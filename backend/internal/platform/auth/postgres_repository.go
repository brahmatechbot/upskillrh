package auth

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) FindUserCredentials(ctx context.Context, email string) (User, string, error) {
	const q = `
SELECT u.id::text, u.email::text, u.full_name, u.status, u.session_version, c.password_hash
FROM users u
JOIN user_password_credentials c ON c.user_id = u.id
WHERE u.email = $1 AND u.deleted_at IS NULL`
	var user User
	var hash string
	if err := r.pool.QueryRow(ctx, q, email).Scan(&user.ID, &user.Email, &user.DisplayName, &user.Status, &user.SessionVersion, &hash); err != nil {
		return User{}, "", err
	}
	return user, hash, nil
}

func (r *PostgresRepository) CreateSession(ctx context.Context, s Session) error {
	const q = `
INSERT INTO auth_sessions (user_id, token_hash, csrf_hash, session_version, remember_me, created_at, last_seen_at, idle_expires_at, absolute_expires_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.pool.Exec(ctx, q, s.UserID, s.TokenHash, s.CSRFHash, s.SessionVersion, s.RememberMe, s.CreatedAt, s.LastSeenAt, s.IdleExpiresAt, s.AbsoluteExpiresAt)
	return err
}

func (r *PostgresRepository) UpdateLastLogin(ctx context.Context, userID string) error {
	_, err := r.pool.Exec(ctx, `UPDATE users SET last_login_at = now(), updated_at = now() WHERE id = $1`, userID)
	if err == pgx.ErrNoRows {
		return nil
	}
	return err
}

func (r *PostgresRepository) CurrentLegalVersions(ctx context.Context) (LegalVersions, error) {
	versions := LegalVersions{}
	q := `SELECT id::text, version FROM policy_versions WHERE policy_type=$1 AND status='active' ORDER BY effective_at DESC LIMIT 1`
	if err := r.pool.QueryRow(ctx, q, "terms_of_use").Scan(&versions.TermsID, &versions.TermsVersion); err != nil {
		return LegalVersions{}, err
	}
	if err := r.pool.QueryRow(ctx, q, "privacy_notice").Scan(&versions.PrivacyID, &versions.PrivacyVersion); err != nil {
		return LegalVersions{}, err
	}
	return versions, nil
}

func (r *PostgresRepository) Segments(ctx context.Context) ([]Segment, error) {
	rows, err := r.pool.Query(ctx, `SELECT code, name FROM industry_segments WHERE status='active' ORDER BY display_order, name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Segment
	for rows.Next() {
		var s Segment
		if err := rows.Scan(&s.Code, &s.Name); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

func (r *PostgresRepository) CreateCandidateRegistration(ctx context.Context, in CandidateRegistrationRecord) (RegistrationResult, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return RegistrationResult{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var userID, profileID string
	if err := tx.QueryRow(ctx, `INSERT INTO users (email, full_name, status, session_version) VALUES ($1,$2,'active',1) RETURNING id::text`, in.Email, in.FullName).Scan(&userID); err != nil {
		return RegistrationResult{}, mapPgError(err)
	}
	if _, err := tx.Exec(ctx, `INSERT INTO user_password_credentials (user_id, password_hash, algorithm) VALUES ($1,$2,'argon2id')`, userID, in.PasswordHash); err != nil {
		return RegistrationResult{}, mapPgError(err)
	}
	if err := tx.QueryRow(ctx, `INSERT INTO candidate_profiles (user_id, status, cpf_encrypted, cpf_blind_index, profile_completion_status) VALUES ($1,'active',$2,$3,'basic_registration') RETURNING id::text`, userID, in.CPFEncrypted, in.CPFBlindIndex).Scan(&profileID); err != nil {
		return RegistrationResult{}, mapPgError(err)
	}
	if err := insertLegalAndConsent(ctx, tx, userID, in.TermsVersionID, in.PrivacyVersionID, in.MarketingGranted, in.RequestID); err != nil {
		return RegistrationResult{}, err
	}
	if _, err := tx.Exec(ctx, `INSERT INTO email_verification_tokens (user_id, token_hash, expires_at) VALUES ($1,$2,$3)`, userID, in.VerificationHash, in.VerificationExp); err != nil {
		return RegistrationResult{}, err
	}
	in.Session.UserID = userID
	if _, err := tx.Exec(ctx, `INSERT INTO auth_sessions (user_id, token_hash, csrf_hash, session_version, remember_me, created_at, last_seen_at, idle_expires_at, absolute_expires_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`, userID, in.Session.TokenHash, in.Session.CSRFHash, in.Session.SessionVersion, in.Session.RememberMe, in.Session.CreatedAt, in.Session.LastSeenAt, in.Session.IdleExpiresAt, in.Session.AbsoluteExpiresAt); err != nil {
		return RegistrationResult{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return RegistrationResult{}, err
	}
	return RegistrationResult{UserID: userID, DisplayName: in.FullName, ContextType: "candidate", EmailVerificationStatus: "pending", NextURL: "/candidate", OrganizationID: "" + profileID[:0]}, nil
}

func (r *PostgresRepository) CreateCompanyRegistration(ctx context.Context, in CompanyRegistrationRecord) (RegistrationResult, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return RegistrationResult{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var userID, orgID string
	if err := tx.QueryRow(ctx, `INSERT INTO users (email, full_name, status, session_version) VALUES ($1,$2,'active',1) RETURNING id::text`, in.Email, in.ResponsibleName).Scan(&userID); err != nil {
		return RegistrationResult{}, mapPgError(err)
	}
	if _, err := tx.Exec(ctx, `INSERT INTO user_password_credentials (user_id, password_hash, algorithm) VALUES ($1,$2,'argon2id')`, userID, in.PasswordHash); err != nil {
		return RegistrationResult{}, mapPgError(err)
	}
	if err := tx.QueryRow(ctx, `INSERT INTO organizations (trade_name, status, employee_range_code, industry_segment_code, other_industry_segment, tax_identifier_encrypted, tax_identifier_blind_index) VALUES ($1,'pending_verification',$2,$3,NULLIF($4,''),$5,$6) RETURNING id::text`, in.TradeName, in.EmployeeRangeCode, in.IndustrySegmentCode, in.OtherIndustrySegment, in.CNPJEncrypted, in.CNPJBlindIndex).Scan(&orgID); err != nil {
		return RegistrationResult{}, mapPgError(err)
	}
	if _, err := tx.Exec(ctx, `INSERT INTO organization_memberships (tenant_id, user_id, status, role_code) VALUES ($1,$2,'active','organization_owner')`, orgID, userID); err != nil {
		return RegistrationResult{}, err
	}
	if err := insertLegalAndConsent(ctx, tx, userID, in.TermsVersionID, in.PrivacyVersionID, in.MarketingGranted, in.RequestID); err != nil {
		return RegistrationResult{}, err
	}
	if _, err := tx.Exec(ctx, `INSERT INTO email_verification_tokens (user_id, token_hash, expires_at) VALUES ($1,$2,$3)`, userID, in.VerificationHash, in.VerificationExp); err != nil {
		return RegistrationResult{}, err
	}
	in.Session.UserID = userID
	if _, err := tx.Exec(ctx, `INSERT INTO auth_sessions (user_id, token_hash, csrf_hash, session_version, remember_me, created_at, last_seen_at, idle_expires_at, absolute_expires_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`, userID, in.Session.TokenHash, in.Session.CSRFHash, in.Session.SessionVersion, in.Session.RememberMe, in.Session.CreatedAt, in.Session.LastSeenAt, in.Session.IdleExpiresAt, in.Session.AbsoluteExpiresAt); err != nil {
		return RegistrationResult{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return RegistrationResult{}, err
	}
	return RegistrationResult{UserID: userID, DisplayName: in.ResponsibleName, ContextType: "company", OrganizationID: orgID, EmailVerificationStatus: "pending", NextURL: "/app"}, nil
}

func insertLegalAndConsent(ctx context.Context, tx pgx.Tx, userID, termsID, privacyID string, marketing bool, requestID string) error {
	if _, err := tx.Exec(ctx, `INSERT INTO user_legal_acceptances (user_id, policy_version_id, acceptance_type, source, request_id) VALUES ($1,$2,'terms_of_use','registration',$3), ($1,$4,'privacy_notice_ack','registration',$3)`, userID, termsID, requestID, privacyID); err != nil {
		return err
	}
	_, err := tx.Exec(ctx, `INSERT INTO marketing_consents (user_id, granted, granted_at, source, policy_version) VALUES ($1,$2,CASE WHEN $2 THEN now() ELSE NULL END,'registration',$3)`, userID, marketing, privacyID)
	return err
}

func (r *PostgresRepository) LoginDestination(ctx context.Context, userID string) (string, error) {
	var exists bool
	if err := r.pool.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM candidate_profiles WHERE user_id=$1 AND deleted_at IS NULL AND status='active')`, userID).Scan(&exists); err == nil && exists {
		return "/candidate", nil
	}
	if err := r.pool.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM organization_memberships WHERE user_id=$1 AND status='active')`, userID).Scan(&exists); err == nil && exists {
		return "/app", nil
	}
	return "/app", nil
}

func (r *PostgresRepository) SessionContext(ctx context.Context, rawToken string) (SessionContext, error) {
	q := `SELECT u.id::text, u.full_name, (u.email_verified_at IS NOT NULL), COALESCE(cp.id::text,''), COALESCE(o.id::text,''), COALESCE(o.trade_name,''), COALESCE(o.status,''), COALESCE(o.employee_range_code,''), COALESCE(seg.name, o.other_industry_segment, '')
FROM auth_sessions s
JOIN users u ON u.id=s.user_id AND u.deleted_at IS NULL AND u.status='active'
LEFT JOIN candidate_profiles cp ON cp.user_id=u.id AND cp.deleted_at IS NULL AND cp.status='active'
LEFT JOIN organization_memberships om ON om.user_id=u.id AND om.status='active'
LEFT JOIN organizations o ON o.id=om.tenant_id AND o.deleted_at IS NULL
LEFT JOIN industry_segments seg ON seg.code=o.industry_segment_code
WHERE s.token_hash=$1 AND s.revoked_at IS NULL AND s.idle_expires_at > now() AND s.absolute_expires_at > now()
LIMIT 1`
	var out SessionContext
	if err := r.pool.QueryRow(ctx, q, hashToken(rawToken)).Scan(&out.UserID, &out.DisplayName, &out.EmailVerified, &out.CandidateProfileID, &out.OrganizationID, &out.OrganizationName, &out.OrganizationStatus, &out.EmployeeRangeCode, &out.IndustrySegmentName); err != nil {
		return SessionContext{}, err
	}
	return out, nil
}

func mapPgError(err error) error {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return err
	}
	c := strings.ToLower(pgErr.ConstraintName)
	switch {
	case strings.Contains(c, "users_email"):
		return ErrDuplicateEmail
	case strings.Contains(c, "candidate_profiles_cpf"):
		return ErrDuplicateCPF
	case strings.Contains(c, "organizations_tax_identifier"):
		return ErrDuplicateCNPJ
	default:
		return err
	}
}

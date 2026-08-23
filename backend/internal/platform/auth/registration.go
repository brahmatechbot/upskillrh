package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

var (
	ErrDuplicateEmail = errors.New("duplicate email")
	ErrDuplicateCPF   = errors.New("duplicate cpf")
	ErrDuplicateCNPJ  = errors.New("duplicate cnpj")
	ErrLegalMissing   = errors.New("legal acceptance missing")
	ErrInvalidCPF     = errors.New("invalid cpf")
	ErrInvalidCNPJ    = errors.New("invalid cnpj")
)

type RegistrationStore interface {
	CreateCandidateRegistration(ctx context.Context, in CandidateRegistrationRecord) (RegistrationResult, error)
	CreateCompanyRegistration(ctx context.Context, in CompanyRegistrationRecord) (RegistrationResult, error)
	CurrentLegalVersions(ctx context.Context) (LegalVersions, error)
	Segments(ctx context.Context) ([]Segment, error)
	SessionContext(ctx context.Context, rawToken string) (SessionContext, error)
}

type LegalVersions struct {
	TermsID, TermsVersion     string
	PrivacyID, PrivacyVersion string
}

type Segment struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type CandidateRegistrationRecord struct {
	FullName         string
	Email            string
	PasswordHash     string
	CPFEncrypted     string
	CPFBlindIndex    string
	TermsVersionID   string
	PrivacyVersionID string
	MarketingGranted bool
	RequestID        string
	Session          Session
	VerificationHash []byte
	VerificationExp  time.Time
}

type CompanyRegistrationRecord struct {
	ResponsibleName      string
	Email                string
	PasswordHash         string
	TradeName            string
	EmployeeRangeCode    string
	IndustrySegmentCode  string
	OtherIndustrySegment string
	CNPJEncrypted        string
	CNPJBlindIndex       string
	TermsVersionID       string
	PrivacyVersionID     string
	MarketingGranted     bool
	RequestID            string
	Session              Session
	VerificationHash     []byte
	VerificationExp      time.Time
}

type RegistrationResult struct {
	UserID                  string    `json:"user_id"`
	DisplayName             string    `json:"display_name"`
	ContextType             string    `json:"context_type"`
	OrganizationID          string    `json:"organization_id,omitempty"`
	EmailVerificationStatus string    `json:"email_verification_status"`
	NextURL                 string    `json:"next_url"`
	SessionRaw              string    `json:"-"`
	SessionExpiresAt        time.Time `json:"-"`
}

type SessionContext struct {
	UserID              string
	DisplayName         string
	EmailVerified       bool
	CandidateProfileID  string
	OrganizationID      string
	OrganizationName    string
	OrganizationStatus  string
	EmployeeRangeCode   string
	IndustrySegmentName string
}

type RegistrationService struct {
	store RegistrationStore
	now   func() time.Time
}

func NewRegistrationService(store RegistrationStore) *RegistrationService {
	return &RegistrationService{store: store, now: time.Now}
}

type RegisterInput struct {
	AccountType          string `json:"account_type"`
	FullName             string `json:"full_name"`
	ResponsibleName      string `json:"responsible_name"`
	Email                string `json:"email"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
	CPF                  string `json:"cpf"`
	TradeName            string `json:"trade_name"`
	EmployeeRangeCode    string `json:"employee_range_code"`
	CNPJ                 string `json:"cnpj"`
	IndustrySegmentCode  string `json:"industry_segment_code"`
	OtherIndustrySegment string `json:"other_industry_segment"`
	AcceptTerms          bool   `json:"accept_terms"`
	AcceptPrivacy        bool   `json:"accept_privacy"`
	MarketingGranted     bool   `json:"marketing_granted"`
}

func (s *RegistrationService) Register(ctx context.Context, in RegisterInput, requestID string) (RegistrationResult, map[string]string, error) {
	if s.store == nil {
		return RegistrationResult{}, nil, ErrAccessUnavailable
	}
	fields := validateRegistrationInput(in)
	if len(fields) > 0 {
		return RegistrationResult{}, fields, nil
	}
	legal, err := s.store.CurrentLegalVersions(ctx)
	if err != nil {
		return RegistrationResult{}, nil, err
	}
	hash, err := HashPassword(in.Password)
	if err != nil {
		return RegistrationResult{}, nil, err
	}
	sessionRaw, session, err := s.newSession("", false)
	if err != nil {
		return RegistrationResult{}, nil, err
	}
	verifyRaw, err := randomToken(32)
	if err != nil {
		return RegistrationResult{}, nil, err
	}
	verificationHash := hashToken(verifyRaw)
	verificationExp := s.now().UTC().Add(48 * time.Hour)

	switch in.AccountType {
	case "candidate":
		cpf := onlyDigits(in.CPF)
		result, err := s.store.CreateCandidateRegistration(ctx, CandidateRegistrationRecord{
			FullName: strings.TrimSpace(in.FullName), Email: normalizeEmail(in.Email), PasswordHash: hash,
			CPFEncrypted: protectIdentifier("cpf", cpf), CPFBlindIndex: blindIndex("cpf", cpf),
			TermsVersionID: legal.TermsID, PrivacyVersionID: legal.PrivacyID, MarketingGranted: in.MarketingGranted,
			RequestID: requestID, Session: session, VerificationHash: verificationHash, VerificationExp: verificationExp,
		})
		if err != nil {
			return RegistrationResult{}, duplicateField(err), err
		}
		result.SessionRaw = sessionRaw
		result.SessionExpiresAt = session.AbsoluteExpiresAt
		return result, nil, nil
	case "company":
		cnpj := onlyDigits(in.CNPJ)
		result, err := s.store.CreateCompanyRegistration(ctx, CompanyRegistrationRecord{
			ResponsibleName: strings.TrimSpace(in.ResponsibleName), Email: normalizeEmail(in.Email), PasswordHash: hash,
			TradeName: strings.TrimSpace(in.TradeName), EmployeeRangeCode: in.EmployeeRangeCode,
			IndustrySegmentCode: in.IndustrySegmentCode, OtherIndustrySegment: strings.TrimSpace(in.OtherIndustrySegment),
			CNPJEncrypted: protectIdentifier("cnpj", cnpj), CNPJBlindIndex: blindIndex("cnpj", cnpj),
			TermsVersionID: legal.TermsID, PrivacyVersionID: legal.PrivacyID, MarketingGranted: in.MarketingGranted,
			RequestID: requestID, Session: session, VerificationHash: verificationHash, VerificationExp: verificationExp,
		})
		if err != nil {
			return RegistrationResult{}, duplicateField(err), err
		}
		result.SessionRaw = sessionRaw
		result.SessionExpiresAt = session.AbsoluteExpiresAt
		return result, nil, nil
	default:
		return RegistrationResult{}, map[string]string{"account_type": "Escolha empresa ou candidato."}, nil
	}
}

func (s *RegistrationService) newSession(userID string, remember bool) (string, Session, error) {
	rawToken, err := randomToken(32)
	if err != nil {
		return "", Session{}, err
	}
	csrfToken, err := randomToken(32)
	if err != nil {
		return "", Session{}, err
	}
	now := s.now().UTC()
	return rawToken, Session{UserID: userID, TokenHash: hashToken(rawToken), CSRFHash: hashToken(csrfToken), SessionVersion: 1, RememberMe: remember, CreatedAt: now, LastSeenAt: now, IdleExpiresAt: now.Add(12 * time.Hour), AbsoluteExpiresAt: now.Add(24 * time.Hour)}, nil
}

func validateRegistrationInput(in RegisterInput) map[string]string {
	fields := map[string]string{}
	if in.AccountType != "candidate" && in.AccountType != "company" {
		fields["account_type"] = "Escolha empresa ou candidato."
	}
	if normalizeEmail(in.Email) == "" || len(in.Email) > 254 || !strings.Contains(in.Email, "@") || strings.ContainsAny(in.Email, " \t\n\r") {
		fields["email"] = "Informe um e-mail válido."
	}
	if len(in.Password) < 8 || len(in.Password) > 1024 || !utf8.ValidString(in.Password) {
		fields["password"] = "A senha deve ter pelo menos 8 caracteres."
	}
	if in.Password != in.PasswordConfirmation {
		fields["password_confirmation"] = "A confirmação deve ser idêntica à senha."
	}
	if !in.AcceptTerms {
		fields["accept_terms"] = "Aceite os Termos de Uso."
	}
	if !in.AcceptPrivacy {
		fields["accept_privacy"] = "Confirme ciência da Política de Privacidade."
	}
	if in.AccountType == "candidate" {
		if strings.TrimSpace(in.FullName) == "" {
			fields["full_name"] = "Informe seu nome completo."
		}
		if !validCPF(in.CPF) {
			fields["cpf"] = "Informe um CPF válido."
		}
	}
	if in.AccountType == "company" {
		if strings.TrimSpace(in.ResponsibleName) == "" {
			fields["responsible_name"] = "Informe o nome do responsável."
		}
		if strings.TrimSpace(in.TradeName) == "" {
			fields["trade_name"] = "Informe o nome da empresa."
		}
		if !validEmployeeRange(in.EmployeeRangeCode) {
			fields["employee_range_code"] = "Selecione uma faixa de funcionários."
		}
		if !validCNPJ(in.CNPJ) {
			fields["cnpj"] = "Informe um CNPJ válido."
		}
		if !validSegment(in.IndustrySegmentCode) {
			fields["industry_segment_code"] = "Selecione um segmento."
		}
		if in.IndustrySegmentCode == "other" && strings.TrimSpace(in.OtherIndustrySegment) == "" {
			fields["other_industry_segment"] = "Descreva o segmento."
		}
		if len(strings.TrimSpace(in.OtherIndustrySegment)) > 80 {
			fields["other_industry_segment"] = "Use até 80 caracteres."
		}
	}
	return fields
}

func duplicateField(err error) map[string]string {
	switch {
	case errors.Is(err, ErrDuplicateEmail):
		return map[string]string{"email": "Este e-mail já possui conta. Entre com a conta existente."}
	case errors.Is(err, ErrDuplicateCPF):
		return map[string]string{"cpf": "CPF já associado a um perfil ativo."}
	case errors.Is(err, ErrDuplicateCNPJ):
		return map[string]string{"cnpj": "CNPJ já associado a uma organização ativa."}
	default:
		return nil
	}
}

func normalizeEmail(v string) string { return strings.ToLower(strings.TrimSpace(v)) }
func onlyDigits(v string) string     { return regexp.MustCompile(`\D`).ReplaceAllString(v, "") }
func validEmployeeRange(v string) bool {
	switch v {
	case "solo", "2_10", "11_50", "51_200", "201_500", "501_1000", "1001_5000", "5000_plus":
		return true
	}
	return false
}
func validSegment(v string) bool {
	switch v {
	case "technology_software", "financial_services_banks", "insurance", "telecommunications", "media_entertainment", "consulting_professional_services", "other":
		return true
	}
	return false
}

func validCPF(v string) bool {
	d := onlyDigits(v)
	if len(d) != 11 {
		return false
	}
	all := true
	for i := 1; i < len(d); i++ {
		if d[i] != d[0] {
			all = false
			break
		}
	}
	if all {
		return false
	}
	calc := func(n int) byte {
		sum := 0
		for i := 0; i < n; i++ {
			sum += int(d[i]-'0') * (n + 1 - i)
		}
		r := sum % 11
		if r < 2 {
			return '0'
		}
		return byte('0' + (11 - r))
	}
	return d[9] == calc(9) && d[10] == calc(10)
}

func validCNPJ(v string) bool {
	d := onlyDigits(v)
	if len(d) != 14 {
		return false
	}
	all := true
	for i := 1; i < len(d); i++ {
		if d[i] != d[0] {
			all = false
			break
		}
	}
	if all {
		return false
	}
	calc := func(weights []int) byte {
		sum := 0
		for i, w := range weights {
			sum += int(d[i]-'0') * w
		}
		r := sum % 11
		if r < 2 {
			return '0'
		}
		return byte('0' + (11 - r))
	}
	return d[12] == calc([]int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}) && d[13] == calc([]int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2})
}

func piiSecret() []byte {
	secret := os.Getenv("UPSKILLRH_PII_SECRET")
	if secret == "" {
		secret = "upskillrh-local-development-secret"
	}
	return []byte(secret)
}
func blindIndex(kind, value string) string {
	mac := hmac.New(sha256.New, piiSecret())
	_, _ = mac.Write([]byte(kind + ":" + value))
	return hex.EncodeToString(mac.Sum(nil))
}
func protectIdentifier(kind, value string) string {
	return fmt.Sprintf("protected:v1:%s:%s", kind, blindIndex(kind+":cipher", value))
}

func sessionCookie(r *http.Request) (string, bool) {
	c, err := r.Cookie("upskillrh_session")
	return func() string {
		if err != nil {
			return ""
		}
		return c.Value
	}(), err == nil && c.Value != ""
}

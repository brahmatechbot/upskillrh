ALTER TABLE users ADD COLUMN IF NOT EXISTS email_verified_at TIMESTAMPTZ NULL;

CREATE TABLE IF NOT EXISTS industry_segments (
  code TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  display_order INTEGER NOT NULL,
  status TEXT NOT NULL DEFAULT 'active',
  CONSTRAINT industry_segments_status_check CHECK (status IN ('active', 'inactive'))
);

INSERT INTO industry_segments (code, name, display_order, status) VALUES
  ('technology_software', 'Tecnologia e Software', 10, 'active'),
  ('financial_services_banks', 'Serviços Financeiros e Bancos', 20, 'active'),
  ('insurance', 'Seguros', 30, 'active'),
  ('telecommunications', 'Telecomunicações', 40, 'active'),
  ('media_entertainment', 'Mídia e Entretenimento', 50, 'active'),
  ('consulting_professional_services', 'Consultoria e Serviços Profissionais', 60, 'active'),
  ('other', 'Outro', 70, 'active')
ON CONFLICT (code) DO UPDATE SET name = EXCLUDED.name, display_order = EXCLUDED.display_order, status = EXCLUDED.status;

CREATE TABLE IF NOT EXISTS organizations (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  trade_name TEXT NOT NULL,
  status TEXT NOT NULL,
  employee_range_code TEXT NOT NULL,
  industry_segment_code TEXT NOT NULL REFERENCES industry_segments(code),
  other_industry_segment TEXT NULL,
  tax_identifier_encrypted TEXT NOT NULL,
  tax_identifier_blind_index TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  deleted_at TIMESTAMPTZ NULL,
  CONSTRAINT organizations_status_check CHECK (status IN ('pending_verification', 'active', 'suspended', 'disabled')),
  CONSTRAINT organizations_employee_range_check CHECK (employee_range_code IN ('solo','2_10','11_50','51_200','201_500','501_1000','1001_5000','5000_plus')),
  CONSTRAINT organizations_other_segment_check CHECK ((industry_segment_code <> 'other') OR (other_industry_segment IS NOT NULL AND length(other_industry_segment) BETWEEN 1 AND 80))
);

CREATE UNIQUE INDEX IF NOT EXISTS organizations_tax_identifier_active_uq ON organizations (tax_identifier_blind_index) WHERE deleted_at IS NULL AND status <> 'disabled';

CREATE TABLE IF NOT EXISTS organization_memberships (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID NOT NULL REFERENCES organizations(id),
  user_id UUID NOT NULL REFERENCES users(id),
  status TEXT NOT NULL,
  role_code TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CONSTRAINT organization_memberships_status_check CHECK (status IN ('active', 'inactive', 'suspended')),
  CONSTRAINT organization_memberships_role_check CHECK (role_code IN ('organization_owner', 'organization_admin', 'member'))
);

CREATE UNIQUE INDEX IF NOT EXISTS organization_memberships_user_tenant_uq ON organization_memberships (tenant_id, user_id);

CREATE TABLE IF NOT EXISTS candidate_profiles (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id),
  status TEXT NOT NULL,
  cpf_encrypted TEXT NOT NULL,
  cpf_blind_index TEXT NOT NULL,
  profile_completion_status TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  deleted_at TIMESTAMPTZ NULL,
  CONSTRAINT candidate_profiles_status_check CHECK (status IN ('active', 'suspended', 'disabled')),
  CONSTRAINT candidate_profiles_completion_check CHECK (profile_completion_status IN ('basic_registration'))
);

CREATE UNIQUE INDEX IF NOT EXISTS candidate_profiles_user_active_uq ON candidate_profiles (user_id) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX IF NOT EXISTS candidate_profiles_cpf_active_uq ON candidate_profiles (cpf_blind_index) WHERE deleted_at IS NULL AND status <> 'disabled';

CREATE TABLE IF NOT EXISTS policy_versions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  policy_type TEXT NOT NULL,
  version TEXT NOT NULL,
  content_url TEXT NOT NULL,
  effective_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  status TEXT NOT NULL,
  CONSTRAINT policy_versions_type_check CHECK (policy_type IN ('terms_of_use', 'privacy_notice')),
  CONSTRAINT policy_versions_status_check CHECK (status IN ('active', 'retired'))
);

CREATE UNIQUE INDEX IF NOT EXISTS policy_versions_active_type_uq ON policy_versions (policy_type) WHERE status = 'active';

INSERT INTO policy_versions (policy_type, version, content_url, effective_at, status)
VALUES
  ('terms_of_use', '2026-08-23', '/termos-de-uso', now(), 'active'),
  ('privacy_notice', '2026-08-23', '/politica-de-privacidade', now(), 'active')
ON CONFLICT DO NOTHING;

CREATE TABLE IF NOT EXISTS user_legal_acceptances (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id),
  policy_version_id UUID NOT NULL REFERENCES policy_versions(id),
  acceptance_type TEXT NOT NULL,
  accepted_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  source TEXT NOT NULL,
  request_id TEXT NOT NULL,
  CONSTRAINT user_legal_acceptances_type_check CHECK (acceptance_type IN ('terms_of_use', 'privacy_notice_ack'))
);

CREATE TABLE IF NOT EXISTS marketing_consents (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id),
  granted BOOLEAN NOT NULL,
  granted_at TIMESTAMPTZ NULL,
  revoked_at TIMESTAMPTZ NULL,
  source TEXT NOT NULL,
  policy_version TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS email_verification_tokens (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id),
  token_hash BYTEA NOT NULL UNIQUE,
  expires_at TIMESTAMPTZ NOT NULL,
  used_at TIMESTAMPTZ NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

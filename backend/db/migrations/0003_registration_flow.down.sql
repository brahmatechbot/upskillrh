DROP TABLE IF EXISTS email_verification_tokens;
DROP TABLE IF EXISTS marketing_consents;
DROP TABLE IF EXISTS user_legal_acceptances;
DROP TABLE IF EXISTS policy_versions;
DROP TABLE IF EXISTS candidate_profiles;
DROP TABLE IF EXISTS organization_memberships;
DROP TABLE IF EXISTS organizations;
DROP TABLE IF EXISTS industry_segments;
ALTER TABLE users DROP COLUMN IF EXISTS email_verified_at;

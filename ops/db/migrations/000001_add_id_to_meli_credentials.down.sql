ALTER TABLE mercadolivre_credentials
DROP CONSTRAINT IF EXISTS mercadolivre_credentials_pkey,
DROP COLUMN IF EXISTS id;
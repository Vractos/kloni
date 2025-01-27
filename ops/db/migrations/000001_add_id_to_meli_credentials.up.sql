ALTER TABLE mercadolivre_credentials ADD COLUMN id UUID;

UPDATE mercadolivre_credentials SET id = gen_random_uuid() WHERE id IS NULL;

ALTER TABLE mercadolivre_credentials ALTER COLUMN id SET NOT NULL;

ALTER TABLE mercadolivre_credentials
DROP CONSTRAINT IF EXISTS mercadolivre_credentials_pkey,
ADD PRIMARY KEY (id);
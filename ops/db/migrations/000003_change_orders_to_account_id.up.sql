ALTER TABLE orders ADD COLUMN account_id UUID;

UPDATE orders o
SET account_id = mc.id
FROM mercadolivre_credentials mc
WHERE o.store_id = mc.owner_id;

ALTER TABLE orders ALTER COLUMN account_id SET NOT NULL;

ALTER TABLE orders
DROP CONSTRAINT IF EXISTS orders_store_id_fkey;

ALTER TABLE orders
ADD CONSTRAINT orders_account_id_fkey
FOREIGN KEY (account_id) REFERENCES mercadolivre_credentials(id);

ALTER TABLE orders DROP COLUMN store_id;

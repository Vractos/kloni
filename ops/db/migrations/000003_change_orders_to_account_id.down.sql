ALTER TABLE orders ADD COLUMN store_id UUID;

ALTER TABLE orders
DROP CONSTRAINT IF EXISTS orders_account_id_fkey;

ALTER TABLE orders
ADD CONSTRAINT orders_store_id_fkey
FOREIGN KEY (store_id) REFERENCES mercadolivre_credentials(owner_id);

UPDATE orders o
SET store_id = mc.owner_id
FROM mercadolivre_credentials mc
WHERE o.account_id = mc.id;

ALTER TABLE orders ALTER COLUMN store_id SET NOT NULL;

ALTER TABLE orders DROP COLUMN account_id;

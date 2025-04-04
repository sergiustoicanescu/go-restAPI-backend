ALTER TABLE customers
ADD CONSTRAINT unique_user UNIQUE (user_id);

ALTER TABLE orders
DROP CONSTRAINT IF EXISTS orders_status_check;

ALTER TABLE orders
ADD CONSTRAINT orders_status_check CHECK (status IN ('pending', 'completed', 'cancelled'));

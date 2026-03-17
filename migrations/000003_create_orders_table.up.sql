CREATE TABLE IF NOT EXISTS orders (
  id BIGSERIAL PRIMARY KEY,
  customer_id BIGINT NOT NULL,
  status text NOT NULL,
  total_amount_in_cents integer NOT NULL,
  created_at TIMESTAMPTZ DEFAULT now(),
  version integer NOT NULL DEFAULT 1,

  CONSTRAINT fk_customer_id FOREIGN KEY (customer_id) REFERENCES users(id)
);


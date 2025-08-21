-- +goose Up
CREATE TABLE IF NOT EXISTS payments (
  order_uid     TEXT PRIMARY KEY
                REFERENCES orders(order_uid) ON DELETE CASCADE,
  transaction   TEXT NOT NULL,
  request_id    TEXT,
  currency      TEXT NOT NULL,
  provider      TEXT NOT NULL,
  amount        INTEGER NOT NULL,
  payment_dt    BIGINT NOT NULL,
  bank          TEXT NOT NULL,
  delivery_cost INTEGER NOT NULL,
  goods_total   INTEGER NOT NULL,
  custom_fee    INTEGER NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS payments;

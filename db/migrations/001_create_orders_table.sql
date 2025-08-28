-- +goose Up
CREATE TABLE IF NOT EXISTS orders (
  order_uid          TEXT PRIMARY KEY,
  track_number       TEXT NOT NULL,
  entry              TEXT NOT NULL,
  locale             TEXT NOT NULL,
  internal_signature TEXT,
  customer_id        TEXT NOT NULL,
  delivery_service   TEXT NOT NULL,
  shardkey           TEXT NOT NULL,
  sm_id              INTEGER  NOT NULL,
  date_created       TIMESTAMPTZ NOT NULL,
  oof_shard          TEXT NOT NULL,
  created_at         TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_orders_created_at_desc
  ON orders (created_at DESC);

-- +goose Down
DROP TABLE IF EXISTS orders;
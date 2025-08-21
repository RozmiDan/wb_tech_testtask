-- +goose Up
CREATE TABLE IF NOT EXISTS items (
  order_uid    TEXT NOT NULL
               REFERENCES orders(order_uid) ON DELETE CASCADE,
  chrt_id      BIGINT NOT NULL,
  track_number TEXT NOT NULL,
  price        INTEGER NOT NULL,
  rid          TEXT NOT NULL,
  name         TEXT NOT NULL,
  sale         INTEGER NOT NULL,
  size         TEXT NOT NULL,
  total_price  INTEGER NOT NULL,
  nm_id        BIGINT NOT NULL,
  brand        TEXT NOT NULL,
  status       INTEGER NOT NULL,
);

-- +goose Down
DROP TABLE IF EXISTS items;

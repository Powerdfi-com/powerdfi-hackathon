CREATE TABLE IF NOT EXISTS trades (
    id               UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    buy_order_id          UUID   REFERENCES orders (id) ON DELETE SET NULL,
    sell_order_id          UUID   REFERENCES orders (id) ON DELETE SET NULL,
    price            NUMERIC          NOT NULL,
    quantity         INTEGER          NOT NULL,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);
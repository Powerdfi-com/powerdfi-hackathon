CREATE TABLE IF NOT EXISTS orders (
    id               UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    user_id          UUID   REFERENCES users (id) ON DELETE SET NULL,
    asset_id          UUID REFERENCES assets (id) ON DELETE CASCADE,
  
    type    CITEXT,
    kind    CITEXT,
    status    CITEXT,
    price            NUMERIC          NOT NULL,
    quantity         INTEGER          NOT NULL,
    inital_quantity  INTEGER          NOT NULL,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT now()

);
CREATE TABLE IF NOT EXISTS admins
(
    id          UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    email          CITEXT UNIQUE,
    name          TEXT,
    password_hash BYTEA NOT NULL,
    role_mask SMALLINT NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ      NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ      NOT NULL DEFAULT now()
);
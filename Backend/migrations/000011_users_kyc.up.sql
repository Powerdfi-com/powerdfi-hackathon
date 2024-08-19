CREATE TABLE IF NOT EXISTS users_kyc (
    id               UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
   user_id          UUID   REFERENCES users (id) ON DELETE SET NULL,
    url          TEXT             NOT NULL DEFAULT '',
    platform          CITEXT,
    reference_id    CITEXT,
    status    CITEXT,
    comment    CITEXT,
    UNIQUE (user_id)
);
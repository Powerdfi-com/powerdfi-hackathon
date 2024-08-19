-- nonces table for authentication
CREATE TABLE IF NOT EXISTS nonces
(
    -- a foreign key isn't used for the user_address as
    -- the account isn't created before the nonce is generated
    user_address CITEXT      NOT NULL UNIQUE,
    message      TEXT        NOT NULL,

    -- action should be 'create' or 'update'
    action       TEXT        NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);
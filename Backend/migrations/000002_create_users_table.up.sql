CREATE TABLE IF NOT EXISTS users
(
    id          UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    address        CITEXT           NOT NULL UNIQUE,
    account_id     CITEXT           NOT NULL UNIQUE,
    email          CITEXT UNIQUE,
    username          CITEXT UNIQUE,
    first_name     CITEXT,
    last_name      CITEXT,
    bio         TEXT             NOT NULL DEFAULT '',
    website     TEXT             NOT NULL DEFAULT '',
    twitter     TEXT             NOT NULL DEFAULT '',
    discord     TEXT             NOT NULL DEFAULT '',
    avatar      TEXT             NOT NULL DEFAULT '',
    
    public_key     CITEXT,
    encrypted_private_key BYTEA,
    
 
    is_verified BOOLEAN          NOT NULL DEFAULT FALSE,
    is_active BOOLEAN          NOT NULL DEFAULT FALSE,

    kyc_registered_date  TIMESTAMPTZ      ,
    created_at  TIMESTAMPTZ      NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ      NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS assets
(
    id            UUID PRIMARY KEY                                   NOT NULL DEFAULT gen_random_uuid(),
    token_id      CITEXT                                     NOT NULL,
    name          TEXT                                               NOT NULL DEFAULT '',
    symbol          CITEXT                                               NOT NULL DEFAULT '',
    category_id      SMALLINT       REFERENCES categories (id),
    blockchain       TEXT             NOT NULL DEFAULT '',

    creator_id    UUID REFERENCES users (id)                         NOT NULL,
    metadata_url  TEXT                                               NOT NULL DEFAULT '',
    
    urls     JSON                                               NOT NULL DEFAULT '[]',
    legal_docs     JSON                                               NOT NULL DEFAULT '[]',
    issuance_docs     JSON                                               NOT NULL DEFAULT '[]',
    signatories     JSON                                               NOT NULL DEFAULT '[]',
   
    description   TEXT                                               NOT NULL DEFAULT '',
    total_supply  INTEGER                                            NOT NULL DEFAULT 1,
    serial_number  INTEGER,


    properties    JSONB                                              NOT NULL DEFAULT '{}',
    status TEXT NOT NULL DEFAULT '',

    is_verified BOOLEAN          NOT NULL DEFAULT FALSE,
    is_minted BOOLEAN          NOT NULL DEFAULT FALSE,

    executed_at    TIMESTAMPTZ                                        ,
    expires_at    TIMESTAMPTZ                                        ,
    created_at    TIMESTAMPTZ                                        NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ                                        NOT NULL DEFAULT now()
);
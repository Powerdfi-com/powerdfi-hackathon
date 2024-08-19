-- CREATE TABLE IF NOT EXISTS activities
-- (
--     id               UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
--     action           TEXT             NOT NULL,               -- list, delist, bid, bid_cancel, sale
--     blockchain       TEXT             NOT NULL,
    
--     asset_id          UUID REFERENCES assets (id) ON DELETE CASCADE,
--     listing_id UUID REFERENCES listings (id) ON DELETE CASCADE,

--     from_address     CITEXT,
--     to_address       CITEXT,
--     price            NUMERIC,
--     currency         CITEXT,
--     quantity         INTEGER,
--     block_number     INTEGER,
--     transaction_hash TEXT,
--     occurred_at      TIMESTAMPTZ      NOT NULL DEFAULT now(), -- used for activities retrieved from the blockchain
--     created_at       TIMESTAMPTZ      NOT NULL DEFAULT now()
-- )


CREATE TABLE IF NOT EXISTS activities
(
    id               UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    action           TEXT             NOT NULL,               -- list, delist, bid, bid_cancel, sale
   
    asset_id          UUID REFERENCES assets (id) ON DELETE CASCADE,


    from_user_id     UUID REFERENCES users (id),
    to_user_id       UUID REFERENCES users (id),
  
    price            NUMERIC,
    currency         CITEXT,
    quantity         INTEGER,

    created_at       TIMESTAMPTZ      NOT NULL DEFAULT now()
)



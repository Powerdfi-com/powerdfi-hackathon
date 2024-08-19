CREATE TABLE IF NOT EXISTS listings
(
    id               UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    type             CITEXT           NOT NULL, -- fixed or auction
    user_id          UUID             NOT NULL REFERENCES users (id),
    asset_id          UUID REFERENCES assets (id) ON DELETE CASCADE,
  
    price            NUMERIC          NOT NULL,
    min_invest_amount      NUMERIC,
    max_invest_amount      NUMERIC,
    min__raise_amount      NUMERIC,
    max_raise_amount       NUMERIC,
   
    currency         JSON                                               NOT NULL DEFAULT '[]',
    quantity         INTEGER          NOT NULL,
   
    start_date       TIMESTAMPTZ,
    end_date         TIMESTAMPTZ,
     
    is_active BOOLEAN          NOT NULL DEFAULT FALSE,
    is_cancelled BOOLEAN          NOT NULL DEFAULT FALSE,

    created_at       TIMESTAMPTZ      NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ      NOT NULL DEFAULT now()
)
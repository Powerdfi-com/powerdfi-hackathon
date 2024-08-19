CREATE TABLE IF NOT EXISTS favourites
(
    user_id          UUID   NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    asset_id          UUID REFERENCES assets (id) ON DELETE CASCADE,
    
    UNIQUE (user_id, asset_id)
);
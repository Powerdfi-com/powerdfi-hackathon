

CREATE TABLE IF NOT EXISTS asset_views
(
    user_id          UUID   REFERENCES users (id) ON DELETE SET NULL,
    asset_id          UUID REFERENCES assets (id) ON DELETE CASCADE,

    UNIQUE (user_id, asset_id)
);
CREATE TABLE IF NOT EXISTS asset_owners (
   user_id          UUID   REFERENCES users (id) ON DELETE SET NULL,
    asset_id          UUID REFERENCES assets (id) ON DELETE CASCADE,
    serial_numbers     JSON                                               NOT NULL DEFAULT '[]',
    UNIQUE (user_id, asset_id)
);
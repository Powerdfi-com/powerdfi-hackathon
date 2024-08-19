CREATE TABLE IF NOT EXISTS notification_prefs
(
    id        UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    user_id   UUID             NOT NULL UNIQUE REFERENCES users (id) ON DELETE CASCADE,

    -- whitelist holds a list of values which correspond
    -- to the IDs of allowed notification types
    whitelist SMALLINT[]       NOT NULL DEFAULT '{}'
);
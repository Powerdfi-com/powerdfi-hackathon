CREATE TABLE IF NOT EXISTS user_notifications
(
    id          UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    type        SMALLINT         NOT NULL REFERENCES notification_types (id),
    user_id     UUID             NOT NULL REFERENCES users (id) ON DELETE CASCADE,

    -- for sales
    -- activity_id UUID             REFERENCES activities (id) ON DELETE SET NULL,
    data     JSON           NOT NULL DEFAULT '{}',
    viewed     BOOLEAN          NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ      NOT NULL DEFAULT now()
);
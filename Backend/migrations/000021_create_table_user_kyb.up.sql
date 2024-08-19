CREATE TABLE IF NOT EXISTS users_kyb (
    id               UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
   user_id          UUID   REFERENCES users (id) ON DELETE SET NULL,
   
    platform          CITEXT,
    reference_id    CITEXT,
    status    CITEXT,
    comment    CITEXT,
    certificate_of_inc    TEXT,
    company_name    TEXT,
    company_location    TEXT,
    company_address    TEXT,
    company_reg_no    TEXT,
    UNIQUE (user_id)
);
ALTER TABLE accounts DROP COLUMN password_hash;

CREATE TABLE account_opaque (
    account_uuid UUID PRIMARY KEY,
    registration_record bytea,
    credential_identifier bytea UNIQUE,
    created_at DATE NOT NULL,
    updated_at DATE NOT NULL
);

CREATE TABLE sessions (
    id UUID PRIMARY KEY,
    account_uuid UUID,
    created_at DATE NOT NULL,
    last_seen_at DATE NOT NULL
);

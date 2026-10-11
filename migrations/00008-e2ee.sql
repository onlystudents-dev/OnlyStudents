CREATE TABLE account_crypto (
    account_uuid UUID PRIMARY KEY,
    wrapped_mk_password BYTEA NOT NULL,
    wrapped_mk_recovery BYTEA NOT NULL,
    recovery_key_hash   BYTEA NOT NULL,
    recovery_key_salt   BYTEA NOT NULL,
    recovery_key_params TEXT  NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

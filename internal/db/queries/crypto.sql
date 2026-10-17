-- name: GetAccountCrypto :one
SELECT * FROM account_crypto WHERE account_uuid = $1;

-- name: CreateAccountCrypto :exec
INSERT INTO account_crypto (account_uuid, wrapped_mk_password, wrapped_mk_recovery, recovery_key_hash, recovery_key_salt, recovery_key_params, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, now(), now());

-- name: UpdateAccountCrypto :execrows
UPDATE account_crypto SET wrapped_mk_password = $1, wrapped_mk_recovery = $2, recovery_key_hash = $3, recovery_key_salt = $4, recovery_key_params = $5, updated_at = now() WHERE account_uuid = $6;

-- name: GetOpaqueRecord :one
SELECT * FROM account_opaque WHERE account_uuid = $1;

-- name: CreateOpaqueRecord :exec
INSERT INTO account_opaque (account_uuid, registration_record, credential_identifier, created_at, updated_at)
VALUES ($1, $2, $3, current_date, current_date)
ON CONFLICT (account_uuid) DO UPDATE
SET registration_record = EXCLUDED.registration_record,
    credential_identifier = EXCLUDED.credential_identifier,
    updated_at = current_date;

-- name: UpdateOpaqueRecord :exec
UPDATE account_opaque SET registration_record = $1, credential_identifier = $2, updated_at = current_date WHERE account_uuid = $3;

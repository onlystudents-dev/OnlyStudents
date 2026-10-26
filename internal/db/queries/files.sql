-- name: UploadFile :execrows
INSERT INTO files (account_uuid, type, size, original_name, uuid, created_at) VALUES($1, $2, $3, $4, $5, now());

-- name: DeleteFile :execrows
DELETE FROM files WHERE uuid = $1 AND account_uuid = $2;

-- name: GetFile :one
SELECT type, size, original_name, uuid, created_at FROM files WHERE uuid = $1;

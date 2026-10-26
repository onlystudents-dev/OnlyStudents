-- +goose Up
SELECT setval(pg_get_serial_sequence('files', 'id'), GREATEST((SELECT COALESCE(MAX(id), 1) FROM files), 1));

-- +goose Down
SELECT 1;

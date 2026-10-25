-- +goose Up
ALTER TABLE homework_submissions ADD COLUMN file_uuid UUID REFERENCES files(uuid);
UPDATE homework_submissions hs SET file_uuid = (SELECT uuid FROM files WHERE id = hs.file_id);
ALTER TABLE homework_submissions DROP COLUMN file_id;

ALTER TABLE absence_reason ADD COLUMN file_uuid UUID REFERENCES files(uuid);
UPDATE absence_reason ar SET file_uuid = (SELECT uuid FROM files WHERE id = ar.file_id);
ALTER TABLE absence_reason DROP COLUMN file_id;

-- +goose Down
ALTER TABLE homework_submissions ADD COLUMN file_id BIGINT DEFAULT NULL REFERENCES files (id);
UPDATE homework_submissions hs SET file_id = (SELECT id FROM files WHERE uuid = hs.file_uuid);
ALTER TABLE homework_submissions DROP COLUMN file_uuid;

ALTER TABLE absence_reason ADD COLUMN file_id BIGINT DEFAULT NULL REFERENCES files (id);
UPDATE absence_reason ar SET file_id = (SELECT id FROM files WHERE uuid = ar.file_uuid);
ALTER TABLE absence_reason DROP COLUMN file_uuid;

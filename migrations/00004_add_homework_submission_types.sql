-- +goose Up
ALTER TABLE homework_submissions ADD COLUMN type VARCHAR(512) NOT NULL DEFAULT 'note';
ALTER TABLE homework_submissions ADD COLUMN note TEXT DEFAULT NULL;
ALTER TABLE homework_submissions ADD COLUMN file_id BIGINT DEFAULT NULL REFERENCES files (id);

UPDATE homework_submissions SET note = COALESCE(content, '') WHERE note IS NULL;

ALTER TABLE homework_submissions DROP COLUMN content;

ALTER TABLE homework_submissions ADD CONSTRAINT chk_homework_submission_type CHECK (
    (type = 'file' AND file_id IS NOT NULL AND note IS NULL) OR
    (type = 'note' AND note IS NOT NULL AND file_id IS NULL)
);

-- +goose Down
ALTER TABLE homework_submissions DROP CONSTRAINT chk_homework_submission_type;

ALTER TABLE homework_submissions ADD COLUMN content TEXT DEFAULT NULL;
UPDATE homework_submissions SET content = note WHERE note IS NOT NULL;

ALTER TABLE homework_submissions DROP COLUMN type;
ALTER TABLE homework_submissions DROP COLUMN note;
ALTER TABLE homework_submissions DROP COLUMN file_id;

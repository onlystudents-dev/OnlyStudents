-- +goose Up
ALTER TABLE homework_submissions DROP CONSTRAINT chk_homework_submission_type;
ALTER TABLE homework_submissions DROP COLUMN type;

ALTER TABLE absence_reason DROP CONSTRAINT chk_absence_reason_type;
ALTER TABLE absence_reason DROP COLUMN type;

-- +goose Down
ALTER TABLE homework_submissions ADD COLUMN type VARCHAR(512) NOT NULL DEFAULT 'note';

UPDATE homework_submissions SET type = 'file' WHERE file_id IS NOT NULL;
UPDATE homework_submissions SET note = NULL WHERE file_id IS NOT NULL;

UPDATE homework_submissions SET note = COALESCE(note, '') WHERE type = 'note';

ALTER TABLE homework_submissions ADD CONSTRAINT chk_homework_submission_type CHECK (
    (type = 'file' AND file_id IS NOT NULL AND note IS NULL) OR
    (type = 'note' AND note IS NOT NULL AND file_id IS NULL)
);

ALTER TABLE absence_reason ADD COLUMN type VARCHAR(512) NOT NULL DEFAULT 'note';

UPDATE absence_reason SET type = 'file' WHERE file_id IS NOT NULL;
UPDATE absence_reason SET note = NULL WHERE file_id IS NOT NULL;

UPDATE absence_reason SET note = COALESCE(note, '') WHERE type = 'note';

ALTER TABLE absence_reason ADD CONSTRAINT chk_absence_reason_type CHECK (
    (type = 'file' AND file_id IS NOT NULL AND note IS NULL) OR
    (type = 'note' AND note IS NOT NULL AND file_id IS NULL)
);

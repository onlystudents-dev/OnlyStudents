-- +goose Up
CREATE TABLE absence_reason (
    id BIGSERIAL PRIMARY KEY,
    absence_id INT UNIQUE NOT NULL REFERENCES absences (id) ON DELETE CASCADE,
    submitted_at DATE NOT NULL,
    type VARCHAR(512) NOT NULL,
    note VARCHAR(512) DEFAULT NULL,
    file_id BIGINT DEFAULT NULL REFERENCES files (id) ON DELETE CASCADE,
    CONSTRAINT chk_absence_reason_type CHECK (
        (type = 'file' AND file_id IS NOT NULL AND note IS NULL) OR
        (type = 'note' AND note IS NOT NULL AND file_id IS NULL)
    )
);

-- +goose Down
DROP TABLE absence_reason;

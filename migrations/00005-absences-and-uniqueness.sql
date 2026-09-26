ALTER TABLE absences
    ADD COLUMN class_subjects_id INT DEFAULT NULL REFERENCES class_subjects (id) ON DELETE RESTRICT;

CREATE INDEX idx_absences_class_subjects_id ON absences (class_subjects_id);

ALTER TABLE class_subjects
    DROP CONSTRAINT uq_class_subjects;

CREATE UNIQUE INDEX uq_class_subjects_subject
    ON class_subjects (class_id, subject_id)
    WHERE NOT custom_subject;

CREATE UNIQUE INDEX uq_class_subjects_custom
    ON class_subjects (class_id, custom_subject_id)
    WHERE custom_subject;

ALTER TABLE homework
    ADD COLUMN assigned_date DATE DEFAULT NULL;

CREATE INDEX idx_homework_assigned
    ON homework (class_subjects_id, COALESCE(assigned_date, due_date));

CREATE UNIQUE INDEX uq_exams_class_subject_date
    ON exams (class_subjects_id, date);
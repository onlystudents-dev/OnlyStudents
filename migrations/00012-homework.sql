CREATE TABLE homework (
    id BIGSERIAL PRIMARY KEY,
    class_subjects_id INT NOT NULL REFERENCES class_subjects (id) ON DELETE CASCADE,
    teacher_id INT NOT NULL REFERENCES teachers (id) ON DELETE RESTRICT,
    title VARCHAR(256) NOT NULL,
    description TEXT DEFAULT NULL,
    due_date DATE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE homework_submissions (
    id BIGSERIAL PRIMARY KEY,
    homework_id INT NOT NULL REFERENCES homework (id) ON DELETE CASCADE,
    student_id INT NOT NULL REFERENCES students (id) ON DELETE CASCADE,
    content TEXT DEFAULT NULL,
    submitted_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    graded_value SMALLINT DEFAULT NULL,
    CONSTRAINT chk_homework_submissions_value CHECK (graded_value BETWEEN 1 AND 5),
    CONSTRAINT uq_homework_submissions UNIQUE (homework_id, student_id)
);

CREATE INDEX idx_homework_class_subjects_id ON homework (class_subjects_id);
CREATE INDEX idx_homework_teacher_id ON homework (teacher_id);
CREATE INDEX idx_homework_due_date ON homework (due_date);
CREATE INDEX idx_homework_submissions_homework_id ON homework_submissions (homework_id);
CREATE INDEX idx_homework_submissions_student_id ON homework_submissions (student_id);

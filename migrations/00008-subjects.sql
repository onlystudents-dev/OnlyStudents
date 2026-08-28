CREATE TABLE subjects (
    id SERIAL PRIMARY KEY,
    school_id INT NOT NULL REFERENCES schools (id) ON DELETE CASCADE,
    name VARCHAR(128) NOT NULL,
    code VARCHAR(16) NOT NULL,
    CONSTRAINT uq_subjects_school_code UNIQUE (school_id, code)
);

CREATE TABLE class_subjects (
    id SERIAL PRIMARY KEY,
    class_id INT NOT NULL REFERENCES classes (id) ON DELETE CASCADE,
    subject_id INT NOT NULL REFERENCES subjects (id) ON DELETE CASCADE,
    teacher_id INT NOT NULL REFERENCES teachers (id) ON DELETE RESTRICT,
    term_id INT NOT NULL REFERENCES terms (id) ON DELETE RESTRICT
);

CREATE INDEX idx_subjects_school_id ON subjects (school_id);
CREATE INDEX idx_class_subjects_class_id ON class_subjects (class_id);
CREATE INDEX idx_class_subjects_subject_id ON class_subjects (subject_id);
CREATE INDEX idx_class_subjects_teacher_id ON class_subjects (teacher_id);
CREATE INDEX idx_class_subjects_term_id ON class_subjects (term_id);

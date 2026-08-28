CREATE TABLE grade_types (
    id SERIAL PRIMARY KEY,
    school_id INT NOT NULL REFERENCES schools (id) ON DELETE CASCADE,
    name VARCHAR(64) NOT NULL,
    weight INT NOT NULL DEFAULT 1,
    CONSTRAINT chk_grade_types_weight CHECK (weight > 0)
);

CREATE TABLE grades (
    id BIGSERIAL PRIMARY KEY,
    student_id INT NOT NULL REFERENCES students (id) ON DELETE CASCADE,
    class_subjects_id INT NOT NULL REFERENCES class_subjects (id) ON DELETE RESTRICT,
    teacher_id INT NOT NULL REFERENCES teachers (id) ON DELETE RESTRICT,
    term_id INT NOT NULL REFERENCES terms (id) ON DELETE RESTRICT,
    grade_type_id INT NOT NULL REFERENCES grade_types (id) ON DELETE RESTRICT,
    value SMALLINT NOT NULL,
    date DATE NOT NULL,
    note VARCHAR(512) DEFAULT NULL,
    CONSTRAINT chk_grades_value CHECK (value BETWEEN 1 AND 5)
);

CREATE TABLE final_grades (
    id BIGSERIAL PRIMARY KEY,
    student_id INT NOT NULL REFERENCES students (id) ON DELETE CASCADE,
    class_subjects_id INT NOT NULL REFERENCES class_subjects (id) ON DELETE RESTRICT,
    term_id INT NOT NULL REFERENCES terms (id) ON DELETE RESTRICT,
    teacher_id INT NOT NULL REFERENCES teachers (id) ON DELETE RESTRICT,
    value SMALLINT NOT NULL,
    CONSTRAINT chk_final_grades_value CHECK (value BETWEEN 1 AND 5),
    CONSTRAINT uq_final_grades UNIQUE (student_id, class_subjects_id, term_id)
);

CREATE INDEX idx_grade_types_school_id ON grade_types (school_id);
CREATE INDEX idx_grades_student_id ON grades (student_id);
CREATE INDEX idx_grades_class_subjects_id ON grades (class_subjects_id);
CREATE INDEX idx_grades_term_id ON grades (term_id);
CREATE INDEX idx_grades_teacher_id ON grades (teacher_id);
CREATE INDEX idx_final_grades_student_id ON final_grades (student_id);
CREATE INDEX idx_final_grades_class_subjects_id ON final_grades (class_subjects_id);
CREATE INDEX idx_final_grades_term_id ON final_grades (term_id);

CREATE TABLE absences (
    id BIGSERIAL PRIMARY KEY,
    student_id INT NOT NULL REFERENCES students (id) ON DELETE CASCADE,
    lesson_id INT DEFAULT NULL REFERENCES lessons (id) ON DELETE SET NULL,
    date DATE NOT NULL,
    type VARCHAR(8) NOT NULL,
    justified BOOLEAN NOT NULL DEFAULT FALSE,
    verified_by INT DEFAULT NULL REFERENCES teachers (id) ON DELETE SET NULL,
    note VARCHAR(512) DEFAULT NULL,
    CONSTRAINT chk_absences_type CHECK (type IN ('absent','tardy')),
    CONSTRAINT uq_absences_student_lesson UNIQUE (student_id, lesson_id)
);

CREATE INDEX idx_absences_student_id ON absences (student_id);
CREATE INDEX idx_absences_lesson_id ON absences (lesson_id);
CREATE INDEX idx_absences_date ON absences (date);

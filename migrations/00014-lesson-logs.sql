CREATE TABLE lesson_logs (
    id BIGSERIAL PRIMARY KEY,
    lesson_id INT NOT NULL REFERENCES lessons (id) ON DELETE CASCADE,
    date DATE NOT NULL,
    teacher_id INT NOT NULL REFERENCES teachers (id) ON DELETE RESTRICT,
    topic VARCHAR(512) DEFAULT NULL,
    homework_id INT DEFAULT NULL REFERENCES homework (id) ON DELETE SET NULL,
    conducted BOOLEAN NOT NULL DEFAULT TRUE,
    CONSTRAINT uq_lesson_logs UNIQUE (lesson_id, date)
);

CREATE INDEX idx_lesson_logs_lesson_id ON lesson_logs (lesson_id);
CREATE INDEX idx_lesson_logs_date ON lesson_logs (date);
CREATE INDEX idx_lesson_logs_teacher_id ON lesson_logs (teacher_id);
CREATE INDEX idx_lesson_logs_homework_id ON lesson_logs (homework_id);

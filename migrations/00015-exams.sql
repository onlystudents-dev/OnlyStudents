CREATE TABLE exams (
    id BIGSERIAL PRIMARY KEY,
    class_subjects_id INT NOT NULL REFERENCES class_subjects (id) ON DELETE CASCADE,
    teacher_id INT NOT NULL REFERENCES teachers (id) ON DELETE RESTRICT,
    title VARCHAR(256) NOT NULL,
    description TEXT DEFAULT NULL,
    date DATE NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    room_id INT DEFAULT NULL REFERENCES rooms (id) ON DELETE SET NULL,
    CONSTRAINT chk_exams_time CHECK (end_time > start_time)
);

CREATE INDEX idx_exams_class_subjects_id ON exams (class_subjects_id);
CREATE INDEX idx_exams_teacher_id ON exams (teacher_id);
CREATE INDEX idx_exams_date ON exams (date);
CREATE INDEX idx_exams_room_id ON exams (room_id);

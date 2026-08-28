CREATE TABLE rooms (
    id SERIAL PRIMARY KEY,
    school_id INT NOT NULL REFERENCES schools (id) ON DELETE CASCADE,
    name VARCHAR(32) NOT NULL,
    capacity INT NOT NULL DEFAULT 0,
    CONSTRAINT chk_rooms_capacity CHECK (capacity >= 0)
);

CREATE TABLE time_slots (
    id SERIAL PRIMARY KEY,
    school_id INT NOT NULL REFERENCES schools (id) ON DELETE CASCADE,
    label VARCHAR(32) NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    position INT NOT NULL,
    CONSTRAINT chk_time_slots_time CHECK (end_time > start_time),
    CONSTRAINT uq_time_slots_school_position UNIQUE (school_id, position)
);

CREATE TABLE lessons (
    id SERIAL PRIMARY KEY,
    class_subjects_id INT NOT NULL REFERENCES class_subjects (id) ON DELETE CASCADE,
    day_of_week SMALLINT NOT NULL,
    time_slot_id INT NOT NULL REFERENCES time_slots (id) ON DELETE RESTRICT,
    room_id INT NOT NULL REFERENCES rooms (id) ON DELETE RESTRICT,
    CONSTRAINT chk_lessons_dow CHECK (day_of_week BETWEEN 1 AND 7),
    CONSTRAINT uq_lessons_class UNIQUE (day_of_week, time_slot_id, class_subjects_id),
    CONSTRAINT uq_lessons_room UNIQUE (day_of_week, time_slot_id, room_id)
);

CREATE INDEX idx_rooms_school_id ON rooms (school_id);
CREATE INDEX idx_time_slots_school_id ON time_slots (school_id);
CREATE INDEX idx_lessons_class_subjects_id ON lessons (class_subjects_id);
CREATE INDEX idx_lessons_time_slot_id ON lessons (time_slot_id);
CREATE INDEX idx_lessons_room_id ON lessons (room_id);

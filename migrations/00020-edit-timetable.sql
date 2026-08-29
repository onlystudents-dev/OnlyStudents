DROP TABLE IF EXISTS time_slots;
DROP TABLE IF EXISTS lessons;
DROP TABLE IF EXISTS calendar_events;

DROP TABLE IF EXISTS subjects;
DROP TABLE IF EXISTS class_subjects;

CREATE TABLE bell_schedule_type (
    id SERIAL PRIMARY KEY,
    school_id INT NOT NULL REFERENCES schools(id),
    name TEXT NOT NULL
);

CREATE TABLE bell_schedule (
    id SERIAL PRIMARY KEY,
    type_id INT NOT NULL REFERENCES bell_schedule_type(id),
    is_lesson bool NOT NULL DEFAULT(true),
    lesson_number INT,
    at_start TIME NOT NULL,
    at_end TIME NOT NULL,
    CONSTRAINT chk_is_lesson CHECK (is_lesson = TRUE OR lesson_number IS NULL)
);

CREATE TABLE subjects (
    id SERIAL PRIMARY KEY,
    subject_name TEXT NOT NULL
);

CREATE TABLE custom_subjects (
    id SERIAL PRIMARY KEY,
    school_id INT NOT NULL REFERENCES schools(id),
    subject_name TEXT NOT NULL
);

CREATE TABLE groups (
    id SERIAL PRIMARY KEY,
    school_id INT NOT NULL REFERENCES schools(id),
    bell_id INT NOT NULL REFERENCES bell_schedule_type(id),
    group_name TEXT NOT NULL
);

CREATE TABLE group_members (
    id SERIAL PRIMARY KEY,
    group_id INT NOT NULL REFERENCES groups(id),
    student_id INT NOT NULL REFERENCES students(id)
);

CREATE TABLE lessons (
    id SERIAL PRIMARY KEY,
    school_id INT NOT NULL REFERENCES schools(id),
    teacher_id INT NOT NULL REFERENCES teachers(id), 
    room_id INT NOT NULL REFERENCES rooms(id),
    group_id INT NOT NULL REFERENCES groups(id),
    custom_subject bool NOT NULL DEFAULT(false),
    subject_id INT REFERENCES subjects(id),
    custom_subject_id INT REFERENCES custom_subjects(id),
    CONSTRAINT chk_subject CHECK ( (custom_subject = FALSE AND subject_id IS NOT NULL AND custom_subject_id IS NULL) OR (custom_subject = TRUE AND custom_subject_id IS NOT NULL AND subject_id IS NULL) )
);

CREATE TABLE base_schedule (
    id SERIAL PRIMARY KEY,
    school_id INT NOT NULL REFERENCES schools(id),
    class_id INT NOT NULL REFERENCES classes(id),
    day_of_week INT NOT NULL,
    lesson_num INT NOT NULL,
    lesson_id INT NOT NULL REFERENCES lessons(id),
    CONSTRAINT chk_day CHECK (day_of_week BETWEEN 1 AND 7),
    UNIQUE(day_of_week, lesson_num, room_id)
);

CREATE TABLE time_table (
    id SERIAL PRIMARY KEY,
    school_id INT NOT NULL REFERENCES schools(id),
    actual_date DATE NOT NULL,
    lesson_num INT NOT NULL,
    lesson_id INT NOT NULL REFERENCES lessons(id),
    is_substitution bool NOT NULL DEFAULT(false),
    substitution_teacher_id INT REFERENCES teachers(id),
    canceled bool DEFAULT(false),
    CONSTRAINT chk_day CHECK (day_of_week BETWEEN 1 AND 7),
    CONSTRAINT chk_substitution CHECK (is_substitution = TRUE or substitution_teacher_id IS NULL),
    UNIQUE(day_of_week, lesson_num, room_id)
);

CREATE TABLE principal (
    id SERIAL PRIMARY KEY,
    school_id INT NOT NULL REFERENCES schools(id),
    teacher_id INT NOT NULL REFERENCES teacher(id)
);

CREATE TABLE permissions (
    id SERIAL PRIMARY KEY,
    school_id INT NOT NULL REFERENCES schools(id),
    teacher_id INT NOT NULL REFERENCES teacher(id),
    permission_id INT NOT NULL REFERENCES permission_type(id)
);

CREATE TABLE permission_type (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL
);

INSERT INTO permission_type (name, description) VALUES ('PRINCIPAL', 'Has all additional permissions.');
INSERT INTO permission_type (name, description) VALUES ('MANAGE_TIMETABLES', 'Can manage all class schedules.');
INSERT INTO permission_type (name, description) VALUES ('MANAGE_BELL_SCHEDULE', 'Can manage all bell schedule type and manage bell schedule.');
INSERT INTO permission_type (name, description) VALUES ('MANAGE_GROUPS', 'Can manage all groups and assign any student to any group.');
INSERT INTO permission_type (name, description) VALUES ('MANAGE_ROOMS', 'Can manage all rooms');
INSERT INTO permission_type (name, description) VALUES ('MANAGE_SUBSTITUTIONS', 'Can manage daily substitutions and cancel lessons.');
INSERT INTO permission_type (name, description) VALUES ('MANAGE_CUSTOM_SUBJECT', 'Can create, edit, or delete custom subjects.');

ALTER TABLE classes ADD COLUMN bell_id INT NOT NULL REFERENCES bell_schedule_type(id);

CREATE TABLE school_type (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE teachers (
    id SERIAL PRIMARY KEY,
    phone_number VARCHAR(256) NOT NULL UNIQUE,
    username VARCHAR(256) NOT NULL,
    birth_first_name VARCHAR(256) NOT NULL,
    birth_last_name VARCHAR(256) NOT NULL,
    birth_date DATE NOT NULL,
    birth_city VARCHAR(100) NOT NULL,
    birth_country CHAR(2) NOT NULL,
    permament_address VARCHAR(256) NOT NULL,
    temporary_address VARCHAR(256) NOT NULL,
    first_name VARCHAR(256) NOT NULL DEFAULT '',
    last_name VARCHAR(256) NOT NULL DEFAULT ''
);

CREATE TABLE schools (
    id SERIAL PRIMARY KEY,
    name VARCHAR(256) NOT NULL,
    zip_code VARCHAR(100) NOT NULL,
    city VARCHAR(100) NOT NULL,
    address_line VARCHAR(256) NOT NULL,
    school_type INT NOT NULL,
    principal_id INT NOT NULL,
    phone_number VARCHAR(30) NOT NULL,
    email_address VARCHAR(255) NOT NULL,
    CONSTRAINT fk_schools_school_type FOREIGN KEY (school_type) REFERENCES school_type (id),
    CONSTRAINT fk_schools_principal FOREIGN KEY (principal_id) REFERENCES teachers (id)
);

CREATE TABLE bell_schedule_type (
    id SERIAL PRIMARY KEY,
    school_id INT NOT NULL REFERENCES schools (id),
    name TEXT NOT NULL
);

CREATE TABLE classes (
    id SERIAL PRIMARY KEY,
    school_id INT NOT NULL,
    name VARCHAR(30) NOT NULL,
    teacher_id INT NOT NULL,
    co_teacher_id INT DEFAULT NULL,
    bell_id INT REFERENCES bell_schedule_type (id),
    CONSTRAINT fk_classes_school FOREIGN KEY (school_id) REFERENCES schools (id),
    CONSTRAINT fk_classes_teacher FOREIGN KEY (teacher_id) REFERENCES teachers (id),
    CONSTRAINT fk_classes_co_teacher FOREIGN KEY (co_teacher_id) REFERENCES teachers (id)
);

CREATE TABLE students (
    id SERIAL PRIMARY KEY,
    id_number INT NOT NULL UNIQUE,
    school_id INT NOT NULL,
    phone_number VARCHAR(30) DEFAULT NULL,
    first_name VARCHAR(256) NOT NULL,
    last_name VARCHAR(256) NOT NULL,
    birth_first_name VARCHAR(256) NOT NULL,
    birth_last_name VARCHAR(256) NOT NULL,
    birth_date DATE NOT NULL,
    birth_city VARCHAR(100) NOT NULL,
    birth_country CHAR(2) NOT NULL,
    mother_birth_first_name VARCHAR(256) NOT NULL,
    mother_birth_last_name VARCHAR(256) NOT NULL,
    classes_id INT NOT NULL,
    permament_address VARCHAR(256) NOT NULL,
    temporary_address VARCHAR(256) NOT NULL,
    tax_number INT DEFAULT NULL,
    ssn_number INT NOT NULL,
    bank_name VARCHAR(256) NOT NULL,
    iban_owner VARCHAR(256) NOT NULL,
    iban_number VARCHAR(34) NOT NULL,
    document_type VARCHAR(50) NOT NULL,
    document_number VARCHAR(50) NOT NULL,
    CONSTRAINT fk_students_school FOREIGN KEY (school_id) REFERENCES schools (id),
    CONSTRAINT fk_students_class FOREIGN KEY (classes_id) REFERENCES classes (id),
    CONSTRAINT uq_students_phone_number UNIQUE (phone_number)
);

CREATE TABLE guardians (
    id SERIAL PRIMARY KEY,
    phone_number VARCHAR(30) NOT NULL UNIQUE,
    first_name VARCHAR(256) NOT NULL,
    last_name VARCHAR(256) NOT NULL,
    birth_first_name VARCHAR(256) NOT NULL,
    birth_last_name VARCHAR(256) NOT NULL,
    birth_date DATE NOT NULL,
    birth_city VARCHAR(100) NOT NULL,
    birth_country CHAR(2) NOT NULL,
    permament_address VARCHAR(256) NOT NULL,
    temporary_address VARCHAR(256) NOT NULL
);

CREATE TABLE guardians_access (
    id SERIAL PRIMARY KEY,
    student_id INT NOT NULL,
    guardian_id INT NOT NULL,
    legal_representative BOOLEAN NOT NULL,
    CONSTRAINT fk_guardians_access_student FOREIGN KEY (student_id) REFERENCES students (id),
    CONSTRAINT fk_guardians_access_guardian FOREIGN KEY (guardian_id) REFERENCES guardians (id)
);

CREATE TABLE teacher_school (
    id SERIAL PRIMARY KEY,
    teacher_id INT NOT NULL,
    school_id INT NOT NULL,
    CONSTRAINT fk_teacher_school_teacher FOREIGN KEY (teacher_id) REFERENCES teachers (id),
    CONSTRAINT fk_teacher_school_school FOREIGN KEY (school_id) REFERENCES schools (id)
);

CREATE TABLE student_citizenships (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    document_type VARCHAR(50) NOT NULL,
    country CHAR(2) NOT NULL,
    CONSTRAINT fk_student_citizenships_user FOREIGN KEY (user_id) REFERENCES students (id)
);

CREATE TABLE guardian_citizenships (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    document_type VARCHAR(50) NOT NULL,
    country CHAR(2) NOT NULL,
    CONSTRAINT fk_guardian_citizenships_user FOREIGN KEY (user_id) REFERENCES guardians (id)
);

CREATE TABLE teacher_citizenships (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    document_type VARCHAR(50) NOT NULL,
    country CHAR(2) NOT NULL,
    CONSTRAINT fk_teacher_citizenships_user FOREIGN KEY (user_id) REFERENCES teachers (id)
);

CREATE TABLE accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    role TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    student_id INT REFERENCES students (id) ON DELETE CASCADE,
    teacher_id INT REFERENCES teachers (id) ON DELETE CASCADE,
    guardian_id INT REFERENCES guardians (id) ON DELETE CASCADE,
    email_address VARCHAR(256),
    email_verified BOOLEAN NOT NULL DEFAULT false,
    preferences jsonb NOT NULL DEFAULT '{}'::jsonb,
    CONSTRAINT chk_accounts_role CHECK (role IN ('student', 'teacher', 'guardian')),
    CONSTRAINT chk_accounts_role_link CHECK (
        (role = 'student' AND student_id IS NOT NULL AND teacher_id IS NULL AND guardian_id IS NULL) OR
        (role = 'teacher' AND teacher_id IS NOT NULL AND student_id IS NULL AND guardian_id IS NULL) OR
        (role = 'guardian' AND guardian_id IS NOT NULL AND student_id IS NULL AND teacher_id IS NULL)
    ),
    CONSTRAINT uq_accounts_email_address UNIQUE (email_address)
);

CREATE TABLE school_years (
    id SERIAL PRIMARY KEY,
    school_id INT NOT NULL REFERENCES schools (id) ON DELETE CASCADE,
    name VARCHAR(64) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT FALSE,
    CONSTRAINT chk_school_years_dates CHECK (end_date > start_date)
);

CREATE TABLE terms (
    id SERIAL PRIMARY KEY,
    school_year_id INT NOT NULL REFERENCES school_years (id) ON DELETE CASCADE,
    name VARCHAR(64) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    grade_deadline DATE NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT FALSE,
    CONSTRAINT chk_terms_dates CHECK (end_date >= start_date AND grade_deadline <= end_date)
);

CREATE TABLE rooms (
    id SERIAL PRIMARY KEY,
    school_id INT NOT NULL REFERENCES schools (id) ON DELETE CASCADE,
    name VARCHAR(32) NOT NULL,
    capacity INT NOT NULL DEFAULT 0,
    CONSTRAINT chk_rooms_capacity CHECK (capacity >= 0)
);

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
    class_subjects_id INT NOT NULL,
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
    class_subjects_id INT NOT NULL,
    term_id INT NOT NULL REFERENCES terms (id) ON DELETE RESTRICT,
    teacher_id INT NOT NULL REFERENCES teachers (id) ON DELETE RESTRICT,
    value SMALLINT NOT NULL,
    CONSTRAINT chk_final_grades_value CHECK (value BETWEEN 1 AND 5),
    CONSTRAINT uq_final_grades UNIQUE (student_id, class_subjects_id, term_id)
);

CREATE TABLE absences (
    id BIGSERIAL PRIMARY KEY,
    student_id INT NOT NULL REFERENCES students (id) ON DELETE CASCADE,
    lesson_id INT DEFAULT NULL,
    date DATE NOT NULL,
    type VARCHAR(8) NOT NULL,
    justified BOOLEAN NOT NULL DEFAULT FALSE,
    verified_by INT DEFAULT NULL REFERENCES teachers (id) ON DELETE SET NULL,
    note VARCHAR(512) DEFAULT NULL,
    CONSTRAINT chk_absences_type CHECK (type IN ('absent','tardy')),
    CONSTRAINT uq_absences_student_lesson UNIQUE (student_id, lesson_id)
);

CREATE TABLE homework (
    id BIGSERIAL PRIMARY KEY,
    class_subjects_id INT NOT NULL,
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

CREATE TABLE announcements (
    id BIGSERIAL PRIMARY KEY,
    author_account_id UUID NOT NULL REFERENCES accounts (id) ON DELETE RESTRICT,
    school_id INT DEFAULT NULL REFERENCES schools (id) ON DELETE CASCADE,
    class_id INT DEFAULT NULL REFERENCES classes (id) ON DELETE CASCADE,
    title VARCHAR(256) NOT NULL,
    body TEXT NOT NULL,
    pinned BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_announcements_scope CHECK (school_id IS NOT NULL OR class_id IS NOT NULL)
);

CREATE TABLE lesson_logs (
    id BIGSERIAL PRIMARY KEY,
    lesson_id INT NOT NULL,
    date DATE NOT NULL,
    teacher_id INT NOT NULL REFERENCES teachers (id) ON DELETE RESTRICT,
    topic VARCHAR(512) DEFAULT NULL,
    homework_id INT DEFAULT NULL REFERENCES homework (id) ON DELETE SET NULL,
    conducted BOOLEAN NOT NULL DEFAULT TRUE,
    CONSTRAINT uq_lesson_logs UNIQUE (lesson_id, date)
);

CREATE TABLE exams (
    id BIGSERIAL PRIMARY KEY,
    class_subjects_id INT NOT NULL,
    teacher_id INT NOT NULL REFERENCES teachers (id) ON DELETE RESTRICT,
    title VARCHAR(256) NOT NULL,
    description TEXT DEFAULT NULL,
    date DATE NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    room_id INT DEFAULT NULL REFERENCES rooms (id) ON DELETE SET NULL,
    CONSTRAINT chk_exams_time CHECK (end_time > start_time)
);

CREATE TABLE bell_schedule (
    id SERIAL PRIMARY KEY,
    school_id INT NOT NULL REFERENCES schools (id),
    type_id INT NOT NULL REFERENCES bell_schedule_type (id),
    lesson_number INT,
    at_start TIME NOT NULL,
    at_end TIME NOT NULL
);

CREATE TABLE subjects (
    id SERIAL PRIMARY KEY,
    subject_name TEXT NOT NULL,
    code TEXT
);

CREATE TABLE custom_subjects (
    id SERIAL PRIMARY KEY,
    school_id INT NOT NULL REFERENCES schools (id),
    subject_name TEXT NOT NULL
);

CREATE TABLE groups (
    id SERIAL PRIMARY KEY,
    school_id INT NOT NULL REFERENCES schools (id),
    bell_id INT NOT NULL REFERENCES bell_schedule_type (id),
    group_name TEXT NOT NULL
);

CREATE TABLE group_members (
    id SERIAL PRIMARY KEY,
    school_id INT NOT NULL REFERENCES schools (id),
    group_id INT NOT NULL REFERENCES groups (id),
    student_id INT NOT NULL REFERENCES students (id)
);

CREATE TABLE base_schedule (
    id SERIAL PRIMARY KEY,
    school_id INT NOT NULL REFERENCES schools (id),
    teacher_id INT NOT NULL REFERENCES teachers (id),
    room_id INT NOT NULL REFERENCES rooms (id),
    day_of_week INT NOT NULL,
    lesson_num INT NOT NULL,
    group_id INT NOT NULL REFERENCES groups (id),
    custom_subject bool NOT NULL DEFAULT(false),
    subject_id INT REFERENCES subjects (id),
    custom_subject_id INT REFERENCES custom_subjects (id),
    CONSTRAINT chk_subject CHECK ( (custom_subject = FALSE AND subject_id IS NOT NULL AND custom_subject_id IS NULL) OR (custom_subject = TRUE AND custom_subject_id IS NOT NULL AND subject_id IS NULL) ),
    CONSTRAINT chk_day CHECK (day_of_week BETWEEN 1 AND 7),
    UNIQUE(day_of_week, lesson_num, room_id)
);

CREATE TABLE time_table (
    id SERIAL PRIMARY KEY,
    school_id INT NOT NULL REFERENCES schools (id),
    teacher_id INT NOT NULL REFERENCES teachers (id),
    room_id INT NOT NULL REFERENCES rooms (id),
    day_of_week INT NOT NULL,
    group_id INT NOT NULL REFERENCES groups (id),
    custom_subject bool NOT NULL DEFAULT(false),
    subject_id INT REFERENCES subjects (id),
    custom_subject_id INT REFERENCES custom_subjects (id),
    actual_date DATE NOT NULL,
    lesson_num INT NOT NULL,
    is_substitution bool NOT NULL DEFAULT(false),
    substitution_teacher_id INT REFERENCES teachers (id),
    canceled bool DEFAULT(false),
    CONSTRAINT chk_day CHECK (day_of_week BETWEEN 1 AND 7),
    CONSTRAINT chk_substitution CHECK (is_substitution = TRUE or substitution_teacher_id IS NULL),
    UNIQUE(school_id, room_id, actual_date, lesson_num, group_id)
);

CREATE TABLE principal (
    id SERIAL PRIMARY KEY,
    school_id INT NOT NULL REFERENCES schools (id),
    teacher_id INT NOT NULL REFERENCES teachers (id)
);

CREATE TABLE permission_type (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL
);

CREATE TABLE permissions (
    id SERIAL PRIMARY KEY,
    school_id INT NOT NULL REFERENCES schools (id),
    teacher_id INT NOT NULL REFERENCES teachers (id),
    permission_id INT NOT NULL REFERENCES permission_type (id)
);

CREATE INDEX idx_students_school_id ON students (school_id);
CREATE INDEX idx_students_classes_id ON students (classes_id);

CREATE INDEX idx_schools_school_type ON schools (school_type);
CREATE INDEX idx_schools_principal_id ON schools (principal_id);

CREATE INDEX idx_guardians_access_student_id ON guardians_access (student_id);
CREATE INDEX idx_guardians_access_guardian_id ON guardians_access (guardian_id);

CREATE INDEX idx_teacher_school_teacher_id ON teacher_school (teacher_id);
CREATE INDEX idx_teacher_school_school_id ON teacher_school (school_id);

CREATE INDEX idx_classes_school_id ON classes (school_id);
CREATE INDEX idx_classes_teacher_id ON classes (teacher_id);
CREATE INDEX idx_classes_co_teacher_id ON classes (co_teacher_id);

CREATE INDEX idx_student_citizenships_user_id ON student_citizenships (user_id);
CREATE INDEX idx_guardian_citizenships_user_id ON guardian_citizenships (user_id);
CREATE INDEX idx_teacher_citizenships_user_id ON teacher_citizenships (user_id);

CREATE INDEX idx_accounts_role ON accounts (role);
CREATE INDEX idx_accounts_student_id ON accounts (student_id);
CREATE INDEX idx_accounts_teacher_id ON accounts (teacher_id);
CREATE INDEX idx_accounts_guardian_id ON accounts (guardian_id);

CREATE INDEX idx_school_years_school_id ON school_years (school_id);
CREATE INDEX idx_terms_school_year_id ON terms (school_year_id);
CREATE UNIQUE INDEX uq_school_years_active ON school_years (school_id) WHERE is_active;
CREATE UNIQUE INDEX uq_terms_active ON terms (school_year_id) WHERE is_active;

CREATE INDEX idx_rooms_school_id ON rooms (school_id);

CREATE INDEX idx_grade_types_school_id ON grade_types (school_id);
CREATE INDEX idx_grades_student_id ON grades (student_id);
CREATE INDEX idx_grades_class_subjects_id ON grades (class_subjects_id);
CREATE INDEX idx_grades_term_id ON grades (term_id);
CREATE INDEX idx_grades_teacher_id ON grades (teacher_id);
CREATE INDEX idx_final_grades_student_id ON final_grades (student_id);
CREATE INDEX idx_final_grades_class_subjects_id ON final_grades (class_subjects_id);
CREATE INDEX idx_final_grades_term_id ON final_grades (term_id);

CREATE INDEX idx_absences_student_id ON absences (student_id);
CREATE INDEX idx_absences_lesson_id ON absences (lesson_id);
CREATE INDEX idx_absences_date ON absences (date);

CREATE INDEX idx_homework_class_subjects_id ON homework (class_subjects_id);
CREATE INDEX idx_homework_teacher_id ON homework (teacher_id);
CREATE INDEX idx_homework_due_date ON homework (due_date);
CREATE INDEX idx_homework_submissions_homework_id ON homework_submissions (homework_id);
CREATE INDEX idx_homework_submissions_student_id ON homework_submissions (student_id);

CREATE INDEX idx_announcements_author_account_id ON announcements (author_account_id);
CREATE INDEX idx_announcements_school_id ON announcements (school_id);
CREATE INDEX idx_announcements_class_id ON announcements (class_id);
CREATE INDEX idx_announcements_created_at ON announcements (created_at);

CREATE INDEX idx_lesson_logs_lesson_id ON lesson_logs (lesson_id);
CREATE INDEX idx_lesson_logs_date ON lesson_logs (date);
CREATE INDEX idx_lesson_logs_teacher_id ON lesson_logs (teacher_id);
CREATE INDEX idx_lesson_logs_homework_id ON lesson_logs (homework_id);

CREATE INDEX idx_exams_class_subjects_id ON exams (class_subjects_id);
CREATE INDEX idx_exams_teacher_id ON exams (teacher_id);
CREATE INDEX idx_exams_date ON exams (date);
CREATE INDEX idx_exams_room_id ON exams (room_id);

INSERT INTO permission_type (name, description) VALUES
    ('PRINCIPAL', 'Has all additional permissions.'),
    ('MANAGE_TIMETABLES', 'Can manage all class schedules.'),
    ('MANAGE_BELL_SCHEDULE', 'Can manage all bell schedule type and manage bell schedule.'),
    ('MANAGE_GROUPS', 'Can manage all groups and assign any student to any group.'),
    ('MANAGE_ROOMS', 'Can manage all rooms'),
    ('MANAGE_SUBSTITUTIONS', 'Can manage daily substitutions and cancel lessons.'),
    ('MANAGE_CUSTOM_SUBJECT', 'Can create, edit, or delete custom subjects.');

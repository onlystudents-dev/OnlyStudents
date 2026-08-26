ALTER TABLE students DROP COLUMN password_hash;
ALTER TABLE students DROP COLUMN password_salt;

ALTER TABLE guardians DROP COLUMN password_hash;
ALTER TABLE guardians DROP COLUMN password_salt;

CREATE TABLE accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    role TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    student_id INT REFERENCES students (id) ON DELETE CASCADE,
    teacher_id INT REFERENCES teachers (id) ON DELETE CASCADE,
    guardian_id INT REFERENCES guardians (id) ON DELETE CASCADE,
    CONSTRAINT chk_accounts_role CHECK (role IN ('student', 'teacher', 'guardian')),
    CONSTRAINT chk_accounts_role_link CHECK (
        (role = 'student' AND student_id IS NOT NULL AND teacher_id IS NULL AND guardian_id IS NULL) OR
        (role = 'teacher' AND teacher_id IS NOT NULL AND student_id IS NULL AND guardian_id IS NULL) OR
        (role = 'guardian' AND guardian_id IS NOT NULL AND student_id IS NULL AND teacher_id IS NULL)
    )
);

CREATE INDEX idx_accounts_role ON accounts (role);
CREATE INDEX idx_accounts_student_id ON accounts (student_id);
CREATE INDEX idx_accounts_teacher_id ON accounts (teacher_id);
CREATE INDEX idx_accounts_guardian_id ON accounts (guardian_id);

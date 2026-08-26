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

ALTER TABLE students
    ADD CONSTRAINT uq_students_email_address UNIQUE (email_address),
    ADD CONSTRAINT uq_students_phone_number UNIQUE (phone_number);

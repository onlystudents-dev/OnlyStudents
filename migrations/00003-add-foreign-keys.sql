INSERT INTO school_type (id, name) VALUES (1, 'Default');

INSERT INTO teachers (
    id, email_address, phone_number, username,
    birth_first_name, birth_last_name, birth_date, birth_city, birth_country,
    permament_address, temporary_address
) VALUES (
    1, 'principal@school.test', '+0000000000', 'principal',
    'Principal', 'Principal', '1970-01-01', 'City', 'HU',
    'Addr', 'Addr'
);

INSERT INTO schools (
    id, name, zip_code, city, address_line,
    school_type, principal_id, phone_number, email_address
) VALUES (
    1, 'Default School', '0000', 'City', 'Addr',
    1, 1, '+0000000000', 'school@school.test'
);

INSERT INTO classes (id, school_id, name, teacher_id) VALUES (1, 1, 'Default Class', 1);

ALTER TABLE schools
    ADD CONSTRAINT fk_schools_school_type FOREIGN KEY (school_type) REFERENCES school_type (id),
    ADD CONSTRAINT fk_schools_principal FOREIGN KEY (principal_id) REFERENCES teachers (id);

ALTER TABLE students
    ADD CONSTRAINT fk_students_school FOREIGN KEY (school_id) REFERENCES schools (id),
    ADD CONSTRAINT fk_students_class FOREIGN KEY (classes_id) REFERENCES classes (id);

ALTER TABLE guardians_access
    ADD CONSTRAINT fk_guardians_access_student FOREIGN KEY (student_id) REFERENCES students (id),
    ADD CONSTRAINT fk_guardians_access_guardian FOREIGN KEY (guardian_id) REFERENCES guardians (id);

ALTER TABLE teacher_school
    ADD CONSTRAINT fk_teacher_school_teacher FOREIGN KEY (teacher_id) REFERENCES teachers (id),
    ADD CONSTRAINT fk_teacher_school_school FOREIGN KEY (school_id) REFERENCES schools (id);

ALTER TABLE classes
    ADD CONSTRAINT fk_classes_school FOREIGN KEY (school_id) REFERENCES schools (id),
    ADD CONSTRAINT fk_classes_teacher FOREIGN KEY (teacher_id) REFERENCES teachers (id),
    ADD CONSTRAINT fk_classes_co_teacher FOREIGN KEY (co_teacher_id) REFERENCES teachers (id);

ALTER TABLE student_citizenships
    ADD CONSTRAINT fk_student_citizenships_user FOREIGN KEY (user_id) REFERENCES students (id);

ALTER TABLE guardian_citizenships
    ADD CONSTRAINT fk_guardian_citizenships_user FOREIGN KEY (user_id) REFERENCES guardians (id);

ALTER TABLE teacher_citizenships
    ADD CONSTRAINT fk_teacher_citizenships_user FOREIGN KEY (user_id) REFERENCES teachers (id);

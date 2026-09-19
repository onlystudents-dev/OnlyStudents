INSERT INTO school_type (id, name) VALUES (1, 'Default');

INSERT INTO teachers (
    id, phone_number, username,
    birth_first_name, birth_last_name, birth_date, birth_city, birth_country,
    permament_address, temporary_address, first_name, last_name
) VALUES (
    1, '+0000000000', 'principal',
    'Principal', 'Principal', '1970-01-01', 'City', 'HU',
    'Addr', 'Addr', 'Principal', 'Principal'
);

INSERT INTO schools (
    id, name, zip_code, city, address_line,
    school_type, principal_id, phone_number, email_address
) VALUES (
    1, 'Default School', '0000', 'City', 'Addr',
    1, 1, '+0000000000', 'school@school.test'
);

INSERT INTO classes (id, school_id, name, teacher_id) VALUES (1, 1, 'Default Class', 1);

INSERT INTO students (
    id_number, school_id, first_name, last_name,
    birth_first_name, birth_last_name, birth_date, birth_city, birth_country,
    mother_birth_first_name, mother_birth_last_name,
    classes_id, permament_address, temporary_address, ssn_number,
    bank_name, iban_owner, iban_number, document_type, document_number
) VALUES (
    123456789, 1, 'Test', 'Student',
    'Test', 'Student', '2005-01-01', 'TestCity', 'HU',
    'MotherTest', 'MotherStudent',
    1, 'PermAddr', 'TempAddr', 987654321,
    'TestBank', 'TestOwner', 'HU12345678901234567890123456', 'ID', 'DOC123'
);

INSERT INTO guardians (
    id, phone_number, first_name, last_name,
    birth_first_name, birth_last_name, birth_date, birth_city, birth_country,
    permament_address, temporary_address
) VALUES (
    1, '+0000000001', 'Test', 'Guardian',
    'Test', 'Guardian', '1980-01-01', 'City', 'HU',
    'Addr', 'Addr'
);

INSERT INTO accounts (role, password_hash, student_id, email_address)
SELECT 'student', '$argon2id$v=19$m=65536,t=1,p=16$uVcQrvROzMvx8pEKN8eMlA$wfXFSmeuUNDWKTZ3edtmeA+VjXec53r1++TE9n5Z5Jg', s.id, NULL
FROM students s
WHERE s.id_number = 123456789
  AND NOT EXISTS (
    SELECT 1 FROM accounts a
    WHERE a.student_id = s.id AND a.role = 'student'
  );

INSERT INTO accounts (role, password_hash, teacher_id, email_address)
SELECT 'teacher', '$argon2id$v=19$m=65536,t=1,p=16$dVZjUXJ2Uk96TXZ4OHBFS044ZU1sQQ$yt9FrwyGnwXP4H7Go14ot4R7DgD8BoEguVNTfmF3mfM', t.id, 'principal@school.test'
FROM teachers t
WHERE t.id = 1
  AND NOT EXISTS (
    SELECT 1 FROM accounts a
    WHERE a.teacher_id = t.id AND a.role = 'teacher'
  );

INSERT INTO accounts (role, password_hash, guardian_id, email_address)
SELECT 'guardian', '$argon2id$v=19$m=65536,t=1,p=16$dVZjUXJ2Uk96TXZ4OHBFS044ZU1sQQ$yt9FrwyGnwXP4H7Go14ot4R7DgD8BoEguVNTfmF3mfM', g.id, 'guardian@school.test'
FROM guardians g
WHERE g.id = 1
  AND NOT EXISTS (
    SELECT 1 FROM accounts a
    WHERE a.guardian_id = g.id AND a.role = 'guardian'
  );

INSERT INTO guardians_access (student_id, guardian_id, legal_representative)
SELECT s.id, 1, TRUE
FROM students s
WHERE s.id_number = 123456789
  AND NOT EXISTS (
    SELECT 1 FROM guardians_access ga
    WHERE ga.student_id = s.id AND ga.guardian_id = 1
  );

INSERT INTO school_years (school_id, name, start_date, end_date, is_active)
VALUES (1, '2025/2026', '2025-09-01', '2026-06-15', TRUE);

INSERT INTO terms (school_year_id, name, start_date, end_date, grade_deadline, is_active)
VALUES (1, '1. félév', '2025-09-01', '2026-01-31', '2026-01-31', TRUE);

INSERT INTO subjects (subject_name)
VALUES
    ('Matematika'),
    ('Magyar nyelv'),
    ('Történelem'),
    ('Angol nyelv'),
    ('Informatika');

INSERT INTO grade_types (school_id, name, weight)
VALUES
    (1, 'Témazáró', 2),
    (1, 'Házi feladat', 1),
    (1, 'Felelés', 1),
    (1, 'Dolgozat', 3);

INSERT INTO grades (student_id, class_subjects_id, teacher_id, term_id, grade_type_id, value, date, note)
VALUES
    -- Matematika
    (1, 1, 1, 1, 1, 2, '2025-10-10', 'Első témazáró'),
    (1, 1, 1, 1, 2, 1, '2025-10-05', 'Házi feladat'),
    (1, 1, 1, 1, 3, 3, '2025-10-15', 'Felelés'),
    (1, 1, 1, 1, 4, 2, '2025-11-20', 'Első dolgozat'),
    -- Magyar nyelv
    (1, 2, 1, 1, 1, 3, '2025-10-12', 'Témazáró'),
    (1, 2, 1, 1, 2, 2, '2025-10-08', 'Házi feladat'),
    (1, 2, 1, 1, 3, 4, '2025-10-20', 'Felelés'),
    (1, 2, 1, 1, 4, 3, '2025-11-25', 'Dolgozat'),
    -- Történelem
    (1, 3, 1, 1, 1, 1, '2025-10-14', 'Témazáró'),
    (1, 3, 1, 1, 2, 2, '2025-10-07', 'Házi feladat'),
    (1, 3, 1, 1, 3, 2, '2025-10-22', 'Felelés'),
    -- Angol nyelv
    (1, 4, 1, 1, 1, 2, '2025-10-16', 'Témazáró'),
    (1, 4, 1, 1, 3, 3, '2025-10-25', 'Felelés'),
    (1, 4, 1, 1, 4, 1, '2025-11-30', 'Dolgozat'),
    -- Informatika
    (1, 5, 1, 1, 2, 1, '2025-10-03', 'Házi feladat'),
    (1, 5, 1, 1, 2, 2, '2025-10-17', 'Házi feladat'),
    (1, 5, 1, 1, 4, 2, '2025-11-15', 'Projekt');

INSERT INTO final_grades (student_id, class_subjects_id, term_id, teacher_id, value)
VALUES
    (1, 1, 1, 1, 2),
    (1, 2, 1, 1, 3),
    (1, 3, 1, 1, 2),
    (1, 4, 1, 1, 2),
    (1, 5, 1, 1, 1);

INSERT INTO permission_type (name, description) VALUES
    ('PRINCIPAL', 'Has all additional permissions.'),
    ('MANAGE_TIMETABLES', 'Can manage all class schedules.'),
    ('MANAGE_BELL_SCHEDULE', 'Can manage all bell schedule type and manage bell schedule.'),
    ('MANAGE_GROUPS', 'Can manage all groups and assign any student to any group.'),
    ('MANAGE_ROOMS', 'Can manage all rooms'),
    ('MANAGE_SUBSTITUTIONS', 'Can manage daily substitutions and cancel lessons.'),
    ('MANAGE_CUSTOM_SUBJECT', 'Can create, edit, or delete custom subjects.');

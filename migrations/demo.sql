INSERT INTO school_type (id, name) VALUES
    (1, 'Default');

INSERT INTO teachers (
    id, phone_number,
    birth_first_name, birth_last_name, birth_date, birth_city, birth_country,
    permament_address, temporary_address, first_name, last_name
) VALUES (
             1, '+0000000000',
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

INSERT INTO bell_schedule_type (
    id,
    school_id,
    name
) VALUES (
             1,
             1,
             'Default'
         );

INSERT INTO classes (
    id, school_id, name, teacher_id, bell_id
) VALUES (
             1, 1, 'Default Class', 1, 1
         );

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

INSERT INTO accounts (role, student_id, email_address)
SELECT 'student', s.id, NULL
FROM students s
WHERE s.id_number = 123456789
  AND NOT EXISTS (
    SELECT 1
    FROM accounts a
    WHERE a.student_id = s.id
      AND a.role = 'student'
);

INSERT INTO accounts (role, teacher_id, email_address)
SELECT 'teacher', t.id, 'principal@school.test'
FROM teachers t
WHERE t.id = 1
  AND NOT EXISTS (
    SELECT 1
    FROM accounts a
    WHERE a.teacher_id = t.id
      AND a.role = 'teacher'
);

INSERT INTO accounts (role, guardian_id, email_address)
SELECT 'guardian', g.id, 'guardian@school.test'
FROM guardians g
WHERE g.id = 1
  AND NOT EXISTS (
    SELECT 1
    FROM accounts a
    WHERE a.guardian_id = g.id
      AND a.role = 'guardian'
);

INSERT INTO guardians_access (
    student_id,
    guardian_id,
    legal_representative
)
SELECT
    s.id,
    1,
    TRUE
FROM students s
WHERE s.id_number = 123456789
  AND NOT EXISTS (
    SELECT 1
    FROM guardians_access ga
    WHERE ga.student_id = s.id
      AND ga.guardian_id = 1
);

INSERT INTO school_years (
    school_id,
    name,
    start_date,
    end_date,
    is_active
) VALUES (
             1,
             '2026/2027',
             '2026-09-01',
             '2027-06-15',
             TRUE
         );

INSERT INTO terms (
    school_year_id,
    name,
    start_date,
    end_date,
    grade_deadline,
    is_active
) VALUES (
             1,
             '1. félév',
             '2026-09-01',
             '2027-01-31',
             '2027-01-31',
             TRUE
         );

INSERT INTO subjects (
    subject_name
) VALUES
      ('Matematika'),
      ('Magyar nyelv'),
      ('Történelem'),
      ('Angol nyelv'),
      ('Informatika');

-- Taught subjects.
INSERT INTO class_subjects (
    id,
    school_id,
    class_id,
    subject_id,
    custom_subject,
    custom_subject_id,
    teacher_id
) VALUES
      -- Matematika
      (1, 1, 1, 1, FALSE, NULL, 1),

      -- Magyar nyelv
      (2, 1, 1, 2, FALSE, NULL, 1),

      -- Történelem
      (3, 1, 1, 3, FALSE, NULL, 1),

      -- Angol nyelv
      (4, 1, 1, 4, FALSE, NULL, 1),

      -- Informatika
      (5, 1, 1, 5, FALSE, NULL, 1);

INSERT INTO grade_types (
    school_id,
    name,
    weight
) VALUES
      (1, 'Témazáró', 2),
      (1, 'Házi feladat', 1),
      (1, 'Felelés', 1),
      (1, 'Dolgozat', 3);

INSERT INTO grades (
    student_id,
    class_subjects_id,
    teacher_id,
    term_id,
    grade_type_id,
    value,
    date,
    note
) VALUES
      -- Matematika
      (1, 1, 1, 1, 1, 2, '2026-09-10', 'Első témazáró'),
      (1, 1, 1, 1, 2, 1, '2026-09-08', 'Házi feladat'),
      (1, 1, 1, 1, 3, 3, '2026-09-14', 'Felelés'),
      (1, 1, 1, 1, 4, 2, '2026-09-18', 'Dolgozat'),

      -- Magyar nyelv
      (1, 2, 1, 1, 1, 3, '2026-09-11', 'Témazáró'),
      (1, 2, 1, 1, 2, 2, '2026-09-09', 'Házi feladat'),
      (1, 2, 1, 1, 3, 4, '2026-09-15', 'Felelés'),
      (1, 2, 1, 1, 4, 3, '2026-09-18', 'Dolgozat'),

      -- Történelem
      (1, 3, 1, 1, 1, 1, '2026-09-14', 'Témazáró'),
      (1, 3, 1, 1, 2, 2, '2026-09-10', 'Házi feladat'),
      (1, 3, 1, 1, 3, 2, '2026-09-17', 'Felelés'),

      -- Angol nyelv
      (1, 4, 1, 1, 1, 2, '2026-09-15', 'Témazáró'),
      (1, 4, 1, 1, 3, 3, '2026-09-16', 'Felelés'),
      (1, 4, 1, 1, 4, 1, '2026-09-18', 'Dolgozat'),

      -- Informatika
      (1, 5, 1, 1, 2, 1, '2026-09-09', 'Házi feladat'),
      (1, 5, 1, 1, 2, 2, '2026-09-16', 'Házi feladat'),
      (1, 5, 1, 1, 4, 2, '2026-09-18', 'Projekt');

INSERT INTO final_grades (
    student_id,
    class_subjects_id,
    term_id,
    teacher_id,
    value
) VALUES
      (1, 1, 1, 1, 2),
      (1, 2, 1, 1, 3),
      (1, 3, 1, 1, 2),
      (1, 4, 1, 1, 2),
      (1, 5, 1, 1, 1);

-- MANAGE_PRINCIPAL
INSERT INTO permissions (
    school_id,
    teacher_id,
    permission_id
) VALUES (
             1,
             1,
             1
         );

-- Test timetable setup

INSERT INTO bell_schedule (
    school_id,
    type_id,
    lesson_number,
    at_start,
    at_end
) VALUES
      (1, 1, 1, '08:00', '08:45'),
      (1, 1, 2, '08:55', '09:40'),
      (1, 1, 3, '09:50', '10:35'),
      (1, 1, 4, '10:45', '11:30'),
      (1, 1, 5, '11:40', '12:25'),
      (1, 1, 6, '12:35', '13:20'),
      (1, 1, 7, '13:30', '14:15'),
      (1, 1, 8, '14:25', '15:10');

INSERT INTO rooms (
    id,
    school_id,
    name,
    capacity
) VALUES (
             1,
             1,
             '101',
             30
         );

INSERT INTO groups (
    id,
    school_id,
    bell_id,
    group_name
) VALUES (
             1,
             1,
             1,
             '9.A'
         );

INSERT INTO group_members (
    school_id,
    group_id,
    student_id
) VALUES (
             1,
             1,
             1
         );

INSERT INTO teacher_school (
    teacher_id,
    school_id
) VALUES (
             1,
             1
         );

-- Weekly base timetable

INSERT INTO base_schedule (
    school_id,
    teacher_id,
    room_id,
    day_of_week,
    lesson_num,
    group_id,
    custom_subject,
    subject_id
) VALUES
      -- Monday
      (1, 1, 1, 1, 1, 1, FALSE, 1),
      (1, 1, 1, 1, 2, 1, FALSE, 2),
      (1, 1, 1, 1, 3, 1, FALSE, 3),
      (1, 1, 1, 1, 4, 1, FALSE, 4),
      (1, 1, 1, 1, 5, 1, FALSE, 5),

      -- Tuesday
      (1, 1, 1, 2, 1, 1, FALSE, 2),
      (1, 1, 1, 2, 2, 1, FALSE, 1),
      (1, 1, 1, 2, 3, 1, FALSE, 4),
      (1, 1, 1, 2, 4, 1, FALSE, 3),

      -- Wednesday
      (1, 1, 1, 3, 1, 1, FALSE, 1),
      (1, 1, 1, 3, 2, 1, FALSE, 5),
      (1, 1, 1, 3, 3, 1, FALSE, 2),
      (1, 1, 1, 3, 4, 1, FALSE, 4),

      -- Thursday
      (1, 1, 1, 4, 1, 1, FALSE, 3),
      (1, 1, 1, 4, 2, 1, FALSE, 1),
      (1, 1, 1, 4, 3, 1, FALSE, 4),
      (1, 1, 1, 4, 4, 1, FALSE, 5),

      -- Friday
      (1, 1, 1, 5, 1, 1, FALSE, 2),
      (1, 1, 1, 5, 2, 1, FALSE, 3),
      (1, 1, 1, 5, 3, 1, FALSE, 1),
      (1, 1, 1, 5, 4, 1, FALSE, 5);

-- Actual timetable for 2026-09-14 through 2026-09-18

INSERT INTO time_table (
    school_id,
    teacher_id,
    room_id,
    day_of_week,
    group_id,
    custom_subject,
    subject_id,
    actual_date,
    lesson_num,
    is_substitution,
    substitution_teacher_id,
    canceled
) VALUES
      -- Monday 2026-09-14
      (1, 1, 1, 1, 1, FALSE, 1, '2026-09-14', 1, FALSE, NULL, FALSE),
      (1, 1, 1, 1, 1, FALSE, 2, '2026-09-14', 2, FALSE, NULL, FALSE),
      (1, 1, 1, 1, 1, FALSE, 3, '2026-09-14', 3, FALSE, NULL, FALSE),
      (1, 1, 1, 1, 1, FALSE, 4, '2026-09-14', 4, FALSE, NULL, FALSE),
      (1, 1, 1, 1, 1, FALSE, 5, '2026-09-14', 5, FALSE, NULL, FALSE),

      -- Tuesday 2026-09-15
      (1, 1, 1, 2, 1, FALSE, 2, '2026-09-15', 1, FALSE, NULL, FALSE),
      (1, 1, 1, 2, 1, FALSE, 1, '2026-09-15', 2, FALSE, NULL, FALSE),
      (1, 1, 1, 2, 1, FALSE, 4, '2026-09-15', 3, FALSE, NULL, FALSE),
      (1, 1, 1, 2, 1, FALSE, 3, '2026-09-15', 4, FALSE, NULL, FALSE),

      -- Wednesday 2026-09-16
      (1, 1, 1, 3, 1, FALSE, 1, '2026-09-16', 1, FALSE, NULL, FALSE),
      (1, 1, 1, 3, 1, FALSE, 5, '2026-09-16', 2, FALSE, NULL, FALSE),
      (1, 1, 1, 3, 1, FALSE, 2, '2026-09-16', 3, FALSE, NULL, FALSE),
      (1, 1, 1, 3, 1, FALSE, 4, '2026-09-16', 4, FALSE, NULL, FALSE),

      -- Thursday 2026-09-17
      (1, 1, 1, 4, 1, FALSE, 3, '2026-09-17', 1, FALSE, NULL, FALSE),
      (1, 1, 1, 4, 1, FALSE, 1, '2026-09-17', 2, FALSE, NULL, FALSE),
      (1, 1, 1, 4, 1, FALSE, 4, '2026-09-17', 3, FALSE, NULL, FALSE),
      (1, 1, 1, 4, 1, FALSE, 5, '2026-09-17', 4, FALSE, NULL, FALSE),

      -- Friday 2026-09-18
      (1, 1, 1, 5, 1, FALSE, 2, '2026-09-18', 1, FALSE, NULL, FALSE),
      (1, 1, 1, 5, 1, FALSE, 3, '2026-09-18', 2, FALSE, NULL, FALSE),
      (1, 1, 1, 5, 1, FALSE, 1, '2026-09-18', 3, FALSE, NULL, FALSE),
      (1, 1, 1, 5, 1, FALSE, 5, '2026-09-18', 4, FALSE, NULL, FALSE),

      -- Friday 2026-09-18
      (1, 1, 1, 6, 1, FALSE, 1, '2026-09-19', 1, FALSE, NULL, FALSE),
      (1, 1, 1, 6, 1, FALSE, 4, '2026-09-19', 2, FALSE, NULL, FALSE),
      (1, 1, 1, 6, 1, FALSE, 3, '2026-09-19', 3, FALSE, NULL, FALSE),
      (1, 1, 1, 6, 1, FALSE, 5, '2026-09-19', 4, FALSE, NULL, FALSE);

-- Homework

INSERT INTO homework (
    class_subjects_id,
    teacher_id,
    title,
    description,
    due_date
) VALUES
      -- Monday
      (1, 1, 'Algebra gyakorló feladatok',
       'Oldd meg a 12–20. feladatokat a munkafüzetből.',
       '2026-09-21'),

      (4, 1, 'English vocabulary',
       'Tanuld meg az Unit 1 szavait.',
       '2026-09-21'),

      -- Tuesday
      (2, 1, 'Nyelvtani feladatlap',
       'A kijelölt nyelvtani feladatok megoldása.',
       '2026-09-22'),

      (1, 1, 'Geometria gyakorlás',
       'Háromszögek és szögek gyakorlása.',
       '2026-09-22'),

      -- Wednesday
      (5, 1, 'Informatika projekt',
       'Készíts egy rövid bemutatót a megadott témáról.',
       '2026-09-23'),

      -- Thursday
      (1, 1, 'Dolgozatra készülés',
       'Geometriai feladatok átnézése, megtanulása',
       '2026-09-24'),

      (4, 1, 'English exercises',
       'Complete exercises 4–8 in the workbook.',
       '2026-09-24'),

      -- Friday
      (1, 1, 'Törtek gyakorlása',
       'Oldd meg a kijelölt törtes feladatokat.',
       '2026-09-25'),

      (5, 1, 'Programozási feladat',
       'Készíts egy egyszerű Java programot.',
       '2026-09-25');

-- Exams

INSERT INTO exams (
    class_subjects_id,
    teacher_id,
    title,
    description,
    date,
    start_time,
    end_time,
    room_id
) VALUES
      -- Monday
      (4, 1, 'Angol témazáró',
       'Unit 1 dolgozat.',
       '2026-09-21',
       '08:00',
       '08:45',
       1),

      -- Tuesday
      (1, 1, 'Matematika dolgozat',
       'Algebra és egyenletek.',
       '2026-09-22',
       '08:55',
       '09:40',
       1),

      -- Wednesday
      (5, 1, 'Informatika számonkérés',
       'Alapvető programozási ismeretek.',
       '2026-09-23',
       '08:55',
       '09:40',
       1),

      -- Thursday
      (3, 1, 'Történelem témazáró',
       'Az ókori Róma.',
       '2026-09-24',
       '08:00',
       '08:45',
       1),

      -- Friday
      (1, 1, 'Matematika témazáró',
       'Geometria.',
       '2026-09-25',
       '09:50',
       '10:35',
       1);

-- Absences

INSERT INTO absences (
    student_id,
    class_subjects_id,
    lesson_id,
    date,
    type,
    justified,
    verified_by,
    note
)
SELECT
    1,
    2,
    bs.id,
    '2026-09-15',
    'absent',
    TRUE,
    1,
    'Betegség miatt nem jelent meg.'
FROM base_schedule bs
WHERE bs.school_id = 1
  AND bs.group_id = 1
  AND bs.day_of_week = 2
  AND bs.lesson_num = 1
  AND NOT EXISTS (
    SELECT 1
    FROM absences a
    WHERE a.student_id = 1
      AND a.lesson_id = bs.id
);

INSERT INTO absences (
    student_id,
    class_subjects_id,
    lesson_id,
    date,
    type,
    justified,
    verified_by,
    note
)
SELECT
    1,
    1,
    bs.id,
    '2026-09-16',
    'absent',
    FALSE,
    NULL,
    NULL
FROM base_schedule bs
WHERE bs.school_id = 1
  AND bs.group_id = 1
  AND bs.day_of_week = 3
  AND bs.lesson_num = 1
  AND NOT EXISTS (
    SELECT 1
    FROM absences a
    WHERE a.student_id = 1
      AND a.lesson_id = bs.id
);

INSERT INTO absences (
    student_id,
    class_subjects_id,
    lesson_id,
    date,
    type,
    justified,
    verified_by,
    note
)
SELECT
    1,
    5,
    bs.id,
    '2026-09-17',
    'tardy',
    FALSE,
    NULL,
    'Késett a tanórára.'
FROM base_schedule bs
WHERE bs.school_id = 1
  AND bs.group_id = 1
  AND bs.day_of_week = 4
  AND bs.lesson_num = 4
  AND NOT EXISTS (
    SELECT 1
    FROM absences a
    WHERE a.student_id = 1
      AND a.lesson_id = bs.id
);

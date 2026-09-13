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

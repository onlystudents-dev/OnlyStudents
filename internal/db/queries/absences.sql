-- name: GetStudentAbsences :many
SELECT
    a.id,
    s.name AS subject,
    s.code AS subject_code,
    CONCAT(t.first_name, ' ', t.last_name) AS teacher,
    a.date,
    a.type,
    a.justified,
    a.note,
    CONCAT(v.first_name, ' ', v.last_name) AS verified_by
FROM absences a
LEFT JOIN lessons l ON l.id = a.lesson_id
LEFT JOIN class_subjects cs ON cs.id = l.class_subjects_id
LEFT JOIN subjects s ON s.id = cs.subject_id
LEFT JOIN teachers t ON t.id = cs.teacher_id
LEFT JOIN teachers v ON v.id = a.verified_by
WHERE a.student_id = $1
ORDER BY a.date, s.name, a.id;

-- name: GetTeacherAbsences :many
SELECT
    a.id,
    CONCAT(st.first_name, ' ', st.last_name) AS student,
    s.name AS subject,
    s.code AS subject_code,
    a.date,
    a.type,
    a.justified,
    a.note,
    CONCAT(v.first_name, ' ', v.last_name) AS verified_by
FROM absences a
JOIN students st ON st.id = a.student_id
LEFT JOIN lessons l ON l.id = a.lesson_id
LEFT JOIN class_subjects cs ON cs.id = l.class_subjects_id
LEFT JOIN subjects s ON s.id = cs.subject_id
LEFT JOIN teachers t ON t.id = cs.teacher_id
LEFT JOIN teachers v ON v.id = a.verified_by
WHERE t.id = $1
ORDER BY a.date, st.last_name, st.first_name, a.id;

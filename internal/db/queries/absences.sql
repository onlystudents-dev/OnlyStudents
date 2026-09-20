-- name: GetStudentAbsences :many
SELECT
    a.id,
    COALESCE(cs.subject_name, s.subject_name) AS subject,
    s.code AS subject_code,
    CONCAT(t.first_name, ' ', t.last_name) AS teacher,
    a.date,
    a.type,
    a.justified,
    a.note,
    CONCAT(v.first_name, ' ', v.last_name) AS verified_by
FROM absences a
LEFT JOIN base_schedule bs ON bs.id = a.lesson_id
LEFT JOIN subjects s ON s.id = bs.subject_id
LEFT JOIN custom_subjects cs ON cs.id = bs.custom_subject_id
LEFT JOIN teachers t ON t.id = bs.teacher_id
LEFT JOIN teachers v ON v.id = a.verified_by
WHERE a.student_id = $1
ORDER BY a.date, COALESCE(cs.subject_name, s.subject_name), a.id;

-- name: GetTeacherAbsences :many
SELECT
    a.id,
    CONCAT(st.first_name, ' ', st.last_name) AS student,
    COALESCE(cs.subject_name, s.subject_name) AS subject,
    s.code AS subject_code,
    a.date,
    a.type,
    a.justified,
    a.note,
    CONCAT(v.first_name, ' ', v.last_name) AS verified_by
FROM absences a
JOIN students st ON st.id = a.student_id
LEFT JOIN base_schedule bs ON bs.id = a.lesson_id
LEFT JOIN subjects s ON s.id = bs.subject_id
LEFT JOIN custom_subjects cs ON cs.id = bs.custom_subject_id
LEFT JOIN teachers t ON t.id = bs.teacher_id
LEFT JOIN teachers v ON v.id = a.verified_by
WHERE t.id = $1
ORDER BY a.date, st.last_name, st.first_name, a.id;
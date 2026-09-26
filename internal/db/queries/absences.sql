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
LEFT JOIN class_subjects csub ON csub.id = a.class_subjects_id
LEFT JOIN subjects s ON s.id = csub.subject_id
LEFT JOIN custom_subjects cs ON cs.id = csub.custom_subject_id
LEFT JOIN base_schedule bs ON bs.id = a.lesson_id
LEFT JOIN teachers t ON t.id = COALESCE(bs.teacher_id, csub.teacher_id)
LEFT JOIN teachers v ON v.id = a.verified_by
WHERE a.student_id = $1
ORDER BY a.date, COALESCE(cs.subject_name, s.subject_name), a.id;

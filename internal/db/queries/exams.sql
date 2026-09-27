-- name: GetStudentExams :many
SELECT
    e.id,
    COALESCE(cs.subject_name, s.subject_name) AS subject,
    s.code AS subject_code,
    t.first_name AS teacher_first_name,
    t.last_name AS teacher_last_name,
    e.title,
    e.description,
    e.date,
    e.start_time,
    e.end_time,
    r.name AS room
FROM exams e
JOIN class_subjects csub ON csub.id = e.class_subjects_id
LEFT JOIN subjects s ON s.id = csub.subject_id
LEFT JOIN custom_subjects cs ON cs.id = csub.custom_subject_id
JOIN teachers t ON t.id = e.teacher_id
LEFT JOIN rooms r ON r.id = e.room_id
WHERE csub.class_id = (
    SELECT st.classes_id
    FROM students st
    WHERE st.id = $1
)
ORDER BY e.date, e.start_time;

-- name: AddExam :exec
INSERT INTO exams (class_subjects_id, teacher_id, title, description, date, start_time, end_time, room_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: EditExam :exec
UPDATE exams SET class_subjects_id = $2, teacher_id = $3, title = $4, description = $5, date = $6, start_time = $7, end_time = $8, room_id = $9 WHERE id = $1;

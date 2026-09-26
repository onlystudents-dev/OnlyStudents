-- name: GetStudentExams :many
SELECT
    e.id,
    s.subject_name AS subject,
    s.code AS subject_code,
    CONCAT(t.first_name, ' ', t.last_name) AS teacher,
    e.title,
    e.description,
    e.date,
    e.start_time,
    e.end_time,
    r.name AS room
FROM exams e
JOIN subjects s ON s.id = e.class_subjects_id
JOIN teachers t ON t.id = e.teacher_id
LEFT JOIN rooms r ON r.id = e.room_id
WHERE e.class_subjects_id IN (
    SELECT bs.subject_id
    FROM base_schedule bs
    JOIN group_members gm ON gm.group_id = bs.group_id
    WHERE gm.student_id = $1
      AND bs.subject_id IS NOT NULL
)
ORDER BY e.date, e.start_time;

-- name: AddExam :exec
INSERT INTO exams (class_subjects_id, teacher_id, title, description, date, start_time, end_time, room_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: EditExam :exec
UPDATE exams SET class_subjects_id = $2, teacher_id = $3, title = $4, description = $5, date = $6, start_time = $7, end_time = $8, room_id = $9 WHERE id = $1;

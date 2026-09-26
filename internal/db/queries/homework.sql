-- name: GetStudentHomework :many
SELECT
    h.id,
    s.subject_name AS subject,
    s.code AS subject_code,
    CONCAT(t.first_name, ' ', t.last_name) AS teacher,
    h.title,
    h.description,
    h.due_date,
    h.created_at,
    hs.submitted_at,
    hs.graded_value
FROM homework h
JOIN subjects s ON s.id = h.class_subjects_id
JOIN teachers t ON t.id = h.teacher_id
LEFT JOIN homework_submissions hs ON hs.homework_id = h.id AND hs.student_id = $1
WHERE h.class_subjects_id IN (
    SELECT bs.subject_id
    FROM base_schedule bs
    JOIN group_members gm ON gm.group_id = bs.group_id
    WHERE gm.student_id = $1
      AND bs.subject_id IS NOT NULL
)
ORDER BY h.due_date, h.created_at;

-- name: AddHomework :exec
INSERT INTO homework (class_subjects_id, teacher_id, title, description, due_date) VALUES ($1, $2, $3, $4, $5);

-- name: EditHomework :exec
UPDATE homework SET class_subjects_id = $2, teacher_id = $3, title = $4, description = $5, due_date = $6 WHERE id = $1;

-- name: AddHomeworkSubmission :exec
INSERT INTO homework_submissions (homework_id, student_id, content) VALUES ($1, $2, $3);

-- name: EditHomeworkSubmission :exec
UPDATE homework_submissions SET content = $2, graded_value = $3 WHERE homework_id = $1 AND student_id = $4;

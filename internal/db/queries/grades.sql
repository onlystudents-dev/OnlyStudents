-- name: GetStudentGrades :many
SELECT
    g.id,
    s.name AS subject,
    s.code AS subject_code,
    CONCAT(t.first_name, ' ', t.last_name) AS teacher,
    ter.name AS term,
    gt.name AS type,
    g.value,
    g.date,
    g.note
FROM grades g
JOIN class_subjects cs ON cs.id = g.class_subjects_id
JOIN subjects s ON s.id = cs.subject_id
JOIN teachers t ON t.id = g.teacher_id
JOIN terms ter ON ter.id = g.term_id
JOIN grade_types gt ON gt.id = g.grade_type_id
WHERE g.student_id = $1
ORDER BY g.term_id, s.name, g.date;

-- name: GetTeacherGrades :many
SELECT
    g.id,
    s.name AS subject,
    s.code AS subject_code,
    CONCAT(t.first_name, ' ', t.last_name) AS teacher,
    ter.name AS term,
    gt.name AS type,
    g.value,
    g.date,
    g.note
FROM grades g
JOIN class_subjects cs ON cs.id = g.class_subjects_id
JOIN subjects s ON s.id = cs.subject_id
JOIN teachers t ON t.id = g.teacher_id
JOIN terms ter ON ter.id = g.term_id
JOIN grade_types gt ON gt.id = g.grade_type_id
WHERE g.teacher_id = $1
ORDER BY g.term_id, s.name, g.date;

-- name: GetStudentFinalGrades :many
SELECT
    fg.id,
    s.name AS subject,
    s.code AS subject_code,
    CONCAT(t.first_name, ' ', t.last_name) AS teacher,
    ter.name AS term,
    fg.value
FROM final_grades fg
JOIN class_subjects cs ON cs.id = fg.class_subjects_id
JOIN subjects s ON s.id = cs.subject_id
JOIN teachers t ON t.id = fg.teacher_id
JOIN terms ter ON ter.id = fg.term_id
WHERE fg.student_id = $1
ORDER BY fg.term_id, s.name;

-- name: GetStudentGrades :many
SELECT
    g.id,
    COALESCE(cs.subject_name, s.subject_name) AS subject,
    s.code AS subject_code,
    t.first_name AS teacher_first_name,
    t.last_name AS teacher_last_name,
    ter.name AS term,
    gt.name AS type,
    g.value,
    g.date,
    g.note
FROM grades g
JOIN class_subjects csub ON csub.id = g.class_subjects_id
LEFT JOIN subjects s ON s.id = csub.subject_id
LEFT JOIN custom_subjects cs ON cs.id = csub.custom_subject_id
JOIN teachers t ON t.id = g.teacher_id
JOIN terms ter ON ter.id = g.term_id
JOIN grade_types gt ON gt.id = g.grade_type_id
WHERE g.student_id = $1
ORDER BY g.term_id, COALESCE(cs.subject_name, s.subject_name), g.date;

-- name: GetStudentFinalGrades :many
SELECT
    fg.id,
    COALESCE(cs.subject_name, s.subject_name) AS subject,
    s.code AS subject_code,
    t.first_name AS teacher_first_name,
    t.last_name AS teacher_last_name,
    ter.name AS term,
    fg.value
FROM final_grades fg
JOIN class_subjects csub ON csub.id = fg.class_subjects_id
LEFT JOIN subjects s ON s.id = csub.subject_id
LEFT JOIN custom_subjects cs ON cs.id = csub.custom_subject_id
JOIN teachers t ON t.id = fg.teacher_id
JOIN terms ter ON ter.id = fg.term_id
WHERE fg.student_id = $1
ORDER BY fg.term_id, COALESCE(cs.subject_name, s.subject_name);

-- name: AddGrade :exec
INSERT INTO grades (student_id, class_subjects_id, teacher_id, term_id, grade_type_id, value, date, note) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: AddFinalGrade :exec
INSERT INTO final_grades (student_id, class_subjects_id, term_id, teacher_id, value) VALUES ($1, $2, $3, $4, $5);

-- name: EditGrade :exec
UPDATE grades SET class_subjects_id = $2, teacher_id = $3, term_id = $4, grade_type_id = $5, value = $6, date = $7, note = $8 WHERE id = $1;

-- name: EditFinalGrade :exec
UPDATE final_grades SET class_subjects_id = $2, term_id = $3, teacher_id = $4, value = $5 WHERE id = $1;

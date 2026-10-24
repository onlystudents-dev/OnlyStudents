-- name: GetTeacherSubjects :many
SELECT DISTINCT
    bs.school_id,
    COALESCE(s.id, cs.id) AS id,
    COALESCE(s.subject_name, cs.subject_name) AS name,
    bs.custom_subject AS custom,
    t.id AS teacher_id,
    COALESCE(NULLIF(t.first_name, ''), t.birth_first_name)::text AS teacher_first_name,
    COALESCE(NULLIF(t.last_name, ''), t.birth_last_name)::text AS teacher_last_name,
    g.id AS group_id,
    g.group_name
FROM base_schedule bs
JOIN teachers t ON t.id = bs.teacher_id
JOIN groups g ON g.id = bs.group_id
LEFT JOIN subjects s ON s.id = bs.subject_id
LEFT JOIN custom_subjects cs ON cs.id = bs.custom_subject_id AND cs.school_id = bs.school_id
WHERE bs.school_id = sqlc.arg(school_id)
  AND bs.teacher_id = sqlc.arg(teacher_id)
ORDER BY name;

-- name: GetStudentSubjects :many
SELECT DISTINCT
    bs.school_id,
    COALESCE(s.id, cs.id) AS id,
    COALESCE(s.subject_name, cs.subject_name) AS name,
    bs.custom_subject AS custom,
    t.id AS teacher_id,
    COALESCE(NULLIF(t.first_name, ''), t.birth_first_name)::text AS teacher_first_name,
    COALESCE(NULLIF(t.last_name, ''), t.birth_last_name)::text AS teacher_last_name,
    g.id AS group_id,
    g.group_name
FROM base_schedule bs
JOIN group_members gm ON gm.group_id = bs.group_id AND gm.school_id = bs.school_id
JOIN teachers t ON t.id = bs.teacher_id
JOIN groups g ON g.id = bs.group_id
LEFT JOIN subjects s ON s.id = bs.subject_id
LEFT JOIN custom_subjects cs ON cs.id = bs.custom_subject_id AND cs.school_id = bs.school_id
WHERE bs.school_id = sqlc.arg(school_id)
  AND gm.student_id = sqlc.arg(student_id)
ORDER BY name;

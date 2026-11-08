-- name: GetTeacher :one
SELECT
    t.*,
    COALESCE(
            (
                SELECT jsonb_object_agg(school_id, permissions)
                FROM (
                         SELECT
                             p.school_id,
                             array_agg(pt.name) AS permissions
                         FROM permissions p
                                  JOIN permission_type pt ON pt.id = p.permission_id
                         WHERE p.teacher_id = t.id
                         GROUP BY p.school_id
                     ) sp
            ),
            '{}'::jsonb
    ) AS permissions
FROM teachers t
WHERE t.id = $1;

-- name: CreateTeacher :one
INSERT INTO teachers (phone_number,birth_first_name,birth_last_name,birth_date,birth_city,birth_country,permament_address,temporary_address,first_name,last_name) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id;

-- name: ListTeacherSchoolMemberships :many
SELECT *, school.name AS school_name FROM teacher_school ts
JOIN schools school ON school.id = ts.school_id
WHERE teacher_id = $1;
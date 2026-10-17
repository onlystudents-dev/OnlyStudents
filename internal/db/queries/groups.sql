-- name: GetStudentClass :one
SELECT c.*
FROM classes c
JOIN student_school ss ON ss.class_id = c.id
WHERE ss.id = sqlc.arg(student_id) AND ss.school_id = sqlc.arg(school_id);

-- name: GetGroupsByStudentID :many
SELECT g.*
FROM groups g
JOIN group_members gm ON gm.group_id = g.id AND gm.school_id = g.school_id
WHERE gm.student_id = $1 AND g.school_id = $2
ORDER BY g.group_name;

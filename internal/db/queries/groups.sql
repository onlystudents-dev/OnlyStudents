-- name: GetClassByID :one
SELECT * FROM classes WHERE id = $1;

-- name: GetGroupsByStudentID :many
SELECT g.*
FROM groups g
JOIN group_members gm ON gm.group_id = g.id AND gm.school_id = g.school_id
WHERE gm.student_id = $1 AND g.school_id = $2
ORDER BY g.group_name;

-- name: GetGuardianChildren :many
SELECT
    s.id,
    s.first_name,
    s.last_name,
    s.school_id,
    s.classes_id AS class_id
FROM guardians_access ga
JOIN students s ON s.id = ga.student_id
WHERE ga.guardian_id = $1
ORDER BY s.last_name, s.first_name;

-- name: CanViewStudent :one
SELECT 1 FROM guardians_access WHERE guardian_id = $1 AND student_id = $2;

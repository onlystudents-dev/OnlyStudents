-- name: GetGuardianChildren :many
SELECT
    s.id,
    s.first_name,
    s.last_name,
    ss.school_id,
    ss.class_id AS class_id,
    ss.id_number
FROM guardians_access ga
JOIN students s ON s.id = ga.student_id
JOIN student_school ss ON ss.student_id = s.id
WHERE ga.guardian_id = $1
ORDER BY s.last_name, s.first_name, ss.school_id;

-- name: CanViewStudent :one
SELECT 1 FROM guardians_access WHERE guardian_id = $1 AND student_id = $2;

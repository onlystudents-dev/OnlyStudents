-- name: ListStudentMemberships :many
SELECT *, school.name AS school_name FROM student_school ss
JOIN schools school ON school.id = ss.school_id
WHERE student_id = $1
ORDER BY school_id;

-- name: ListStudentMembershipsBySchool :many
SELECT
    ss.id,
    ss.student_id,
    ss.school_id,
    ss.class_id,
    ss.id_number,
    s.first_name,
    s.last_name
FROM student_school ss
JOIN students s ON s.id = ss.student_id
WHERE ss.school_id = $1
ORDER BY s.last_name, s.first_name, ss.school_id;

-- name: GetStudentMembership :one
SELECT * FROM student_school WHERE student_id = $1 AND school_id = $2;

-- name: CreateStudentMembership :one
INSERT INTO student_school (student_id, school_id, class_id, id_number)
SELECT sqlc.arg(student_id), sqlc.arg(school_id), sqlc.arg(class_id), sqlc.arg(id_number)
WHERE EXISTS (
    SELECT 1
    FROM classes c
    WHERE c.id = sqlc.arg(class_id)
      AND c.school_id = sqlc.arg(school_id)
)
RETURNING id;

-- name: DeleteStudentMembership :execrows
DELETE FROM student_school WHERE student_id = $1 AND school_id = $2;

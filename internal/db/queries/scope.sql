-- name: CheckPermission :one
SELECT EXISTS (SELECT 1 FROM permissions p JOIN permission_type pt ON pt.id = p.permission_id WHERE (pt.name = $1 OR pt.name = 'PRINCIPAL') AND p.teacher_id = $2 AND p.school_id = $3);

-- name: IsTeacherSchoolMember :one
SELECT EXISTS (SELECT 1 from teacher_school WHERE teacher_id = $1 AND school_id = $2);

-- name: CheckPermission :one
SELECT EXISTS (SELECT 1 FROM permissions p JOIN permission_type pt ON pt.id = p.permission_id WHERE (pt.name = $1 OR pt.name = 'PRINCIPAL') AND p.teacher_id = $2 AND p.school_id = $3);

-- name: IsTeacherSchoolMember :one
SELECT EXISTS (SELECT 1 from teacher_school WHERE teacher_id = $1 AND school_id = $2);

-- name: TeacherOwnsClassSubject :one
SELECT EXISTS (
    SELECT 1
    FROM class_subjects cs
    WHERE cs.id = $1
      AND cs.school_id = $2
      AND cs.teacher_id = $3
);

-- name: TeacherClassSubjectHasStudent :one
SELECT EXISTS (
    SELECT 1
    FROM class_subjects cs
    JOIN students st ON st.classes_id = cs.class_id
    WHERE cs.id = $1
      AND cs.school_id = $2
      AND cs.teacher_id = $3
      AND st.id = $4
);

-- name: TeacherTeachesStudent :one
SELECT EXISTS (
    SELECT 1
    FROM class_subjects cs
    JOIN students st ON st.classes_id = cs.class_id
    WHERE cs.school_id = $1
      AND cs.teacher_id = $2
      AND st.id = $3
);

-- name: GetStudentAbsences :many
SELECT
    a.id,
    COALESCE(cs.subject_name, s.subject_name) AS subject,
    s.code AS subject_code,
    t.first_name AS teacher_first_name,
    t.last_name AS teacher_last_name,
    a.date,
    a.type,
    a.justified,
    a.note,
    CONCAT(v.first_name, ' ', v.last_name) AS verified_by
FROM absences a
LEFT JOIN class_subjects csub ON csub.id = a.class_subjects_id
LEFT JOIN subjects s ON s.id = csub.subject_id
LEFT JOIN custom_subjects cs ON cs.id = csub.custom_subject_id
LEFT JOIN base_schedule bs ON bs.id = a.lesson_id
LEFT JOIN teachers t ON t.id = COALESCE(bs.teacher_id, csub.teacher_id)
LEFT JOIN teachers v ON v.id = a.verified_by
WHERE a.student_id = sqlc.arg(student_id)
  AND (
    csub.school_id = sqlc.arg(school_id)
    OR (csub.id IS NULL AND bs.school_id = sqlc.arg(school_id))
  )
ORDER BY a.date, COALESCE(cs.subject_name, s.subject_name), a.id;

-- name: GetTeacherAbsences :many
SELECT
    a.id,
    a.student_id,
    st.first_name,
    st.last_name,
    csub.id AS class_subjects_id,
    COALESCE(cs.subject_name, s.subject_name) AS subject,
    s.code AS subject_code,
    a.date,
    a.type,
    a.justified,
    a.note,
    a.verified_by
FROM absences a
JOIN class_subjects csub ON csub.id = a.class_subjects_id
JOIN students st ON st.id = a.student_id
LEFT JOIN subjects s ON s.id = csub.subject_id
LEFT JOIN custom_subjects cs ON cs.id = csub.custom_subject_id
WHERE csub.school_id = $1
  AND csub.teacher_id = $2
ORDER BY a.date, COALESCE(cs.subject_name, s.subject_name), a.id;

-- name: TeacherAddAbsence :execrows
INSERT INTO absences (student_id, class_subjects_id, date, type, note)
SELECT ss.student_id, cs.id, sqlc.arg(date), sqlc.arg(type), sqlc.arg(note)
FROM class_subjects cs
JOIN student_school ss ON ss.class_id = cs.class_id AND ss.school_id = cs.school_id
WHERE cs.id = sqlc.arg(class_subjects_id)
  AND cs.school_id = sqlc.arg(school_id)
  AND cs.teacher_id = sqlc.arg(teacher_id)
  AND ss.student_id = sqlc.arg(student_id)
  AND sqlc.arg(type) IN ('absent', 'tardy');

-- name: TeacherEditAbsence :execrows
UPDATE absences a
SET class_subjects_id = $2,
    date = $3,
    type = $4,
    note = $5
WHERE a.id = $1
  AND $4 IN ('absent', 'tardy')
  AND EXISTS (
      SELECT 1
      FROM class_subjects cs
      WHERE cs.id = a.class_subjects_id
        AND cs.school_id = $6
        AND cs.teacher_id = $7
  )
  AND EXISTS (
      SELECT 1
      FROM class_subjects cs
      WHERE cs.id = $2
        AND cs.school_id = $6
        AND cs.teacher_id = $7
  );

-- name: TeacherVerifyAbsence :execrows
UPDATE absences a
SET justified = sqlc.arg(justified),
    verified_by = CASE WHEN sqlc.arg(justified) THEN sqlc.arg(teacher_id) ELSE NULL END
WHERE a.id = sqlc.arg(id)
  AND EXISTS (
      SELECT 1
      FROM class_subjects cs
      WHERE cs.id = a.class_subjects_id
        AND cs.school_id = sqlc.arg(school_id)
        AND cs.teacher_id = sqlc.arg(teacher_id)
  );

-- name: TeacherDeleteAbsence :execrows
DELETE FROM absences a
WHERE a.id = $1
  AND EXISTS (
      SELECT 1
      FROM class_subjects cs
      WHERE cs.id = a.class_subjects_id
        AND cs.school_id = $2
        AND cs.teacher_id = $3
  );

-- name: GetStudentGrades :many
SELECT
    g.id,
    COALESCE(cs.subject_name, s.subject_name) AS subject,
    s.code AS subject_code,
    t.first_name AS teacher_first_name,
    t.last_name AS teacher_last_name,
    ter.name AS term,
    gt.name AS type,
    g.value,
    g.date,
    g.note
FROM grades g
JOIN class_subjects csub ON csub.id = g.class_subjects_id
LEFT JOIN subjects s ON s.id = csub.subject_id
LEFT JOIN custom_subjects cs ON cs.id = csub.custom_subject_id
JOIN teachers t ON t.id = g.teacher_id
JOIN terms ter ON ter.id = g.term_id
JOIN grade_types gt ON gt.id = g.grade_type_id
WHERE g.student_id = sqlc.arg(student_id)
  AND csub.school_id = sqlc.arg(school_id)
ORDER BY g.term_id, COALESCE(cs.subject_name, s.subject_name), g.date;

-- name: GetStudentFinalGrades :many
SELECT
    fg.id,
    COALESCE(cs.subject_name, s.subject_name) AS subject,
    s.code AS subject_code,
    t.first_name AS teacher_first_name,
    t.last_name AS teacher_last_name,
    ter.name AS term,
    fg.value
FROM final_grades fg
JOIN class_subjects csub ON csub.id = fg.class_subjects_id
LEFT JOIN subjects s ON s.id = csub.subject_id
LEFT JOIN custom_subjects cs ON cs.id = csub.custom_subject_id
JOIN teachers t ON t.id = fg.teacher_id
JOIN terms ter ON ter.id = fg.term_id
WHERE fg.student_id = sqlc.arg(student_id)
  AND csub.school_id = sqlc.arg(school_id)
ORDER BY fg.term_id, COALESCE(cs.subject_name, s.subject_name);

-- name: GetTeacherGrades :many
SELECT
    g.id,
    g.student_id,
    st.first_name,
    st.last_name,
    csub.id AS class_subjects_id,
    COALESCE(cs.subject_name, s.subject_name) AS subject,
    s.code AS subject_code,
    g.term_id,
    ter.name AS term,
    g.grade_type_id,
    gt.name AS type,
    g.value,
    g.date,
    g.note
FROM grades g
JOIN class_subjects csub ON csub.id = g.class_subjects_id
JOIN students st ON st.id = g.student_id
LEFT JOIN subjects s ON s.id = csub.subject_id
LEFT JOIN custom_subjects cs ON cs.id = csub.custom_subject_id
JOIN terms ter ON ter.id = g.term_id
JOIN grade_types gt ON gt.id = g.grade_type_id
WHERE csub.school_id = $1
  AND csub.teacher_id = $2
ORDER BY g.date, COALESCE(cs.subject_name, s.subject_name), g.id;

-- name: GetTeacherFinalGrades :many
SELECT
    fg.id,
    fg.student_id,
    st.first_name,
    st.last_name,
    csub.id AS class_subjects_id,
    COALESCE(cs.subject_name, s.subject_name) AS subject,
    s.code AS subject_code,
    fg.term_id,
    ter.name AS term,
    fg.value
FROM final_grades fg
JOIN class_subjects csub ON csub.id = fg.class_subjects_id
JOIN students st ON st.id = fg.student_id
LEFT JOIN subjects s ON s.id = csub.subject_id
LEFT JOIN custom_subjects cs ON cs.id = csub.custom_subject_id
JOIN terms ter ON ter.id = fg.term_id
WHERE csub.school_id = $1
  AND csub.teacher_id = $2
ORDER BY fg.term_id, COALESCE(cs.subject_name, s.subject_name), fg.id;

-- name: TeacherAddGrade :execrows
INSERT INTO grades (student_id, class_subjects_id, teacher_id, term_id, grade_type_id, value, date, note)
SELECT ss.student_id, cs.id, sqlc.arg(teacher_id), sqlc.arg(term_id), sqlc.arg(grade_type_id), sqlc.arg(value), sqlc.arg(date), sqlc.arg(note)
FROM class_subjects cs
JOIN student_school ss ON ss.classes_id = cs.class_id AND ss.school_id = cs.school_id
WHERE cs.id = sqlc.arg(class_subjects_id)
  AND cs.school_id = sqlc.arg(school_id)
  AND cs.teacher_id = sqlc.arg(teacher_id)
  AND ss.student_id = sqlc.arg(student_id)
  AND EXISTS (SELECT 1 FROM grade_types gt WHERE gt.id = sqlc.arg(grade_type_id) AND gt.school_id = sqlc.arg(school_id))
  AND EXISTS (SELECT 1 FROM terms t JOIN school_years sy ON sy.id = t.school_year_id WHERE t.id = sqlc.arg(term_id) AND sy.school_id = sqlc.arg(school_id));

-- name: TeacherAddFinalGrade :execrows
INSERT INTO final_grades (student_id, class_subjects_id, term_id, teacher_id, value)
SELECT ss.student_id, cs.id, sqlc.arg(term_id), sqlc.arg(teacher_id), sqlc.arg(value)
FROM class_subjects cs
JOIN student_school ss ON ss.classes_id = cs.class_id AND ss.school_id = cs.school_id
WHERE cs.id = sqlc.arg(class_subjects_id)
  AND cs.school_id = sqlc.arg(school_id)
  AND cs.teacher_id = sqlc.arg(teacher_id)
  AND ss.student_id = sqlc.arg(student_id)
  AND EXISTS (SELECT 1 FROM terms t JOIN school_years sy ON sy.id = t.school_year_id WHERE t.id = sqlc.arg(term_id) AND sy.school_id = sqlc.arg(school_id));

-- name: TeacherEditGrade :execrows
UPDATE grades g
SET class_subjects_id = $2,
    term_id = $3,
    grade_type_id = $4,
    value = $5,
    date = $6,
    note = $7
WHERE g.id = $1
  AND EXISTS (
      SELECT 1
      FROM class_subjects cs
      WHERE cs.id = g.class_subjects_id
        AND cs.school_id = $8
        AND cs.teacher_id = $9
  )
  AND EXISTS (
      SELECT 1
      FROM class_subjects cs
      WHERE cs.id = $2
        AND cs.school_id = $8
        AND cs.teacher_id = $9
  )
  AND EXISTS (SELECT 1 FROM grade_types gt WHERE gt.id = $4 AND gt.school_id = $8)
  AND EXISTS (SELECT 1 FROM terms t JOIN school_years sy ON sy.id = t.school_year_id WHERE t.id = $3 AND sy.school_id = $8);

-- name: TeacherEditFinalGrade :execrows
UPDATE final_grades fg
SET class_subjects_id = $2,
    term_id = $3,
    value = $4
WHERE fg.id = $1
  AND EXISTS (
      SELECT 1
      FROM class_subjects cs
      WHERE cs.id = fg.class_subjects_id
        AND cs.school_id = $5
        AND cs.teacher_id = $6
  )
  AND EXISTS (
      SELECT 1
      FROM class_subjects cs
      WHERE cs.id = $2
        AND cs.school_id = $5
        AND cs.teacher_id = $6
  )
  AND EXISTS (SELECT 1 FROM terms t JOIN school_years sy ON sy.id = t.school_year_id WHERE t.id = $3 AND sy.school_id = $5);

-- name: TeacherDeleteGrade :execrows
DELETE FROM grades g
WHERE g.id = $1
  AND EXISTS (
      SELECT 1
      FROM class_subjects cs
      WHERE cs.id = g.class_subjects_id
        AND cs.school_id = $2
        AND cs.teacher_id = $3
  );

-- name: TeacherDeleteFinalGrade :execrows
DELETE FROM final_grades fg
WHERE fg.id = $1
  AND EXISTS (
      SELECT 1
      FROM class_subjects cs
      WHERE cs.id = fg.class_subjects_id
        AND cs.school_id = $2
        AND cs.teacher_id = $3
  );

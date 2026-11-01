-- name: GetStudentGrades :many
SELECT
    g.id,
    COALESCE(cs.subject_name, s.subject_name) AS subject,
    s.code AS subject_code,
    t.first_name AS teacher_first_name,
    t.last_name AS teacher_last_name,
    ter.name AS term,
    gt.name AS type,
    gt.weight AS weight,
    g.value,
    g.date,
    g.note
FROM grades g
         JOIN lesson_grades lg ON lg.id = g.lesson_grade_id
         JOIN class_subjects csub ON csub.id = lg.class_subjects_id
         LEFT JOIN subjects s ON s.id = csub.subject_id
         LEFT JOIN custom_subjects cs ON cs.id = csub.custom_subject_id
         JOIN teachers t ON t.id = g.teacher_id
         JOIN terms ter ON ter.id = g.term_id
         JOIN grade_types gt ON gt.id = lg.type
WHERE g.student_id = sqlc.arg(student_id)
  AND csub.school_id = sqlc.arg(school_id)
ORDER BY g.term_id, COALESCE(cs.subject_name, s.subject_name), g.date;

-- name: GetStudentGradesBySubject :many
SELECT
    g.id,
    COALESCE(cs.subject_name, s.subject_name) AS subject,
    s.code AS subject_code,
    t.first_name AS teacher_first_name,
    t.last_name AS teacher_last_name,
    ter.name AS term,
    gt.name AS type,
    gt.weight AS weight,
    g.value,
    g.date,
    g.note
FROM grades g
         JOIN lesson_grades lg ON lg.id = g.lesson_grade_id
         JOIN class_subjects csub ON csub.id = lg.class_subjects_id
         LEFT JOIN subjects s ON s.id = csub.subject_id
         LEFT JOIN custom_subjects cs ON cs.id = csub.custom_subject_id
         JOIN teachers t ON t.id = g.teacher_id
         JOIN terms ter ON ter.id = g.term_id
         JOIN grade_types gt ON gt.id = lg.type
WHERE g.student_id = sqlc.arg(student_id)
  AND csub.school_id = sqlc.arg(school_id)
  AND (
    (NOT sqlc.arg(custom_subject)::bool AND csub.subject_id = sqlc.arg(subject_id)::int)
        OR
    (sqlc.arg(custom_subject)::bool AND csub.custom_subject_id = sqlc.arg(subject_id)::int)
    )
ORDER BY g.term_id, g.date;

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
    g.lesson_grade_id,
    gt.id AS grade_type_id,
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
         JOIN lesson_grades lg ON lg.id = g.lesson_grade_id
         JOIN grade_types gt ON gt.id = lg.type
WHERE csub.school_id = $1
  AND csub.teacher_id = $2
ORDER BY g.date, COALESCE(cs.subject_name, s.subject_name), g.id;

-- name: GetTeacherGradesByLesson :many
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
    g.lesson_grade_id,
    gt.id AS grade_type_id,
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
         JOIN lesson_grades lg ON lg.id = g.lesson_grade_id
         JOIN grade_types gt ON gt.id = lg.type
WHERE csub.school_id = $1
  AND csub.teacher_id = $2
  AND g.lesson_grade_id = $3
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

-- name: GetTeacherGradeTypes :many
SELECT
    id,
    name,
    weight
FROM grade_types
WHERE school_id = $1
ORDER BY weight;

-- name: TeacherAddGrade :execrows
INSERT INTO grades (student_id, teacher_id, term_id, lesson_grade_id, value, date, note)
SELECT ss.student_id, sqlc.arg(teacher_id), sqlc.arg(term_id), lg.id, sqlc.arg(value), sqlc.arg(date), sqlc.arg(note)
FROM class_subjects cs
         JOIN student_school ss ON ss.class_id = cs.class_id AND ss.school_id = cs.school_id
         JOIN lesson_grades lg ON lg.id = sqlc.arg(lesson_grade_id) AND lg.class_subjects_id = cs.id
WHERE cs.id = lg.class_subjects_id
  AND cs.school_id = sqlc.arg(school_id)
  AND cs.teacher_id = sqlc.arg(teacher_id)
  AND ss.student_id = sqlc.arg(student_id)
  AND EXISTS (SELECT 1 FROM grade_types gt WHERE gt.id = lg.type AND gt.school_id = sqlc.arg(school_id))
  AND EXISTS (SELECT 1 FROM terms t JOIN school_years sy ON sy.id = t.school_year_id WHERE t.id = sqlc.arg(term_id) AND sy.school_id = sqlc.arg(school_id));

-- name: TeacherAddFinalGrade :execrows
INSERT INTO final_grades (student_id, class_subjects_id, term_id, teacher_id, value)
SELECT ss.student_id, cs.id, sqlc.arg(term_id), sqlc.arg(teacher_id), sqlc.arg(value)
FROM class_subjects cs
         JOIN student_school ss ON ss.class_id = cs.class_id AND ss.school_id = cs.school_id
WHERE cs.id = sqlc.arg(class_subjects_id)
  AND cs.school_id = sqlc.arg(school_id)
  AND cs.teacher_id = sqlc.arg(teacher_id)
  AND ss.student_id = sqlc.arg(student_id)
  AND EXISTS (SELECT 1 FROM terms t JOIN school_years sy ON sy.id = t.school_year_id WHERE t.id = sqlc.arg(term_id) AND sy.school_id = sqlc.arg(school_id));

-- name: TeacherEditGrade :execrows
UPDATE grades g
SET term_id = $3,
    lesson_grade_id = $4,
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
             JOIN lesson_grades lg ON lg.class_subjects_id = cs.id
             JOIN grade_types gt ON gt.id = lg.type
    WHERE cs.id = $2
      AND cs.school_id = $8
      AND cs.teacher_id = $9
      AND lg.id = $4
      AND gt.school_id = $8
)
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

-- name: GetTeacherLessonGrades :many
SELECT
    lg.id,
    lg.class_subjects_id,
    lg.actual_date,
    lg.title,
    lg.description,
    lg.group_id,
    lg.type AS grade_type_id,
    gt.name AS type,
    gt.weight
FROM lesson_grades lg
         JOIN class_subjects cs ON cs.id = lg.class_subjects_id
         JOIN grade_types gt ON gt.id = lg.type
WHERE cs.school_id = $1
  AND cs.teacher_id = $2
ORDER BY lg.actual_date, lg.id;

-- name: GetTeacherLessonGradesBySubject :many
SELECT
    lg.id,
    lg.class_subjects_id,
    lg.actual_date,
    lg.title,
    lg.description,
    lg.group_id,
    lg.type AS grade_type_id,
    gt.name AS type,
    gt.weight
FROM lesson_grades lg
         JOIN class_subjects cs ON cs.id = lg.class_subjects_id
         JOIN grade_types gt ON gt.id = lg.type
WHERE cs.school_id = $1
  AND cs.teacher_id = $2
  AND cs.id = $3
  AND lg.actual_date = $4
ORDER BY lg.actual_date, lg.id;

-- name: TeacherAddLessonGrade :execrows
INSERT INTO lesson_grades (class_subjects_id, actual_date, title, description, type, group_id)
SELECT cs.id, sqlc.arg(actual_date), sqlc.arg(title), sqlc.arg(description), gt.id, sqlc.narg(group_id)
FROM class_subjects cs
         JOIN grade_types gt ON gt.id = sqlc.arg(type)
WHERE cs.id = sqlc.arg(class_subjects_id)
  AND cs.school_id = sqlc.arg(school_id)
  AND cs.teacher_id = sqlc.arg(teacher_id)
  AND gt.school_id = sqlc.arg(school_id)
  AND (
    sqlc.narg(group_id)::int IS NULL
        OR EXISTS (
            SELECT 1
            FROM groups gr
            WHERE gr.id = sqlc.narg(group_id)::int
        )
    );

-- name: TeacherEditLessonGrade :execrows
UPDATE lesson_grades lg
SET title = $2,
    description = $3,
    type = $4
WHERE lg.id = $1
  AND EXISTS (
    SELECT 1
    FROM class_subjects cs
    WHERE cs.id = lg.class_subjects_id
      AND cs.school_id = $5
      AND cs.teacher_id = $6
)
  AND EXISTS (
    SELECT 1
    FROM grade_types gt
    WHERE gt.id = $4
      AND gt.school_id = $5
);

-- name: TeacherDeleteLessonGrade :execrows
DELETE FROM lesson_grades lg
WHERE lg.id = $1
  AND EXISTS (
    SELECT 1
    FROM class_subjects cs
    WHERE cs.id = lg.class_subjects_id
      AND cs.school_id = $2
      AND cs.teacher_id = $3
);
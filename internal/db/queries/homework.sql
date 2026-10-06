-- name: GetStudentHomework :many
SELECT
    h.id,
    COALESCE(cs.subject_name, s.subject_name) AS subject,
    s.code AS subject_code,
    t.first_name AS teacher_first_name,
    t.last_name AS teacher_last_name,
    h.title,
    h.description,
    h.due_date,
    h.created_at,
    hs.submitted_at,
    hs.graded_value
FROM homework h
JOIN class_subjects csub ON csub.id = h.class_subjects_id
LEFT JOIN subjects s ON s.id = csub.subject_id
LEFT JOIN custom_subjects cs ON cs.id = csub.custom_subject_id
JOIN teachers t ON t.id = h.teacher_id
LEFT JOIN homework_submissions hs ON hs.homework_id = h.id AND hs.student_id = sqlc.arg(student_id)
WHERE csub.school_id = sqlc.arg(school_id)
  AND csub.class_id = (
      SELECT ss.classes_id
      FROM student_school ss
      WHERE ss.student_id = sqlc.arg(student_id)
        AND ss.school_id = sqlc.arg(school_id)
  )
ORDER BY h.due_date, h.created_at;

-- name: GetTeacherHomework :many
SELECT
    h.id,
    csub.id AS class_subjects_id,
    COALESCE(cs.subject_name, s.subject_name) AS subject,
    s.code AS subject_code,
    h.title,
    h.description,
    h.due_date,
    h.created_at,
    (SELECT COUNT(*) FROM homework_submissions hs WHERE hs.homework_id = h.id) AS submission_count,
    (SELECT COUNT(*) FROM homework_submissions hs WHERE hs.homework_id = h.id AND hs.submitted_at::date <= h.due_date) AS on_time_submission_count
FROM homework h
JOIN class_subjects csub ON csub.id = h.class_subjects_id
LEFT JOIN subjects s ON s.id = csub.subject_id
LEFT JOIN custom_subjects cs ON cs.id = csub.custom_subject_id
WHERE csub.school_id = $1
  AND csub.teacher_id = $2
ORDER BY h.due_date, h.created_at;

-- name: TeacherAddHomework :execrows
INSERT INTO homework (class_subjects_id, teacher_id, title, description, due_date)
SELECT cs.id, sqlc.arg(teacher_id), sqlc.arg(title), sqlc.arg(description), sqlc.arg(due_date)
FROM class_subjects cs
WHERE cs.id = sqlc.arg(class_subjects_id)
  AND cs.school_id = sqlc.arg(school_id)
  AND cs.teacher_id = sqlc.arg(teacher_id);

-- name: TeacherEditHomework :execrows
UPDATE homework h
SET class_subjects_id = $2,
    title = $3,
    description = $4,
    due_date = $5
WHERE h.id = $1
  AND EXISTS (
      SELECT 1
      FROM class_subjects cs
      WHERE cs.id = h.class_subjects_id
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

-- name: TeacherDeleteHomework :execrows
DELETE FROM homework h
WHERE h.id = $1
  AND EXISTS (
      SELECT 1
      FROM class_subjects cs
      WHERE cs.id = h.class_subjects_id
        AND cs.school_id = $2
        AND cs.teacher_id = $3
  );

-- name: StudentUpsertHomeworkSubmission :execrows
INSERT INTO homework_submissions (homework_id, student_id, content)
SELECT sqlc.arg(homework_id), sqlc.arg(student_id), sqlc.arg(content)
WHERE EXISTS (
    SELECT 1
    FROM homework h
    JOIN class_subjects csub ON csub.id = h.class_subjects_id
    WHERE h.id = sqlc.arg(homework_id)
      AND csub.school_id = sqlc.arg(school_id)
      AND csub.class_id = (
          SELECT ss.classes_id
          FROM student_school ss
          WHERE ss.student_id = sqlc.arg(student_id)
            AND ss.school_id = sqlc.arg(school_id)
      )
)
ON CONFLICT (homework_id, student_id) DO UPDATE
SET content = EXCLUDED.content,
    submitted_at = NOW();

-- name: TeacherUpdateHomeworkSubmission :execrows
UPDATE homework_submissions hs
SET content = $1,
    graded_value = $2
WHERE hs.homework_id = $3
  AND hs.student_id = $4
  AND EXISTS (
      SELECT 1
      FROM homework h
      JOIN class_subjects cs ON cs.id = h.class_subjects_id
      WHERE h.id = hs.homework_id
        AND cs.school_id = $5
        AND cs.teacher_id = $6
  );

-- name: GetTeacherHomeworkSubmissions :many
SELECT
    hs.id,
    st.first_name,
    st.last_name,
    hs.content,
    hs.submitted_at,
    hs.graded_value
FROM homework_submissions hs
JOIN students st ON st.id = hs.student_id
WHERE hs.homework_id = $1
  AND EXISTS (
      SELECT 1
      FROM homework h
      JOIN class_subjects cs ON cs.id = h.class_subjects_id
      WHERE h.id = hs.homework_id
        AND cs.school_id = $2
        AND cs.teacher_id = $3
  )
ORDER BY hs.submitted_at, hs.student_id;

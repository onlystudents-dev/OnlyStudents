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
    hs.graded_value,
    hs.note AS submission_note,
    f.uuid AS submission_file_uuid,
    f.original_name AS submission_file_name,
    f.size AS submission_file_size,
    f.type AS submission_file_type
FROM homework h
         JOIN class_subjects csub ON csub.id = h.class_subjects_id
         LEFT JOIN subjects s ON s.id = csub.subject_id
         LEFT JOIN custom_subjects cs ON cs.id = csub.custom_subject_id
         JOIN teachers t ON t.id = h.teacher_id
         LEFT JOIN homework_submissions hs ON hs.homework_id = h.id AND hs.student_id = sqlc.arg(student_id)
         LEFT JOIN files f ON f.uuid = hs.file_uuid
WHERE csub.school_id = sqlc.arg(school_id)
  AND csub.class_id = (
    SELECT ss.class_id
    FROM student_school ss
    WHERE ss.student_id = sqlc.arg(student_id)
      AND ss.school_id = sqlc.arg(school_id)
)
ORDER BY h.due_date, h.created_at;

-- name: GetStudentHomeworkByID :one
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
    hs.graded_value,
    hs.note AS submission_note,
    f.uuid AS submission_file_uuid,
    f.original_name AS submission_file_name,
    f.size AS submission_file_size,
    f.type AS submission_file_type
FROM homework h
         JOIN class_subjects csub ON csub.id = h.class_subjects_id
         LEFT JOIN subjects s ON s.id = csub.subject_id
         LEFT JOIN custom_subjects cs ON cs.id = csub.custom_subject_id
         JOIN teachers t ON t.id = h.teacher_id
         LEFT JOIN homework_submissions hs ON hs.homework_id = h.id AND hs.student_id = sqlc.arg(student_id)
         LEFT JOIN files f ON f.uuid = hs.file_uuid
WHERE csub.school_id = sqlc.arg(school_id)
  AND csub.class_id = (
    SELECT ss.class_id
    FROM student_school ss
    WHERE ss.student_id = sqlc.arg(student_id)
      AND ss.school_id = sqlc.arg(school_id)
) AND h.id = $1
ORDER BY h.due_date, h.created_at;

-- name: TeacherAddHomework :execrows
INSERT INTO homework (class_subjects_id, teacher_id, title, description, due_date, assigned_date)
SELECT cs.id, sqlc.arg(teacher_id), sqlc.arg(title), sqlc.arg(description), sqlc.arg(due_date), sqlc.arg(assigned_date)
FROM class_subjects cs
WHERE cs.id = sqlc.arg(class_subjects_id)
  AND cs.school_id = sqlc.arg(school_id)
  AND cs.teacher_id = sqlc.arg(teacher_id);

-- name: GetTeacherHomework :one
SELECT
    h.id,
    h.class_subjects_id,
    h.title,
    h.description,
    h.due_date,
    h.created_at,
    COALESCE(cs.subject_name, s.subject_name) AS subject
FROM homework h
JOIN class_subjects csub ON csub.id = h.class_subjects_id
LEFT JOIN subjects s ON s.id = csub.subject_id
LEFT JOIN custom_subjects cs ON cs.id = csub.custom_subject_id
WHERE h.id = $1
  AND csub.school_id = $2
  AND csub.teacher_id = $3;

-- name: TeacherEditHomework :execrows
UPDATE homework h
SET class_subjects_id = $2,
    title = $3,
    description = $4,
    due_date = $5,
    assigned_date = $6
WHERE h.id = $1
  AND EXISTS (
      SELECT 1
      FROM class_subjects cs
      WHERE cs.id = h.class_subjects_id
        AND cs.school_id = $7
        AND cs.teacher_id = $8
  )
  AND EXISTS (
      SELECT 1
      FROM class_subjects cs
      WHERE cs.id = $2
        AND cs.school_id = $7
        AND cs.teacher_id = $8
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
INSERT INTO homework_submissions (homework_id, student_id, note, file_uuid)
SELECT sqlc.arg(homework_id)::int, sqlc.arg(student_id), sqlc.arg(note), sqlc.arg(file_uuid)
WHERE EXISTS (
    SELECT 1
    FROM homework h
    JOIN class_subjects csub ON csub.id = h.class_subjects_id
    WHERE h.id = sqlc.arg(homework_id)
      AND csub.school_id = sqlc.arg(school_id)
      AND csub.class_id = (
          SELECT ss.class_id
          FROM student_school ss
          WHERE ss.student_id = sqlc.arg(student_id)
            AND ss.school_id = sqlc.arg(school_id)
      )
)
ON CONFLICT (homework_id, student_id) DO UPDATE
SET note = EXCLUDED.note,
    file_uuid = EXCLUDED.file_uuid,
    submitted_at = NOW();

-- name: TeacherGradeHomeworkSubmission :execrows
UPDATE homework_submissions hs
SET graded_value = sqlc.arg(graded_value)
WHERE hs.homework_id = sqlc.arg(homework_id)
  AND hs.student_id = sqlc.arg(student_id)
  AND EXISTS (
      SELECT 1
      FROM homework h
      JOIN class_subjects cs ON cs.id = h.class_subjects_id
      WHERE h.id = hs.homework_id
        AND cs.school_id = sqlc.arg(school_id)
        AND cs.teacher_id = sqlc.arg(teacher_id)
  );

-- name: GetTeacherHomeworkSubmissions :many
SELECT
    hs.id,
    st.first_name,
    st.last_name,
    h.due_date,
    f.original_name AS file_name,
    f.size AS file_size,
    f.type AS file_type,
    f.uuid AS file_uuid,
    hs.note,
    hs.submitted_at,
    hs.graded_value
FROM homework_submissions hs
JOIN homework h ON h.id = hs.homework_id
JOIN students st ON st.id = hs.student_id
LEFT JOIN files f ON f.uuid = hs.file_uuid
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

-- name: GetTeacherUngradedHomeworkSubmissions :many
SELECT
    hs.id,
    hs.homework_id,
    h.title AS homework_title,
    h.due_date,
    csub.class_id AS class_subjects_id,
    COALESCE(cs.subject_name, s.subject_name) AS subject,
    st.id AS student_id,
    st.first_name,
    st.last_name,
    f.original_name AS file_name,
    f.size AS file_size,
    f.type AS file_type,
    f.uuid AS file_uuid,
    hs.note,
    hs.submitted_at
FROM homework_submissions hs
         JOIN homework h ON h.id = hs.homework_id
         JOIN class_subjects csub ON csub.id = h.class_subjects_id
         LEFT JOIN subjects s ON s.id = csub.subject_id
         LEFT JOIN custom_subjects cs ON cs.id = csub.custom_subject_id
         JOIN students st ON st.id = hs.student_id
         LEFT JOIN files f ON f.uuid = hs.file_uuid
WHERE hs.graded_value IS NULL
  AND csub.school_id = $1
  AND csub.teacher_id = $2
ORDER BY hs.submitted_at, hs.student_id;

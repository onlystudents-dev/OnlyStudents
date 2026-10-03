-- name: GetStudentExams :many
SELECT
    e.id,
    COALESCE(cs.subject_name, s.subject_name) AS subject,
    s.code AS subject_code,
    t.first_name AS teacher_first_name,
    t.last_name AS teacher_last_name,
    e.title,
    e.description,
    e.date,
    e.start_time,
    e.end_time,
    r.name AS room
FROM exams e
JOIN class_subjects csub ON csub.id = e.class_subjects_id
LEFT JOIN subjects s ON s.id = csub.subject_id
LEFT JOIN custom_subjects cs ON cs.id = csub.custom_subject_id
JOIN teachers t ON t.id = e.teacher_id
LEFT JOIN rooms r ON r.id = e.room_id
WHERE csub.class_id = (
    SELECT st.classes_id
    FROM students st
    WHERE st.id = $1
)
ORDER BY e.date, e.start_time;

-- name: GetTeacherExams :many
SELECT
    e.id,
    csub.id AS class_subjects_id,
    COALESCE(cs.subject_name, s.subject_name) AS subject,
    s.code AS subject_code,
    e.title,
    e.description,
    e.date,
    e.start_time,
    e.end_time,
    r.id AS room_id,
    r.name AS room
FROM exams e
JOIN class_subjects csub ON csub.id = e.class_subjects_id
LEFT JOIN subjects s ON s.id = csub.subject_id
LEFT JOIN custom_subjects cs ON cs.id = csub.custom_subject_id
LEFT JOIN rooms r ON r.id = e.room_id
WHERE csub.school_id = $1
  AND csub.teacher_id = $2
ORDER BY e.date, e.start_time;

-- name: TeacherAddExam :execrows
INSERT INTO exams (class_subjects_id, teacher_id, title, description, date, start_time, end_time, room_id)
SELECT cs.id, sqlc.arg(teacher_id), sqlc.arg(title), sqlc.arg(description), sqlc.arg(date), sqlc.arg(start_time), sqlc.arg(end_time), sqlc.arg(room_id)
FROM class_subjects cs
WHERE cs.id = sqlc.arg(class_subjects_id)
  AND cs.school_id = sqlc.arg(school_id)
  AND cs.teacher_id = sqlc.arg(teacher_id)
  AND (sqlc.arg(room_id) IS NULL OR EXISTS (SELECT 1 FROM rooms r WHERE r.id = sqlc.arg(room_id) AND r.school_id = sqlc.arg(school_id)));

-- name: TeacherEditExam :execrows
UPDATE exams e
SET class_subjects_id = $2,
    title = $3,
    description = $4,
    date = $5,
    start_time = $6,
    end_time = $7,
    room_id = $8
WHERE e.id = $1
  AND EXISTS (
      SELECT 1
      FROM class_subjects cs
      WHERE cs.id = e.class_subjects_id
        AND cs.school_id = $9
        AND cs.teacher_id = $10
  )
  AND EXISTS (
      SELECT 1
      FROM class_subjects cs
      WHERE cs.id = $2
        AND cs.school_id = $9
        AND cs.teacher_id = $10
  )
  AND ($8 IS NULL OR EXISTS (SELECT 1 FROM rooms r WHERE r.id = $8 AND r.school_id = $9));

-- name: TeacherDeleteExam :execrows
DELETE FROM exams e
WHERE e.id = $1
  AND EXISTS (
      SELECT 1
      FROM class_subjects cs
      WHERE cs.id = e.class_subjects_id
        AND cs.school_id = $2
        AND cs.teacher_id = $3
  );

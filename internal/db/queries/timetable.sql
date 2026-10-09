-- name: CreateBellScheduleType :execrows
INSERT INTO bell_schedule_type (school_id, name) VALUES ($1, $2);

-- name: DeleteBellScheduleType :execrows
DELETE FROM bell_schedule_type WHERE school_id = $1 AND id = $2;

-- name: EditBellScheduleType :execrows
UPDATE bell_schedule_type SET name = $1 WHERE id = $2 AND school_id = $3;

-- name: ReadBellScheduleType :many
SELECT id, name FROM bell_schedule_type WHERE school_id = $1;

-- name: CreateLessonTime :execrows
INSERT INTO bell_schedule (school_id, type_id, lesson_number, at_start, at_end)
SELECT sqlc.arg(school_id), sqlc.arg(type_id), sqlc.arg(lesson_number), sqlc.arg(at_start), sqlc.arg(at_end)
FROM schools s
WHERE s.id = sqlc.arg(school_id)
  AND EXISTS (SELECT 1 FROM bell_schedule_type bst WHERE bst.id = sqlc.arg(type_id) AND bst.school_id = sqlc.arg(school_id));

-- name: DeleteLessonTime :execrows
DELETE FROM bell_schedule WHERE school_id = $1 AND id = $2;

-- name: EditLessonTime :execrows
UPDATE bell_schedule set lesson_number = $1, at_start = $2, at_end = $3 WHERE id = $4 AND school_id = $5;

-- name: ReadLessonTime :many
SELECT * FROM bell_schedule WHERE school_id = $1 AND type_id = $2;

-- name: CreateCustomSubject :execrows
INSERT INTO custom_subjects (school_id, subject_name) VALUES ($1, $2);

-- name: EditCustomSubject :execrows
UPDATE custom_subjects SET subject_name = $1 WHERE school_id = $2 AND id = $3;

-- name: DeleteCustomSubject :execrows
DELETE FROM custom_subjects WHERE school_id = $1 AND id = $2;

-- name: ReadCustomSubject :many
SELECT id, subject_name FROM custom_subjects WHERE school_id = $1;

-- name: CreateRoom :execrows
INSERT INTO rooms (school_id, name, capacity) VALUES ($1, $2, $3);

-- name: EditRoom :execrows
UPDATE rooms SET name = $1, capacity = $2 WHERE school_id = $3 AND id = $4;

-- name: DeleteRoom :execrows
DELETE FROM rooms WHERE school_id = $1 AND id = $2;

-- name: ReadRoom :many
SELECT id, name, capacity FROM rooms WHERE school_id = $1;

-- name: CreateGroup :execrows
INSERT INTO groups (school_id, bell_id, group_name)
SELECT sqlc.arg(school_id), sqlc.arg(bell_id), sqlc.arg(group_name)
FROM schools s
WHERE s.id = sqlc.arg(school_id)
  AND EXISTS (SELECT 1 FROM bell_schedule_type bst WHERE bst.id = sqlc.arg(bell_id) AND bst.school_id = sqlc.arg(school_id));

-- name: EditGroup :execrows
UPDATE groups SET group_name = $1 WHERE school_id = $2 AND id = $3;

-- name: DeleteGroup :execrows
DELETE FROM groups WHERE school_id = $1 AND id = $2;

-- name: ReadGroup :many
SELECT * FROM groups WHERE school_id = $1;

-- name: InsertStudentToGroup :execrows
INSERT INTO group_members (school_id, group_id, student_id)
SELECT sqlc.arg(school_id), sqlc.arg(group_id), sqlc.arg(student_id)
FROM schools s
WHERE s.id = sqlc.arg(school_id)
  AND EXISTS (SELECT 1 FROM groups g WHERE g.id = sqlc.arg(group_id) AND g.school_id = sqlc.arg(school_id))
  AND EXISTS (SELECT 1 FROM student_school ss WHERE ss.student_id = sqlc.arg(student_id) AND ss.school_id = sqlc.arg(school_id));

-- name: DeleteStudentFromGroup :execrows
DELETE FROM group_members WHERE group_id = $1 AND student_id = $2 AND school_id = $3;

-- name: ReadListOfStudents :many
SELECT s.id, s.first_name, s.last_name FROM students AS s INNER JOIN group_members AS g_m ON g_m.student_id = s.id WHERE g_m.group_id = sqlc.arg(group_id) AND g_m.school_id = sqlc.arg(school_id) AND EXISTS (SELECT 1 FROM student_school ss WHERE ss.student_id = s.id AND ss.school_id = sqlc.arg(school_id));

-- name: CreateBaseSchedule :execrows
INSERT INTO base_schedule (school_id, teacher_id, room_id, day_of_week, lesson_num, group_id, custom_subject, subject_id, custom_subject_id)
SELECT sqlc.arg(school_id), sqlc.arg(teacher_id), sqlc.arg(room_id), sqlc.arg(day_of_week), sqlc.arg(lesson_num), sqlc.arg(group_id), sqlc.arg(custom_subject), sqlc.arg(subject_id), sqlc.arg(custom_subject_id)
FROM schools s
WHERE s.id = sqlc.arg(school_id)
  AND EXISTS (SELECT 1 FROM teacher_school ts WHERE ts.teacher_id = sqlc.arg(teacher_id) AND ts.school_id = sqlc.arg(school_id))
  AND EXISTS (SELECT 1 FROM rooms r WHERE r.id = sqlc.arg(room_id) AND r.school_id = sqlc.arg(school_id))
  AND EXISTS (SELECT 1 FROM groups g WHERE g.id = sqlc.arg(group_id) AND g.school_id = sqlc.arg(school_id))
  AND (NOT sqlc.arg(custom_subject) OR EXISTS (SELECT 1 FROM custom_subjects cs WHERE cs.id = sqlc.arg(custom_subject_id) AND cs.school_id = sqlc.arg(school_id)));

-- name: DeleteBaseSchedule :execrows
DELETE FROM base_schedule WHERE id = $1 AND school_id = $2;

-- name: UpdateBaseSchedule :execrows
UPDATE base_schedule bs
SET teacher_id = sqlc.arg(teacher_id),
    day_of_week = sqlc.arg(day_of_week),
    lesson_num = sqlc.arg(lesson_num),
    room_id = sqlc.arg(room_id),
    group_id = sqlc.arg(group_id),
    custom_subject = sqlc.arg(custom_subject),
    subject_id = sqlc.arg(subject_id),
    custom_subject_id = sqlc.arg(custom_subject_id)
WHERE bs.id = sqlc.arg(id)
  AND bs.school_id = sqlc.arg(school_id)
  AND EXISTS (SELECT 1 FROM teacher_school ts WHERE ts.teacher_id = sqlc.arg(teacher_id) AND ts.school_id = sqlc.arg(school_id))
  AND EXISTS (SELECT 1 FROM rooms r WHERE r.id = sqlc.arg(room_id) AND r.school_id = sqlc.arg(school_id))
  AND EXISTS (SELECT 1 FROM groups g WHERE g.id = sqlc.arg(group_id) AND g.school_id = sqlc.arg(school_id))
  AND (NOT sqlc.arg(custom_subject) OR EXISTS (SELECT 1 FROM custom_subjects cs WHERE cs.id = sqlc.arg(custom_subject_id) AND cs.school_id = sqlc.arg(school_id)));

-- name: ReadBaseScheduleClass :many
SELECT DISTINCT bs.*, COALESCE(sub.subject_name, cs.subject_name) AS subject_name, EXISTS (SELECT 1 FROM exams e JOIN class_subjects csub ON csub.id = e.class_subjects_id WHERE csub.class_id = sqlc.arg(class_id) AND csub.subject_id = bs.subject_id) AS has_exam, EXISTS (SELECT 1 FROM homework hw JOIN class_subjects csub ON csub.id = hw.class_subjects_id WHERE csub.class_id = sqlc.arg(class_id) AND csub.subject_id = bs.subject_id) AS has_homework, teach.first_name AS teacher_first_name, teach.last_name AS teacher_last_name FROM base_schedule bs JOIN groups g ON bs.group_id = g.id JOIN group_members gm ON g.id = gm.group_id JOIN student_school ss ON ss.student_id = gm.student_id AND ss.school_id = sqlc.arg(school_id) LEFT JOIN subjects sub ON sub.id = bs.subject_id LEFT JOIN custom_subjects cs ON cs.id = bs.custom_subject_id LEFT JOIN teachers teach ON teach.id = bs.teacher_id WHERE ss.class_id = sqlc.arg(class_id) AND bs.school_id = sqlc.arg(school_id) AND gm.school_id = sqlc.arg(school_id);

-- name: ReadBaseScheduleGroup :many
SELECT bs.*, COALESCE(sub.subject_name, cs.subject_name) AS subject_name, EXISTS (SELECT 1 FROM exams e JOIN class_subjects csub ON csub.id = e.class_subjects_id WHERE csub.subject_id = bs.subject_id AND csub.class_id IN (SELECT ss.class_id FROM group_members gm JOIN student_school ss ON ss.student_id = gm.student_id AND ss.school_id = sqlc.arg(school_id) WHERE gm.group_id = bs.group_id)) AS has_exam, EXISTS (SELECT 1 FROM homework hw JOIN class_subjects csub ON csub.id = hw.class_subjects_id WHERE csub.subject_id = bs.subject_id AND csub.class_id IN (SELECT ss.class_id FROM group_members gm JOIN student_school ss ON ss.student_id = gm.student_id AND ss.school_id = sqlc.arg(school_id) WHERE gm.group_id = bs.group_id)) AS has_homework, teach.first_name AS teacher_first_name, teach.last_name AS teacher_last_name FROM base_schedule bs LEFT JOIN subjects sub ON sub.id = bs.subject_id LEFT JOIN custom_subjects cs ON cs.id = bs.custom_subject_id LEFT JOIN teachers teach ON teach.id = bs.teacher_id WHERE bs.group_id = sqlc.arg(group_id) AND bs.school_id = sqlc.arg(school_id);

-- name: ReadRealTimeTable :many
SELECT d.actual_date::date AS actual_date, COALESCE(t.room_id, b.room_id) AS room_id, b.lesson_num, b.day_of_week, COALESCE(t.substitution_teacher_id, t.teacher_id, b.teacher_id) AS effective_teacher_id, COALESCE(t.group_id, b.group_id) AS group_id, COALESCE(t.school_id, b.school_id) AS school_id, COALESCE(t.custom_subject, b.custom_subject) AS custom_subject, COALESCE(t.subject_id, b.subject_id) AS subject_id, COALESCE(t.custom_subject_id, b.custom_subject_id) AS custom_subject_id, COALESCE(sub.subject_name, cs.subject_name) AS subject_name, COALESCE(t.is_substitution, FALSE) AS is_substitution, COALESCE(t.canceled, FALSE) AS canceled, EXISTS (SELECT 1 FROM exams e JOIN class_subjects csub ON csub.id = e.class_subjects_id WHERE csub.subject_id = COALESCE(t.subject_id, b.subject_id) AND csub.class_id IN (SELECT ss.class_id FROM group_members gm JOIN student_school ss ON ss.student_id = gm.student_id AND ss.school_id = sqlc.arg(school_id) WHERE gm.group_id = COALESCE(t.group_id, b.group_id)) AND e.date = d.actual_date::date) AS has_exam, EXISTS (SELECT 1 FROM homework hw JOIN class_subjects csub ON csub.id = hw.class_subjects_id WHERE csub.subject_id = COALESCE(t.subject_id, b.subject_id) AND csub.class_id IN (SELECT ss.class_id FROM group_members gm JOIN student_school ss ON ss.student_id = gm.student_id AND ss.school_id = sqlc.arg(school_id) WHERE gm.group_id = COALESCE(t.group_id, b.group_id)) AND hw.due_date = d.actual_date::date) AS has_homework, teach.first_name AS teacher_first_name, teach.last_name AS teacher_last_name FROM base_schedule b JOIN generate_series(CAST(sqlc.arg(start_date) AS date), CAST(sqlc.arg(end_date) AS date), interval '1 day') AS d(actual_date) ON EXTRACT(ISODOW FROM d.actual_date) = b.day_of_week LEFT JOIN time_table t ON t.school_id = b.school_id AND t.group_id = b.group_id AND t.day_of_week = b.day_of_week AND t.lesson_num = b.lesson_num AND t.actual_date = d.actual_date::date LEFT JOIN subjects sub ON sub.id = COALESCE(t.subject_id, b.subject_id) LEFT JOIN custom_subjects cs ON cs.id = COALESCE(t.custom_subject_id, b.custom_subject_id) LEFT JOIN teachers teach ON teach.id = COALESCE(t.substitution_teacher_id, t.teacher_id, b.teacher_id) WHERE b.school_id = sqlc.arg(school_id) AND b.group_id IN ( SELECT g.id FROM groups g JOIN group_members gm ON g.id = gm.group_id JOIN student_school ss ON ss.student_id = gm.student_id AND ss.school_id = sqlc.arg(school_id) WHERE ss.class_id = sqlc.arg(class_id)) ORDER BY d.actual_date, b.lesson_num;

-- name: ReadMyRealTimeTable :many
SELECT d.actual_date::date AS actual_date, COALESCE(t.room_id, b.room_id) AS room_id, b.lesson_num, b.day_of_week, COALESCE(t.substitution_teacher_id, t.teacher_id, b.teacher_id) AS effective_teacher_id, COALESCE(t.group_id, b.group_id) AS group_id, COALESCE(t.school_id, b.school_id) AS school_id, COALESCE(t.custom_subject, b.custom_subject) AS custom_subject, COALESCE(t.subject_id, b.subject_id) AS subject_id, COALESCE(t.custom_subject_id, b.custom_subject_id) AS custom_subject_id, COALESCE(sub.subject_name, cs.subject_name) AS subject_name, COALESCE(t.is_substitution, FALSE) AS is_substitution, COALESCE(t.canceled, FALSE) AS canceled, EXISTS (SELECT 1 FROM exams e JOIN class_subjects csub ON csub.id = e.class_subjects_id WHERE csub.subject_id = COALESCE(t.subject_id, b.subject_id) AND csub.class_id IN (SELECT ss.class_id FROM group_members gm JOIN student_school ss ON ss.student_id = gm.student_id AND ss.school_id = sqlc.arg(school_id) WHERE gm.group_id = COALESCE(t.group_id, b.group_id)) AND e.date = d.actual_date::date) AS has_exam, EXISTS (SELECT 1 FROM homework hw JOIN class_subjects csub ON csub.id = hw.class_subjects_id WHERE csub.subject_id = COALESCE(t.subject_id, b.subject_id) AND csub.class_id IN (SELECT ss.class_id FROM group_members gm JOIN student_school ss ON ss.student_id = gm.student_id AND ss.school_id = sqlc.arg(school_id) WHERE gm.group_id = COALESCE(t.group_id, b.group_id)) AND hw.due_date = d.actual_date::date) AS has_homework, teach.first_name AS teacher_first_name, teach.last_name AS teacher_last_name FROM base_schedule b JOIN generate_series(CAST(sqlc.arg(start_date) AS date), CAST(sqlc.arg(end_date) AS date), interval '1 day') AS d(actual_date) ON EXTRACT(ISODOW FROM d.actual_date) = b.day_of_week LEFT JOIN time_table t ON t.school_id = b.school_id AND t.group_id = b.group_id AND t.day_of_week = b.day_of_week AND t.lesson_num = b.lesson_num AND t.actual_date = d.actual_date::date LEFT JOIN subjects sub ON sub.id = COALESCE(t.subject_id, b.subject_id) LEFT JOIN custom_subjects cs ON cs.id = COALESCE(t.custom_subject_id, b.custom_subject_id) LEFT JOIN teachers teach ON teach.id = COALESCE(t.substitution_teacher_id, t.teacher_id, b.teacher_id) WHERE b.school_id = sqlc.arg(school_id) AND (b.teacher_id = sqlc.arg(teacher_id) OR t.substitution_teacher_id = sqlc.arg(teacher_id)) ORDER BY d.actual_date, b.lesson_num;

-- name: ReadTeacherLessonTimes :many
SELECT DISTINCT bs.* FROM bell_schedule bs
JOIN groups g ON g.bell_id = bs.type_id AND g.school_id = bs.school_id
JOIN base_schedule b ON b.group_id = g.id
WHERE b.school_id = $1 AND b.teacher_id = $2;

-- name: AddCanceledLesson :execrows
INSERT INTO time_table (school_id, teacher_id, room_id, day_of_week, group_id, custom_subject, custom_subject_id, subject_id, actual_date, lesson_num, canceled)
SELECT sqlc.arg(school_id), sqlc.arg(teacher_id), sqlc.arg(room_id), sqlc.arg(day_of_week), sqlc.arg(group_id), sqlc.arg(custom_subject), sqlc.arg(custom_subject_id), sqlc.arg(subject_id), sqlc.arg(actual_date), sqlc.arg(lesson_num), true
FROM schools s
WHERE s.id = sqlc.arg(school_id)
  AND EXISTS (SELECT 1 FROM teacher_school ts WHERE ts.teacher_id = sqlc.arg(teacher_id) AND ts.school_id = sqlc.arg(school_id))
  AND EXISTS (SELECT 1 FROM rooms r WHERE r.id = sqlc.arg(room_id) AND r.school_id = sqlc.arg(school_id))
  AND EXISTS (SELECT 1 FROM groups g WHERE g.id = sqlc.arg(group_id) AND g.school_id = sqlc.arg(school_id))
  AND (NOT sqlc.arg(custom_subject) OR EXISTS (SELECT 1 FROM custom_subjects cs WHERE cs.id = sqlc.arg(custom_subject_id) AND cs.school_id = sqlc.arg(school_id)))
ON CONFLICT (school_id, room_id, actual_date, lesson_num, group_id) DO UPDATE SET canceled = true;

-- name: RemoveCanceledLesson :execrows
INSERT INTO time_table (school_id, teacher_id, room_id, day_of_week, group_id, custom_subject, custom_subject_id, subject_id, actual_date, lesson_num, canceled)
SELECT sqlc.arg(school_id), sqlc.arg(teacher_id), sqlc.arg(room_id), sqlc.arg(day_of_week), sqlc.arg(group_id), sqlc.arg(custom_subject), sqlc.arg(custom_subject_id), sqlc.arg(subject_id), sqlc.arg(actual_date), sqlc.arg(lesson_num), false
FROM schools s
WHERE s.id = sqlc.arg(school_id)
  AND EXISTS (SELECT 1 FROM teacher_school ts WHERE ts.teacher_id = sqlc.arg(teacher_id) AND ts.school_id = sqlc.arg(school_id))
  AND EXISTS (SELECT 1 FROM rooms r WHERE r.id = sqlc.arg(room_id) AND r.school_id = sqlc.arg(school_id))
  AND EXISTS (SELECT 1 FROM groups g WHERE g.id = sqlc.arg(group_id) AND g.school_id = sqlc.arg(school_id))
  AND (NOT sqlc.arg(custom_subject) OR EXISTS (SELECT 1 FROM custom_subjects cs WHERE cs.id = sqlc.arg(custom_subject_id) AND cs.school_id = sqlc.arg(school_id)))
ON CONFLICT (school_id, room_id, actual_date, lesson_num, group_id) DO UPDATE SET canceled = false;

-- name: ManageSubsitutionLesson :execrows
UPDATE time_table t
SET teacher_id = sqlc.arg(teacher_id),
    room_id = sqlc.arg(room_id),
    day_of_week = sqlc.arg(day_of_week),
    group_id = sqlc.arg(group_id),
    custom_subject = sqlc.arg(custom_subject),
    custom_subject_id = sqlc.arg(custom_subject_id),
    subject_id = sqlc.arg(subject_id),
    actual_date = sqlc.arg(actual_date),
    lesson_num = sqlc.arg(lesson_num),
    is_substitution = sqlc.arg(is_substitution),
    substitution_teacher_id = sqlc.arg(substitution_teacher_id)
WHERE t.school_id = sqlc.arg(school_id)
  AND t.id = sqlc.arg(id)
  AND EXISTS (SELECT 1 FROM teacher_school ts WHERE ts.teacher_id = sqlc.arg(teacher_id) AND ts.school_id = sqlc.arg(school_id))
  AND EXISTS (SELECT 1 FROM rooms r WHERE r.id = sqlc.arg(room_id) AND r.school_id = sqlc.arg(school_id))
  AND EXISTS (SELECT 1 FROM groups g WHERE g.id = sqlc.arg(group_id) AND g.school_id = sqlc.arg(school_id))
  AND (NOT sqlc.arg(custom_subject) OR EXISTS (SELECT 1 FROM custom_subjects cs WHERE cs.id = sqlc.arg(custom_subject_id) AND cs.school_id = sqlc.arg(school_id)))
  AND (NOT sqlc.arg(is_substitution) OR EXISTS (SELECT 1 FROM teacher_school ts2 WHERE ts2.teacher_id = sqlc.arg(substitution_teacher_id) AND ts2.school_id = sqlc.arg(school_id)));

-- name: CreateRealTimeLesson :execrows
INSERT INTO time_table (school_id, teacher_id, room_id, day_of_week, group_id, custom_subject, custom_subject_id, subject_id, actual_date, lesson_num)
SELECT sqlc.arg(school_id), sqlc.arg(teacher_id), sqlc.arg(room_id), sqlc.arg(day_of_week), sqlc.arg(group_id), sqlc.arg(custom_subject), sqlc.arg(custom_subject_id), sqlc.arg(subject_id), sqlc.arg(actual_date), sqlc.arg(lesson_num)
FROM schools s
WHERE s.id = sqlc.arg(school_id)
  AND EXISTS (SELECT 1 FROM teacher_school ts WHERE ts.teacher_id = sqlc.arg(teacher_id) AND ts.school_id = sqlc.arg(school_id))
  AND EXISTS (SELECT 1 FROM rooms r WHERE r.id = sqlc.arg(room_id) AND r.school_id = sqlc.arg(school_id))
  AND EXISTS (SELECT 1 FROM groups g WHERE g.id = sqlc.arg(group_id) AND g.school_id = sqlc.arg(school_id))
  AND (NOT sqlc.arg(custom_subject) OR EXISTS (SELECT 1 FROM custom_subjects cs WHERE cs.id = sqlc.arg(custom_subject_id) AND cs.school_id = sqlc.arg(school_id)));

-- name: UpdateRealTimeLesson :execrows
UPDATE time_table t
SET teacher_id = sqlc.arg(teacher_id),
    room_id = sqlc.arg(room_id),
    day_of_week = sqlc.arg(day_of_week),
    group_id = sqlc.arg(group_id),
    custom_subject = sqlc.arg(custom_subject),
    custom_subject_id = sqlc.arg(custom_subject_id),
    subject_id = sqlc.arg(subject_id),
    actual_date = sqlc.arg(actual_date),
    lesson_num = sqlc.arg(lesson_num)
WHERE t.school_id = sqlc.arg(school_id)
  AND t.id = sqlc.arg(id)
  AND EXISTS (SELECT 1 FROM teacher_school ts WHERE ts.teacher_id = sqlc.arg(teacher_id) AND ts.school_id = sqlc.arg(school_id))
  AND EXISTS (SELECT 1 FROM rooms r WHERE r.id = sqlc.arg(room_id) AND r.school_id = sqlc.arg(school_id))
  AND EXISTS (SELECT 1 FROM groups g WHERE g.id = sqlc.arg(group_id) AND g.school_id = sqlc.arg(school_id))
  AND (NOT sqlc.arg(custom_subject) OR EXISTS (SELECT 1 FROM custom_subjects cs WHERE cs.id = sqlc.arg(custom_subject_id) AND cs.school_id = sqlc.arg(school_id)));

-- name: DeleteRealTimeLesson :execrows
DELETE FROM time_table WHERE school_id = $1 AND id = $2;

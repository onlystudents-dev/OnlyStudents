-- name: CreateBellScheduleType :exec
INSERT INTO bell_schedule_type (school_id, name) VALUES ($1, $2);

-- name: DeleteBellScheduleType :exec
DELETE FROM bell_schedule_type WHERE school_id = $1 AND id = $2;

-- name: EditBellScheduleType :exec
UPDATE bell_schedule_type SET name = $1 WHERE id = $2 AND school_id = $3;

-- name: ReadBellScheduleType :many
SELECT id, name FROM bell_schedule_type WHERE school_id = $1;

-- name: CreateLessonTime :exec
INSERT INTO bell_schedule (school_id, type_id, lesson_number, at_start, at_end) VALUES ($1, $2, $3, $4, $5);

-- name: DeleteLessonTime :exec
DELETE FROM bell_schedule WHERE school_id = $1 AND id = $2;

-- name: EditLessonTime :exec
UPDATE bell_schedule set lesson_number = $1, at_start = $2, at_end = $3 WHERE id = $4 AND school_id = $5;

-- name: ReadLessonTime :many
SELECT * FROM bell_schedule WHERE school_id = $1 AND type_id = $2;

-- name: CreateCustomSubject :exec
INSERT INTO custom_subjects (school_id, subject_name) VALUES ($1, $2);

-- name: EditCustomSubject :exec
UPDATE custom_subjects SET subject_name = $1 WHERE school_id = $2 AND id = $3;

-- name: DeleteCustomSubject :exec
DELETE FROM custom_subjects WHERE school_id = $1 AND id = $2;

-- name: ReadCustomSubject :many
SELECT id, subject_name FROM custom_subjects WHERE school_id = $1;

-- name: CreateRoom :exec
INSERT INTO rooms (school_id, name, capacity) VALUES ($1, $2, $3);

-- name: EditRoom :exec
UPDATE rooms SET name = $1, capacity = $2 WHERE school_id = $3 AND id = $4;

-- name: DeleteRoom :exec
DELETE FROM rooms WHERE school_id = $1 AND id = $2;

-- name: ReadRoom :many
SELECT id, name, capacity FROM rooms WHERE school_id = $1;

-- name: CreateGroup :exec
INSERT INTO groups (school_id, bell_id, group_name) VALUES ($1, $2, $3);

-- name: EditGroup :exec
UPDATE groups SET group_name = $1 WHERE school_id = $2 AND id = $3;

-- name: DeleteGroup :exec
DELETE FROM groups WHERE school_id = $1 AND id = $2;

-- name: ReadGroup :many
SELECT * FROM groups WHERE school_id = $1;

-- name: InsertStudentToGroup :exec
INSERT INTO group_members (school_id, group_id, student_id) VALUES ($1, $2, $3);

-- name: DeleteStudentFromGroup :exec
DELETE FROM group_members WHERE group_id = $1 AND student_id = $2 AND school_id = $3;

-- name: ReadListOfStudents :many
SELECT s.id, s.first_name, s.last_name FROM students AS s INNER JOIN group_members AS g_m ON g_m.student_id = s.id WHERE g_m.group_id = $1 AND g_m.school_id = $2;

-- name: CreateBaseSchedule :exec
INSERT INTO base_schedule (school_id, teacher_id, room_id, day_of_week, lesson_num, group_id, custom_subject, subject_id, custom_subject_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);

-- name: DeleteBaseSchedule :exec
DELETE FROM base_schedule WHERE id = $1 AND school_id = $2;

-- name: UpdateBaseSchedule :exec
UPDATE base_schedule SET teacher_id = $1, day_of_week = $2, lesson_num = $3, room_id = $4, group_id = $5, custom_subject = $6, subject_id = $7, custom_subject_id = $8 WHERE id = $9 AND school_id = $10;

-- name: ReadBaseScheduleClass :many
SELECT DISTINCT bs.* FROM base_schedule bs JOIN groups g ON bs.group_id = g.id JOIN group_members gm ON g.id = gm.group_id JOIN students s ON gm.student_id = s.id WHERE s.classes_id = $1 AND s.school_id = $2;

-- name: ReadBaseScheduleGroup :many
SELECT * FROM base_schedule WHERE group_id = $1 AND school_id = $2;

-- name: ReadRealTimeTable :many
SELECT  COALESCE(t.room_id, b.room_id) AS room_id, COALESCE(t.lesson_num, b.lesson_num) AS lesson_num, COALESCE(t.day_of_week, b.day_of_week) AS day_of_week, COALESCE(t.substitution_teacher_id, t.teacher_id, b.teacher_id) AS effective_teacher_id, COALESCE(t.group_id, b.group_id) AS group_id, COALESCE(t.school_id, b.school_id) AS school_id, COALESCE(t.custom_subject, b.custom_subject) AS custom_subject, COALESCE(t.subject_id, b.subject_id) AS subject_id, COALESCE(t.custom_subject_id, b.custom_subject_id) AS custom_subject_id, COALESCE(t.is_substitution, FALSE) AS is_substitution, COALESCE(t.canceled, FALSE) AS canceled FROM base_schedule b LEFT JOIN time_table t ON t.school_id = b.school_id AND t.group_id = b.group_id AND t.day_of_week = b.day_of_week AND t.lesson_num = b.lesson_num AND t.actual_date BETWEEN $1 AND $2 WHERE b.group_id IN ( SELECT g.id FROM groups g JOIN group_members gm ON g.id = gm.group_id JOIN students s ON gm.student_id = s.id WHERE s.classes_id = $3 AND s.school_id = $4);

-- name: AddCanceledLesson :exec
INSERT INTO time_table (school_id, teacher_id, room_id, day_of_week, group_id, custom_subject, custom_subject_id, subject_id, actual_date, lesson_num, canceled) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, true) ON CONFLICT (school_id, room_id, actual_date, lesson_num, group_id) DO UPDATE SET canceled = true;

-- name: RemoveCanceledLesson :exec
INSERT INTO time_table (school_id, teacher_id, room_id, day_of_week, group_id, custom_subject, custom_subject_id, subject_id, actual_date, lesson_num, canceled) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, false) ON CONFLICT (school_id, room_id, actual_date, lesson_num, group_id) DO UPDATE SET canceled = false;

-- name: ManageSubsitutionLesson :exec
UPDATE time_table SET teacher_id = $2, room_id = $3, day_of_week = $4, group_id = $5, custom_subject = $6, custom_subject_id = $7, subject_id = $8, actual_date = $9, lesson_num = $10, is_substitution = $11, substitution_teacher_id = $12 WHERE school_id = $1 AND id = $13;

-- name: CreateRealTimeLesson :exec
INSERT INTO time_table (school_id, teacher_id, room_id, day_of_week, group_id, custom_subject, custom_subject_id, subject_id, actual_date, lesson_num) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);

-- name: UpdateRealTimeLesson :exec
UPDATE time_table SET teacher_id = $1, room_id = $2, day_of_week = $3, group_id = $4, custom_subject = $5, custom_subject_id = $6, subject_id = $7, actual_date = $8, lesson_num = $9 WHERE school_id = $10 AND id = $11;

-- name: DeleteRealTimeLesson :exec
DELETE FROM time_table WHERE school_id = $1 AND id = $2;

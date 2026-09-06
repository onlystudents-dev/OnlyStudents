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
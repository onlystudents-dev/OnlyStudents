-- name: CreateBellScheduleType :exec
INSERT INTO bell_schedule_type (school_id, name) VALUES ($1, $2);

-- name: DeleteBellScheduleType :exec
DELETE FROM bell_schedule_type WHERE school_id = $1 AND id = $2;

-- name: EditBellScheduleType :exec
UPDATE bell_schedule_type SET name = $1 WHERE id = $2 AND school_id = $3;

-- name: ReadBellScheduleType :many
SELECT id, name FROM bell_schedule_type WHERE school_id = $1;

-- name: CreateCustomSubject :exec
INSERT INTO custom_subjects (school_id, subject_name) VALUES ($1, $2);

-- name: EditCustomSubject :exec
UPDATE custom_subjects SET subject_name = $1 WHERE school_id = $2 AND id = $3;

-- name: DeleteCustomSubject :exec
DELETE FROM custom_subjects WHERE school_id = $1 AND id = $2;

-- name: ReadCustomSubject :many
SELECT id, subject_name FROM custom_subjects WHERE school_id = $1;
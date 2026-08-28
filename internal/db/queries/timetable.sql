-- name: CreateBellScheduleType :exec
INSERT INTO bell_schedule_type (school_id, name) VALUES ($1, $2);
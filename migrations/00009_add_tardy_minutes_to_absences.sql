-- +goose Up
ALTER TABLE absences ADD COLUMN tardy_minutes INT DEFAULT NULL;

-- +goose Down
ALTER TABLE absences DROP COLUMN tardy_minutes;

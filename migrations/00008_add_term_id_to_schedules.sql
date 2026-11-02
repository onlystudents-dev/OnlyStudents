-- +goose Up
-- Resolve the term that belongs to a school. When a date is supplied, the term
-- whose date range contains that date wins; otherwise the active term is used.
-- +goose StatementBegin
CREATE FUNCTION resolve_term_id(p_school_id INT, p_date DATE DEFAULT NULL)
RETURNS INT
LANGUAGE sql
STABLE
AS $$
    SELECT t.id
    FROM terms t
    JOIN school_years sy ON sy.id = t.school_year_id
    WHERE sy.school_id = p_school_id
    ORDER BY
        CASE WHEN p_date IS NOT NULL AND p_date BETWEEN t.start_date AND t.end_date THEN 0 ELSE 1 END,
        t.is_active DESC,
        sy.is_active DESC,
        t.start_date DESC,
        t.id DESC
    LIMIT 1;
$$;
-- +goose StatementEnd

ALTER TABLE base_schedule ADD COLUMN term_id INT REFERENCES terms (id) ON DELETE RESTRICT;
ALTER TABLE time_table ADD COLUMN term_id INT REFERENCES terms (id) ON DELETE RESTRICT;

UPDATE base_schedule bs
SET term_id = resolve_term_id(bs.school_id, NULL)
WHERE bs.term_id IS NULL;

UPDATE time_table t
SET term_id = resolve_term_id(t.school_id, t.actual_date)
WHERE t.term_id IS NULL;

ALTER TABLE base_schedule ALTER COLUMN term_id SET NOT NULL;
ALTER TABLE time_table ALTER COLUMN term_id SET NOT NULL;

-- A base schedule is now term scoped, so the same room/lesson slot may exist
-- in more than one term.
ALTER TABLE base_schedule DROP CONSTRAINT base_schedule_day_of_week_lesson_num_room_id_key;
ALTER TABLE base_schedule ADD CONSTRAINT uq_base_schedule_term_slot UNIQUE (term_id, day_of_week, lesson_num, room_id);

CREATE INDEX idx_base_schedule_term_id ON base_schedule (term_id);
CREATE INDEX idx_time_table_term_id ON time_table (term_id);

-- Keep inserts/updates working for callers that do not know the term yet.
-- +goose StatementBegin
CREATE FUNCTION fill_schedule_term_id()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
    IF TG_TABLE_NAME = 'time_table' THEN
        IF NEW.term_id IS NULL OR NEW.term_id = 0 OR (TG_OP = 'UPDATE' AND NEW.actual_date IS DISTINCT FROM OLD.actual_date) THEN
            NEW.term_id := resolve_term_id(NEW.school_id, NEW.actual_date);
        END IF;
    ELSIF NEW.term_id IS NULL OR NEW.term_id = 0 THEN
        NEW.term_id := resolve_term_id(NEW.school_id, NULL);
    END IF;

    RETURN NEW;
END;
$$;
-- +goose StatementEnd

CREATE TRIGGER trg_base_schedule_term_id
BEFORE INSERT OR UPDATE ON base_schedule
FOR EACH ROW EXECUTE FUNCTION fill_schedule_term_id();

CREATE TRIGGER trg_time_table_term_id
BEFORE INSERT OR UPDATE ON time_table
FOR EACH ROW EXECUTE FUNCTION fill_schedule_term_id();

-- +goose Down
DROP TRIGGER trg_time_table_term_id ON time_table;
DROP TRIGGER trg_base_schedule_term_id ON base_schedule;
DROP FUNCTION fill_schedule_term_id();

DROP INDEX idx_time_table_term_id;
DROP INDEX idx_base_schedule_term_id;

ALTER TABLE base_schedule DROP CONSTRAINT uq_base_schedule_term_slot;
ALTER TABLE base_schedule ADD CONSTRAINT base_schedule_day_of_week_lesson_num_room_id_key UNIQUE (day_of_week, lesson_num, room_id);

ALTER TABLE time_table DROP COLUMN term_id;
ALTER TABLE base_schedule DROP COLUMN term_id;

DROP FUNCTION resolve_term_id(INT, DATE);

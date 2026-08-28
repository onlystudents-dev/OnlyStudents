CREATE TABLE calendar_events (
    id BIGSERIAL PRIMARY KEY,
    school_id INT DEFAULT NULL REFERENCES schools (id) ON DELETE CASCADE,
    class_id INT DEFAULT NULL REFERENCES classes (id) ON DELETE CASCADE,
    title VARCHAR(256) NOT NULL,
    description TEXT DEFAULT NULL,
    start TIMESTAMPTZ NOT NULL,
    "end" TIMESTAMPTZ DEFAULT NULL,
    all_day BOOLEAN NOT NULL DEFAULT FALSE,
    event_type VARCHAR(32) NOT NULL,
    CONSTRAINT chk_calendar_events_time CHECK ("end" IS NULL OR "end" >= start),
    CONSTRAINT chk_calendar_events_scope CHECK (school_id IS NOT NULL OR class_id IS NOT NULL)
);

CREATE INDEX idx_calendar_events_school_id ON calendar_events (school_id);
CREATE INDEX idx_calendar_events_class_id ON calendar_events (class_id);
CREATE INDEX idx_calendar_events_start ON calendar_events (start);

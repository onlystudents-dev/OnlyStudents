CREATE TABLE announcements (
    id BIGSERIAL PRIMARY KEY,
    author_account_id UUID NOT NULL REFERENCES accounts (id) ON DELETE RESTRICT,
    school_id INT DEFAULT NULL REFERENCES schools (id) ON DELETE CASCADE,
    class_id INT DEFAULT NULL REFERENCES classes (id) ON DELETE CASCADE,
    title VARCHAR(256) NOT NULL,
    body TEXT NOT NULL,
    pinned BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_announcements_scope CHECK (school_id IS NOT NULL OR class_id IS NOT NULL)
);

CREATE INDEX idx_announcements_author_account_id ON announcements (author_account_id);
CREATE INDEX idx_announcements_school_id ON announcements (school_id);
CREATE INDEX idx_announcements_class_id ON announcements (class_id);
CREATE INDEX idx_announcements_created_at ON announcements (created_at);

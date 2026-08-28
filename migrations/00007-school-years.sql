CREATE TABLE school_years (
    id SERIAL PRIMARY KEY,
    school_id INT NOT NULL REFERENCES schools (id) ON DELETE CASCADE,
    name VARCHAR(64) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT FALSE,
    CONSTRAINT chk_school_years_dates CHECK (end_date > start_date)
);

CREATE TABLE terms (
    id SERIAL PRIMARY KEY,
    school_year_id INT NOT NULL REFERENCES school_years (id) ON DELETE CASCADE,
    name VARCHAR(64) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    grade_deadline DATE NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT FALSE,
    CONSTRAINT chk_terms_dates CHECK (end_date >= start_date AND grade_deadline <= end_date)
);

CREATE INDEX idx_school_years_school_id ON school_years (school_id);
CREATE INDEX idx_terms_school_year_id ON terms (school_year_id);
CREATE UNIQUE INDEX uq_school_years_active ON school_years (school_id) WHERE is_active;
CREATE UNIQUE INDEX uq_terms_active ON terms (school_year_id) WHERE is_active;

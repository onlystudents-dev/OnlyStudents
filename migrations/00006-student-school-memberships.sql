CREATE TABLE student_school (
    id SERIAL PRIMARY KEY,
    student_id INT NOT NULL REFERENCES students (id) ON DELETE CASCADE,
    school_id  INT NOT NULL REFERENCES schools (id)  ON DELETE CASCADE,
    class_id INT NOT NULL REFERENCES classes (id)  ON DELETE CASCADE,
    id_number  INT NOT NULL,
    CONSTRAINT uq_student_school          UNIQUE (student_id, school_id),
    CONSTRAINT uq_student_school_id_number UNIQUE (school_id, id_number)
);
CREATE INDEX idx_student_school_school_id  ON student_school (school_id);
CREATE INDEX idx_student_school_class_id ON student_school (class_id);

INSERT INTO student_school (student_id, school_id, class_id, id_number)
SELECT s.id, s.school_id, s.class_id, s.id_number FROM students s;

ALTER TABLE students DROP CONSTRAINT fk_students_school;
ALTER TABLE students DROP CONSTRAINT fk_students_class;
DROP INDEX idx_students_school_id;
DROP INDEX idx_students_class_id;
ALTER TABLE students DROP COLUMN school_id;
ALTER TABLE students DROP COLUMN class_id;
ALTER TABLE students DROP COLUMN id_number;

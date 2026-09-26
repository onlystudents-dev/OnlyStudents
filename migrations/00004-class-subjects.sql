CREATE TABLE class_subjects (
    id SERIAL PRIMARY KEY,
    school_id INT NOT NULL REFERENCES schools(id),
    class_id INT NOT NULL REFERENCES classes(id),
    subject_id INT REFERENCES subjects(id),
    custom_subject bool NOT NULL DEFAULT FALSE,
    custom_subject_id INT REFERENCES custom_subjects(id),
    teacher_id INT NOT NULL REFERENCES teachers(id),
    CONSTRAINT chk_class_subjects CHECK (
      (custom_subject = FALSE AND subject_id IS NOT NULL AND custom_subject_id IS NULL) OR
      (custom_subject = TRUE  AND custom_subject_id IS NOT NULL AND subject_id IS NULL)),
    CONSTRAINT uq_class_subjects UNIQUE (class_id, subject_id, custom_subject_id)
);

ALTER TABLE grades
    ADD CONSTRAINT fk_grades_class_subjects
    FOREIGN KEY (class_subjects_id) REFERENCES class_subjects (id) ON DELETE RESTRICT;

ALTER TABLE final_grades
    ADD CONSTRAINT fk_final_grades_class_subjects
    FOREIGN KEY (class_subjects_id) REFERENCES class_subjects (id) ON DELETE RESTRICT;

ALTER TABLE exams
    ADD CONSTRAINT fk_exams_class_subjects
    FOREIGN KEY (class_subjects_id) REFERENCES class_subjects (id) ON DELETE CASCADE;

ALTER TABLE homework
    ADD CONSTRAINT fk_homework_class_subjects
    FOREIGN KEY (class_subjects_id) REFERENCES class_subjects (id) ON DELETE CASCADE;

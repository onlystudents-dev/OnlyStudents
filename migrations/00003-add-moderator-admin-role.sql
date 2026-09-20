ALTER TABLE accounts DROP CONSTRAINT chk_accounts_role_link;

ALTER TABLE accounts ADD CONSTRAINT chk_accounts_role_link CHECK (
    (role = 'student' AND student_id IS NOT NULL AND teacher_id IS NULL AND guardian_id IS NULL) OR
    (role = 'teacher' AND teacher_id IS NOT NULL AND student_id IS NULL AND guardian_id IS NULL) OR
    (role = 'guardian' AND guardian_id IS NOT NULL AND student_id IS NULL AND teacher_id IS NULL) OR
    (role = 'moderator' AND guardian_id IS NULL AND student_id IS NULL AND teacher_id IS NULL) OR
    (role = 'admin' AND guardian_id IS NULL AND student_id IS NULL AND teacher_id IS NULL)
)

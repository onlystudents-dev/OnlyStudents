ALTER TABLE accounts ADD COLUMN email_address VARCHAR(256);

UPDATE accounts a
SET email_address = s.email_address
FROM students s
WHERE a.role = 'student' AND a.student_id = s.id;

UPDATE accounts a
SET email_address = t.email_address
FROM teachers t
WHERE a.role = 'teacher' AND a.teacher_id = t.id;

UPDATE accounts a
SET email_address = g.email_address
FROM guardians g
WHERE a.role = 'guardian' AND a.guardian_id = g.id;

ALTER TABLE accounts ADD CONSTRAINT uq_accounts_email_address UNIQUE (email_address);

ALTER TABLE students DROP COLUMN email_address;
ALTER TABLE teachers DROP COLUMN email_address;
ALTER TABLE guardians DROP COLUMN email_address;

INSERT INTO accounts (role, password_hash, student_id)
SELECT 'student', '$argon2id$v=19$m=65536,t=1,p=16$uVcQrvROzMvx8pEKN8eMlA$wfXFSmeuUNDWKTZ3edtmeA+VjXec53r1++TE9n5Z5Jg', s.id
FROM students s
WHERE s.id_number = 123456789
  AND NOT EXISTS (
    SELECT 1 FROM accounts a
    WHERE a.student_id = s.id AND a.role = 'student'
);
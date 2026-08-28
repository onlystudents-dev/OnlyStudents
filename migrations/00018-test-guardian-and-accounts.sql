INSERT INTO guardians (
    id, email_address, phone_number, first_name, last_name,
    birth_first_name, birth_last_name, birth_date, birth_city, birth_country,
    permament_address, temporary_address
) VALUES (
    1, 'guardian@school.test', '+0000000001', 'Test', 'Guardian',
    'Test', 'Guardian', '1980-01-01', 'City', 'HU',
    'Addr', 'Addr'
)
ON CONFLICT (id) DO NOTHING;

INSERT INTO accounts (role, password_hash, teacher_id)
SELECT 'teacher', '$argon2id$v=19$m=65536,t=1,p=16$dVZjUXJ2Uk96TXZ4OHBFS044ZU1sQQ$yt9FrwyGnwXP4H7Go14ot4R7DgD8BoEguVNTfmF3mfM', t.id
FROM teachers t
WHERE t.id = 1
  AND NOT EXISTS (
    SELECT 1 FROM accounts a
    WHERE a.teacher_id = t.id AND a.role = 'teacher'
  );

INSERT INTO accounts (role, password_hash, guardian_id)
SELECT 'guardian', '$argon2id$v=19$m=65536,t=1,p=16$dVZjUXJ2Uk96TXZ4OHBFS044ZU1sQQ$yt9FrwyGnwXP4H7Go14ot4R7DgD8BoEguVNTfmF3mfM', g.id
FROM guardians g
WHERE g.id = 1
  AND NOT EXISTS (
    SELECT 1 FROM accounts a
    WHERE a.guardian_id = g.id AND a.role = 'guardian'
  );

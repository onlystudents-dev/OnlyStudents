INSERT INTO guardians_access (student_id, guardian_id, legal_representative)
SELECT s.id, 1, TRUE
FROM students s
WHERE s.id_number = 123456789
  AND NOT EXISTS (
    SELECT 1 FROM guardians_access ga
    WHERE ga.student_id = s.id AND ga.guardian_id = 1
  );
INSERT INTO students (
	id_number, school_id, first_name, last_name,
	birth_first_name, birth_last_name, birth_date, birth_city, birth_country,
	mother_birth_first_name, mother_birth_last_name,
	classes_id, permament_address, temporary_address, ssn_number,
	bank_name, iban_owner, iban_number, document_type, document_number,
	password_hash, password_salt
) VALUES (
	123456789, 1, 'Test', 'Student',
	'Test', 'Student', '2005-01-01', 'TestCity', 'HU',
	'MotherTest', 'MotherStudent',
	1, 'PermAddr', 'TempAddr', 987654321,
	'TestBank', 'TestOwner', 'HU12345678901234567890123456', 'ID', 'DOC123',
	'hash', 'salt'
);

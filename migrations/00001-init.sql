CREATE TABLE schools (
	id SERIAL PRIMARY KEY,
	name VARCHAR(256) NOT NULL,
	zip_code VARCHAR(100) NOT NULL,
    city VARCHAR(100) NOT NULL,
    address_line VARCHAR(256) NOT NULL,
	school_type INT NOT NULL,
    principal_id INT NOT NULL,
    phone_number VARCHAR(30) NOT NULL,
    email_address VARCHAR(255) NOT NULL
);

CREATE TABLE school_type (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE students (
    id SERIAL PRIMARY KEY,
    id_number INT NOT NULL UNIQUE,
    school_id INT NOT NULL,
    phone_number VARCHAR(30) DEFAULT(NULL),
    email_address VARCHAR(255) DEFAULT(NULL),
    first_name VARCHAR(256) NOT NULL,
    last_name VARCHAR(256) NOT NULL,
    birth_first_name VARCHAR(256) NOT NULL,
    birth_last_name VARCHAR(256) NOT NULL,
    birth_date DATE NOT NULL,
    birth_city VARCHAR(100) NOT NULL,
    birth_country CHAR(2) NOT NULL,
    mother_birth_first_name VARCHAR(256) NOT NULL,
    mother_birth_last_name VARCHAR(256) NOT NULL,
    classes_id INT NOT NULL,
    permament_address VARCHAR(256) NOT NULL,
    temporary_address VARCHAR(256) NOT NULL,
    tax_number INT DEFAULT(NULL),
    ssn_number INT NOT NULL,
    bank_name VARCHAR(256) NOT NULL,
    iban_owner VARCHAR(256) NOT NULL,
    iban_number VARCHAR(34) NOT NULL,
    document_type VARCHAR(50) NOT NULL,
    document_number VARCHAR(50) NOT NULL,
    password_hash TEXT NOT NULL,
    password_salt TEXT NOT NULL
);

CREATE TABLE guardians_access (
    id SERIAL PRIMARY KEY,
    student_id INT NOT NULL,
    guardian_id INT NOT NULL,
    legal_representative BOOLEAN NOT NULL
);

CREATE TABLE guardians (
    id SERIAL PRIMARY KEY,
    email_address VARCHAR(256) NOT NULL UNIQUE,
    phone_number VARCHAR(30) NOT NULL UNIQUE,
    first_name VARCHAR(256) NOT NULL,
    last_name VARCHAR(256) NOT NULL,
    birth_first_name VARCHAR(256) NOT NULL,
    birth_last_name VARCHAR(256) NOT NULL,
    birth_date DATE NOT NULL,
    birth_city VARCHAR(100) NOT NULL,
    birth_country CHAR(2) NOT NULL,
    permament_address VARCHAR(256) NOT NULL,
    temporary_address VARCHAR(256) NOT NULL,
    password_hash TEXT NOT NULL,
    password_salt TEXT NOT NULL
);

CREATE TABLE teachers (
    id SERIAL PRIMARY KEY,
    email_address VARCHAR(256) NOT NULL UNIQUE,
    phone_number VARCHAR(256) NOT NULL UNIQUE,
    username VARCHAR(256) NOT NULL,
    birth_first_name VARCHAR(256) NOT NULL,
    birth_last_name VARCHAR(256) NOT NULL,
    birth_date DATE NOT NULL,
    birth_city VARCHAR(100) NOT NULL,
    birth_country CHAR(2) NOT NULL,
    permament_address VARCHAR(256) NOT NULL,
    temporary_address VARCHAR(256) NOT NULL
);

CREATE TABLE teacher_school (
    id SERIAL PRIMARY KEY,
    teacher_id INT NOT NULL,
    school_id INT NOT NULL
);

CREATE TABLE classes (
    id SERIAL PRIMARY KEY,
    school_id INT NOT NULL,
    name VARCHAR(30) NOT NULL,
    teacher_id INT NOT NULL,
    co_teacher_id INT DEFAULT(NULL)
);

CREATE TABLE student_citizenships (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    document_type VARCHAR(50) NOT NULL,
    country CHAR(2) NOT NULL
);

CREATE TABLE guardian_citizenships (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    document_type VARCHAR(50) NOT NULL,
    country CHAR(2) NOT NULL
);

CREATE TABLE teacher_citizenships (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    document_type VARCHAR(50) NOT NULL,
    country CHAR(2) NOT NULL
);

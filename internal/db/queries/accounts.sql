-- name: CreateStudentAccount :exec
INSERT INTO accounts (role, password_hash, student_id) VALUES ("student", $1, $2);

-- name: CreateGuardianAccount :exec
INSERT INTO accounts (role, password_hash, guardian_id) VALUES ("guardian", $1, $2);

-- name: CreateTeacherAccount :exec
INSERT INTO accounts (role, password_hash, teacher_id) VALUES ("teacher", $1, $2);

-- name: GetAccountByStudentID :one
SELECT * FROM accounts WHERE student_id = $1 AND role = "student";

-- name: GetAccountByGuardianID :one
SELECT * FROM accounts WHERE guardian_id = $1 AND role = "guardian";

-- name: GetAccountByTeacherID :one
SELECT * FROM accounts WHERE teacher_id = $1 AND role = "teacher";

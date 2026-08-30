-- name: CreateStudentAccount :exec
INSERT INTO accounts (role, password_hash, student_id) VALUES ("student", $1, $2);

-- name: CreateGuardianAccount :exec
INSERT INTO accounts (role, password_hash, guardian_id) VALUES ("guardian", $1, $2);

-- name: CreateTeacherAccount :exec
INSERT INTO accounts (role, password_hash, teacher_id) VALUES ("teacher", $1, $2);

-- name: GetAccountByStudentID :one
SELECT * FROM accounts WHERE student_id = $1 AND role = 'student';

-- name: GetAccountByGuardianID :one
SELECT * FROM accounts WHERE guardian_id = $1 AND role = 'guardian';

-- name: GetAccountByTeacherID :one
SELECT * FROM accounts WHERE teacher_id = $1 AND role = 'teacher';

-- name: ResetPasswordStudent :exec
UPDATE accounts SET password_hash = $1 WHERE student_id = $2 AND role = 'student';

-- name: ResetPasswordTeacher :exec
UPDATE accounts SET password_hash = $1 WHERE teacher_id = $2 AND role = 'teacher';

-- name: ResetPasswordGuardian :exec
UPDATE accounts SET password_hash = $1 WHERE guardian_id = $2 AND role = 'guardian';

-- name: ChangePasswordStudent :exec
UPDATE accounts SET password_hash = $1 WHERE student_id = $2 AND role = 'student';

-- name: ChangePasswordTeacher :exec
UPDATE  accounts SET password_hash = $1 WHERE teacher_id = $2 AND role = 'teacher';

-- name: ChangePasswordGuardian :exec
UPDATE  accounts SET password_hash = $1 WHERE  teacher_id = $2 AND role = 'guardian';

-- name: UpdateEmailStudent :exec
UPDATE accounts SET email_address = $1, email_verified = false WHERE student_id = $2 AND role = 'student';

-- name: UpdaEmailteGuardian :exec
UPDATE accounts SET email_address = $1, email_verfied = false WHERE guardian_id $2 AND role = 'guardian';

-- name: UpdateEmailTeacher :exec
UPDATE accounts SET email_address = $1, email_verfied = false WHERE teacher_id $2 AND role = 'teacher';

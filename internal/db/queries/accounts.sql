-- name: CreateAccount :exec
INSERT INTO accounts (id,role,student_id,teacher_id,guardian_id,email_address,email_verified,theme,lang,pfp_url,nickname) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);

-- name: GetAccountByStudentID :one
SELECT * FROM accounts WHERE student_id = $1 AND role = 'student';

-- name: GetAccountByGuardianID :one
SELECT * FROM accounts WHERE guardian_id = $1 AND role = 'guardian';

-- name: GetAccountByTeacherID :one
SELECT * FROM accounts WHERE teacher_id = $1 AND role = 'teacher';

-- name: GetAccountByUUID :one
SELECT * FROM accounts WHERE id = $1;

-- name: UpdateEmailStudent :exec
UPDATE accounts SET email_address = $1, email_verified = false WHERE student_id = $2 AND role = 'student';

-- name: UpdateEmailGuardian :exec
UPDATE accounts SET email_address = $1, email_verified = false WHERE guardian_id = $2 AND role = 'guardian';

-- name: UpdateEmailTeacher :exec
UPDATE accounts SET email_address = $1, email_verified = false WHERE teacher_id = $2 AND role = 'teacher';

-- name: UpdateNicknameStudent :exec
UPDATE accounts SET nickname = $1 WHERE student_id = $2 AND role = 'student';

-- name: UpdateNicknameGuardian :exec
UPDATE accounts SET nickname = $1 WHERE guardian_id = $2 AND role = 'guardian';

-- name: UpdateNicknameTeacher :exec
UPDATE accounts SET nickname = $1 WHERE teacher_id = $2 AND role = 'teacher';

-- name: UpdateThemeStudent :exec
UPDATE accounts SET theme = $1 WHERE student_id = $2 AND role = 'student';

-- name: UpdateThemeGuardian :exec
UPDATE accounts SET theme = $1 WHERE guardian_id = $2 AND role = 'guardian';

-- name: UpdateThemeTeacher :exec
UPDATE accounts SET theme = $1 WHERE teacher_id = $2 AND role = 'teacher';

-- name: UpdateLangStudent :exec
UPDATE accounts SET lang = $1 WHERE student_id = $2 AND role = 'student';

-- name: UpdateLangGuardian :exec
UPDATE accounts SET lang = $1 WHERE guardian_id = $2 AND role = 'guardian';

-- name: UpdateLangTeacher :exec
UPDATE accounts SET lang = $1 WHERE teacher_id = $2 AND role = 'teacher';

-- name: UpdatePreferencesStudent :exec
UPDATE accounts SET preferences = $1 WHERE student_id = $2 AND role = 'student';

-- name: UpdatePreferencesGuardian :exec
UPDATE accounts SET preferences = $1 WHERE guardian_id = $2 AND role = 'guardian';

-- name: UpdatePreferencesTeacher :exec
UPDATE accounts SET preferences = $1 WHERE teacher_id = $2 AND role = 'teacher';

-- name: VerifyEmailStudent :exec
UPDATE accounts SET email_verified = true WHERE student_id = $1 AND role = 'student';

-- name: VerifyEmailGuardian :exec
UPDATE accounts SET email_verified = true WHERE guardian_id = $1 AND role = 'guardian';

-- name: VerifyEmailTeacher :exec
UPDATE accounts SET email_verified = true WHERE teacher_id = $1 AND role = 'teacher';

package meapi

import (
	db_queries "onlystudents/internal/db/store"
	"onlystudents/internal/helpers"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type changePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

func ChangePassword(c fiber.Ctx, pool *pgxpool.Pool, session_store *helpers.SessionStore, cache_store *helpers.CacheStore) error {
	session_token := c.Cookies("session_token", "")
	if session_token == "" {
		return c.SendStatus(401)
	}

	session_data, err := session_store.Get(c, session_token)
	if err != nil {
		return c.SendStatus(401)
	}

	var req changePasswordRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.SendStatus(400)
	}

	if req.CurrentPassword == "" || req.NewPassword == "" {
		return c.SendStatus(400)
	}

	queries := db_queries.New(pool)

	account, err := cache_store.CacheOrGetAccount(c.Context(), *queries, session_data.Role, session_data.AccountID, helpers.GetInt32EnvFallback("ACCOUNT_CACHE_TTL", 5*60, 604800))
	if err != nil {
		return c.SendStatus(401)
	}

	if !helpers.Argon2Verify(req.CurrentPassword, account.PasswordHash) {
		return c.SendStatus(401)
	}

	new_password_hash, err := helpers.Argon2HashPassword(req.NewPassword)
	if err != nil {
		return c.SendStatus(500)
	}

	account_id := pgtype.Int4{Int32: session_data.AccountID, Valid: true}

	switch session_data.Role {
	case "student":
		err = queries.ChangePasswordStudent(c.Context(), db_queries.ChangePasswordStudentParams{PasswordHash: new_password_hash, StudentID: account_id})
	case "guardian":
		err = queries.ChangePasswordGuardian(c.Context(), db_queries.ChangePasswordGuardianParams{PasswordHash: new_password_hash, GuardianID: account_id})
	case "teacher":
		err = queries.ChangePasswordTeacher(c.Context(), db_queries.ChangePasswordTeacherParams{PasswordHash: new_password_hash, TeacherID: account_id})
	default:
		return c.SendStatus(400)
	}

	if err != nil {
		return c.SendStatus(500)
	}
	// invalidate cached account object which has the old password hash
	cache_store.InvalidateCachedAccount(c.Context(), session_data.Role, session_data.AccountID)
	return c.SendStatus(200)
}

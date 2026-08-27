package v1

import (
	"context"
	db_queries "onlystudents/internal/db/store"
	"onlystudents/internal/helpers"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type loginRequest struct {
	User     int32  `json:"user"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func Login(c fiber.Ctx, pool *pgxpool.Pool, session_store *helpers.SessionStore) error {
	var req loginRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.SendStatus(400)
	}

	if req.User <= 0 || req.Password == "" || req.Role == "" {
		return c.SendStatus(400)
	}

	user := pgtype.Int4{Int32: req.User, Valid: true}

	queries := db_queries.New(pool)

	var account db_queries.Account
	var err error

	switch req.Role {
	case "student":
		account, err = queries.GetAccountByStudentID(context.Background(), user)
	case "guardian":
		account, err = queries.GetAccountByGuardianID(context.Background(), user)
	case "teacher":
		account, err = queries.GetAccountByTeacherID(context.Background(), user)
	default:
		return c.SendStatus(400)
	}

	if err != nil {
		return c.SendStatus(401)
	}

	if !helpers.Argon2Verify(req.Password, account.PasswordHash) {
		return c.SendStatus(401)
	}

	session_token, err := session_store.Create(context.Background(), req.User, req.Role)
	if err != nil {
		return c.SendStatus(500)
	}

	duration := time.Duration(helpers.GetUintEnvFallback("SESSION_TTL", 3600)) * time.Second

	c.Cookie(&fiber.Cookie{
		Name:     "session_token",
		Value:    session_token,
		Expires:  time.Now().Add(duration),
		HTTPOnly: true,
		Secure:   helpers.GetEnvFallback("APP_ENV", "development") == "production",
		SameSite: "Lax",
	})

	return c.JSON(helpers.SessionData{
		Role:      req.Role,
		AccountID: req.User,
	})
}

func Logout(c fiber.Ctx, pool *pgxpool.Pool, session_store *helpers.SessionStore) error {
	token := c.Cookies("session_token", "")

	if token == "" {
		return c.SendStatus(401)
	}

	err := session_store.Delete(context.Background(), token)

	if err != nil {
		return c.SendStatus(401)
	}

	c.ClearCookie("session_token")

	return c.Redirect().To("/")
}

func Me(c fiber.Ctx, pool *pgxpool.Pool, session_store *helpers.SessionStore) error {
	session_token := c.Cookies("session_token", "")
	if session_token == "" {
		return c.SendStatus(401)
	}

	session_data, err := session_store.Get(c, session_token)
	if err != nil {
		return c.SendStatus(401)
	}

	return c.JSON(session_data)
}

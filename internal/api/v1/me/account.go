package meapi

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	db_queries "onlystudents/internal/db/store"
	"onlystudents/internal/helpers"
	"onlystudents/internal/mail"
	"time"

	"github.com/goccy/go-json"
	"github.com/redis/go-redis/v9"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

const pgUniqueViolation = "23505"

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation
}

func generateCode(length int) (string, error) {
	const charset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	code := make([]byte, length)
	for i := range code {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		code[i] = charset[n.Int64()]
	}
	return string(code), nil
}

func ChangePassword(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	_, ok := c.Locals("session").(helpers.SessionData)

	if !ok {
		return c.SendStatus(401)
	}

	// IMPORTANT TODO: OPAQUE password change
	return c.SendStatus(501)
}

type changeNicknameRequest struct {
	Nickname string `json:"nickname"`
}

type changeThemeRequest struct {
	Theme string `json:"theme"`
}

type changeLangRequest struct {
	Lang string `json:"lang"`
}

type changeEmailRequest struct {
	Password string `json:"password"`
	NewEmail string `json:"new_email"`
}

func ChangeTheme(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	session_data, ok := c.Locals("session").(helpers.SessionData)

	if !ok {
		return c.SendStatus(401)
	}

	var req changeThemeRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.SendStatus(400)
	}

	queries := db_queries.New(pool)

	pgAccountID := pgtype.Int4{Int32: session_data.AccountID, Valid: true}
	var err error
	switch session_data.Role {
	case "student":
		err = queries.UpdateThemeStudent(c.Context(), db_queries.UpdateThemeStudentParams{Theme: req.Theme, StudentID: pgAccountID})
	case "teacher":
		err = queries.UpdateThemeTeacher(c.Context(), db_queries.UpdateThemeTeacherParams{Theme: req.Theme, TeacherID: pgAccountID})
	case "guardian":
		err = queries.UpdateThemeGuardian(c.Context(), db_queries.UpdateThemeGuardianParams{Theme: req.Theme, GuardianID: pgAccountID})
	default:
		return c.SendStatus(400)
	}

	if err != nil {
		if isUniqueViolation(err) {
			return c.SendStatus(400)
		}
		slog.Error("change nickname error", "err", err)
		return c.SendStatus(500)
	}

	helpers.InvalidateCachedAccount(c.Context(), rdb, session_data.Role, session_data.AccountID)
	return c.SendStatus(200)
}

func ChangeLang(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	session_data, ok := c.Locals("session").(helpers.SessionData)

	if !ok {
		return c.SendStatus(401)
	}

	var req changeLangRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.SendStatus(400)
	}

	queries := db_queries.New(pool)

	pgAccountID := pgtype.Int4{Int32: session_data.AccountID, Valid: true}
	var err error
	switch session_data.Role {
	case "student":
		err = queries.UpdateLangStudent(c.Context(), db_queries.UpdateLangStudentParams{Lang: req.Lang, StudentID: pgAccountID})
	case "teacher":
		err = queries.UpdateLangTeacher(c.Context(), db_queries.UpdateLangTeacherParams{Lang: req.Lang, TeacherID: pgAccountID})
	case "guardian":
		err = queries.UpdateLangGuardian(c.Context(), db_queries.UpdateLangGuardianParams{Lang: req.Lang, GuardianID: pgAccountID})
	default:
		return c.SendStatus(400)
	}

	if err != nil {
		if isUniqueViolation(err) {
			return c.SendStatus(400)
		}
		slog.Error("change nickname error", "err", err)
		return c.SendStatus(500)
	}

	helpers.InvalidateCachedAccount(c.Context(), rdb, session_data.Role, session_data.AccountID)
	return c.SendStatus(200)
}

func ChangeNickname(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	session_data, ok := c.Locals("session").(helpers.SessionData)

	if !ok {
		return c.SendStatus(401)
	}

	var req changeNicknameRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.SendStatus(400)
	}

	queries := db_queries.New(pool)

	pgAccountID := pgtype.Int4{Int32: session_data.AccountID, Valid: true}
	var err error
	switch session_data.Role {
	case "student":
		err = queries.UpdateNicknameStudent(c.Context(), db_queries.UpdateNicknameStudentParams{Nickname: req.Nickname, StudentID: pgAccountID})
	case "teacher":
		err = queries.UpdateNicknameTeacher(c.Context(), db_queries.UpdateNicknameTeacherParams{Nickname: req.Nickname, TeacherID: pgAccountID})
	case "guardian":
		err = queries.UpdateNicknameGuardian(c.Context(), db_queries.UpdateNicknameGuardianParams{Nickname: req.Nickname, GuardianID: pgAccountID})
	default:
		return c.SendStatus(400)
	}

	if err != nil {
		if isUniqueViolation(err) {
			return c.SendStatus(400)
		}
		slog.Error("change nickname error", "err", err)
		return c.SendStatus(500)
	}

	helpers.InvalidateCachedAccount(c.Context(), rdb, session_data.Role, session_data.AccountID)
	return c.SendStatus(200)
}

func ChangeEmail(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	session_data, ok := c.Locals("session").(helpers.SessionData)

	if !ok {
		return c.SendStatus(401)
	}

	var req changeEmailRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.SendStatus(400)
	}

	if req.Password == "" || req.NewEmail == "" {
		return c.SendStatus(400)
	}

	queries := db_queries.New(pool)
	_, err := helpers.CacheOrGetAccount(c.Context(), rdb, *queries, session_data.Role, session_data.AccountID, helpers.GetInt32EnvFallback("ACCOUNT_CACHE_TTL", 5*60, 604800))
	if err != nil {
		return c.SendStatus(401)
	}

	// IMPORTANT TODO: Add OPAQUE auth
	return c.SendStatus(501)

	// pgEmail := pgtype.Text{String: req.NewEmail, Valid: true}
	// pgAccountID := pgtype.Int4{Int32: session_data.AccountID, Valid: true}
	// switch session_data.Role {
	// case "student":
	// 	err = queries.UpdateEmailStudent(c.Context(), db_queries.UpdateEmailStudentParams{EmailAddress: pgEmail, StudentID: pgAccountID})
	// case "teacher":
	// 	err = queries.UpdateEmailTeacher(c.Context(), db_queries.UpdateEmailTeacherParams{EmailAddress: pgEmail, TeacherID: pgAccountID})
	// case "guardian":
	// 	err = queries.UpdateEmailGuardian(c.Context(), db_queries.UpdateEmailGuardianParams{EmailAddress: pgEmail, GuardianID: pgAccountID})
	// default:
	// 	return c.SendStatus(400)
	// }

	// if err != nil {
	// 	if isUniqueViolation(err) {
	// 		return c.SendStatus(400)
	// 	}
	// 	slog.Error("change email error", "err", err)
	// 	return c.SendStatus(500)
	// }

	// helpers.InvalidateCachedAccount(c.Context(), rdb, session_data.Role, session_data.AccountID)
	// return c.SendStatus(200)
}

type VerifyEmailData struct {
	VerifyCode       string
	Name             string
	ExpiresInMinutes int
}

func VerifyEmailRequest(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	session_data, ok := c.Locals("session").(helpers.SessionData)

	if !ok {
		return c.SendStatus(401)
	}

	queries := db_queries.New(pool)
	account, err := helpers.CacheOrGetAccount(c.Context(), rdb, *queries, session_data.Role, session_data.AccountID, helpers.GetInt32EnvFallback("ACCOUNT_CACHE_TTL", 5*60, 604800))

	if err != nil {
		return c.SendStatus(401)
	}

	if !account.EmailAddress.Valid || account.EmailAddress.String == "" {
		return c.SendStatus(500)
	}

	if account.EmailVerified {
		return c.SendStatus(200)
	}

	code, e := generateCode(helpers.GetIntEnvFallback("EMAIL_VERIFY_CODE_LEN", 10, 64))
	if e != nil {
		return c.SendStatus(500)
	}

	cacheKey := fmt.Sprintf("pending_email_verify_%s", code)
	pending := helpers.SessionData{AccountID: session_data.AccountID, Role: session_data.Role}
	dataJSON, e := json.Marshal(pending)
	if e != nil {
		return c.SendStatus(500)
	}

	expiresInMinutes := helpers.GetInt64EnvFallback("EMAIL_VERIFY_CODE_TTL", 15, 1440)
	ttl := time.Duration(expiresInMinutes) * time.Minute
	if e := rdb.Set(c.Context(), cacheKey, string(dataJSON), ttl).Err(); e != nil {
		return c.SendStatus(500)
	}

	email := account.EmailAddress.String

	// #nosec G118
	go func() {
		// intentionally detached to avoid timing-based user enumeration
		bg, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		mailer, e := mail.NewFromEnv()
		if e != nil {
			slog.Error("email verify mailer init failed", "err", e)
			return
		}
		if e := mailer.SendTemplateContext(bg, email, "Verify your email", "email_verify.html", VerifyEmailData{VerifyCode: code, Name: account.Role, ExpiresInMinutes: int(expiresInMinutes)}); e != nil {
			slog.Error("email verify email failed", "account", email, "err", e)
		}
	}()

	return c.SendStatus(200)
}

type VerifyEmailConfirmRequest struct {
	Code string `json:"code"`
}

func VerifyEmailConfirm(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	session_data, ok := c.Locals("session").(helpers.SessionData)

	if !ok {
		return c.SendStatus(401)
	}

	var req VerifyEmailConfirmRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.SendStatus(400)
	}

	if req.Code == "" {
		return c.SendStatus(400)
	}

	cacheKey := fmt.Sprintf("pending_email_verify_%s", req.Code)
	dataJSON, err := rdb.GetDel(c.Context(), cacheKey).Result()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "INVALID_VERIFICATION_CODE",
		})
	}

	var pending helpers.SessionData
	if err := json.Unmarshal([]byte(dataJSON), &pending); err != nil {
		slog.Error("json unmarshal error", "err", err)
		return c.SendStatus(500)
	}

	if pending.AccountID != session_data.AccountID || pending.Role != session_data.Role {
		return c.SendStatus(401)
	}

	queries := db_queries.New(pool)
	pgAccountID := pgtype.Int4{Int32: session_data.AccountID, Valid: true}
	switch session_data.Role {
	case "student":
		err = queries.VerifyEmailStudent(c.Context(), pgAccountID)
	case "teacher":
		err = queries.VerifyEmailTeacher(c.Context(), pgAccountID)
	case "guardian":
		err = queries.VerifyEmailGuardian(c.Context(), pgAccountID)
	default:
		return c.SendStatus(400)
	}

	if err != nil {
		slog.Error("verify email error", "err", err)
		return c.SendStatus(500)
	}

	helpers.InvalidateCachedAccount(c.Context(), rdb, session_data.Role, session_data.AccountID)
	return c.SendStatus(200)
}

package meapi

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"html"
	"log/slog"
	"math/big"
	db_queries "onlystudents/internal/db/store"
	"onlystudents/internal/helpers"
	"onlystudents/internal/mail"
	"slices"
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
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	// IMPORTANT TODO: OPAQUE password change
	return c.SendStatus(fiber.StatusNotImplemented)
}

type Preferences struct {
	Nickname         string `json:"nickname"`
	Lang             string `json:"lang"`
	Theme            string `json:"theme"`
	PfpURL           string `json:"pfp_url"`
	TimeTableDisplay int8   `json:"timetable_display"`
	TimeTableNext    bool   `json:"timetable_next"`
}

type UpdatePreferencesRequest struct {
	NewPreferences Preferences `json:"new_preferences"`
}

type changeEmailRequest struct {
	NewEmail string `json:"new_email"`
}

func UpdatePreferences(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	session_data, ok := c.Locals("session").(helpers.SessionData)

	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	var req UpdatePreferencesRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if !slices.Contains(helpers.Themes, req.NewPreferences.Theme) {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if !slices.Contains(helpers.Languages, req.NewPreferences.Lang) {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if len(req.NewPreferences.Nickname) < 8 || len(req.NewPreferences.Nickname) > 64 || html.EscapeString(req.NewPreferences.Nickname) != req.NewPreferences.Nickname {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	preferences_jsonb, marshal_err := json.Marshal(Preferences{
		Nickname:         req.NewPreferences.Nickname,
		Lang:             req.NewPreferences.Lang,
		Theme:            req.NewPreferences.Theme,
		PfpURL:           req.NewPreferences.PfpURL,
		TimeTableDisplay: req.NewPreferences.TimeTableDisplay,
		TimeTableNext:    req.NewPreferences.TimeTableNext,
	})

	if marshal_err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	queries := db_queries.New(pool)

	pgAccountID := pgtype.Int4{Int32: session_data.AccountID, Valid: true}
	var err error
	switch session_data.Role {
	case "student":
		err = queries.UpdatePreferencesStudent(c.Context(), db_queries.UpdatePreferencesStudentParams{Preferences: preferences_jsonb, StudentID: pgAccountID})
	case "teacher":
		err = queries.UpdatePreferencesTeacher(c.Context(), db_queries.UpdatePreferencesTeacherParams{Preferences: preferences_jsonb, TeacherID: pgAccountID})
	case "guardian":
		err = queries.UpdatePreferencesGuardian(c.Context(), db_queries.UpdatePreferencesGuardianParams{Preferences: preferences_jsonb, GuardianID: pgAccountID})
	default:
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if err != nil {
		slog.Error("update preferences error", "err", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	helpers.InvalidateCachedAccount(c.Context(), rdb, session_data.Role, session_data.AccountID)
	return c.SendStatus(fiber.StatusOK)
}

func ChangeEmail(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	session_data, ok := c.Locals("session").(helpers.SessionData)

	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	var req changeEmailRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	queries := db_queries.New(pool)
	_, err := helpers.CacheOrGetAccount(c.Context(), rdb, *queries, session_data.Role, session_data.AccountID, helpers.GetInt32EnvFallback("ACCOUNT_CACHE_TTL", 5*60, 604800))
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	pgEmail := pgtype.Text{String: req.NewEmail, Valid: true}
	pgAccountID := pgtype.Int4{Int32: session_data.AccountID, Valid: true}
	switch session_data.Role {
	case "student":
		err = queries.UpdateEmailStudent(c.Context(), db_queries.UpdateEmailStudentParams{EmailAddress: pgEmail, StudentID: pgAccountID})
	case "teacher":
		err = queries.UpdateEmailTeacher(c.Context(), db_queries.UpdateEmailTeacherParams{EmailAddress: pgEmail, TeacherID: pgAccountID})
	case "guardian":
		err = queries.UpdateEmailGuardian(c.Context(), db_queries.UpdateEmailGuardianParams{EmailAddress: pgEmail, GuardianID: pgAccountID})
	default:
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if err != nil {
		if isUniqueViolation(err) {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		slog.Error("change email error", "err", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	helpers.InvalidateCachedAccount(c.Context(), rdb, session_data.Role, session_data.AccountID)
	return c.SendStatus(fiber.StatusOK)
}

type VerifyEmailData struct {
	VerifyCode       string
	Name             string
	ExpiresInMinutes int
}

func VerifyEmailRequest(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	session_data, ok := c.Locals("session").(helpers.SessionData)

	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	queries := db_queries.New(pool)
	account, err := helpers.CacheOrGetAccount(c.Context(), rdb, *queries, session_data.Role, session_data.AccountID, helpers.GetInt32EnvFallback("ACCOUNT_CACHE_TTL", 5*60, 604800))

	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	if !account.EmailAddress.Valid || account.EmailAddress.String == "" {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	if account.EmailVerified {
		return c.SendStatus(fiber.StatusOK)
	}

	code, e := generateCode(helpers.GetIntEnvFallback("EMAIL_VERIFY_CODE_LEN", 10, 64))
	if e != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	cacheKey := fmt.Sprintf("pending_email_verify_%s", code)
	pending := helpers.SessionData{AccountID: session_data.AccountID, Role: session_data.Role}
	dataJSON, e := json.Marshal(pending)
	if e != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	expiresInMinutes := helpers.GetInt64EnvFallback("EMAIL_VERIFY_CODE_TTL", 15, 1440)
	ttl := time.Duration(expiresInMinutes) * time.Minute
	if e := rdb.Set(c.Context(), cacheKey, string(dataJSON), ttl).Err(); e != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
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

	return c.SendStatus(fiber.StatusOK)
}

type VerifyEmailConfirmRequest struct {
	Code string `json:"code"`
}

func VerifyEmailConfirm(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	session_data, ok := c.Locals("session").(helpers.SessionData)

	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	var req VerifyEmailConfirmRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if req.Code == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	cacheKey := fmt.Sprintf("pending_email_verify_%s", req.Code)
	dataJSON, err := rdb.GetDel(c.Context(), cacheKey).Result()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "INVALID_VERIFICATION_CODE",
		})
	}

	var pending helpers.SessionData
	if err := json.Unmarshal([]byte(dataJSON), &pending); err != nil {
		slog.Error("json unmarshal error", "err", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	if pending.AccountID != session_data.AccountID || pending.Role != session_data.Role {
		return c.SendStatus(fiber.StatusUnauthorized)
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
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if err != nil {
		slog.Error("verify email error", "err", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	helpers.InvalidateCachedAccount(c.Context(), rdb, session_data.Role, session_data.AccountID)
	return c.SendStatus(fiber.StatusOK)
}

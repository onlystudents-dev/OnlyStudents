package meapi

import (
	"errors"
	"log/slog"
	db_queries "onlystudents/internal/db/store"
	"onlystudents/internal/helpers"
	opaquepkg "onlystudents/internal/opaque"
	"time"

	"github.com/bytemare/opaque"
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type ChangePasswordInitRequest struct {
	StartLoginRequest opaquepkg.Message `json:"start_login_request"`
}

type RekeyInitRequest struct {
	StartEnrollRequest opaquepkg.Message `json:"start_enroll_request"`
}

type RekeyFinishRequest struct {
	CryptoGrantHandle   string            `json:"crypto_grant_handle"`
	WrappedMKPassword   string            `json:"wrapped_mk_password"`
	FinishEnrollRequest opaquepkg.Message `json:"finish_enroll_request"`
}

func ChangePasswordInit(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client, opaque_server *opaque.Server) error {
	session_data, ok := c.Locals("session").(helpers.SessionData)

	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "INVALID_SESSION"})
	}

	var req ChangePasswordInitRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	var record *opaque.ClientRecord

	queries := db_queries.New(pool)
	account, err := helpers.CacheOrGetAccount(c.Context(), rdb, *queries, session_data.Role, session_data.AccountID, helpers.GetInt32EnvFallback("ACCOUNT_CACHE_TTL", 5*60, 604800))

	switch {
	case err == nil:
		opaque_record, rerr := queries.GetOpaqueRecord(c.Context(), account.ID)
		if rerr != nil {
			record = opaquepkg.FakeRecord
		} else {
			reg, derr := opaque_server.Deserialize.RegistrationRecord(opaque_record.RegistrationRecord)

			if derr != nil {
				slog.Error("corrupt opaque record", "err", derr)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "SERVER_ERROR"})
			}
			record = &opaque.ClientRecord{
				CredentialIdentifier: []byte(account.ID.String()),
				RegistrationRecord:   reg,
			}
		}
	case errors.Is(err, pgx.ErrNoRows):
		record = opaquepkg.FakeRecord
	default:
		slog.Error("login init account lookup", "err", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "SERVER_ERROR"})
	}

	ke2, handle, err := opaquepkg.ChangePasswordInit(c.Context(), opaque_server, rdb, record, req.StartLoginRequest)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	return c.JSON(loginInitResponse{LoginResponse: ke2, LoginHandle: handle})
}

func ChangePasswordFinish(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client, opaque_server *opaque.Server) error {
	var req loginFinishRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	if req.LoginHandle == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	sessionToken, crypto_grant_handle, err := opaquepkg.ChangePasswordFinish(c.Context(), opaque_server, rdb, pool, req.LoginHandle, req.FinishLoginRequest)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "WRONG_CREDENTIALS"})
	}

	duration := time.Duration(helpers.GetInt64EnvFallback("SESSION_TTL", 3600, 2592000)) * time.Second

	c.Cookie(&fiber.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  time.Now().Add(duration),
		HTTPOnly: true,
		Secure:   helpers.GetEnvFallback("APP_ENV", "development") == "production",
		SameSite: "Strict",
	})

	return c.JSON(fiber.Map{
		"crypto_grant_handle": crypto_grant_handle,
		"session_token":       sessionToken,
	})
}

func RekeyInit(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client, opaque_server *opaque.Server) error {
	session_data, ok := c.Locals("session").(helpers.SessionData)

	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "INVALID_SESSION"})
	}

	var req RekeyInitRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	ke2, err := opaquepkg.RekeyInit(c.Context(), opaque_server, session_data.AccountUUID, req.StartEnrollRequest)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	return c.JSON(enrollInitResponse{EnrollResponse: ke2})
}

func RekeyFinish(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client, opaque_server *opaque.Server) error {
	session_data, ok := c.Locals("session").(helpers.SessionData)

	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "INVALID_SESSION"})
	}

	var req RekeyFinishRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	DecodedWrappedMKPassword, DecodedWrappedMKPasswordErr := b64decode(req.WrappedMKPassword)

	if DecodedWrappedMKPasswordErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	err := opaquepkg.RekeyFinish(c.Context(), opaque_server, rdb, pool, session_data.AccountUUID, req.CryptoGrantHandle, DecodedWrappedMKPassword, req.FinishEnrollRequest)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "WRONG_CREDENTIALS"})
	}

	return c.SendStatus(fiber.StatusOK)
}

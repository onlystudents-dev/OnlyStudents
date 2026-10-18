package meapi

import (
	"encoding/base64"
	"log/slog"
	db_queries "onlystudents/internal/db/store"
	"onlystudents/internal/helpers"
	opaquepkg "onlystudents/internal/opaque"

	"github.com/goccy/go-json"
	"github.com/google/uuid"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type enrollE2EERequest struct {
	WrappedMKPassword string `json:"wrapped_mk_password"`
	WrappedMKRecovery string `json:"wrapped_mk_recovery"`
	RecoveryKeyHash   string `json:"recovery_key_hash"`
	RecoveryKeySalt   string `json:"recovery_key_salt"`
	RecoveryKeyParams string `json:"recovery_key_params"`
}

type E2EERewrapRequest struct {
	CryptoGrantHandle string `json:"crypto_grant_handle"`
	WrappedMKPassword string `json:"wrapped_mk_password"`
}

func b64decode(base64_string string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(base64_string)
}

func b64encode(base64_bytes []byte) string {
	return base64.StdEncoding.EncodeToString(base64_bytes)
}

func GetE2EE(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	session_data, ok := c.Locals("session").(helpers.SessionData)

	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "INVALID_SESSION"})
	}

	queries := db_queries.New(pool)

	account_uuid, aerr := uuid.Parse(session_data.AccountUUID)

	if aerr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "SERVER_ERROR"})
	}

	e2ee_data, err := queries.GetAccountCrypto(c.Context(), pgtype.UUID{Bytes: account_uuid, Valid: true})

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "NO_CRYPTO"})
	}

	base64_encoded_wrapped_mk := b64encode(e2ee_data.WrappedMkPassword)

	if base64_encoded_wrapped_mk == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "SERVER_ERROR"})
	}

	return c.JSON(fiber.Map{
		"wrapped_mk_password": base64_encoded_wrapped_mk,
	})
}

func EnrollE2EE(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	session_data, ok := c.Locals("session").(helpers.SessionData)

	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "INVALID_SESSION"})
	}

	var req enrollE2EERequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	if req.WrappedMKPassword == "" || req.WrappedMKRecovery == "" || req.RecoveryKeyHash == "" || req.RecoveryKeySalt == "" || req.RecoveryKeyParams == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	account_uuid, aerr := uuid.Parse(session_data.AccountUUID)

	if aerr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	DecodedWrappedMKPassword, DecodedWrappedMKPasswordErr := b64decode(req.WrappedMKPassword)

	if DecodedWrappedMKPasswordErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	DecodedWrappedMKRecovery, DecodedWrappedMKRecoveryErr := b64decode(req.WrappedMKRecovery)

	if DecodedWrappedMKRecoveryErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	DecodedRecoveryKeyHash, DecodedRecoveryKeyHashErr := b64decode(req.RecoveryKeyHash)

	if DecodedRecoveryKeyHashErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	DecodedRecoveryKeySalt, DecodedRecoveryKeySaltErr := b64decode(req.RecoveryKeySalt)

	if DecodedRecoveryKeySaltErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	queries := db_queries.New(pool)

	err := queries.CreateAccountCrypto(c.Context(), db_queries.CreateAccountCryptoParams{
		AccountUuid:       pgtype.UUID{Bytes: account_uuid, Valid: true},
		WrappedMkPassword: DecodedWrappedMKPassword,
		WrappedMkRecovery: DecodedWrappedMKRecovery,
		RecoveryKeyHash:   DecodedRecoveryKeyHash,
		RecoveryKeySalt:   DecodedRecoveryKeySalt,
		RecoveryKeyParams: req.RecoveryKeyParams,
	})
	if err != nil {
		if helpers.IsUniqueViolation(err) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "CONFLICT"})
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "SERVER_ERROR"})
		}
	}

	return c.SendStatus(fiber.StatusOK)
}

func E2EERewrap(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	session_data, ok := c.Locals("session").(helpers.SessionData)

	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "INVALID_SESSION"})
	}

	var req E2EERewrapRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	if req.WrappedMKPassword == "" || req.CryptoGrantHandle == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	raw, err := rdb.Get(c.Context(), "crypto_grant:"+req.CryptoGrantHandle).Bytes()

	if err != nil {
		slog.Error("crypto grant rdb get", "err", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "SERVER_ERROR"})
	}

	var crypto_grant opaquepkg.EnrollState

	unmarshal_err := json.Unmarshal(raw, &crypto_grant)

	if unmarshal_err != nil {
		slog.Error("json unmarshall error", "err", unmarshal_err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "SERVER_ERROR"})
	}

	if crypto_grant.AccountUUID != session_data.AccountUUID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "FORBIDDEN"})
	}

	parsed_account_uuid, err := uuid.Parse(crypto_grant.AccountUUID)

	if err != nil {
		slog.Error("account uuid parse", "err", err)
		return c.SendStatus(500)
	}

	DecodedWrappedMKPassword, DecodedWrappedMKPasswordErr := b64decode(req.WrappedMKPassword)

	if DecodedWrappedMKPasswordErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	queries := db_queries.New(pool)

	rows, err := queries.UpdateWrappedMKPassword(c.Context(), db_queries.UpdateWrappedMKPasswordParams{
		AccountUuid:       pgtype.UUID{Bytes: parsed_account_uuid, Valid: true},
		WrappedMkPassword: DecodedWrappedMKPassword,
	})

	if err != nil {
		slog.Error("update wrapped mk password", "err", err)
		if helpers.IsUniqueViolation(err) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "CONFLICT"})
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "SERVER_ERROR"})
		}
	}

	if rows != 1 {
		slog.Error("unexpected updated row count", "rows", rows)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	if err := rdb.Del(c.Context(), "crypto_grant:"+req.CryptoGrantHandle).Err(); err != nil {
		slog.Error("crypto grant rdb del", "err", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "SERVER_ERROR"})
	}

	return c.SendStatus(fiber.StatusOK)
}

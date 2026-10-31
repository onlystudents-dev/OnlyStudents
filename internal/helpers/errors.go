package helpers

import (
	"context"
	"errors"
	"log/slog"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type ErrorType struct {
	error_code  string
	status_code int
}

var db_errors = map[string]ErrorType{
	"23505": {error_code: "CONFLICT", status_code: 409},            // unique_violation
	"23503": {error_code: "BAD_REQUEST", status_code: 400},         // foreign_key_violation
	"23502": {error_code: "BAD_REQUEST", status_code: 400},         // not_null_violation
	"23514": {error_code: "BAD_REQUEST", status_code: 400},         // check_violation
	"23P01": {error_code: "CONFLICT", status_code: 409},            // exclusion_violation (same as unique)
	"22012": {error_code: "BAD_REQUEST", status_code: 400},         // division_by_zero
	"22P02": {error_code: "BAD_REQUEST", status_code: 400},         // invalid_text_representation
	"22001": {error_code: "BAD_REQUEST", status_code: 400},         // string_data_right_truncation
	"22003": {error_code: "BAD_REQUEST", status_code: 400},         // numeric_value_out_of_range
	"42501": {error_code: "FORBIDDEN", status_code: 403},           // insufficient_privilege
	"40001": {error_code: "SERVICE_UNAVAILABLE", status_code: 503}, // serialization_failure
	"40P01": {error_code: "SERVICE_UNAVAILABLE", status_code: 503}, // deadlock_detected
	"55P03": {error_code: "SERVICE_UNAVAILABLE", status_code: 503}, // lock_not_available
	"53300": {error_code: "SERVICE_UNAVAILABLE", status_code: 503}, // too many connections
	"53100": {error_code: "SERVICE_UNAVAILABLE", status_code: 503}, // disk full
	"53200": {error_code: "SERVICE_UNAVAILABLE", status_code: 503}, // out of memory
	"08006": {error_code: "SERVICE_UNAVAILABLE", status_code: 503}, // connection failure
	"08000": {error_code: "SERVICE_UNAVAILABLE", status_code: 503}, // connection failure
	"57P01": {error_code: "SERVICE_UNAVAILABLE", status_code: 503}, // admin shutdown
	"42P01": {error_code: "SERVER_ERROR", status_code: 500},        // undefined_table (a bug, not client fault)
	"42703": {error_code: "SERVER_ERROR", status_code: 500},        // undefined_column
}

func MapDBError(err error) (string, int) {
	slog.Debug("DB Error", "error", err.Error())

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return "NOT_FOUND", 404
	case errors.Is(err, pgx.ErrTooManyRows):
		return "SERVER_ERROR", 500
	case errors.Is(err, context.DeadlineExceeded):
		return "GATEWAY_TIMEOUT", 504
	case errors.Is(err, context.Canceled):
		return "REQUEST_CANCELLED", 499
	}

	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return "SERVER_ERROR", 500
	}

	error_type, ok := db_errors[pgErr.Code]

	if !ok {
		return "SERVER_ERROR", 500
	}

	return error_type.error_code, error_type.status_code
}

func ErrorByStatusCode(c fiber.Ctx, status_code int) error {
	switch status_code {
	case 400:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	case 401:
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "UNAUTHORIZED"})
	case 403:
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "FORBIDDEN"})
	case 500:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "SERVER_ERROR"})
	case 501:
		return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{"error": "NOT_IMPLEMENTED"})
	default:
		return c.SendStatus(status_code)
	}
}

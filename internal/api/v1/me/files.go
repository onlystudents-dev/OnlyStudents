package meapi

import (
	"log/slog"
	db_queries "onlystudents/internal/db/store"
	"onlystudents/internal/helpers"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

const (
	KiB = 1 << 10
)

type GetFileRequest struct {
	file_uuid string
}

func FileUpload(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	sessionData, ok := c.Locals("session").(helpers.SessionData)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "INVALID_SESSION"})
	}

	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "MISSING_FILE"})
	}

	cat, err := helpers.ValidateFileAndCategorize(file)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "FILE_TYPE_NOT_ALLOWED"})
	}

	fileUuid, err := uuid.NewRandom()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "SERVER_ERROR"})
	}

	if float64(file.Size)/float64(KiB) > float64(helpers.GetInt64EnvFallback("MAX_FILE_SIZE", 10, 99999)) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "FILE_TOO_BIG"})
	}

	if helpers.GetEnvFallback("FILE_UPLOAD_BACKEND", "directory") == "directory" {
		err := c.SaveFile(file, filepath.Join(helpers.GetEnvFallback("FILE_UPLOADS_DIR", "./uploads"), fileUuid.String()))
		if err != nil {
			slog.Error(err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "SERVER_ERROR"})
		}
	} else {
		// TODO: add S3-compatible file storage support
		return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{"error": "NOT_IMPLEMENTED"})
	}

	accountUuid, err := uuid.Parse(sessionData.AccountUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "SERVER_ERROR"})
	}

	queries := db_queries.New(pool)
	rows_changed, err := queries.UploadFile(c.Context(), db_queries.UploadFileParams{
		AccountUuid:  pgtype.UUID{Bytes: accountUuid, Valid: true},
		Type:         string(cat),
		Size:         file.Size,
		OriginalName: pgtype.Text{String: file.Filename, Valid: true},
		Uuid:         pgtype.UUID{Bytes: fileUuid, Valid: true},
	})

	if err != nil {
		slog.Error(err.Error())
		error_code, status_code := helpers.MapDBError(err)
		return c.Status(status_code).JSON(fiber.Map{"error": error_code})
	}

	if rows_changed != 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	return c.JSON(fiber.Map{"file_uuid": fileUuid.String()})
}

func FileDelete(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	sessionData, ok := c.Locals("session").(helpers.SessionData)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "INVALID_SESSION"})
	}

	accountUuid, err := uuid.Parse(sessionData.AccountUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "SERVER_ERROR"})
	}

	file_uuid_str := c.Params("file_uuid", "")

	file_uuid, file_uuid_err := uuid.Parse(file_uuid_str)

	if file_uuid_err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	queries := db_queries.New(pool)
	rows_changed, err := queries.DeleteFile(c.Context(), db_queries.DeleteFileParams{
		Uuid:        pgtype.UUID{Bytes: file_uuid, Valid: true},
		AccountUuid: pgtype.UUID{Bytes: accountUuid, Valid: true},
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "SERVER_ERROR"})
	}

	// forbidden because deletefile matches based on account uuid, and if its the wrong one, changed rows will be 0
	if rows_changed != 1 {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "FORBIDDEN"})
	}

	if helpers.GetEnvFallback("FILE_UPLOAD_BACKEND", "directory") == "directory" {
		root, err := os.OpenRoot(helpers.GetEnvFallback("FILE_UPLOADS_DIR", "./uploads"))

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "SERVER_ERROR"})
		}

		if err := root.Remove(file_uuid_str); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "SERVER_ERROR"})
		}
	} else {
		// TODO: add S3-compatible file storage support
		return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{"error": "NOT_IMPLEMENTED"})
	}

	return c.SendStatus(fiber.StatusOK)
}

func FetchFile(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	if helpers.GetEnvFallback("FILE_UPLOAD_BACKEND", "directory") != "directory" {
		// TODO: add S3-compatible file storage support
		return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{"error": "NOT_IMPLEMENTED"})
	}

	fs := os.DirFS(helpers.GetEnvFallback("FILE_UPLOADS_DIR", "./uploads"))

	file_uuid := c.Params("file_uuid", "")

	_, err := uuid.Parse(file_uuid)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	file, err := fs.Open(file_uuid)
	defer file.Close()

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "NOT_FOUND"})
	}

	return c.SendFile(file_uuid, fiber.SendFile{
		FS:       fs,
		Compress: true,
		MaxAge:   999999, // you cannot modify a file, and reuploads generate a new UUID
	})
}

func GetFile(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	file_uuid_str := c.Params("file_uuid", "")

	file_uuid, err := uuid.Parse(file_uuid_str)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	queries := db_queries.New(pool)

	data, err := helpers.CacheOrGetFile(c.Context(), rdb, *queries, file_uuid, helpers.GetInt32EnvFallback("FILE_CACHE_TTL", 5*60, 604800))

	if err != nil {
		error_code, status_code := helpers.MapDBError(err)
		return c.Status(status_code).JSON(fiber.Map{"error": error_code})
	}

	return c.JSON(data)
}

package meapi

import (
	"log/slog"
	"math"
	db_queries "onlystudents/internal/db/store"
	"onlystudents/internal/helpers"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mojatter/s2"
	"github.com/redis/go-redis/v9"
)

const (
	KiB = 1 << 10
	MiB = 1 << 20
)

type DeleteFileRequest struct {
	FileUUID string `json:"file_uuid"`
}

type File struct {
	Type            string `json:"type"`
	Size            int64  `json:"size"`
	HasOriginalName bool   `json:"has_original_name"`
	OriginalName    string `json:"original_name"`
	Uuid            string `json:"uuid"`
	CreatedAt       int64  `json:"created_at"`
}

func GetMyFiles(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	sessionData, ok := c.Locals("session").(helpers.SessionData)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "INVALID_SESSION"})
	}

	accountUuid, err := uuid.Parse(sessionData.AccountUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "SERVER_ERROR"})
	}

	queries := db_queries.New(pool)

	my_files, err := queries.GetMyFiles(c.Context(), pgtype.UUID{Bytes: accountUuid, Valid: true})

	summaries := []File{}

	for _, file := range my_files {
		summaries = append(summaries, File{
			Type:            file.Type,
			Size:            file.Size,
			HasOriginalName: file.OriginalName.Valid,
			OriginalName:    file.OriginalName.String,
			Uuid:            file.Uuid.String(),
			CreatedAt:       file.CreatedAt.Time.Unix(),
		})
	}

	return c.JSON(summaries)
}

func FileUpload(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client, s2_client s2.Storage) error {
	sessionData, ok := c.Locals("session").(helpers.SessionData)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "INVALID_SESSION"})
	}

	accountUuid, err := uuid.Parse(sessionData.AccountUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "SERVER_ERROR"})
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

	if float64(file.Size)/float64(KiB) > float64(helpers.GetInt64EnvFallback("MAX_FILE_SIZE", 10240, 99999)) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "FILE_TOO_BIG"})
	}

	file_reader, file_read_err := file.Open()

	if file_read_err != nil {
		slog.Debug("multipart file open failed", "file_read_err", file_read_err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "SERVER_ERROR"})
	}

	queries := db_queries.New(pool)

	total_file_size, err := queries.GetTotalFileSize(c.Context(), pgtype.UUID{Bytes: accountUuid, Valid: true})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "SERVER_ERROR"})
	}

	if (total_file_size+file.Size)/MiB >= helpers.GetInt64EnvFallback("USER_STORAGE_LIMIT", 500, 100000) {
		return c.Status(fiber.StatusInsufficientStorage).JSON(fiber.Map{"error": "INSUFFICIENT_ACCOUNT_STORAGE"})
	}

	s2_put_err := s2_client.Put(c.Context(), s2.NewObjectReader(fileUuid.String(), file_reader, uint64(min(file.Size, math.MaxInt64)))) // #nosec G115 -- value clamped to <= max, which fits uint64
	if s2_put_err != nil {
		slog.Error("file provider put failed", "s2_put_err", s2_put_err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "SERVER_ERROR"})
	}

	accountUuid, uuid_err := uuid.Parse(sessionData.AccountUUID)
	if uuid_err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "SERVER_ERROR"})
	}

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

func FileDelete(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client, s2_client s2.Storage) error {
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
		error_code, status_code := helpers.MapDBError(err)
		return c.Status(status_code).JSON(fiber.Map{"error": error_code})
	}

	// forbidden because deletefile matches based on account uuid, and if its the wrong one, changed rows will be 0
	if rows_changed != 1 {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "FORBIDDEN"})
	}

	s2_delete_err := s2_client.Delete(c.Context(), file_uuid_str)
	if s2_delete_err != nil {
		slog.Error("file provider delete failed", "s2_delete_err", s2_delete_err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "SERVER_ERROR"})
	}

	return c.SendStatus(fiber.StatusOK)
}

func FetchFile(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client, s2_client s2.Storage) error {
	file_uuid := c.Params("file_uuid", "") // NOTE: in the URL param its :file_uuid, do not use body
	_, err := uuid.Parse(file_uuid)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	object, err := s2_client.Get(c.Context(), file_uuid)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "NOT_FOUND"})
	}

	object_reader, err := object.Open()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "SERVER_ERROR"})
	}

	return c.SendStream(object_reader)
}

func GetFile(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	file_uuid_str := c.Params("file_uuid", "") // NOTE: in the URL param its :file_uuid, do not use body

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

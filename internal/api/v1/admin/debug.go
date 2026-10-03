package adminapi

import (
	"fmt"
	"io/fs"
	"log/slog"
	db_queries "onlystudents/internal/db/store"
	"onlystudents/internal/helpers"
	"os"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type LogFile struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type AdminStatusData struct {
	DBPing         int64  `json:"db_ping"`
	RedisStatus    string `json:"redis_ping"`
	DBQuerySuccess bool   `json:"db_query_success"`
	SchoolCount    int64  `json:"school_count"`
	AccountCount   int64  `json:"account_count"`
	StudentCount   int64  `json:"student_count"`
	TeacherCount   int64  `json:"teacher_count"`
	GuardianCount  int64  `json:"guardian_count"`
}

func AdminLogs(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	_, ok := c.Locals("session").(helpers.SessionData)

	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "UNAUTHORIZED"})
	}

	log_dir := helpers.LogDir()

	root, err := os.OpenRoot(log_dir)
	if err != nil {
		slog.Error("Error opening logs folder", "err", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "SERVER_ERROR"})
	}
	defer root.Close()

	root_fs := root.FS()

	log_json := []LogFile{}

	err = fs.WalkDir(root_fs, ".", func(path string, entry fs.DirEntry, err error) error {
		if err != nil || entry.IsDir() {
			return err
		}

		content, err := fs.ReadFile(root_fs, path)
		if err != nil {
			return err
		}

		log_json = append(log_json, LogFile{Name: path, Content: string(content)})

		return nil
	})

	if err != nil {
		slog.Error("Error reading logs folder", "err", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "SERVER_ERROR"})
	}

	return c.JSON(log_json)
}

func GetDBPing(c fiber.Ctx, pool *pgxpool.Pool) int64 {
	var db_ping int64

	start := time.Now().UnixMilli()
	err := pool.Ping(c.Context())

	if err != nil {
		db_ping = -1
	} else {
		db_ping = time.Now().UnixMilli() - start
	}

	return db_ping
}

func GetRedisStatus(c fiber.Ctx, rdb *redis.Client) string {
	var redis_db_status string

	start := time.Now().UnixMilli()
	_, err := rdb.Ping(c.Context()).Result()

	if err != nil {
		redis_db_status = fmt.Sprintf("Dragonfly returned error: %s in %d ms", err, time.Now().UnixMilli()-start)
	} else {
		redis_db_status = fmt.Sprintf("PONG in %d ms", time.Now().UnixMilli()-start)
	}

	return redis_db_status
}

func AdminStatus(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	_, ok := c.Locals("session").(helpers.SessionData)

	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "UNAUTHORIZED"})
	}

	queries := db_queries.New(pool)

	debug_data, err := queries.GetAdminDebugData(c.Context())

	if err != nil {
		slog.Error(fmt.Sprintf("Error getting debug data from Postgres DB: %s", err))
		return c.JSON(AdminStatusData{
			DBPing:         GetDBPing(c, pool),
			RedisStatus:    GetRedisStatus(c, rdb),
			DBQuerySuccess: false,
		})
	} else {
		return c.JSON(AdminStatusData{
			DBPing:         GetDBPing(c, pool),
			RedisStatus:    GetRedisStatus(c, rdb),
			DBQuerySuccess: true,
			SchoolCount:    debug_data.SchoolCount,
			AccountCount:   debug_data.AccountCount,
			StudentCount:   debug_data.StudentCount,
			TeacherCount:   debug_data.TeacherCount,
			GuardianCount:  debug_data.GuardianCount,
		})
	}
}

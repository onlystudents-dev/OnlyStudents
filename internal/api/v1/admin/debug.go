package adminapi

import (
	"io/fs"
	"log/slog"
	"onlystudents/internal/helpers"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type LogFile struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type runSQLRequest struct {
	Cmd string `json:"cmd"`
}

func AdminLogs(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	_, ok := c.Locals("session").(helpers.SessionData)

	if !ok {
		return c.SendStatus(401)
	}

	log_dir := helpers.LogDir()

	root, err := os.OpenRoot(log_dir)
	if err != nil {
		slog.Error("Error opening logs folder", "err", err)
		return c.SendStatus(500)
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
		return c.SendStatus(500)
	}

	return c.JSON(log_json)
}

func AdminRunSQL(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	_, ok := c.Locals("session").(helpers.SessionData)

	if !ok {
		return c.SendStatus(401)
	}

	var req runSQLRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.SendStatus(400)
	}

	cmd_tag, err := pool.Exec(c.Context(), req.Cmd)

	if err != nil {
		return c.Status(400).SendString(err.Error())
	} else {
		return c.Status(200).SendString(cmd_tag.String())
	}
}

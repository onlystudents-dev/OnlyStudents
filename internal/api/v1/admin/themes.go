package adminapi

import (
	"onlystudents/internal/helpers"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func AdminSetTheme(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var body struct {
		Name   string            `json:"name"`
		Colors map[string]string `json:"colors"`
	}

	if err := c.Bind().JSON(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	theme := helpers.Theme{
		Name:   body.Name,
		Colors: body.Colors,
	}

	if err := helpers.SetTheme(c, pool, rdb, theme.Name, theme.Colors); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "SERVER_ERROR"})
	}

	return c.SendStatus(fiber.StatusOK)
}

func AdminDeleteTheme(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	err := helpers.DeleteTheme(c, pool, rdb, string(c.Body()))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "SERVER_ERROR"})
	}

	return c.SendStatus(fiber.StatusOK)
}

func AdminGetThemes(c fiber.Ctx, pool *pgxpool.Pool) error {
	themes, err := helpers.GetThemes(c, pool)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "SERVER_ERROR"})
	}

	return c.JSON(themes)
}

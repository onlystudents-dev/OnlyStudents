package v1

import (
	"context"
	db_queries "onlystudents/internal/db/store"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
)

func School(c fiber.Ctx, pool *pgxpool.Pool) error {
	queries := db_queries.New(pool)

	schools, err := queries.ListSchools(context.Background())

	if err != nil {
		return c.SendStatus(500)
	}

	summaries := make([]fiber.Map, 0, len(schools))
	for _, s := range schools {
		summaries = append(summaries, fiber.Map{
			"name": s.Name,
			"id":   s.ID,
		})
	}

	return c.JSON(summaries)
}

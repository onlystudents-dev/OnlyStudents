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

func Teachers(c fiber.Ctx, pool *pgxpool.Pool) error {
	queries := db_queries.New(pool)
	teachers, err := queries.ListTeachers(context.Background())

	if err != nil {
		return c.SendStatus(500)
	}

	return c.JSON(teachers)
}

func Guardians(c fiber.Ctx, pool *pgxpool.Pool) error {
	queries := db_queries.New(pool)
	guardians, err := queries.ListGuardians(context.Background())

	if err != nil {
		return c.SendStatus(500)
	}

	return c.JSON(guardians)
}

func Students(c fiber.Ctx, pool *pgxpool.Pool) error {
	queries := db_queries.New(pool)
	students, err := queries.ListStudents(context.Background())

	if err != nil {
		return c.SendStatus(500)
	}

	return c.JSON(students)
}

package timetable

import (
	db_queries "onlystudents/internal/db/store"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CreateBellScheduleTypeRequest struct {
	Name string `json:"name"`
}

type EditBellScheduleTypeRequest struct {
	Id   int32  `json:"id"`
	Name string `json:"name"`
}

type DeleteBellScheduleTypeRequest struct {
	Id int32 `json:"id"`
}

func CreateBellScheduleType(c fiber.Ctx, pool *pgxpool.Pool) error {
	var req CreateBellScheduleTypeRequest

	token := c.Cookies("session_token")
	if token == "" {
		return c.SendStatus(401)
	}

	if err := c.Bind().Body(&req); err != nil {
		return c.SendStatus(400)
	}

	if req.Name == "" {
		return c.SendStatus(400)
	}

	queries := db_queries.New(pool)

	params := db_queries.CreateBellScheduleTypeParams{
		SchoolID: 1,
		Name:     req.Name,
	}

	err := queries.CreateBellScheduleType(c.Context(), params)
	if err != nil {
		return c.SendStatus(500)
	}

	return c.SendStatus(200)
}

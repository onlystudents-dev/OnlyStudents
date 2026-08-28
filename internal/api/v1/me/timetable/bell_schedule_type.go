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

	if err := c.Bind().Body(&req); err != nil {
		return c.SendStatus(400)
	}

	if req.Name == "" {
		return c.SendStatus(400)
	}

	queries := db_queries.New(pool)

}

package meapi

import (
	"errors"
	db_queries "onlystudents/internal/db/store"
	"onlystudents/internal/helpers"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type ReadLessonTimeRequest struct {
	TypeID int32 `json:"type_id" query:"type_id"`
}

func parseDateRange(c fiber.Ctx) (time.Time, time.Time, error) {
	start_date := c.Query("start_date", "none")

	if start_date == "none" {
		return time.Time{}, time.Time{}, errors.New("start date not found")
	}

	start_parsed, err := time.Parse(time.DateOnly, start_date)

	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	end_date := c.Query("end_date", "none")

	if end_date == "none" {
		return time.Time{}, time.Time{}, errors.New("end date not found")
	}

	end_parsed, err := time.Parse(time.DateOnly, end_date)

	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	return start_parsed, end_parsed, nil
}

func ReadMyRealTimeTable(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	scope, err := helpers.ResolveReaderScope(c, pool, rdb)

	if err != nil {
		return c.SendStatus(401)
	}

	start, end, err := parseDateRange(c)
	if err != nil {
		return c.SendStatus(400)
	}

	queries := db_queries.New(pool)

	data, err := queries.ReadRealTimeTable(c.Context(), db_queries.ReadRealTimeTableParams{
		SchoolID:     scope.SchoolID,
		ClassesID:    scope.ClassID,
		ActualDate:   pgtype.Date{Time: start, Valid: true},
		ActualDate_2: pgtype.Date{Time: end, Valid: true},
	})

	if err != nil {
		return c.SendStatus(500)
	}

	return c.JSON(data)
}

func ReadBaseSchedule(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	scope, err := helpers.ResolveReaderScope(c, pool, rdb)

	if err != nil {
		return c.SendStatus(401)
	}

	queries := db_queries.New(pool)

	data, err := queries.ReadBaseScheduleClass(c.Context(), db_queries.ReadBaseScheduleClassParams{
		SchoolID:  scope.SchoolID,
		ClassesID: scope.ClassID,
	})

	if err != nil {
		return c.SendStatus(500)
	}

	return c.JSON(data)
}

func ReadLessonTime(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req ReadLessonTimeRequest

	if err := c.Bind().Query(&req); err != nil {
		return c.SendStatus(400)
	}

	scope, err := helpers.ResolveReaderScope(c, pool, rdb)

	if err != nil {
		return c.SendStatus(401)
	}

	queries := db_queries.New(pool)

	data, err := queries.ReadLessonTime(c.Context(), db_queries.ReadLessonTimeParams{
		SchoolID: scope.SchoolID,
		TypeID:   req.TypeID,
	})

	if err != nil {
		return c.SendStatus(500)
	}

	return c.JSON(data)
}

func ReadRoom(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	scope, err := helpers.ResolveReaderScope(c, pool, rdb)

	if err != nil {
		return c.SendStatus(401)
	}

	queries := db_queries.New(pool)

	data, err := queries.ReadRoom(c.Context(), scope.SchoolID)
	if err != nil {
		return c.SendStatus(500)
	}

	return c.JSON(data)
}

func ReadCustomSubject(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	scope, err := helpers.ResolveReaderScope(c, pool, rdb)

	if err != nil {
		return c.SendStatus(401)
	}

	queries := db_queries.New(pool)

	data, err := queries.ReadCustomSubject(c.Context(), scope.SchoolID)
	if err != nil {
		return c.SendStatus(500)
	}

	return c.JSON(data)
}

func ReadBellScheduleType(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	scope, err := helpers.ResolveReaderScope(c, pool, rdb)

	if err != nil {
		return c.SendStatus(401)
	}

	queries := db_queries.New(pool)

	data, err := queries.ReadBellScheduleType(c.Context(), scope.SchoolID)
	if err != nil {
		return c.SendStatus(500)
	}

	return c.JSON(data)
}

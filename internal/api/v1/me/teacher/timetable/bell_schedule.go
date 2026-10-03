package timetable

import (
	db_queries "onlystudents/internal/db/store"
	"onlystudents/internal/helpers"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
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

type CreateLessonTimeRequest struct {
	TypeID       int32     `json:"type_id"`
	LessonNumber int32     `json:"lesson_number"`
	Start        time.Time `json:"lesson_start"`
	End          time.Time `json:"lesson_stop"`
}

type EditLessonTimeRequest struct {
	Id           int32     `json:"id"`
	LessonNumber int32     `json:"lesson_number"`
	Start        time.Time `json:"lesson_start"`
	End          time.Time `json:"lesson_stop"`
}

type DeleteLessonTimeRequest struct {
	Id int32 `json:"id"`
}

type ReadLessonTimeRequest struct {
	TypeID int32 `json:"type_id" query:"type_id"`
}

func CreateBellScheduleType(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req CreateBellScheduleTypeRequest

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	scope, status_code := helpers.ResolveTeacherCapabilityScope(c, pool, rdb, "MANAGE_BELL_SCHEDULE")

	if status_code != fiber.StatusOK {
		return c.SendStatus(status_code)
	}

	queries := db_queries.New(pool)

	params := db_queries.CreateBellScheduleTypeParams{
		SchoolID: scope.SchoolID,
		Name:     req.Name,
	}

	err := queries.CreateBellScheduleType(c.Context(), params)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "SERVER_ERROR"})
	}

	return c.SendStatus(fiber.StatusOK)
}

func DeleteBellScheduleType(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req DeleteBellScheduleTypeRequest

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	scope, status_code := helpers.ResolveTeacherCapabilityScope(c, pool, rdb, "MANAGE_BELL_SCHEDULE")

	if status_code != fiber.StatusOK {
		return c.SendStatus(status_code)
	}

	queries := db_queries.New(pool)

	params := db_queries.DeleteBellScheduleTypeParams{
		SchoolID: scope.SchoolID,
		ID:       req.Id,
	}

	err := queries.DeleteBellScheduleType(c.Context(), params)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "SERVER_ERROR"})
	}

	return c.SendStatus(fiber.StatusOK)
}

func EditBellScheduleType(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req EditBellScheduleTypeRequest

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	scope, status_code := helpers.ResolveTeacherCapabilityScope(c, pool, rdb, "MANAGE_BELL_SCHEDULE")

	if status_code != fiber.StatusOK {
		return c.SendStatus(status_code)
	}

	queries := db_queries.New(pool)

	params := db_queries.EditBellScheduleTypeParams{
		Name:     req.Name,
		SchoolID: scope.SchoolID,
		ID:       req.Id,
	}

	err := queries.EditBellScheduleType(c.Context(), params)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "SERVER_ERROR"})
	}

	return c.SendStatus(fiber.StatusOK)
}

func ReadBellScheduleType(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	scope, status_code := helpers.ResolveTeacherCapabilityScope(c, pool, rdb, "MANAGE_BELL_SCHEDULE")

	if status_code != fiber.StatusOK {
		return c.SendStatus(status_code)
	}

	queries := db_queries.New(pool)

	listoftype, err := queries.ReadBellScheduleType(c.Context(), int32(scope.SchoolID))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "SERVER_ERROR"})
	}

	return c.JSON(listoftype)
}

func CreateLessonTime(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req CreateLessonTimeRequest

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	if req.TypeID == 0 || req.LessonNumber <= 0 || req.Start.IsZero() || req.End.IsZero() {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	scope, status_code := helpers.ResolveTeacherCapabilityScope(c, pool, rdb, "MANAGE_BELL_SCHEDULE")

	if status_code != fiber.StatusOK {
		return c.SendStatus(status_code)
	}

	queries := db_queries.New(pool)

	startMicro := int64(req.Start.Hour()*3600+req.Start.Minute()*60+req.Start.Second()*1)*1_000_000 + int64(req.Start.Nanosecond()/1000)
	stopMicro := int64(req.End.Hour()*3600+req.End.Minute()*60+req.End.Second()*1)*1_000_000 + int64(req.End.Nanosecond()/1000)

	params := db_queries.CreateLessonTimeParams{
		SchoolID:     int32(scope.SchoolID),
		TypeID:       req.TypeID,
		LessonNumber: pgtype.Int4{Int32: req.LessonNumber, Valid: true},
		AtStart:      pgtype.Time{Microseconds: startMicro, Valid: true},
		AtEnd:        pgtype.Time{Microseconds: stopMicro, Valid: true},
	}

	err := queries.CreateLessonTime(c.Context(), params)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "SERVER_ERROR"})
	}

	return c.SendStatus(fiber.StatusOK)
}

func DeleteLessonTime(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req DeleteLessonTimeRequest

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	if req.Id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	scope, status_code := helpers.ResolveTeacherCapabilityScope(c, pool, rdb, "MANAGE_BELL_SCHEDULE")

	if status_code != fiber.StatusOK {
		return c.SendStatus(status_code)
	}

	queries := db_queries.New(pool)

	params := db_queries.DeleteLessonTimeParams{
		SchoolID: scope.SchoolID,
		ID:       req.Id,
	}

	err := queries.DeleteLessonTime(c.Context(), params)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "SERVER_ERROR"})
	}

	return c.SendStatus(fiber.StatusOK)
}

func EditLessonTime(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req EditLessonTimeRequest

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	if req.Id == 0 || req.LessonNumber <= 0 || req.Start.IsZero() || req.End.IsZero() {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	scope, status_code := helpers.ResolveTeacherCapabilityScope(c, pool, rdb, "MANAGE_BELL_SCHEDULE")

	if status_code != fiber.StatusOK {
		return c.SendStatus(status_code)
	}

	queries := db_queries.New(pool)

	startMicro := int64(req.Start.Hour()*3600+req.Start.Minute()*60+req.Start.Second()*1)*1_000_000 + int64(req.Start.Nanosecond()/1000)
	stopMicro := int64(req.End.Hour()*3600+req.End.Minute()*60+req.End.Second()*1)*1_000_000 + int64(req.End.Nanosecond()/1000)

	params := db_queries.EditLessonTimeParams{
		SchoolID:     int32(scope.SchoolID),
		ID:           req.Id,
		LessonNumber: pgtype.Int4{Int32: req.LessonNumber, Valid: true},
		AtStart:      pgtype.Time{Microseconds: startMicro, Valid: true},
		AtEnd:        pgtype.Time{Microseconds: stopMicro, Valid: true},
	}

	err := queries.EditLessonTime(c.Context(), params)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "SERVER_ERROR"})
	}

	return c.SendStatus(fiber.StatusOK)
}

func ReadLessonTime(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req ReadLessonTimeRequest

	if err := c.Bind().Query(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	if req.TypeID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	scope, status_code := helpers.ResolveTeacherCapabilityScope(c, pool, rdb, "MANAGE_BELL_SCHEDULE")

	if status_code != fiber.StatusOK {
		return c.SendStatus(status_code)
	}

	queries := db_queries.New(pool)

	params := db_queries.ReadLessonTimeParams{
		SchoolID: scope.SchoolID,
		TypeID:   req.TypeID,
	}

	data, err := queries.ReadLessonTime(c.Context(), params)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "SERVER_ERROR"})
	}

	return c.JSON(data)
}

package timetable

import (
	db_queries "onlystudents/internal/db/store"
	"onlystudents/internal/helpers"
	"strconv"
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
	Type_Id      int32     `json:"type_id"`
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
	Type_id int32 `json:"type_id"`
}

func CreateBellScheduleType(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req CreateBellScheduleTypeRequest

	school_id, err := strconv.ParseInt(c.Get("X-School"), 10, 32)
	if err != nil {
		return c.SendStatus(500)
	}

	if req.Name == "" {
		return c.SendStatus(400)
	}

	has_permission := helpers.CheckPermission(c, pool, rdb, "MANAGE_BELL_SCHEDULE", school_id)

	if has_permission == false {
		return c.SendStatus(401)
	}

	queries := db_queries.New(pool)

	params := db_queries.CreateBellScheduleTypeParams{
		SchoolID: int32(school_id),
		Name:     req.Name,
	}

	err = queries.CreateBellScheduleType(c.Context(), params)
	if err != nil {
		return c.SendStatus(500)
	}

	return c.SendStatus(200)
}

func DeleteBellScheduleType(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req DeleteBellScheduleTypeRequest

	school_id, err := strconv.ParseInt(c.Get("X-School"), 10, 32)

	if err != nil {
		c.SendStatus(500)
	}

	if req.Id == 0 {
		c.SendStatus(400)
	}

	has_permission := helpers.CheckPermission(c, pool, rdb, "MANAGE_BELL_SCHEDULE", school_id)

	if has_permission == false {
		return c.SendStatus(401)
	}

	queries := db_queries.New(pool)

	params := db_queries.DeleteBellScheduleTypeParams{
		SchoolID: int32(school_id),
		ID:       req.Id,
	}

	err = queries.DeleteBellScheduleType(c.Context(), params)

	if err != nil {
		return c.SendStatus(500)
	}

	return c.SendStatus(200)
}

func EditBellScheduleType(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req EditBellScheduleTypeRequest

	school_id, err := strconv.ParseInt(c.Get("X-School"), 10, 32)

	if err != nil {
		return c.SendStatus(500)
	}

	if req.Id == 0 || req.Name == "" {
		c.SendStatus(400)
	}

	has_permission := helpers.CheckPermission(c, pool, rdb, "MANAGE_BELL_SCHEDULE", school_id)

	if has_permission == false {
		return c.SendStatus(401)
	}

	queries := db_queries.New(pool)

	params := db_queries.EditBellScheduleTypeParams{
		Name:     req.Name,
		SchoolID: int32(school_id),
		ID:       req.Id,
	}

	err = queries.EditBellScheduleType(c.Context(), params)

	if err != nil {
		return c.SendStatus(500)
	}

	return c.SendStatus(200)
}

func ReadBellScheduleType(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	school_id, err := strconv.ParseInt(c.Get("X-School"), 10, 32)

	if err != nil {
		return c.SendStatus(500)
	}

	if school_id == 0 {
		return c.SendStatus(500)
	}

	has_permission := helpers.CheckPermission(c, pool, rdb, "MANAGE_BELL_SCHEDULE", school_id)

	if has_permission == false {
		return c.SendStatus(401)
	}

	queries := db_queries.New(pool)

	listoftype, err := queries.ReadBellScheduleType(c.Context(), int32(school_id))

	if err != nil {
		return c.SendStatus(500)
	}

	return c.JSON(listoftype)
}

func CreateLessonTime(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req CreateLessonTimeRequest

	school_id, err := strconv.ParseInt(c.Get("X-School"), 10, 32)

	if err != nil {
		return c.SendStatus(500)
	}

	if school_id == 0 {
		return c.SendStatus(400)
	}

	if req.LessonNumber >= 0 || req.Start.IsZero() || req.End.IsZero() {
		return c.SendStatus(500)
	}

	has_permission := helpers.CheckPermission(c, pool, rdb, "MANAGE_BELL_SCHEDULE", school_id)

	if has_permission == false {
		return c.SendStatus(401)
	}

	queries := db_queries.New(pool)

	startMicro := int64(req.Start.Hour()*3600+req.Start.Minute()*60+req.Start.Second()*1)*1_000_000 + int64(req.Start.Nanosecond()/1000)
	stopMicro := int64(req.End.Hour()*3600+req.End.Minute()*60+req.End.Second()*1)*1_000_000 + int64(req.End.Nanosecond()/1000)

	params := db_queries.CreateLessonTimeParams{
		SchoolID:     int32(school_id),
		TypeID:       req.Type_Id,
		LessonNumber: pgtype.Int4{Int32: req.LessonNumber, Valid: true},
		AtStart:      pgtype.Time{Microseconds: startMicro, Valid: true},
		AtEnd:        pgtype.Time{Microseconds: stopMicro, Valid: true},
	}

	err = queries.CreateLessonTime(c.Context(), params)

	if err != nil {
		return c.SendStatus(500)
	}

	return c.SendStatus(200)
}

func DeleteLessonTime(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req DeleteLessonTimeRequest

	school_id, err := strconv.ParseInt(c.Get("X-School"), 10, 32)

	if err != nil {
		return c.SendStatus(500)
	}

	if school_id == 0 {
		return c.SendStatus(500)
	}

	if req.Id == 0 {
		return c.SendStatus(400)
	}

	has_permission := helpers.CheckPermission(c, pool, rdb, "MANAGE_BELL_SCHEDULE", school_id)

	if has_permission == false {
		return c.SendStatus(401)
	}

	queries := db_queries.New(pool)

	params := db_queries.DeleteLessonTimeParams{
		SchoolID: int32(school_id),
		ID:       req.Id,
	}

	err = queries.DeleteLessonTime(c.Context(), params)

	if err != nil {
		return c.SendStatus(500)
	}

	return c.SendStatus(200)
}

func EditLessonTime(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req EditLessonTimeRequest

	school_id, err := strconv.ParseInt(c.Get("X-School"), 10, 32)

	if err != nil {
		return c.SendStatus(500)
	}

	if school_id == 0 {
		return c.SendStatus(500)
	}

	if req.Id == 0 || req.LessonNumber >= 0 || req.Start.IsZero() || req.End.IsZero() {
		return c.SendStatus(400)
	}

	has_permission := helpers.CheckPermission(c, pool, rdb, "MANAGE_BELL_SCHEDULE", school_id)

	if has_permission == false {
		return c.SendStatus(401)
	}

	queries := db_queries.New(pool)

	startMicro := int64(req.Start.Hour()*3600+req.Start.Minute()*60+req.Start.Second()*1)*1_000_000 + int64(req.Start.Nanosecond()/1000)
	stopMicro := int64(req.End.Hour()*3600+req.End.Minute()*60+req.End.Second()*1)*1_000_000 + int64(req.End.Nanosecond()/1000)

	params := db_queries.EditLessonTimeParams{
		SchoolID:     int32(school_id),
		ID:           req.Id,
		LessonNumber: pgtype.Int4{Int32: req.LessonNumber, Valid: true},
		AtStart:      pgtype.Time{Microseconds: startMicro, Valid: true},
		AtEnd:        pgtype.Time{Microseconds: stopMicro, Valid: true},
	}

	err = queries.EditLessonTime(c.Context(), params)

	if err != nil {
		return c.SendStatus(500)
	}

	return c.SendStatus(200)
}

func ReadLessonTime(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req ReadLessonTimeRequest

	school_id, err := strconv.ParseInt(c.Get("X-School"), 10, 32)

	if err != nil {
		return c.SendStatus(500)
	}

	if school_id == 0 {
		return c.SendStatus(500)
	}

	if req.Type_id == 0 {
		return c.SendStatus(500)
	}

	has_permission := helpers.CheckPermission(c, pool, rdb, "MANAGE_BELL_SCHEDULE", school_id)

	if has_permission == false {
		return c.SendStatus(401)
	}

	queries := db_queries.New(pool)

	params := db_queries.ReadLessonTimeParams{
		SchoolID: int32(school_id),
		TypeID:   req.Type_id,
	}

	data, err := queries.ReadLessonTime(c.Context(), params)

	if err != nil {
		return c.SendStatus(500)
	}

	return c.JSON(data)
}

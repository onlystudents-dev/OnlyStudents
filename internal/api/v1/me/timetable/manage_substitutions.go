package timetable

import (
	db_queries "onlystudents/internal/db/store"
	"onlystudents/internal/helpers"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type CanceledLessonRequest struct {
	Date            pgtype.Date `json:"date"`
	LessonNumber    int32       `json:"lesson_number"`
	TeacherId       int32       `json:"teacher_id"`
	RoomId          int32       `json:"room_id"`
	DayOfWeek       int32       `json:"day_of_week"`
	GroupId         int32       `json:"group_id"`
	isCustomSubject bool        `json:"is_custom_subject"`
	CustomSubjectId int32       `json:"custom_subject_id"`
	SubjectId       int32       `json:"subject_id"`
}

func AddCanceledLesson(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {

	var req CanceledLessonRequest

	token := c.Cookies("session-token")

	school_id, err := strconv.ParseInt(c.Get("X-School"), 10, 32)

	if err != nil {
		return c.SendStatus(500)
	}

	if token == "" {
		return c.SendStatus(401)
	}

	if school_id == 0 {
		return c.SendStatus(500)
	}

	if req.Date.Time.IsZero() || req.LessonNumber < 0 || req.TeacherId == 0 || req.RoomId == 0 || req.DayOfWeek == 0 || req.GroupId == 0 || (req.isCustomSubject == false && req.SubjectId == 0) || (req.isCustomSubject == true && req.CustomSubjectId == 0) {
		return c.SendStatus(400)
	}

	has_permission := helpers.CheckPermission(c.Context(), pool, rdb, token, "MANAGE_SUBSTITUTIONS", school_id)

	if has_permission == false {
		return c.SendStatus(401)
	}

	queries := db_queries.New(pool)

	params := db_queries.AddCanceledLessonParams{
		SchoolID:        int32(school_id),
		TeacherID:       req.TeacherId,
		RoomID:          req.RoomId,
		DayOfWeek:       req.DayOfWeek,
		GroupID:         req.GroupId,
		CustomSubject:   req.isCustomSubject,
		CustomSubjectID: pgtype.Int4{Int32: req.CustomSubjectId, Valid: true},
		SubjectID:       pgtype.Int4{Int32: req.SubjectId, Valid: true},
		ActualDate:      req.Date,
		LessonNum:       req.LessonNumber,
	}

	err = queries.AddCanceledLesson(c.Context(), params)

	if err != nil {
		return c.SendStatus(500)
	}

	return c.SendStatus(200)
}

func RemoveCanceledLesson(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req CanceledLessonRequest

	token := c.Cookies("session-token")

	school_id, err := strconv.ParseInt(c.Get("X-School"), 10, 32)

	if err != nil {
		return c.SendStatus(500)
	}

	if token == "" {
		return c.SendStatus(401)
	}

	if school_id == 0 {
		return c.SendStatus(500)
	}

	has_permission := helpers.CheckPermission(c.Context(), pool, rdb, token, "MANAGE_SUBSTITUTIONS", school_id)

	if has_permission == false {
		return c.SendStatus(401)
	}

	queries := db_queries.New(pool)

	params := db_queries.RemoveCanceledLessonParams{
		SchoolID:        int32(school_id),
		TeacherID:       req.TeacherId,
		RoomID:          req.RoomId,
		DayOfWeek:       req.DayOfWeek,
		GroupID:         req.GroupId,
		CustomSubject:   req.isCustomSubject,
		CustomSubjectID: pgtype.Int4{Int32: req.CustomSubjectId, Valid: true},
		SubjectID:       pgtype.Int4{Int32: req.SubjectId, Valid: true},
		ActualDate:      req.Date,
		LessonNum:       req.LessonNumber,
	}

	err = queries.RemoveCanceledLesson(c.Context(), params)

	if err != nil {
		return c.SendStatus(500)
	}

	return c.SendStatus(200)
}

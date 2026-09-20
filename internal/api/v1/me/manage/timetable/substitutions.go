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
	IsCustomSubject bool        `json:"is_custom_subject"`
	CustomSubjectId int32       `json:"custom_subject_id"`
	SubjectId       int32       `json:"subject_id"`
}

type SubstitutionsLessonRequest struct {
	ID                    int32       `json:"id"`
	Date                  pgtype.Date `json:"date"`
	LessonNumber          int32       `json:"lesson_number"`
	TeacherId             int32       `json:"teacher_id"`
	RoomId                int32       `json:"room_id"`
	DayOfWeek             int32       `json:"day_of_week"`
	GroupId               int32       `json:"group_id"`
	IsCustomSubject       bool        `json:"is_custom_subject"`
	CustomSubjectId       int32       `json:"custom_subject_id"`
	SubjectId             int32       `json:"subject_id"`
	IsSubstitution        bool        `json:"is_substitution"`
	SubstitutionTeacherId int32       `json:"substitution_teacher_id"`
}

func AddCanceledLesson(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req CanceledLessonRequest

	if err := c.Bind().Body(&req); err != nil {
		return c.SendStatus(400)
	}

	school_id, err := strconv.ParseInt(c.Get("X-School"), 10, 32)

	if err != nil {
		return c.SendStatus(400)
	}

	if school_id == 0 {
		return c.SendStatus(400)
	}

	if req.Date.Time.IsZero() || req.LessonNumber <= 0 || req.TeacherId == 0 || req.RoomId == 0 || req.DayOfWeek == 0 || req.GroupId == 0 || (req.IsCustomSubject == false && req.SubjectId == 0) || (req.IsCustomSubject == true && req.CustomSubjectId == 0) {
		return c.SendStatus(400)
	}

	has_permission := helpers.CheckPermission(c, pool, rdb, "MANAGE_SUBSTITUTIONS", school_id)

	if !has_permission {
		return c.SendStatus(403)
	}

	queries := db_queries.New(pool)

	var custom_subject_id pgtype.Int4
	var subject_id pgtype.Int4

	if req.IsCustomSubject {
		custom_subject_id = pgtype.Int4{Int32: req.CustomSubjectId, Valid: true}
		subject_id = pgtype.Int4{Valid: false}
	} else {
		custom_subject_id = pgtype.Int4{Valid: false}
		subject_id = pgtype.Int4{Int32: req.SubjectId, Valid: true}
	}

	params := db_queries.AddCanceledLessonParams{
		SchoolID:        int32(school_id),
		TeacherID:       req.TeacherId,
		RoomID:          req.RoomId,
		DayOfWeek:       req.DayOfWeek,
		GroupID:         req.GroupId,
		CustomSubject:   req.IsCustomSubject,
		CustomSubjectID: custom_subject_id,
		SubjectID:       subject_id,
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

	if err := c.Bind().Body(&req); err != nil {
		return c.SendStatus(400)
	}

	school_id, err := strconv.ParseInt(c.Get("X-School"), 10, 32)

	if err != nil {
		return c.SendStatus(400)
	}

	if school_id == 0 {
		return c.SendStatus(400)
	}

	if req.Date.Time.IsZero() || req.LessonNumber <= 0 || req.TeacherId == 0 || req.RoomId == 0 || req.DayOfWeek == 0 || req.GroupId == 0 || (req.IsCustomSubject == false && req.SubjectId == 0) || (req.IsCustomSubject == true && req.CustomSubjectId == 0) {
		return c.SendStatus(400)
	}

	has_permission := helpers.CheckPermission(c, pool, rdb, "MANAGE_SUBSTITUTIONS", school_id)

	if !has_permission {
		return c.SendStatus(403)
	}

	queries := db_queries.New(pool)

	var custom_subject_id pgtype.Int4
	var subject_id pgtype.Int4

	if req.IsCustomSubject {
		custom_subject_id = pgtype.Int4{Int32: req.CustomSubjectId, Valid: true}
		subject_id = pgtype.Int4{Valid: false}
	} else {
		custom_subject_id = pgtype.Int4{Valid: false}
		subject_id = pgtype.Int4{Int32: req.SubjectId, Valid: true}
	}

	params := db_queries.RemoveCanceledLessonParams{
		SchoolID:        int32(school_id),
		TeacherID:       req.TeacherId,
		RoomID:          req.RoomId,
		DayOfWeek:       req.DayOfWeek,
		GroupID:         req.GroupId,
		CustomSubject:   req.IsCustomSubject,
		CustomSubjectID: custom_subject_id,
		SubjectID:       subject_id,
		ActualDate:      req.Date,
		LessonNum:       req.LessonNumber,
	}

	err = queries.RemoveCanceledLesson(c.Context(), params)

	if err != nil {
		return c.SendStatus(500)
	}

	return c.SendStatus(200)
}

func UpdateSubstitution(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req SubstitutionsLessonRequest

	if err := c.Bind().Body(&req); err != nil {
		return c.SendStatus(400)
	}

	school_id, err := strconv.ParseInt(c.Get("X-School"), 10, 32)

	if err != nil {
		return c.SendStatus(400)
	}

	if school_id == 0 {
		return c.SendStatus(400)
	}

	if req.ID == 0 || req.Date.Time.IsZero() || req.LessonNumber <= 0 || req.TeacherId == 0 || req.RoomId == 0 || req.DayOfWeek == 0 || req.GroupId == 0 || (req.IsCustomSubject == false && req.SubjectId == 0) || (req.IsCustomSubject == true && req.CustomSubjectId == 0 || !req.IsSubstitution || req.SubstitutionTeacherId == 0) {
		return c.SendStatus(400)
	}

	has_permission := helpers.CheckPermission(c, pool, rdb, "MANAGE_SUBSTITUTIONS", school_id)

	if !has_permission {
		return c.SendStatus(403)
	}

	queries := db_queries.New(pool)

	var custom_subject_id pgtype.Int4
	var subject_id pgtype.Int4

	if req.IsCustomSubject {
		custom_subject_id = pgtype.Int4{Int32: req.CustomSubjectId, Valid: true}
		subject_id = pgtype.Int4{Valid: false}
	} else {
		custom_subject_id = pgtype.Int4{Valid: false}
		subject_id = pgtype.Int4{Int32: req.SubjectId, Valid: true}
	}

	var substitution_teacher_id pgtype.Int4

	if req.IsSubstitution {
		substitution_teacher_id = pgtype.Int4{Int32: req.SubstitutionTeacherId, Valid: true}
	} else {
		substitution_teacher_id = pgtype.Int4{Valid: false}
	}

	params := db_queries.ManageSubsitutionLessonParams{
		SchoolID:              int32(school_id),
		TeacherID:             req.TeacherId,
		RoomID:                req.RoomId,
		DayOfWeek:             req.DayOfWeek,
		GroupID:               req.GroupId,
		CustomSubject:         req.IsCustomSubject,
		CustomSubjectID:       custom_subject_id,
		SubjectID:             subject_id,
		ActualDate:            req.Date,
		LessonNum:             req.LessonNumber,
		IsSubstitution:        req.IsSubstitution,
		SubstitutionTeacherID: substitution_teacher_id,
		ID:                    req.ID,
	}

	err = queries.ManageSubsitutionLesson(c.Context(), params)

	if err != nil {
		return c.SendStatus(500)
	}

	return c.SendStatus(200)
}

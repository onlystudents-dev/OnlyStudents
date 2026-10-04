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

type CanceledLessonRequest struct {
	Date            int64 `json:"date"`
	LessonNumber    int32 `json:"lesson_number"`
	TeacherId       int32 `json:"teacher_id"`
	RoomId          int32 `json:"room_id"`
	DayOfWeek       int32 `json:"day_of_week"`
	GroupId         int32 `json:"group_id"`
	IsCustomSubject bool  `json:"is_custom_subject"`
	CustomSubjectId int32 `json:"custom_subject_id"`
	SubjectId       int32 `json:"subject_id"`
}

type SubstitutionsLessonRequest struct {
	ID                    int32 `json:"id"`
	Date                  int64 `json:"date"`
	LessonNumber          int32 `json:"lesson_number"`
	TeacherId             int32 `json:"teacher_id"`
	RoomId                int32 `json:"room_id"`
	DayOfWeek             int32 `json:"day_of_week"`
	GroupId               int32 `json:"group_id"`
	IsCustomSubject       bool  `json:"is_custom_subject"`
	CustomSubjectId       int32 `json:"custom_subject_id"`
	SubjectId             int32 `json:"subject_id"`
	IsSubstitution        bool  `json:"is_substitution"`
	SubstitutionTeacherId int32 `json:"substitution_teacher_id"`
}

func AddCanceledLesson(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req CanceledLessonRequest

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	if req.Date == 0 || req.LessonNumber <= 0 || req.TeacherId == 0 || req.RoomId == 0 || req.DayOfWeek == 0 || req.GroupId == 0 || (req.IsCustomSubject == false && req.SubjectId == 0) || (req.IsCustomSubject == true && req.CustomSubjectId == 0) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	scope, status_code := helpers.ResolveTeacherCapabilityScope(c, pool, rdb, "MANAGE_SUBSTITUTIONS")

	if status_code != fiber.StatusOK {
		return helpers.ErrorByStatusCode(c, status_code)
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
		SchoolID:        scope.SchoolID,
		TeacherID:       req.TeacherId,
		RoomID:          req.RoomId,
		DayOfWeek:       req.DayOfWeek,
		GroupID:         req.GroupId,
		CustomSubject:   req.IsCustomSubject,
		CustomSubjectID: custom_subject_id,
		SubjectID:       subject_id,
		ActualDate:      pgtype.Date{Time: time.Unix(req.Date, 0), Valid: true},
		LessonNum:       req.LessonNumber,
	}

	err := queries.AddCanceledLesson(c.Context(), params)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "SERVER_ERROR"})
	}

	return c.SendStatus(fiber.StatusOK)
}

func RemoveCanceledLesson(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req CanceledLessonRequest

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	if req.Date == 0 || req.LessonNumber <= 0 || req.TeacherId == 0 || req.RoomId == 0 || req.DayOfWeek == 0 || req.GroupId == 0 || (req.IsCustomSubject == false && req.SubjectId == 0) || (req.IsCustomSubject == true && req.CustomSubjectId == 0) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	scope, status_code := helpers.ResolveTeacherCapabilityScope(c, pool, rdb, "MANAGE_SUBSTITUTIONS")

	if status_code != fiber.StatusOK {
		return helpers.ErrorByStatusCode(c, status_code)
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
		SchoolID:        scope.SchoolID,
		TeacherID:       req.TeacherId,
		RoomID:          req.RoomId,
		DayOfWeek:       req.DayOfWeek,
		GroupID:         req.GroupId,
		CustomSubject:   req.IsCustomSubject,
		CustomSubjectID: custom_subject_id,
		SubjectID:       subject_id,
		ActualDate:      pgtype.Date{Time: time.Unix(req.Date, 0), Valid: true},
		LessonNum:       req.LessonNumber,
	}

	err := queries.RemoveCanceledLesson(c.Context(), params)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "SERVER_ERROR"})
	}

	return c.SendStatus(fiber.StatusOK)
}

func UpdateSubstitution(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req SubstitutionsLessonRequest

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	if req.ID == 0 || req.Date == 0 || req.LessonNumber <= 0 || req.TeacherId == 0 || req.RoomId == 0 || req.DayOfWeek == 0 || req.GroupId == 0 || (req.IsCustomSubject == false && req.SubjectId == 0) || (req.IsCustomSubject == true && req.CustomSubjectId == 0 || !req.IsSubstitution || req.SubstitutionTeacherId == 0) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	scope, status_code := helpers.ResolveTeacherCapabilityScope(c, pool, rdb, "MANAGE_SUBSTITUTIONS")

	if status_code != fiber.StatusOK {
		return helpers.ErrorByStatusCode(c, status_code)
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
		SchoolID:              scope.SchoolID,
		TeacherID:             req.TeacherId,
		RoomID:                req.RoomId,
		DayOfWeek:             req.DayOfWeek,
		GroupID:               req.GroupId,
		CustomSubject:         req.IsCustomSubject,
		CustomSubjectID:       custom_subject_id,
		SubjectID:             subject_id,
		ActualDate:            pgtype.Date{Time: time.Unix(req.Date, 0), Valid: true},
		LessonNum:             req.LessonNumber,
		IsSubstitution:        req.IsSubstitution,
		SubstitutionTeacherID: substitution_teacher_id,
		ID:                    req.ID,
	}

	err := queries.ManageSubsitutionLesson(c.Context(), params)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "SERVER_ERROR"})
	}

	return c.SendStatus(fiber.StatusOK)
}

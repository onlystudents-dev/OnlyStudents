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

type CreateBaseScheduleLessonRequest struct {
	TeacherId       int32 `json:"teacher_id"`
	DayOfWeek       int32 `json:"day_of_week"`
	LessonNumber    int32 `json:"lesson_number"`
	RoomId          int32 `json:"room_id"`
	GroupId         int32 `json:"group_id"`
	IsCustomSubject bool  `json:"is_custom_subject"`
	SubjectId       int32 `json:"subject_id"`
	CustomSubjectId int32 `json:"custom_subject_id"`
}

type DeleteBaseScheduleLessonRequest struct {
	Id int32 `json:"id"`
}

type EditBaseScheduleLessonRequest struct {
	Id              int32 `json:"id"`
	TeacherId       int32 `json:"teacher_id"`
	DayOfWeek       int32 `json:"day_of_week"`
	LessonNumber    int32 `json:"lesson_number"`
	RoomId          int32 `json:"room_id"`
	GroupId         int32 `json:"group_id"`
	IsCustomSubject bool  `json:"is_custom_subject"`
	SubjectId       int32 `json:"subject_id"`
	CustomSubjectId int32 `json:"custom_subject_id"`
}

type ReadBaseScheduleClassRequest struct {
	ClassId int32 `json:"class_id" query:"class_id"`
}

type ReadBaseScheduleGroupRequest struct {
	GroupId int32 `json:"group_id" query:"group_id"`
}

type ReadRealTimetableRequest struct {
	ClassId int32  `json:"class_id" query:"class_id"`
	Start   string `json:"start_date" query:"start_date"`
	End     string `json:"end_date" query:"end_date"`
}

type CreateRealTimeLessonRequest struct {
	TeacherId       int32       `json:"teacher_id"`
	RoomId          int32       `json:"room_id"`
	DayOfWeek       int32       `json:"day_of_week"`
	GroupId         int32       `json:"group_id"`
	IsCustomSubject bool        `json:"is_custom_subject"`
	CustomSubjectId int32       `json:"custom_subject_id"`
	SubjectId       int32       `json:"subject_id"`
	ActualDate      pgtype.Date `json:"actual_date"`
	LessonNumber    int32       `json:"lesson_number"`
}

type UpdateRealTimeLessonRequest struct {
	Id              int32       `json:"id"`
	TeacherId       int32       `json:"teacher_id"`
	RoomId          int32       `json:"room_id"`
	DayOfWeek       int32       `json:"day_of_week"`
	GroupId         int32       `json:"group_id"`
	IsCustomSubject bool        `json:"is_custom_subject"`
	CustomSubjectId int32       `json:"custom_subject_id"`
	SubjectId       int32       `json:"subject_id"`
	ActualDate      pgtype.Date `json:"actual_date"`
	LessonNumber    int32       `json:"lesson_number"`
}

type DeleteRealTimeLessonRequest struct {
	Id int32 `json:"id"`
}

func CreateBaseSchedule(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req CreateBaseScheduleLessonRequest

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

	if req.TeacherId == 0 || req.DayOfWeek == 0 || req.LessonNumber <= 0 || req.RoomId == 0 || req.GroupId == 0 || (req.IsCustomSubject == true && req.CustomSubjectId == 0) || (req.IsCustomSubject == false && req.SubjectId == 0) {
		return c.SendStatus(400)
	}

	has_permission := helpers.CheckPermission(c, pool, rdb, "MANAGE_TIMETABLES", school_id)

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

	params := db_queries.CreateBaseScheduleParams{
		SchoolID:        int32(school_id),
		TeacherID:       req.TeacherId,
		RoomID:          req.RoomId,
		DayOfWeek:       req.DayOfWeek,
		GroupID:         req.GroupId,
		CustomSubject:   req.IsCustomSubject,
		CustomSubjectID: custom_subject_id,
		SubjectID:       subject_id,
	}

	err = queries.CreateBaseSchedule(c.Context(), params)

	if err != nil {
		return c.SendStatus(500)
	}

	return c.SendStatus(200)
}

func DeleteBaseSchedule(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req DeleteBaseScheduleLessonRequest

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

	if req.Id == 0 {
		return c.SendStatus(400)
	}

	has_permission := helpers.CheckPermission(c, pool, rdb, "MANAGE_TIMETABLES", school_id)

	if !has_permission {
		return c.SendStatus(403)
	}

	queries := db_queries.New(pool)

	params := db_queries.DeleteBaseScheduleParams{
		ID:       req.Id,
		SchoolID: int32(school_id),
	}

	err = queries.DeleteBaseSchedule(c.Context(), params)

	if err != nil {
		return c.SendStatus(500)
	}

	return c.SendStatus(200)
}

func UpdateBaseSchedule(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req EditBaseScheduleLessonRequest

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

	if req.Id == 0 || req.TeacherId == 0 || req.DayOfWeek == 0 || req.LessonNumber <= 0 || req.RoomId == 0 || req.GroupId == 0 || (req.IsCustomSubject == true && req.CustomSubjectId == 0) || (req.IsCustomSubject == false && req.SubjectId == 0) {
		return c.SendStatus(400)
	}

	has_permission := helpers.CheckPermission(c, pool, rdb, "MANAGE_TIMETABLES", school_id)

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

	params := db_queries.UpdateBaseScheduleParams{
		SchoolID:        int32(school_id),
		TeacherID:       req.TeacherId,
		RoomID:          req.RoomId,
		DayOfWeek:       req.DayOfWeek,
		GroupID:         req.GroupId,
		CustomSubject:   req.IsCustomSubject,
		CustomSubjectID: custom_subject_id,
		SubjectID:       subject_id,
	}

	err = queries.UpdateBaseSchedule(c.Context(), params)

	if err != nil {
		return c.SendStatus(500)
	}

	return c.SendStatus(200)
}

func ReadBaseScheduleClass(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req ReadBaseScheduleClassRequest

	if err := c.Bind().Query(&req); err != nil {
		return c.SendStatus(400)
	}

	school_id, err := strconv.ParseInt(c.Get("X-School"), 10, 32)

	if err != nil {
		return c.SendStatus(400)
	}

	if school_id == 0 {
		return c.SendStatus(400)
	}

	if req.ClassId == 0 {
		return c.SendStatus(400)
	}

	has_permission := helpers.CheckPermission(c, pool, rdb, "MANAGE_TIMETABLES", school_id)

	if !has_permission {
		return c.SendStatus(403)
	}

	queries := db_queries.New(pool)

	params := db_queries.ReadBaseScheduleClassParams{
		ClassesID: req.ClassId,
		SchoolID:  int32(school_id),
	}

	data, err := queries.ReadBaseScheduleClass(c.Context(), params)

	if err != nil {
		return c.SendStatus(500)
	}

	return c.JSON(data)
}

func ReadBaseScheduleGroup(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req ReadBaseScheduleGroupRequest

	if err := c.Bind().Query(&req); err != nil {
		return c.SendStatus(400)
	}

	school_id, err := strconv.ParseInt(c.Get("X-School"), 10, 32)

	if err != nil {
		return c.SendStatus(400)
	}

	if school_id == 0 || req.GroupId == 0 {
		return c.SendStatus(400)
	}

	has_permission := helpers.CheckPermission(c, pool, rdb, "MANAGE_TIMETABLES", school_id)

	if !has_permission {
		return c.SendStatus(403)
	}

	queries := db_queries.New(pool)

	params := db_queries.ReadBaseScheduleGroupParams{
		SchoolID: int32(school_id),
		GroupID:  req.GroupId,
	}

	data, err := queries.ReadBaseScheduleGroup(c.Context(), params)

	if err != nil {
		return c.SendStatus(500)
	}

	return c.JSON(data)
}

func ReadRealTimeTable(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req ReadRealTimetableRequest

	if err := c.Bind().Query(&req); err != nil {
		return c.SendStatus(400)
	}

	school_id, err := strconv.ParseInt(c.Get("X-School"), 10, 32)

	if err != nil {
		return c.SendStatus(400)
	}

	if school_id == 0 || req.ClassId == 0 || req.Start == "" || req.End == "" {
		return c.SendStatus(400)
	}

	start, err := time.Parse("2006-01-02", req.Start)

	if err != nil {
		return c.SendStatus(400)
	}

	end, err := time.Parse("2006-01-02", req.End)

	if err != nil {
		return c.SendStatus(400)
	}

	if start.Unix() >= end.Unix() {
		return c.SendStatus(400)
	}

	has_permission := helpers.CheckPermission(c, pool, rdb, "MANAGE_TIMETABLES", school_id)

	if !has_permission {
		return c.SendStatus(403)
	}

	queries := db_queries.New(pool)

	params := db_queries.ReadRealTimeTableParams{
		SchoolID:     int32(school_id),
		ClassesID:    req.ClassId,
		ActualDate:   pgtype.Date{Time: start, Valid: true},
		ActualDate_2: pgtype.Date{Time: end, Valid: true},
	}

	data, err := queries.ReadRealTimeTable(c.Context(), params)

	if err != nil {
		return c.SendStatus(500)
	}

	return c.JSON(data)
}

func CreateRealTimeLesson(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req CreateRealTimeLessonRequest

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

	if req.TeacherId == 0 || req.RoomId == 0 || req.DayOfWeek == 0 || req.GroupId == 0 || (req.IsCustomSubject == false && req.SubjectId == 0) || (req.IsCustomSubject == true && req.CustomSubjectId == 0) || req.ActualDate.Time.IsZero() || req.LessonNumber <= 0 {
		return c.SendStatus(400)
	}

	has_permission := helpers.CheckPermission(c, pool, rdb, "MANAGE_TIMETABLES", school_id)

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

	params := db_queries.CreateRealTimeLessonParams{
		SchoolID:        int32(school_id),
		TeacherID:       req.TeacherId,
		RoomID:          req.RoomId,
		DayOfWeek:       req.DayOfWeek,
		GroupID:         req.GroupId,
		CustomSubject:   req.IsCustomSubject,
		CustomSubjectID: custom_subject_id,
		SubjectID:       subject_id,
		ActualDate:      req.ActualDate,
		LessonNum:       req.LessonNumber,
	}

	err = queries.CreateRealTimeLesson(c.Context(), params)

	if err != nil {
		return c.SendStatus(500)
	}

	return c.SendStatus(200)
}

func UpdateRealTimeLesson(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req UpdateRealTimeLessonRequest

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

	if req.Id == 0 || req.TeacherId == 0 || req.RoomId == 0 || req.DayOfWeek == 0 || req.GroupId == 0 || (req.IsCustomSubject == false && req.SubjectId == 0) || (req.IsCustomSubject == true && req.CustomSubjectId == 0) || req.ActualDate.Time.IsZero() || req.LessonNumber <= 0 {
		return c.SendStatus(400)
	}

	has_permission := helpers.CheckPermission(c, pool, rdb, "MANAGE_TIMETABLES", school_id)

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

	params := db_queries.UpdateRealTimeLessonParams{
		SchoolID:        int32(school_id),
		TeacherID:       req.TeacherId,
		RoomID:          req.RoomId,
		DayOfWeek:       req.DayOfWeek,
		GroupID:         req.GroupId,
		CustomSubject:   req.IsCustomSubject,
		CustomSubjectID: custom_subject_id,
		SubjectID:       subject_id,
		ActualDate:      req.ActualDate,
		LessonNum:       req.LessonNumber,
	}

	err = queries.UpdateRealTimeLesson(c.Context(), params)

	if err != nil {
		return c.SendStatus(500)
	}

	return c.SendStatus(200)
}

func DeleteRealTimeLesson(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req DeleteRealTimeLessonRequest

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

	if req.Id == 0 {
		return c.SendStatus(400)
	}

	has_permission := helpers.CheckPermission(c, pool, rdb, "MANAGE_TIMETABLES", school_id)

	if !has_permission {
		return c.SendStatus(403)
	}

	queries := db_queries.New(pool)

	params := db_queries.DeleteRealTimeLessonParams{
		SchoolID: int32(school_id),
		ID:       req.Id,
	}

	err = queries.DeleteRealTimeLesson(c.Context(), params)

	if err != nil {
		return c.SendStatus(500)
	}

	return c.SendStatus(200)
}

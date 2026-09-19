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

type CreateBaseScheduleLessonRequest struct {
	TeacherId       int32 `json:"teacher_id"`
	DayOfWeek       int32 `json:"day_of_week"`
	LessonNumber    int32 `json:"lesson_number"`
	RoomId          int32 `json:"room_id"`
	GroupId         int32 `json:"group_id"`
	IsCustomSubject bool  `json:"is_custom_subject"`
	SubjectId       int32 `json:"subject_id"`
	CustomSubjectId int32 `json:"CustomSubject"`
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
	CustomSubjectId int32 `json:"CustomSubject"`
}

type ReadBaseScheduleClassRequest struct {
	ClassId int32 `json:"class_id"`
}

type ReadBaseScheduleGroupRequest struct {
	GroupId int32 `json:"group_id"`
}

type ReadRealTimetableRequest struct {
	ClassId int32       `json:"class_id"`
	Start   pgtype.Date `json:"Start_date"`
	End     pgtype.Date `json:"End_date"`
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

	school_id, err := strconv.ParseInt(c.Get("X-School"), 10, 32)

	if err != nil {
		return c.SendStatus(500)
	}

	if school_id == 0 {
		return c.SendStatus(500)
	}

	if req.TeacherId == 0 || req.DayOfWeek == 0 || req.LessonNumber < 0 || req.RoomId == 0 || req.GroupId == 0 || (req.IsCustomSubject == true && req.CustomSubjectId == 0) || (req.IsCustomSubject == false && req.SubjectId == 0) {
		return c.SendStatus(400)
	}

	has_permission := helpers.CheckPermission(c, pool, rdb, "MANAGE_TIMETABLES", school_id)

	if !has_permission {
		return c.SendStatus(401)
	}

	queries := db_queries.New(pool)

	params := db_queries.CreateBaseScheduleParams{
		SchoolID:        int32(school_id),
		TeacherID:       req.TeacherId,
		RoomID:          req.RoomId,
		DayOfWeek:       req.DayOfWeek,
		LessonNum:       req.LessonNumber,
		GroupID:         req.GroupId,
		CustomSubject:   req.IsCustomSubject,
		SubjectID:       pgtype.Int4{Int32: req.SubjectId, Valid: true},
		CustomSubjectID: pgtype.Int4{Int32: req.CustomSubjectId, Valid: true},
	}

	err = queries.CreateBaseSchedule(c.Context(), params)

	if err != nil {
		return c.SendStatus(500)
	}

	return c.SendStatus(200)
}

func DeleteBaseSchedule(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req DeleteBaseScheduleLessonRequest

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

	has_permission := helpers.CheckPermission(c, pool, rdb, "MANAGE_TIMETABLES", school_id)

	if !has_permission {
		return c.SendStatus(401)
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

	school_id, err := strconv.ParseInt(c.Get("X-School"), 10, 32)

	if err != nil {
		return c.SendStatus(500)
	}

	if school_id == 0 {
		return c.SendStatus(500)
	}

	if req.Id == 0 || req.TeacherId == 0 || req.DayOfWeek == 0 || req.LessonNumber < 0 || req.RoomId == 0 || req.GroupId == 0 || (req.IsCustomSubject == true && req.CustomSubjectId == 0) || (req.IsCustomSubject == false && req.SubjectId == 0) {
		return c.SendStatus(400)
	}

	has_permission := helpers.CheckPermission(c, pool, rdb, "MANAGE_TIMETABLES", school_id)

	if !has_permission {
		return c.SendStatus(401)
	}

	queries := db_queries.New(pool)

	params := db_queries.UpdateBaseScheduleParams{
		ID:              req.Id,
		SchoolID:        int32(school_id),
		TeacherID:       req.TeacherId,
		DayOfWeek:       req.DayOfWeek,
		LessonNum:       req.LessonNumber,
		RoomID:          req.RoomId,
		GroupID:         req.GroupId,
		CustomSubject:   req.IsCustomSubject,
		CustomSubjectID: pgtype.Int4{Int32: req.CustomSubjectId, Valid: true},
		SubjectID:       pgtype.Int4{Int32: req.SubjectId, Valid: true},
	}

	err = queries.UpdateBaseSchedule(c.Context(), params)

	if err != nil {
		return c.SendStatus(500)
	}

	return c.SendStatus(200)
}

func ReadBaseScheduleClass(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req ReadBaseScheduleClassRequest

	school_id, err := strconv.ParseInt(c.Get("X-School"), 10, 32)

	if err != nil {
		return c.SendStatus(500)
	}

	if school_id == 0 {
		return c.SendStatus(500)
	}

	if req.ClassId == 0 {
		return c.SendStatus(400)
	}

	has_permission := helpers.CheckPermission(c, pool, rdb, "MANAGE_TIMETABLES", school_id)

	if !has_permission {
		return c.SendStatus(401)
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

	school_id, err := strconv.ParseInt(c.Get("X-School"), 10, 32)

	if err != nil {
		return c.SendStatus(500)
	}

	if school_id == 0 {
		return c.SendStatus(500)
	}

	has_permission := helpers.CheckPermission(c, pool, rdb, "MANAGE_TIMETABLES", school_id)

	if !has_permission {
		return c.SendStatus(401)
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

	school_id, err := strconv.ParseInt(c.Get("X-School"), 10, 32)

	if err != nil {
		return c.SendStatus(500)
	}

	if school_id == 0 {
		return c.SendStatus(500)
	}

	has_permission := helpers.CheckPermission(c, pool, rdb, "MANAGE_TIMETABLES", school_id)

	if !has_permission {
		return c.SendStatus(401)
	}

	queries := db_queries.New(pool)

	params := db_queries.ReadRealTimeTableParams{
		SchoolID:     int32(school_id),
		ClassesID:    req.ClassId,
		ActualDate:   req.Start,
		ActualDate_2: req.End,
	}

	data, err := queries.ReadRealTimeTable(c.Context(), params)

	if err != nil {
		return c.SendStatus(500)
	}

	return c.JSON(data)
}

func CreateRealTimeLesson(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req CreateRealTimeLessonRequest

	school_id, err := strconv.ParseInt(c.Get("X-School"), 10, 32)

	if err != nil {
		return c.SendStatus(500)
	}

	if school_id == 0 {
		return c.SendStatus(500)
	}

	if req.TeacherId == 0 || req.RoomId == 0 || req.DayOfWeek == 0 || req.GroupId == 0 || (req.IsCustomSubject == false && req.SubjectId == 0) || (req.IsCustomSubject == true && req.CustomSubjectId == 0) || req.ActualDate.Time.IsZero() || req.LessonNumber == 0 {
		return c.SendStatus(400)
	}

	has_permission := helpers.CheckPermission(c, pool, rdb, "MANAGE_TIMETABLES", school_id)

	if !has_permission {
		return c.SendStatus(401)
	}

	queries := db_queries.New(pool)

	params := db_queries.CreateRealTimeLessonParams{
		SchoolID:        int32(school_id),
		TeacherID:       req.TeacherId,
		RoomID:          req.RoomId,
		DayOfWeek:       req.DayOfWeek,
		GroupID:         req.GroupId,
		CustomSubject:   req.IsCustomSubject,
		CustomSubjectID: pgtype.Int4{Int32: req.CustomSubjectId, Valid: true},
		SubjectID:       pgtype.Int4{Int32: req.SubjectId, Valid: true},
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

	school_id, err := strconv.ParseInt(c.Get("X-School"), 10, 32)

	if err != nil {
		return c.SendStatus(500)
	}

	if school_id == 0 {
		return c.SendStatus(500)
	}

	if req.Id == 0 || req.TeacherId == 0 || req.RoomId == 0 || req.DayOfWeek == 0 || req.GroupId == 0 || (req.IsCustomSubject == false && req.SubjectId == 0) || (req.IsCustomSubject == true && req.CustomSubjectId == 0) || req.ActualDate.Time.IsZero() || req.LessonNumber == 0 {
		return c.SendStatus(400)
	}

	has_permission := helpers.CheckPermission(c, pool, rdb, "MANAGE_TIMETABLES", school_id)

	if !has_permission {
		return c.SendStatus(401)
	}

	queries := db_queries.New(pool)

	params := db_queries.UpdateRealTimeLessonParams{
		ID:              req.Id,
		SchoolID:        int32(school_id),
		TeacherID:       req.TeacherId,
		RoomID:          req.RoomId,
		DayOfWeek:       req.DayOfWeek,
		GroupID:         req.GroupId,
		CustomSubject:   req.IsCustomSubject,
		CustomSubjectID: pgtype.Int4{Int32: req.CustomSubjectId, Valid: true},
		SubjectID:       pgtype.Int4{Int32: req.SubjectId, Valid: true},
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
	var req DeleteLessonTimeRequest

	school_id, err := strconv.ParseInt(c.Get("X-School"), 10, 32)

	if err != nil {
		return c.SendStatus(500)
	}

	if school_id == 0 {
		return c.SendStatus(500)
	}

	if req.Id == 0 {
		return c.SendStatus(401)
	}

	has_permission := helpers.CheckPermission(c, pool, rdb, "MANAGE_TIMETABLES", school_id)

	if !has_permission {
		return c.SendStatus(401)
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

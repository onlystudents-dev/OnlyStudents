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

type BaseScheduleLessonResponse struct {
	ID               int32       `json:"id"`
	SchoolID         int32       `json:"school_id"`
	TeacherID        int32       `json:"teacher_id"`
	TeacherFirstName pgtype.Text `json:"teacher_first_name"`
	TeacherLastName  pgtype.Text `json:"teacher_last_name"`
	RoomID           int32       `json:"room_id"`
	DayOfWeek        int32       `json:"day_of_week"`
	LessonNum        int32       `json:"lesson_num"`
	GroupID          int32       `json:"group_id"`
	CustomSubject    bool        `json:"custom_subject"`
	SubjectID        pgtype.Int4 `json:"subject_id"`
	CustomSubjectID  pgtype.Int4 `json:"custom_subject_id"`
	SubjectName      string      `json:"subject_name"`
	HasExam          bool        `json:"has_exam"`
	HasHomework      bool        `json:"has_homework"`
}

type RealTimeLessonResponse struct {
	ActualDate         pgtype.Date `json:"actual_date"`
	RoomID             int32       `json:"room_id"`
	LessonNum          int32       `json:"lesson_num"`
	DayOfWeek          int32       `json:"day_of_week"`
	EffectiveTeacherID int32       `json:"effective_teacher_id"`
	TeacherFirstName   pgtype.Text `json:"teacher_first_name"`
	TeacherLastName    pgtype.Text `json:"teacher_last_name"`
	GroupID            int32       `json:"group_id"`
	SchoolID           int32       `json:"school_id"`
	CustomSubject      bool        `json:"custom_subject"`
	SubjectID          pgtype.Int4 `json:"subject_id"`
	CustomSubjectID    pgtype.Int4 `json:"custom_subject_id"`
	SubjectName        string      `json:"subject_name"`
	HasExam            bool        `json:"has_exam"`
	HasHomework        bool        `json:"has_homework"`
	IsSubstitution     bool        `json:"is_substitution"`
	Canceled           bool        `json:"canceled"`
}

func CreateBaseSchedule(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req CreateBaseScheduleLessonRequest

	if err := c.Bind().Body(&req); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if req.TeacherId == 0 || req.DayOfWeek == 0 || req.LessonNumber <= 0 || req.RoomId == 0 || req.GroupId == 0 || (req.IsCustomSubject == true && req.CustomSubjectId == 0) || (req.IsCustomSubject == false && req.SubjectId == 0) {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	scope, status_code := helpers.ResolveTeacherCapabilityScope(c, pool, rdb, "MANAGE_TIMETABLES")

	if status_code != fiber.StatusOK {
		return c.SendStatus(status_code)
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
		SchoolID:        scope.SchoolID,
		TeacherID:       req.TeacherId,
		RoomID:          req.RoomId,
		DayOfWeek:       req.DayOfWeek,
		GroupID:         req.GroupId,
		CustomSubject:   req.IsCustomSubject,
		CustomSubjectID: custom_subject_id,
		SubjectID:       subject_id,
	}

	err := queries.CreateBaseSchedule(c.Context(), params)

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusOK)
}

func DeleteBaseSchedule(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req DeleteBaseScheduleLessonRequest

	if err := c.Bind().Body(&req); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if req.Id == 0 {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	scope, status_code := helpers.ResolveTeacherCapabilityScope(c, pool, rdb, "MANAGE_TIMETABLES")

	if status_code != fiber.StatusOK {
		return c.SendStatus(status_code)
	}

	queries := db_queries.New(pool)

	params := db_queries.DeleteBaseScheduleParams{
		ID:       req.Id,
		SchoolID: scope.SchoolID,
	}

	err := queries.DeleteBaseSchedule(c.Context(), params)

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusOK)
}

func UpdateBaseSchedule(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req EditBaseScheduleLessonRequest

	if err := c.Bind().Body(&req); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if req.Id == 0 || req.TeacherId == 0 || req.DayOfWeek == 0 || req.LessonNumber <= 0 || req.RoomId == 0 || req.GroupId == 0 || (req.IsCustomSubject == true && req.CustomSubjectId == 0) || (req.IsCustomSubject == false && req.SubjectId == 0) {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	scope, status_code := helpers.ResolveTeacherCapabilityScope(c, pool, rdb, "MANAGE_TIMETABLES")

	if status_code != fiber.StatusOK {
		return c.SendStatus(status_code)
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
		SchoolID:        scope.SchoolID,
		TeacherID:       req.TeacherId,
		RoomID:          req.RoomId,
		DayOfWeek:       req.DayOfWeek,
		GroupID:         req.GroupId,
		CustomSubject:   req.IsCustomSubject,
		CustomSubjectID: custom_subject_id,
		SubjectID:       subject_id,
	}

	err := queries.UpdateBaseSchedule(c.Context(), params)

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusOK)
}

func ReadBaseScheduleClass(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req ReadBaseScheduleClassRequest

	if err := c.Bind().Query(&req); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if req.ClassId == 0 {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	scope, status_code := helpers.ResolveTeacherCapabilityScope(c, pool, rdb, "MANAGE_TIMETABLES")

	if status_code != fiber.StatusOK {
		return c.SendStatus(status_code)
	}

	queries := db_queries.New(pool)

	params := db_queries.ReadBaseScheduleClassParams{
		ClassesID: req.ClassId,
		SchoolID:  int32(scope.SchoolID),
	}

	data, err := queries.ReadBaseScheduleClass(c.Context(), params)

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	lessons := make([]BaseScheduleLessonResponse, 0, len(data))

	for _, row := range data {
		lessons = append(lessons, BaseScheduleLessonResponse{
			ID:               row.ID,
			SchoolID:         row.SchoolID,
			TeacherID:        row.TeacherID,
			TeacherFirstName: row.TeacherFirstName,
			TeacherLastName:  row.TeacherLastName,
			RoomID:           row.RoomID,
			DayOfWeek:        row.DayOfWeek,
			LessonNum:        row.LessonNum,
			GroupID:          row.GroupID,
			CustomSubject:    row.CustomSubject,
			SubjectID:        row.SubjectID,
			CustomSubjectID:  row.CustomSubjectID,
			SubjectName:      row.SubjectName,
			HasExam:          row.HasExam,
			HasHomework:      row.HasHomework,
		})
	}

	return c.JSON(lessons)
}

func ReadBaseScheduleGroup(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req ReadBaseScheduleGroupRequest

	if err := c.Bind().Query(&req); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	school_id, err := strconv.ParseInt(c.Get("X-School"), 10, 32)

	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if school_id == 0 || req.GroupId == 0 {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	scope, status_code := helpers.ResolveTeacherCapabilityScope(c, pool, rdb, "MANAGE_TIMETABLES")

	if status_code != fiber.StatusOK {
		return c.SendStatus(status_code)
	}

	queries := db_queries.New(pool)

	params := db_queries.ReadBaseScheduleGroupParams{
		SchoolID: scope.SchoolID,
		GroupID:  req.GroupId,
	}

	data, err := queries.ReadBaseScheduleGroup(c.Context(), params)

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	lessons := make([]BaseScheduleLessonResponse, 0, len(data))

	for _, row := range data {
		lessons = append(lessons, BaseScheduleLessonResponse{
			ID:               row.ID,
			SchoolID:         row.SchoolID,
			TeacherID:        row.TeacherID,
			TeacherFirstName: row.TeacherFirstName,
			TeacherLastName:  row.TeacherLastName,
			RoomID:           row.RoomID,
			DayOfWeek:        row.DayOfWeek,
			LessonNum:        row.LessonNum,
			GroupID:          row.GroupID,
			CustomSubject:    row.CustomSubject,
			SubjectID:        row.SubjectID,
			CustomSubjectID:  row.CustomSubjectID,
			SubjectName:      row.SubjectName,
			HasExam:          row.HasExam,
			HasHomework:      row.HasHomework,
		})
	}

	return c.JSON(lessons)
}

func ReadRealTimeTable(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req ReadRealTimetableRequest

	if err := c.Bind().Query(&req); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	school_id, err := strconv.ParseInt(c.Get("X-School"), 10, 32)

	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if school_id == 0 || req.ClassId == 0 || req.Start == "" || req.End == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	start, err := time.Parse("2006-01-02", req.Start)

	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	end, err := time.Parse("2006-01-02", req.End)

	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if start.Unix() >= end.Unix() {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	scope, status_code := helpers.ResolveTeacherCapabilityScope(c, pool, rdb, "MANAGE_TIMETABLES")

	if status_code != fiber.StatusOK {
		return c.SendStatus(status_code)
	}

	queries := db_queries.New(pool)

	params := db_queries.ReadRealTimeTableParams{
		SchoolID:  int32(scope.SchoolID),
		ClassesID: req.ClassId,
		StartDate: pgtype.Date{Time: start, Valid: true},
		EndDate:   pgtype.Date{Time: end, Valid: true},
	}

	data, err := queries.ReadRealTimeTable(c.Context(), params)

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	lessons := make([]RealTimeLessonResponse, 0, len(data))

	for _, row := range data {
		lessons = append(lessons, RealTimeLessonResponse{
			ActualDate:         row.ActualDate,
			RoomID:             row.RoomID,
			LessonNum:          row.LessonNum,
			DayOfWeek:          row.DayOfWeek,
			EffectiveTeacherID: row.EffectiveTeacherID,
			TeacherFirstName:   row.TeacherFirstName,
			TeacherLastName:    row.TeacherLastName,
			GroupID:            row.GroupID,
			SchoolID:           row.SchoolID,
			CustomSubject:      row.CustomSubject,
			SubjectID:          row.SubjectID,
			CustomSubjectID:    row.CustomSubjectID,
			SubjectName:        row.SubjectName,
			HasExam:            row.HasExam,
			HasHomework:        row.HasHomework,
			IsSubstitution:     row.IsSubstitution,
			Canceled:           row.Canceled,
		})
	}

	return c.JSON(lessons)
}

func CreateRealTimeLesson(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req CreateRealTimeLessonRequest

	if err := c.Bind().Body(&req); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if req.TeacherId == 0 || req.RoomId == 0 || req.DayOfWeek == 0 || req.GroupId == 0 || (req.IsCustomSubject == false && req.SubjectId == 0) || (req.IsCustomSubject == true && req.CustomSubjectId == 0) || req.ActualDate.Time.IsZero() || req.LessonNumber <= 0 {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	scope, status_code := helpers.ResolveTeacherCapabilityScope(c, pool, rdb, "MANAGE_TIMETABLES")

	if status_code != fiber.StatusOK {
		return c.SendStatus(status_code)
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
		SchoolID:        scope.SchoolID,
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

	err := queries.CreateRealTimeLesson(c.Context(), params)

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusOK)
}

func UpdateRealTimeLesson(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req UpdateRealTimeLessonRequest

	if err := c.Bind().Body(&req); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if req.Id == 0 || req.TeacherId == 0 || req.RoomId == 0 || req.DayOfWeek == 0 || req.GroupId == 0 || (req.IsCustomSubject == false && req.SubjectId == 0) || (req.IsCustomSubject == true && req.CustomSubjectId == 0) || req.ActualDate.Time.IsZero() || req.LessonNumber <= 0 {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	scope, status_code := helpers.ResolveTeacherCapabilityScope(c, pool, rdb, "MANAGE_TIMETABLES")

	if status_code != fiber.StatusOK {
		return c.SendStatus(status_code)
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
		SchoolID:        scope.SchoolID,
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

	err := queries.UpdateRealTimeLesson(c.Context(), params)

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusOK)
}

func DeleteRealTimeLesson(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req DeleteRealTimeLessonRequest

	if err := c.Bind().Body(&req); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if req.Id == 0 {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	scope, status_code := helpers.ResolveTeacherCapabilityScope(c, pool, rdb, "MANAGE_TIMETABLES")

	if status_code != fiber.StatusOK {
		return c.SendStatus(status_code)
	}

	queries := db_queries.New(pool)

	params := db_queries.DeleteRealTimeLessonParams{
		SchoolID: scope.SchoolID,
		ID:       req.Id,
	}

	err := queries.DeleteRealTimeLesson(c.Context(), params)

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusOK)
}

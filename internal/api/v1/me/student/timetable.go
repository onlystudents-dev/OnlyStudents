package studentapi

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

type TimetableLesson struct {
	ActualDate          int64  `json:"actual_date"`
	RoomID              int32  `json:"room_id"`
	LessonNum           int32  `json:"lesson_num"`
	DayOfWeek           int32  `json:"day_of_week"`
	EffectiveTeacherID  int32  `json:"effective_teacher_id"`
	HasTeacherFirstName bool   `json:"has_teacher_first_name"`
	TeacherFirstName    string `json:"teacher_first_name"`
	HasTeacherLastName  bool   `json:"has_teacher_last_name"`
	TeacherLastName     string `json:"teacher_last_name"`
	GroupID             int32  `json:"group_id"`
	SchoolID            int32  `json:"school_id"`
	CustomSubject       bool   `json:"custom_subject"`
	HasSubjectID        bool   `json:"has_subject_id"`
	SubjectID           int32  `json:"subject_id"`
	HasCustomSubjectID  bool   `json:"has_custom_subject_id"`
	CustomSubjectID     int32  `json:"custom_subject_id"`
	SubjectName         string `json:"subject_name"`
	HasExam             bool   `json:"has_exam"`
	HasHomework         bool   `json:"has_homework"`
	IsSubstitution      bool   `json:"is_substitution"`
	Canceled            bool   `json:"canceled"`
}

type TimetableBaseLesson struct {
	ID                  int32  `json:"id"`
	SchoolID            int32  `json:"school_id"`
	TeacherID           int32  `json:"teacher_id"`
	HasTeacherFirstName bool   `json:"has_teacher_first_name"`
	TeacherFirstName    string `json:"teacher_first_name"`
	HasTeacherLastName  bool   `json:"has_teacher_last_name"`
	TeacherLastName     string `json:"teacher_last_name"`
	RoomID              int32  `json:"room_id"`
	DayOfWeek           int32  `json:"day_of_week"`
	LessonNum           int32  `json:"lesson_num"`
	GroupID             int32  `json:"group_id"`
	CustomSubject       bool   `json:"custom_subject"`
	HasSubjectID        bool   `json:"has_subject_id"`
	SubjectID           int32  `json:"subject_id"`
	HasCustomSubjectID  bool   `json:"has_custom_subject_id"`
	CustomSubjectID     int32  `json:"custom_subject_id"`
	SubjectName         string `json:"subject_name"`
	HasExam             bool   `json:"has_exam"`
	HasHomework         bool   `json:"has_homework"`
}

type TimetableBellSchedule struct {
	ID              int32 `json:"id"`
	SchoolID        int32 `json:"school_id"`
	TypeID          int32 `json:"type_id"`
	HasLessonNumber bool  `json:"has_lesson_number"`
	LessonNumber    int32 `json:"lesson_number"`
	AtStart         int64 `json:"at_start"`
	AtEnd           int64 `json:"at_end"`
}

type TimetableRoom struct {
	ID       int32  `json:"id"`
	Name     string `json:"name"`
	Capacity int32  `json:"capacity"`
}

type TimetableCustomSubject struct {
	ID          int32  `json:"id"`
	SubjectName string `json:"subject_name"`
}

type TimetableBellScheduleType struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
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
	scope, err := helpers.ResolveMeScope(c, pool, rdb)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "UNAUTHORIZED"})
	}

	start, end, err := parseDateRange(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	queries := db_queries.New(pool)

	data, err := queries.ReadRealTimeTable(c.Context(), db_queries.ReadRealTimeTableParams{
		SchoolID:  scope.SchoolID,
		ClassesID: scope.ClassID,
		StartDate: pgtype.Date{Time: start, Valid: true},
		EndDate:   pgtype.Date{Time: end, Valid: true},
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "SERVER_ERROR"})
	}

	lessons := make([]TimetableLesson, 0, len(data))

	for _, row := range data {
		lessons = append(lessons, TimetableLesson{
			ActualDate:          row.ActualDate.Time.Unix(),
			RoomID:              row.RoomID,
			LessonNum:           row.LessonNum,
			DayOfWeek:           row.DayOfWeek,
			EffectiveTeacherID:  row.EffectiveTeacherID,
			HasTeacherFirstName: row.TeacherFirstName.Valid,
			HasTeacherLastName:  row.TeacherLastName.Valid,
			TeacherFirstName:    row.TeacherFirstName.String,
			TeacherLastName:     row.TeacherLastName.String,
			GroupID:             row.GroupID,
			SchoolID:            row.SchoolID,
			CustomSubject:       row.CustomSubject,
			HasSubjectID:        row.SubjectID.Valid,
			SubjectID:           row.SubjectID.Int32,
			HasCustomSubjectID:  row.CustomSubjectID.Valid,
			CustomSubjectID:     row.CustomSubjectID.Int32,
			SubjectName:         row.SubjectName,
			HasExam:             row.HasExam,
			HasHomework:         row.HasHomework,
			IsSubstitution:      row.IsSubstitution,
			Canceled:            row.Canceled,
		})
	}

	return c.JSON(lessons)
}

func ReadBaseSchedule(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	scope, err := helpers.ResolveMeScope(c, pool, rdb)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "UNAUTHORIZED"})
	}

	queries := db_queries.New(pool)

	data, err := queries.ReadBaseScheduleClass(c.Context(), db_queries.ReadBaseScheduleClassParams{
		SchoolID: scope.SchoolID,
		ClassID:  scope.ClassID,
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "SERVER_ERROR"})
	}

	lessons := make([]TimetableBaseLesson, 0, len(data))

	for _, row := range data {
		lessons = append(lessons, TimetableBaseLesson{
			ID:                  row.ID,
			SchoolID:            row.SchoolID,
			TeacherID:           row.TeacherID,
			HasTeacherFirstName: row.TeacherFirstName.Valid,
			HasTeacherLastName:  row.TeacherLastName.Valid,
			TeacherFirstName:    row.TeacherFirstName.String,
			TeacherLastName:     row.TeacherLastName.String,
			RoomID:              row.RoomID,
			DayOfWeek:           row.DayOfWeek,
			LessonNum:           row.LessonNum,
			GroupID:             row.GroupID,
			CustomSubject:       row.CustomSubject,
			HasSubjectID:        row.SubjectID.Valid,
			SubjectID:           row.SubjectID.Int32,
			HasCustomSubjectID:  row.CustomSubjectID.Valid,
			CustomSubjectID:     row.CustomSubjectID.Int32,
			SubjectName:         row.SubjectName,
			HasExam:             row.HasExam,
			HasHomework:         row.HasHomework,
		})
	}

	return c.JSON(lessons)
}

func ReadLessonTime(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req ReadLessonTimeRequest

	if err := c.Bind().Query(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	scope, err := helpers.ResolveMeScope(c, pool, rdb)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "UNAUTHORIZED"})
	}

	queries := db_queries.New(pool)

	data, err := queries.ReadLessonTime(c.Context(), db_queries.ReadLessonTimeParams{
		SchoolID: scope.SchoolID,
		TypeID:   req.TypeID,
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "SERVER_ERROR"})
	}

	lesson_times := make([]TimetableBellSchedule, 0, len(data))

	for _, row := range data {
		lesson_times = append(lesson_times, TimetableBellSchedule{
			ID:              row.ID,
			SchoolID:        row.SchoolID,
			TypeID:          row.TypeID,
			HasLessonNumber: row.LessonNumber.Valid,
			LessonNumber:    row.LessonNumber.Int32,
			AtStart:         row.AtStart.Microseconds / 1000000,
			AtEnd:           row.AtEnd.Microseconds / 1000000,
		})
	}

	return c.JSON(lesson_times)
}

func ReadRoom(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	scope, err := helpers.ResolveMeScope(c, pool, rdb)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "UNAUTHORIZED"})
	}

	queries := db_queries.New(pool)

	data, err := queries.ReadRoom(c.Context(), scope.SchoolID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "SERVER_ERROR"})
	}

	rooms := make([]TimetableRoom, 0, len(data))

	for _, row := range data {
		rooms = append(rooms, TimetableRoom{
			ID:       row.ID,
			Name:     row.Name,
			Capacity: row.Capacity,
		})
	}

	return c.JSON(rooms)
}

func ReadCustomSubject(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	scope, err := helpers.ResolveMeScope(c, pool, rdb)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "UNAUTHORIZED"})
	}

	queries := db_queries.New(pool)

	data, err := queries.ReadCustomSubject(c.Context(), scope.SchoolID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "SERVER_ERROR"})
	}

	custom_subjects := make([]TimetableCustomSubject, 0, len(data))

	for _, row := range data {
		custom_subjects = append(custom_subjects, TimetableCustomSubject{
			ID:          row.ID,
			SubjectName: row.SubjectName,
		})
	}

	return c.JSON(custom_subjects)
}

func ReadBellScheduleType(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	scope, err := helpers.ResolveMeScope(c, pool, rdb)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "UNAUTHORIZED"})
	}

	queries := db_queries.New(pool)

	data, err := queries.ReadBellScheduleType(c.Context(), scope.SchoolID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "SERVER_ERROR"})
	}

	bell_schedule_types := make([]TimetableBellScheduleType, 0, len(data))

	for _, row := range data {
		bell_schedule_types = append(bell_schedule_types, TimetableBellScheduleType{
			ID:   row.ID,
			Name: row.Name,
		})
	}

	return c.JSON(bell_schedule_types)
}

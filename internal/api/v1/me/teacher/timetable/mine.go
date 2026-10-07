package timetable

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

type TimetableRoom struct {
	ID       int32  `json:"id"`
	Name     string `json:"name"`
	Capacity int32  `json:"capacity"`
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
	scope, status_code := helpers.ResolveTeacherScope(c, pool, rdb)

	if status_code != fiber.StatusOK {
		return helpers.ErrorByStatusCode(c, status_code)
	}

	start, end, err := parseDateRange(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	queries := db_queries.New(pool)

	data, err := queries.ReadMyRealTimeTable(c.Context(), db_queries.ReadMyRealTimeTableParams{
		SchoolID:  scope.SchoolID,
		TeacherID: scope.TeacherID,
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

func ReadMyLessonTimes(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	scope, status_code := helpers.ResolveTeacherScope(c, pool, rdb)
	if status_code != fiber.StatusOK {
		return helpers.ErrorByStatusCode(c, status_code)
	}

	queries := db_queries.New(pool)

	data, err := queries.ReadTeacherLessonTimes(c.Context(), db_queries.ReadTeacherLessonTimesParams{
		SchoolID:  scope.SchoolID,
		TeacherID: scope.TeacherID,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "SERVER_ERROR"})
	}

	lessonTimes := make([]BellScheduleSummary, 0, len(data))
	for _, row := range data {
		lessonTimes = append(lessonTimes, convertBellSchedule(row))
	}

	return c.JSON(lessonTimes)
}

func ReadMyRooms(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	scope, status_code := helpers.ResolveTeacherScope(c, pool, rdb)

	if status_code != fiber.StatusOK {
		return helpers.ErrorByStatusCode(c, status_code)
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

package timetable

import (
	"context"
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
	HasExamID           bool   `json:"has_exam_id"`
	ExamID              int64  `json:"exam_id"`
	HasHomeworkID       bool   `json:"has_homework_id"`
	HomeworkID          int64  `json:"homework_id"`
	IsSubstitution      bool   `json:"is_substitution"`
	Canceled            bool   `json:"canceled"`
}

type TimetableRoom struct {
	ID       int32  `json:"id"`
	Name     string `json:"name"`
	Capacity int32  `json:"capacity"`
}

func convertTimeTableLesson(row db_queries.ReadMyRealTimeTableRow) TimetableLesson {
	return TimetableLesson{
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
		HasExamID:           row.ExamID > 0,
		ExamID:              row.ExamID,
		HasHomeworkID:       row.HomeworkID > 0,
		HomeworkID:          row.HomeworkID,
		IsSubstitution:      row.IsSubstitution,
		Canceled:            row.Canceled,
	}
}

func convertStudentsInMyGroup(row db_queries.ReadStudentsInMyGroupRow) StudentInGroupSummary {
	return StudentInGroupSummary{
		ID:        row.ID,
		FirstName: row.FirstName,
		LastName:  row.LastName,
	}
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

// uncached because changes a lot
func ReadMyRealTimeTable(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return helpers.TeacherSummary(c, pool, rdb, "REAL_TIMETABLE_CACHE_TTL", "", nil, func(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, req bool, scope helpers.TeacherScope, ttl int32) ([]db_queries.ReadMyRealTimeTableRow, error) {
		start, end, err := parseDateRange(c)

		if err != nil {
			return []db_queries.ReadMyRealTimeTableRow{}, err
		}

		return queries.ReadMyRealTimeTable(c.Context(), db_queries.ReadMyRealTimeTableParams{
			SchoolID:  scope.SchoolID,
			TeacherID: scope.TeacherID,
			StartDate: pgtype.Date{Time: start, Valid: true},
			EndDate:   pgtype.Date{Time: end, Valid: true},
		})
	}, convertTimeTableLesson)
}

func ReadMyLessonTimes(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return helpers.TeacherSummary(c, pool, rdb, "LESSON_TIME_CACHE_TTL", "", nil, func(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, req bool, scope helpers.TeacherScope, ttl int32) ([]db_queries.BellSchedule, error) {
		return helpers.CacheOrGetTeacherMyLessonTimes(c.Context(), rdb, queries, scope, ttl)
	}, convertBellSchedule)
}

func ReadMyGroups(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return helpers.TeacherSummary(c, pool, rdb, "GROUPS_CACHE_TTL", "", nil, func(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, req bool, scope helpers.TeacherScope, ttl int32) ([]db_queries.Group, error) {
		return helpers.CacheOrGetTeacherMyGroups(ctx, rdb, queries, scope, ttl)
	}, convertGroup)
}

func ReadStudentsInMyGroup(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return helpers.TeacherSummary(c, pool, rdb, "GROUP_STUDENTS_CACHE_TTL", "", func(req ReadStudentGroupRequest) bool { return req.GroupID <= 0 }, func(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, req ReadStudentGroupRequest, scope helpers.TeacherScope, ttl int32) ([]db_queries.ReadStudentsInMyGroupRow, error) {
		return helpers.CacheOrGetStudentsInMyGroup(ctx, rdb, queries, req.GroupID, scope, ttl)
	}, convertStudentsInMyGroup)
}

package studentapi

import (
	"context"
	"errors"
	"fmt"
	db_queries "onlystudents/internal/db/store"
	"onlystudents/internal/helpers"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type ReadLessonTimeRequest struct {
	TypeID int32 `query:"type_id"`
}

type ReadRealTimetableRequest struct {
	StartDate string `query:"start_date"`
	EndDate   string `query:"end_date"`
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
	HasExamID           bool   `json:"has_exam_id"`
	ExamID              int64  `json:"exam_id"`
	HasHomeworkID       bool   `json:"has_homework_id"`
	HomeworkID          int64  `json:"homework_id"`
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

func convertRealTimeTable(row db_queries.ReadRealTimeTableRow) TimetableLesson {
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

func convertBaseLesson(row db_queries.ReadBaseScheduleClassRow) TimetableBaseLesson {
	return TimetableBaseLesson{
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
	}
}

func convertLessonTime(row db_queries.BellSchedule) TimetableBellSchedule {
	return TimetableBellSchedule{
		ID:              row.ID,
		SchoolID:        row.SchoolID,
		TypeID:          row.TypeID,
		HasLessonNumber: row.LessonNumber.Valid,
		LessonNumber:    row.LessonNumber.Int32,
		AtStart:         row.AtStart.Microseconds / 1000000,
		AtEnd:           row.AtEnd.Microseconds / 1000000,
	}
}

func convertRoom(row db_queries.ReadRoomRow) TimetableRoom {
	return TimetableRoom{
		ID:       row.ID,
		Name:     row.Name,
		Capacity: row.Capacity,
	}
}

func convertCustomSubject(row db_queries.ReadCustomSubjectRow) TimetableCustomSubject {
	return TimetableCustomSubject{
		ID:          row.ID,
		SubjectName: row.SubjectName,
	}
}

func convertBellScheduleType(row db_queries.ReadBellScheduleTypeRow) TimetableBellScheduleType {
	return TimetableBellScheduleType{
		ID:   row.ID,
		Name: row.Name,
	}
}

func ReadMyRealTimeTable(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return helpers.StudentSummary(c, pool, rdb, "REAL_TIMETABLE_CACHE_TTL", false,
		func(req ReadRealTimetableRequest) bool {
			_, start_date_err := time.Parse(time.DateOnly, req.StartDate)
			if start_date_err != nil {
				return true
			}

			_, end_date_err := time.Parse(time.DateOnly, req.EndDate)
			if end_date_err != nil {
				return true
			}

			return false
		},
		func(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, req ReadRealTimetableRequest, scope helpers.MeScope, ttl int32) ([]db_queries.ReadRealTimeTableRow, error) {
			start, end, err := parseDateRange(c)
			if err != nil {
				return []db_queries.ReadRealTimeTableRow{}, err
			}

			key := fmt.Sprintf("real_timetable:v2:%d:%d:%d:%d", scope.SchoolID, scope.StudentID, start.Unix(), end.Unix())
			return helpers.CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.ReadRealTimeTableRow, error) {
				return queries.ReadRealTimeTable(ctx, db_queries.ReadRealTimeTableParams{
					SchoolID:  scope.SchoolID,
					ClassID:   scope.ClassID,
					StartDate: pgtype.Date{Time: start, Valid: true},
					EndDate:   pgtype.Date{Time: end, Valid: true},
				})
			})
		}, convertRealTimeTable)
}

func ReadBaseSchedule(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return helpers.StudentSummary(c, pool, rdb, "BASE_SCHEDULE_CACHE_TTL", false, nil, func(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, req bool, scope helpers.MeScope, ttl int32) ([]db_queries.ReadBaseScheduleClassRow, error) {
		return helpers.CacheOrGetStudentBaseScheduleClass(ctx, rdb, queries, scope, ttl)
	}, convertBaseLesson)
}

func ReadLessonTime(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return helpers.StudentSummary(c, pool, rdb, "LESSON_TIME_CACHE_TTL", false, func(req ReadLessonTimeRequest) bool {
		return req.TypeID <= 0
	}, func(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, req ReadLessonTimeRequest, scope helpers.MeScope, ttl int32) ([]db_queries.BellSchedule, error) {
		key := fmt.Sprintf("lesson_time:%d:%d", scope.SchoolID, req.TypeID)
		return helpers.CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.BellSchedule, error) {
			return queries.ReadLessonTime(ctx, db_queries.ReadLessonTimeParams{
				SchoolID: scope.SchoolID,
				TypeID:   req.TypeID,
			})
		})
	}, convertLessonTime)
}

func ReadRoom(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return helpers.StudentSummary(c, pool, rdb, "ROOM_CACHE_TTL", false, nil, func(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, req bool, scope helpers.MeScope, ttl int32) ([]db_queries.ReadRoomRow, error) {
		return helpers.CacheOrGetStudentRooms(ctx, rdb, queries, scope, ttl)
	}, convertRoom)
}

func ReadCustomSubject(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return helpers.StudentSummary(c, pool, rdb, "CUSTOM_SUBJECT_CACHE_TTL", false, nil, func(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, req bool, scope helpers.MeScope, ttl int32) ([]db_queries.ReadCustomSubjectRow, error) {
		return helpers.CacheOrGetStudentCustomSubjects(ctx, rdb, queries, scope, ttl)
	}, convertCustomSubject)
}

func ReadBellScheduleType(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return helpers.StudentSummary(c, pool, rdb, "BELL_SCHEDULE_TYPE_CACHE_TTL", false, nil, func(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, req bool, scope helpers.MeScope, ttl int32) ([]db_queries.ReadBellScheduleTypeRow, error) {
		return helpers.CacheOrGetStudentBellScheduleTypes(ctx, rdb, queries, scope, ttl)
	}, convertBellScheduleType)
}

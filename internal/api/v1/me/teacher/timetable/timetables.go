package timetable

import (
	"context"
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
	TeacherId       int32 `json:"teacher_id"`
	RoomId          int32 `json:"room_id"`
	DayOfWeek       int32 `json:"day_of_week"`
	GroupId         int32 `json:"group_id"`
	IsCustomSubject bool  `json:"is_custom_subject"`
	CustomSubjectId int32 `json:"custom_subject_id"`
	SubjectId       int32 `json:"subject_id"`
	ActualDate      int64 `json:"actual_date"`
	LessonNumber    int32 `json:"lesson_number"`
}

type UpdateRealTimeLessonRequest struct {
	Id              int32 `json:"id"`
	TeacherId       int32 `json:"teacher_id"`
	RoomId          int32 `json:"room_id"`
	DayOfWeek       int32 `json:"day_of_week"`
	GroupId         int32 `json:"group_id"`
	IsCustomSubject bool  `json:"is_custom_subject"`
	CustomSubjectId int32 `json:"custom_subject_id"`
	SubjectId       int32 `json:"subject_id"`
	ActualDate      int64 `json:"actual_date"`
	LessonNumber    int32 `json:"lesson_number"`
}

type DeleteRealTimeLessonRequest struct {
	Id int32 `json:"id"`
}

type BaseScheduleLessonResponse struct {
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

type RealTimeLessonResponse struct {
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

func convertBaseScheduleClass(row db_queries.ReadBaseScheduleClassRow) BaseScheduleLessonResponse {
	return BaseScheduleLessonResponse{
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
	}
}

func convertBaseScheduleGroup(row db_queries.ReadBaseScheduleGroupRow) BaseScheduleLessonResponse {
	return BaseScheduleLessonResponse{
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
	}
}

func CreateBaseSchedule(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return TeacherTimeTableModify(c, pool, rdb, "MANAGE_TIMETABLES",
		func(req CreateBaseScheduleLessonRequest) bool {
			return req.TeacherId <= 0 || req.DayOfWeek <= 0 || req.LessonNumber <= 0 || req.RoomId <= 0 || req.GroupId <= 0 || (req.IsCustomSubject == true && req.CustomSubjectId <= 0) || (req.IsCustomSubject == false && req.SubjectId <= 0)
		},
		func(ctx context.Context, teacher_scope helpers.TeacherScope, queries *db_queries.Queries, req CreateBaseScheduleLessonRequest) (int64, error) {
			var custom_subject_id pgtype.Int4
			var subject_id pgtype.Int4

			if req.IsCustomSubject {
				custom_subject_id = pgtype.Int4{Int32: req.CustomSubjectId, Valid: true}
				subject_id = pgtype.Int4{Valid: false}
			} else {
				custom_subject_id = pgtype.Int4{Valid: false}
				subject_id = pgtype.Int4{Int32: req.SubjectId, Valid: true}
			}

			return queries.CreateBaseSchedule(c.Context(), db_queries.CreateBaseScheduleParams{
				SchoolID:        teacher_scope.SchoolID,
				TeacherID:       req.TeacherId,
				RoomID:          req.RoomId,
				DayOfWeek:       req.DayOfWeek,
				GroupID:         req.GroupId,
				CustomSubject:   req.IsCustomSubject,
				CustomSubjectID: custom_subject_id,
				SubjectID:       subject_id,
			})
		})
}

func DeleteBaseSchedule(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return TeacherTimeTableModify(c, pool, rdb, "MANAGE_TIMETABLES",
		func(req DeleteBaseScheduleLessonRequest) bool {
			return req.Id <= 0
		},
		func(ctx context.Context, teacher_scope helpers.TeacherScope, queries *db_queries.Queries, req DeleteBaseScheduleLessonRequest) (int64, error) {
			return queries.DeleteBaseSchedule(c.Context(), db_queries.DeleteBaseScheduleParams{
				ID:       req.Id,
				SchoolID: teacher_scope.SchoolID,
			})
		})
}

func UpdateBaseSchedule(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return TeacherTimeTableModify(c, pool, rdb, "MANAGE_TIMETABLES",
		func(req EditBaseScheduleLessonRequest) bool {
			return req.Id <= 0 || req.TeacherId <= 0 || req.DayOfWeek <= 0 || req.LessonNumber <= 0 || req.RoomId <= 0 || req.GroupId <= 0 || (req.IsCustomSubject == true && req.CustomSubjectId <= 0) || (req.IsCustomSubject == false && req.SubjectId <= 0)
		},
		func(ctx context.Context, teacher_scope helpers.TeacherScope, queries *db_queries.Queries, req EditBaseScheduleLessonRequest) (int64, error) {
			var custom_subject_id pgtype.Int4
			var subject_id pgtype.Int4

			if req.IsCustomSubject {
				custom_subject_id = pgtype.Int4{Int32: req.CustomSubjectId, Valid: true}
				subject_id = pgtype.Int4{Valid: false}
			} else {
				custom_subject_id = pgtype.Int4{Valid: false}
				subject_id = pgtype.Int4{Int32: req.SubjectId, Valid: true}
			}

			return queries.UpdateBaseSchedule(c.Context(), db_queries.UpdateBaseScheduleParams{
				SchoolID:        teacher_scope.SchoolID,
				TeacherID:       req.TeacherId,
				RoomID:          req.RoomId,
				DayOfWeek:       req.DayOfWeek,
				GroupID:         req.GroupId,
				CustomSubject:   req.IsCustomSubject,
				CustomSubjectID: custom_subject_id,
				ID:              req.Id,
				SubjectID:       subject_id,
			})
		})
}

func ReadBaseScheduleClass(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return TeacherTimeTableSummaryByID(c, pool, rdb, "MANAGE_TIMETABLES", "BASE_SCHEDULE_CLASS_CACHE_TTL",
		func(req ReadBaseScheduleClassRequest) bool { return req.ClassId <= 0 },
		func(req ReadBaseScheduleClassRequest) int32 { return req.ClassId },
		helpers.CacheOrGetBaseScheduleClass, convertBaseScheduleClass)
}

func ReadBaseScheduleGroup(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return TeacherTimeTableSummaryByID(c, pool, rdb, "MANAGE_TIMETABLES", "BASE_SCHEDULE_GROUP_CACHE_TTL",
		func(req ReadBaseScheduleGroupRequest) bool { return req.GroupId <= 0 },
		func(req ReadBaseScheduleGroupRequest) int32 { return req.GroupId },
		helpers.CacheOrGetBaseScheduleGroup, convertBaseScheduleGroup)
}

// TODO: figure this out
func ReadRealTimeTable(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req ReadRealTimetableRequest

	if err := c.Bind().Query(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	school_id, err := strconv.ParseInt(c.Get("X-School"), 10, 32)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	if school_id <= 0 || req.ClassId <= 0 || req.Start == "" || req.End == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	start, err := time.Parse("2006-01-02", req.Start)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	end, err := time.Parse("2006-01-02", req.End)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	if start.Unix() >= end.Unix() {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	scope, status_code := helpers.ResolveTeacherCapabilityScope(c, pool, rdb, "MANAGE_TIMETABLES")

	if status_code != fiber.StatusOK {
		return helpers.ErrorByStatusCode(c, status_code)
	}

	queries := db_queries.New(pool)

	params := db_queries.ReadRealTimeTableParams{
		SchoolID:  int32(scope.SchoolID),
		ClassID:   req.ClassId,
		StartDate: pgtype.Date{Time: start, Valid: true},
		EndDate:   pgtype.Date{Time: end, Valid: true},
	}

	data, err := queries.ReadRealTimeTable(c.Context(), params)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "SERVER_ERROR"})
	}

	lessons := make([]RealTimeLessonResponse, 0, len(data))

	for _, row := range data {
		lessons = append(lessons, RealTimeLessonResponse{
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

func CreateRealTimeLesson(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return TeacherTimeTableModify(c, pool, rdb, "MANAGE_TIMETABLES",
		func(req CreateRealTimeLessonRequest) bool {
			return req.TeacherId <= 0 || req.RoomId <= 0 || req.DayOfWeek <= 0 || req.GroupId <= 0 || (req.IsCustomSubject == false && req.SubjectId <= 0) || (req.IsCustomSubject == true && req.CustomSubjectId <= 0) || req.ActualDate <= 0 || req.LessonNumber <= 0
		},
		func(ctx context.Context, teacher_scope helpers.TeacherScope, queries *db_queries.Queries, req CreateRealTimeLessonRequest) (int64, error) {
			var custom_subject_id pgtype.Int4
			var subject_id pgtype.Int4

			if req.IsCustomSubject {
				custom_subject_id = pgtype.Int4{Int32: req.CustomSubjectId, Valid: true}
				subject_id = pgtype.Int4{Valid: false}
			} else {
				custom_subject_id = pgtype.Int4{Valid: false}
				subject_id = pgtype.Int4{Int32: req.SubjectId, Valid: true}
			}

			return queries.CreateRealTimeLesson(c.Context(), db_queries.CreateRealTimeLessonParams{
				SchoolID:        teacher_scope.SchoolID,
				TeacherID:       req.TeacherId,
				RoomID:          req.RoomId,
				DayOfWeek:       req.DayOfWeek,
				GroupID:         req.GroupId,
				CustomSubject:   req.IsCustomSubject,
				CustomSubjectID: custom_subject_id,
				SubjectID:       subject_id,
				ActualDate:      pgtype.Date{Time: time.Unix(req.ActualDate, 0), Valid: true},
				LessonNum:       req.LessonNumber,
			})
		})
}

func UpdateRealTimeLesson(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return TeacherTimeTableModify(c, pool, rdb, "MANAGE_TIMETABLES",
		func(req UpdateRealTimeLessonRequest) bool {
			return req.Id <= 0 || req.TeacherId <= 0 || req.RoomId <= 0 || req.DayOfWeek <= 0 || req.GroupId <= 0 || (req.IsCustomSubject == false && req.SubjectId <= 0) || (req.IsCustomSubject == true && req.CustomSubjectId <= 0) || req.ActualDate <= 0 || req.LessonNumber <= 0
		},
		func(ctx context.Context, teacher_scope helpers.TeacherScope, queries *db_queries.Queries, req UpdateRealTimeLessonRequest) (int64, error) {
			var custom_subject_id pgtype.Int4
			var subject_id pgtype.Int4

			if req.IsCustomSubject {
				custom_subject_id = pgtype.Int4{Int32: req.CustomSubjectId, Valid: true}
				subject_id = pgtype.Int4{Valid: false}
			} else {
				custom_subject_id = pgtype.Int4{Valid: false}
				subject_id = pgtype.Int4{Int32: req.SubjectId, Valid: true}
			}

			return queries.UpdateRealTimeLesson(c.Context(), db_queries.UpdateRealTimeLessonParams{
				ID:              req.Id,
				SchoolID:        teacher_scope.SchoolID,
				TeacherID:       req.TeacherId,
				RoomID:          req.RoomId,
				DayOfWeek:       req.DayOfWeek,
				GroupID:         req.GroupId,
				CustomSubject:   req.IsCustomSubject,
				CustomSubjectID: custom_subject_id,
				SubjectID:       subject_id,
				ActualDate:      pgtype.Date{Time: time.Unix(req.ActualDate, 0), Valid: true},
				LessonNum:       req.LessonNumber,
			})
		})
}

func DeleteRealTimeLesson(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return TeacherTimeTableModify(c, pool, rdb, "MANAGE_TIMETABLES",
		func(req DeleteRealTimeLessonRequest) bool {
			return req.Id <= 0
		},
		func(ctx context.Context, teacher_scope helpers.TeacherScope, queries *db_queries.Queries, req DeleteRealTimeLessonRequest) (int64, error) {
			return queries.DeleteRealTimeLesson(c.Context(), db_queries.DeleteRealTimeLessonParams{
				SchoolID: teacher_scope.SchoolID,
				ID:       req.Id,
			})
		})
}

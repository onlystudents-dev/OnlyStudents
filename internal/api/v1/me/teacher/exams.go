package teacherapi

import (
	"context"
	db_queries "onlystudents/internal/db/store"
	"onlystudents/internal/helpers"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type AddExamRequest struct {
	Title           string `json:"title"`
	HasDescription  bool   `json:"has_description"`
	Description     string `json:"description"`
	Date            int64  `json:"date"`
	StartTime       int64  `json:"start_time"`
	EndTime         int64  `json:"end_time"`
	HasRoomID       bool   `json:"has_room_id"`
	RoomID          int32  `json:"room_id"`
	ClassSubjectsID int32  `json:"class_subjects_id"`
}

type EditExamRequest struct {
	ExamID          int64  `json:"exam_id"`
	Title           string `json:"title"`
	HasDescription  bool   `json:"has_description"`
	Description     string `json:"description"`
	Date            int64  `json:"date"`
	StartTime       int64  `json:"start_time"`
	EndTime         int64  `json:"end_time"`
	HasRoomID       bool   `json:"has_room_id"`
	RoomID          int32  `json:"room_id"`
	ClassSubjectsID int32  `json:"class_subjects_id"`
}

type RemoveExamRequest struct {
	ExamID int64 `json:"exam_id"`
}

type ExamSummary struct {
	ID             int64  `json:"id"`
	Subject        string `json:"subject"`
	HasSubjectCode bool   `json:"has_subject_code"`
	SubjectCode    string `json:"subject_code"`
	Title          string `json:"title"`
	HasDescription bool   `json:"has_description"`
	Description    string `json:"description"`
	Date           int64  `json:"date"`
	StartTime      int64  `json:"start_time"`
	EndTime        int64  `json:"end_time"`
	HasRoom        bool   `json:"has_room"`
	Room           string `json:"room"`
}

func convertExam(row db_queries.GetTeacherExamsRow) ExamSummary {
	return ExamSummary{
		ID:             row.ID,
		Subject:        row.Subject,
		HasSubjectCode: row.SubjectCode.Valid,
		SubjectCode:    row.SubjectCode.String,
		Title:          row.Title,
		HasDescription: row.Description.Valid,
		Description:    row.Description.String,
		Date:           row.Date.Time.Unix(),
		StartTime:      row.StartTime.Microseconds / 1000000,
		EndTime:        row.EndTime.Microseconds / 1000000,
		HasRoom:        row.Room.Valid,
		Room:           row.Room.String,
	}
}

func ListExams(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return TeacherSummary(c, pool, rdb, "EXAMS_CACHE_TTL", helpers.CacheOrGetTeacherExams, convertExam)
}

func AddExam(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return TeacherModify[AddExamRequest](c, pool, rdb, true,
		func(req AddExamRequest) int32 { return req.ClassSubjectsID },
		func(req AddExamRequest) bool {
			return req.Title == "" || (req.HasDescription && req.Description == "") || req.Date == 0 || req.StartTime == 0 || req.EndTime == 0 || (req.HasRoomID && req.RoomID == 0) || req.ClassSubjectsID == 0
		},
		func(ctx context.Context, teacher_scope helpers.TeacherScope, queries *db_queries.Queries, req AddExamRequest) (int64, error) {
			return queries.TeacherAddExam(ctx, db_queries.TeacherAddExamParams{
				TeacherID:       teacher_scope.TeacherID,
				SchoolID:        teacher_scope.SchoolID,
				Date:            pgtype.Date{Time: time.Unix(req.Date, 0), Valid: true},
				Title:           req.Title,
				Description:     pgtype.Text{String: req.Description, Valid: req.HasDescription},
				StartTime:       pgtype.Time{Microseconds: req.StartTime * 1000000, Valid: true},
				EndTime:         pgtype.Time{Microseconds: req.EndTime * 1000000, Valid: true},
				RoomID:          pgtype.Int4{Int32: req.RoomID, Valid: req.HasRoomID},
				ClassSubjectsID: req.ClassSubjectsID,
			})
		})
}

func RemoveExam(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return TeacherModify[RemoveExamRequest](c, pool, rdb, false,
		func(req RemoveExamRequest) int32 { return 0 },
		func(req RemoveExamRequest) bool { return req.ExamID == 0 },
		func(ctx context.Context, teacher_scope helpers.TeacherScope, queries *db_queries.Queries, req RemoveExamRequest) (int64, error) {
			return queries.TeacherDeleteExam(ctx, db_queries.TeacherDeleteExamParams{
				ID:        req.ExamID,
				TeacherID: teacher_scope.TeacherID,
				SchoolID:  teacher_scope.SchoolID,
			})
		})
}

func EditExam(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return TeacherModify[EditExamRequest](c, pool, rdb, true,
		func(req EditExamRequest) int32 { return req.ClassSubjectsID },
		func(req EditExamRequest) bool {
			return req.ExamID == 0 || req.Title == "" || (req.HasDescription && req.Description == "") || req.Date == 0 || req.StartTime == 0 || req.EndTime == 0 || (req.HasRoomID && req.RoomID == 0) || req.ClassSubjectsID == 0
		},
		func(ctx context.Context, teacher_scope helpers.TeacherScope, queries *db_queries.Queries, req EditExamRequest) (int64, error) {
			return queries.TeacherEditExam(ctx, db_queries.TeacherEditExamParams{
				ID:              req.ExamID,
				TeacherID:       teacher_scope.TeacherID,
				SchoolID:        teacher_scope.SchoolID,
				Date:            pgtype.Date{Time: time.Unix(req.Date, 0), Valid: true},
				Title:           req.Title,
				Description:     pgtype.Text{String: req.Description, Valid: req.HasDescription},
				StartTime:       pgtype.Time{Microseconds: req.StartTime * 1000000, Valid: true},
				EndTime:         pgtype.Time{Microseconds: req.EndTime * 1000000, Valid: true},
				RoomID:          pgtype.Int4{Int32: req.RoomID, Valid: req.HasRoomID},
				ClassSubjectsID: req.ClassSubjectsID,
			})
		})
}

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

type AddGradeRequest struct {
	TermID          int32  `json:"term_id"`
	GradeTypeID     int32  `json:"grade_type_id"`
	Value           int16  `json:"value"`
	Date            int64  `json:"date"`
	HasNote         bool   `json:"has_note"`
	Note            string `json:"note"`
	Title           string `json:"title"`
	HasDescription  bool   `json:"has_description"`
	Description     string `json:"description"`
	StudentID       int32  `json:"student_id"`
	ClassSubjectsID int32  `json:"class_subjects_id"`
}

type AddFinalGradeRequest struct {
	TermID          int32 `json:"term_id"`
	Value           int16 `json:"value"`
	StudentID       int32 `json:"student_id"`
	ClassSubjectsID int32 `json:"class_subjects_id"`
}

type EditGradeRequest struct {
	GradeID         int64  `json:"grade_id"`
	TermID          int32  `json:"term_id"`
	GradeTypeID     int32  `json:"grade_type_id"`
	Value           int16  `json:"value"`
	Date            int64  `json:"date"`
	HasNote         bool   `json:"has_note"`
	Note            string `json:"note"`
	Title           string `json:"title"`
	HasDescription  bool   `json:"has_description"`
	Description     string `json:"description"`
	ClassSubjectsID int32  `json:"class_subjects_id"`
}

type EditFinalGradeRequest struct {
	FinalGradeID    int64 `json:"grade_id"`
	TermID          int32 `json:"term_id"`
	GradeTypeID     int32 `json:"grade_type_id"`
	Value           int16 `json:"value"`
	ClassSubjectsID int32 `json:"class_subjects_id"`
}

type RemoveGradeRequest struct {
	GradeID int64 `json:"grade_id"`
}

type RemoveFinalGradeRequest struct {
	FinalGradeID int64 `json:"final_grade_id"`
}

type GradeSummary struct {
	ID             int64  `json:"id"`
	Subject        string `json:"subject"`
	HasSubjectCode bool   `json:"has_subject_code"`
	SubjectCode    string `json:"subject_code"`
	Term           string `json:"term"`
	Type           string `json:"type"`
	Value          int16  `json:"value"`
	Date           int64  `json:"date"`
	HasNote        bool   `json:"has_note"`
	Note           string `json:"note"`
}

type FinalGradeSummary struct {
	ID             int64  `json:"id"`
	Subject        string `json:"subject"`
	HasSubjectCode bool   `json:"has_subject_code"`
	SubjectCode    string `json:"subject_code"`
	Term           string `json:"term"`
	Value          int16  `json:"value"`
}

func convertGrade(row db_queries.GetTeacherGradesRow) GradeSummary {
	return GradeSummary{
		ID:             row.ID,
		Subject:        row.Subject,
		HasSubjectCode: row.SubjectCode.Valid,
		SubjectCode:    row.SubjectCode.String,
		Term:           row.Term,
		Type:           row.Type,
		Value:          row.Value,
		Date:           row.Date.Time.Unix(),
		HasNote:        row.Note.Valid,
		Note:           row.Note.String,
	}
}

func convertFinalGrade(row db_queries.GetTeacherFinalGradesRow) FinalGradeSummary {
	return FinalGradeSummary{
		ID:             row.ID,
		Subject:        row.Subject,
		HasSubjectCode: row.SubjectCode.Valid,
		SubjectCode:    row.SubjectCode.String,
		Term:           row.Term,
		Value:          row.Value,
	}
}

func ListGrades(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return TeacherSummary(c, pool, rdb, "GRADES_CACHE_TTL", helpers.CacheOrGetTeacherGrades, convertGrade)
}

func ListFinalGrades(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return TeacherSummary(c, pool, rdb, "FINAL_GRADES_CACHE_TTL", helpers.CacheOrGetTeacherFinalGrades, convertFinalGrade)
}

func AddGrade(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return TeacherModify[AddGradeRequest](c, pool, rdb, true,
		func(req AddGradeRequest) int32 { return req.ClassSubjectsID },
		func(req AddGradeRequest) bool {
			return req.TermID == 0 || req.GradeTypeID == 0 || req.Value == 0 || req.Date == 0 || (req.HasNote && req.Note == "") || req.Title == "" || (req.HasDescription && req.Description == "") || req.ClassSubjectsID == 0
		},
		func(ctx context.Context, teacher_scope helpers.TeacherScope, queries *db_queries.Queries, req AddGradeRequest) (int64, error) {
			return queries.TeacherAddGrade(ctx, db_queries.TeacherAddGradeParams{
				TeacherID:       teacher_scope.TeacherID,
				TermID:          req.TermID,
				GradeTypeID:     req.GradeTypeID,
				ClassSubjectsID: req.ClassSubjectsID,
				Value:           req.Value,
				Date:            pgtype.Date{Time: time.Unix(req.Date, 0), Valid: true},
				Note:            pgtype.Text{String: req.Note, Valid: req.HasNote},
				SchoolID:        teacher_scope.SchoolID,
				StudentID:       req.StudentID,
			})
		})
}

func AddFinalGrade(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return TeacherModify[AddFinalGradeRequest](c, pool, rdb, true,
		func(req AddFinalGradeRequest) int32 { return req.ClassSubjectsID },
		func(req AddFinalGradeRequest) bool {
			return req.TermID == 0 || req.Value == 0 || req.ClassSubjectsID == 0
		},
		func(ctx context.Context, teacher_scope helpers.TeacherScope, queries *db_queries.Queries, req AddFinalGradeRequest) (int64, error) {
			return queries.TeacherAddFinalGrade(ctx, db_queries.TeacherAddFinalGradeParams{
				TermID:          req.TermID,
				TeacherID:       teacher_scope.TeacherID,
				Value:           req.Value,
				ClassSubjectsID: req.ClassSubjectsID,
				SchoolID:        teacher_scope.SchoolID,
				StudentID:       req.StudentID,
			})
		})
}

func RemoveGrade(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return TeacherModify[RemoveGradeRequest](c, pool, rdb, false,
		func(req RemoveGradeRequest) int32 { return 0 },
		func(req RemoveGradeRequest) bool { return req.GradeID == 0 },
		func(ctx context.Context, teacher_scope helpers.TeacherScope, queries *db_queries.Queries, req RemoveGradeRequest) (int64, error) {
			return queries.TeacherDeleteGrade(ctx, db_queries.TeacherDeleteGradeParams{
				ID:        req.GradeID,
				TeacherID: teacher_scope.TeacherID,
				SchoolID:  teacher_scope.SchoolID,
			})
		})
}

func EditGrade(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return TeacherModify[EditGradeRequest](c, pool, rdb, true,
		func(req EditGradeRequest) int32 { return req.ClassSubjectsID },
		func(req EditGradeRequest) bool {
			return req.GradeID == 0 || req.TermID == 0 || req.GradeTypeID == 0 || req.Value == 0 || req.Date == 0 || (req.HasNote && req.Note == "") || req.Title == "" || (req.HasDescription && req.Description == "") || req.ClassSubjectsID == 0
		},
		func(ctx context.Context, teacher_scope helpers.TeacherScope, queries *db_queries.Queries, req EditGradeRequest) (int64, error) {
			return queries.TeacherEditGrade(ctx, db_queries.TeacherEditGradeParams{
				ID:              req.GradeID,
				TeacherID:       teacher_scope.TeacherID,
				TermID:          req.TermID,
				GradeTypeID:     req.GradeTypeID,
				ClassSubjectsID: req.ClassSubjectsID,
				Value:           req.Value,
				Date:            pgtype.Date{Time: time.Unix(req.Date, 0), Valid: true},
				Note:            pgtype.Text{String: req.Note, Valid: req.HasNote},
				SchoolID:        teacher_scope.SchoolID,
			})
		})
}

func EditFinalGrade(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return TeacherModify[EditFinalGradeRequest](c, pool, rdb, true,
		func(req EditFinalGradeRequest) int32 { return req.ClassSubjectsID },
		func(req EditFinalGradeRequest) bool {
			return req.FinalGradeID == 0 || req.TermID == 0 || req.Value == 0 || req.ClassSubjectsID == 0
		},
		func(ctx context.Context, teacher_scope helpers.TeacherScope, queries *db_queries.Queries, req EditFinalGradeRequest) (int64, error) {
			return queries.TeacherEditFinalGrade(ctx, db_queries.TeacherEditFinalGradeParams{
				ID:              req.FinalGradeID,
				ClassSubjectsID: req.ClassSubjectsID,
				TermID:          req.TermID,
				Value:           req.Value,
				SchoolID:        teacher_scope.SchoolID,
				TeacherID:       teacher_scope.TeacherID,
			})
		})
}

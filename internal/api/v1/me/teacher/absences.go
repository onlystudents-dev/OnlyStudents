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

type AddAbsenceRequest struct {
	Date            int64  `json:"date"`
	Type            string `json:"type"`
	HasNote         bool   `json:"has_note"`
	Note            string `json:"note"`
	StudentID       int32  `json:"student_id"`
	ClassSubjectsID int32  `json:"class_subjects_id"`
}

type EditAbsenceRequest struct {
	AbsenceID       int64  `json:"absence_id"`
	Date            int64  `json:"date"`
	Type            string `json:"type"`
	HasNote         bool   `json:"has_note"`
	Note            string `json:"note"`
	StudentID       int32  `json:"student_id"`
	ClassSubjectsID int32  `json:"class_subjects_id"`
}

type RemoveAbsenceRequest struct {
	AbsenceID int64 `json:"absence_id"`
}

type ListAbsencesRequest struct{}

type AbsenceSummary struct {
	ID             int64       `json:"id"`
	Subject        string      `json:"subject"`
	HasSubjectCode bool        `json:"has_subject_code"`
	SubjectCode    string      `json:"subject_code"`
	Date           int64       `json:"date"`
	Type           string      `json:"type"`
	Justified      bool        `json:"justified"`
	HasNote        bool        `json:"has_note"`
	Note           string      `json:"note"`
	VerifiedBy     interface{} `json:"verified_by"`
}

func convertAbsence(row db_queries.GetTeacherAbsencesRow) AbsenceSummary {
	return AbsenceSummary{
		ID:             row.ID,
		Subject:        row.Subject,
		HasSubjectCode: row.SubjectCode.Valid,
		SubjectCode:    row.SubjectCode.String,
		Date:           row.Date.Time.Unix(),
		Type:           row.Type,
		Justified:      row.Justified,
		HasNote:        row.Note.Valid,
		Note:           row.Note.String,
		VerifiedBy:     row.VerifiedBy,
	}
}

func ListAbsences(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return TeacherSummary[db_queries.GetTeacherAbsencesRow, AbsenceSummary](c, pool, rdb, "ABSENCES_CACHE_TTL", helpers.CacheOrGetTeacherAbsences, convertAbsence)
}

func AddAbsence(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return TeacherModify[AddAbsenceRequest](c, pool, rdb, true,
		func(req AddAbsenceRequest) int32 { return req.ClassSubjectsID },
		func(req AddAbsenceRequest) bool {
			return req.Date == 0 || req.Type == "" || req.StudentID == 0 || req.ClassSubjectsID == 0 || (req.HasNote && req.Note == "")
		},
		func(ctx context.Context, teacher_scope helpers.TeacherScope, queries *db_queries.Queries, req AddAbsenceRequest) (int64, error) {
			return queries.TeacherAddAbsence(ctx, db_queries.TeacherAddAbsenceParams{
				Date:            pgtype.Date{Time: time.Unix(req.Date, 0), Valid: true},
				Type:            req.Type,
				Note:            pgtype.Text{String: req.Note, Valid: req.HasNote},
				TeacherID:       teacher_scope.TeacherID,
				SchoolID:        teacher_scope.SchoolID,
				StudentID:       req.StudentID,
				ClassSubjectsID: req.ClassSubjectsID,
			})
		})
}

func RemoveAbsence(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return TeacherModify[RemoveAbsenceRequest](c, pool, rdb, false,
		func(req RemoveAbsenceRequest) int32 { return 0 },
		func(req RemoveAbsenceRequest) bool { return req.AbsenceID == 0 },
		func(ctx context.Context, teacher_scope helpers.TeacherScope, queries *db_queries.Queries, req RemoveAbsenceRequest) (int64, error) {
			return queries.TeacherDeleteAbsence(ctx, db_queries.TeacherDeleteAbsenceParams{
				ID:        req.AbsenceID,
				TeacherID: teacher_scope.TeacherID,
				SchoolID:  teacher_scope.SchoolID,
			})
		})
}

func EditAbsence(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return TeacherModify[EditAbsenceRequest](c, pool, rdb, true,
		func(req EditAbsenceRequest) int32 { return req.ClassSubjectsID },
		func(req EditAbsenceRequest) bool {
			return req.AbsenceID == 0 || req.Date == 0 || req.Type == "" || req.StudentID == 0 || req.ClassSubjectsID == 0 || (req.HasNote && req.Note == "")
		},
		func(ctx context.Context, teacher_scope helpers.TeacherScope, queries *db_queries.Queries, req EditAbsenceRequest) (int64, error) {
			return queries.TeacherEditAbsence(ctx, db_queries.TeacherEditAbsenceParams{
				ID:              req.AbsenceID,
				Date:            pgtype.Date{Time: time.Unix(req.Date, 0), Valid: true},
				Type:            req.Type,
				Note:            pgtype.Text{String: req.Note, Valid: req.HasNote},
				TeacherID:       teacher_scope.TeacherID,
				SchoolID:        teacher_scope.SchoolID,
				ClassSubjectsID: pgtype.Int4{Int32: req.ClassSubjectsID, Valid: true},
			})
		})
}

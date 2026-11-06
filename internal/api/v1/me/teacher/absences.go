package teacherapi

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

type ListAbsencesByLessonRequest struct {
	ClassSubjectsID int32 `json:"class_subjects_id"`
	Date            int64 `json:"date"`
}

type ChangeType string

const (
	ChangeDelete ChangeType = "DELETE"
	ChangeAdd    ChangeType = "ADD"
	ChangeUpdate ChangeType = "UPDATE"
)

type AbsenceKind string

const (
	KindAbsent AbsenceKind = "absent"
	KindTardy  AbsenceKind = "tardy"
)

type AbsenceChange struct {
	Type        ChangeType   `json:"type"`
	StudentID   int32        `json:"student_id"`
	ID          *int64       `json:"id,omitempty"`
	AbsenceType *AbsenceKind `json:"absence_type,omitempty"`
}

type ChangeAbsencesRequest struct {
	AbsenceChanges  []AbsenceChange `json:"absence_changes"`
	ClassSubjectsID int32           `json:"class_subjects_id"`
	Date            int64           `json:"date"`
}

type AbsenceSummary struct {
	ID             int64       `json:"id"`
	StudentID      int32       `json:"student_id"`
	Subject        string      `json:"subject"`
	HasSubjectCode bool        `json:"has_subject_code"`
	SubjectCode    string      `json:"subject_code"`
	Date           int64       `json:"date"`
	Type           string      `json:"type"`
	Justified      bool        `json:"justified"`
	TardyMinutes   int32       `json:"tardy_minutes"`
	FileName       string      `json:"file_name"`
	FileType       string      `json:"file_type"`
	FileSize       int64       `json:"file_size"`
	FileUUID       string      `json:"file_uuid"`
	ReasonNote     string      `json:"reason_note"`
	Note           string      `json:"note"`
	VerifiedBy     interface{} `json:"verified_by"`
}

func convertAbsence(row db_queries.GetTeacherAbsencesRow) AbsenceSummary {
	return AbsenceSummary{
		ID:             row.ID,
		StudentID:      row.StudentID,
		Subject:        row.Subject,
		HasSubjectCode: row.SubjectCode.Valid,
		SubjectCode:    row.SubjectCode.String,
		Date:           row.Date.Time.Unix(),
		Type:           row.Type,
		Justified:      row.Justified,
		TardyMinutes:   row.TardyMinutes.Int32,
		FileName:       row.FileName.String,
		FileType:       row.FileType.String,
		FileSize:       row.FileSize.Int64,
		FileUUID:       row.FileUuid.String(),
		ReasonNote:     row.ReasonNote.String,
		Note:           row.Note.String,
		VerifiedBy:     row.VerifiedBy,
	}
}

func convertAbsenceByLesson(row db_queries.GetTeacherAbsencesByLessonRow) AbsenceSummary {
	return AbsenceSummary{
		ID:             row.ID,
		StudentID:      row.StudentID,
		Subject:        row.Subject,
		HasSubjectCode: row.SubjectCode.Valid,
		SubjectCode:    row.SubjectCode.String,
		Date:           row.Date.Time.Unix(),
		Type:           row.Type,
		Justified:      row.Justified,
		TardyMinutes:   row.TardyMinutes.Int32,
		FileName:       row.FileName.String,
		FileType:       row.FileType.String,
		FileSize:       row.FileSize.Int64,
		FileUUID:       row.FileUuid.String(),
		ReasonNote:     row.ReasonNote.String,
		Note:           row.Note.String,
		VerifiedBy:     row.VerifiedBy,
	}
}

func ListAbsences(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return helpers.TeacherSummary(c, pool, rdb, "ABSENCES_CACHE_TTL", "", nil, func(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, req bool, scope helpers.TeacherScope, ttl int32) ([]db_queries.GetTeacherAbsencesRow, error) {
		return helpers.CacheOrGetTeacherAbsences(ctx, rdb, queries, scope, ttl)
	}, convertAbsence)
}

func ListAbsencesForLesson(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return helpers.TeacherSummary(c, pool, rdb, "ABSENCES_CACHE_TTL", "", func(req ListAbsencesByLessonRequest) bool { return req.ClassSubjectsID <= 0 || req.Date <= 0 }, func(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, req ListAbsencesByLessonRequest, scope helpers.TeacherScope, ttl int32) ([]db_queries.GetTeacherAbsencesByLessonRow, error) {
		return helpers.CacheOrGetTeacherAbsencesByLesson(ctx, rdb, queries, scope, ttl, req.ClassSubjectsID, pgtype.Date{Time: time.Unix(req.Date, 0), Valid: true})
	}, convertAbsenceByLesson)
}

func ChangeAbsencesValidate(req ChangeAbsencesRequest) bool {
	if req.Date <= 0 || req.ClassSubjectsID <= 0 || len(req.AbsenceChanges) == 0 {
		return true
	}

	for _, absence_change := range req.AbsenceChanges {
		switch absence_change.Type {
		case ChangeDelete:
			if absence_change.ID == nil {
				return true
			}
		case ChangeAdd:
			if absence_change.AbsenceType == nil {
				return true
			}
		case ChangeUpdate:
			if absence_change.AbsenceType == nil || absence_change.ID == nil {
				return true
			}
		default:
			return true
		}
	}

	return false
}

func ChangeAbsences(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return helpers.TeacherModify(c, pool, rdb,
		"",
		func(req ChangeAbsencesRequest) bool {
			return ChangeAbsencesValidate(req)
		},
		func(ctx context.Context, teacher_scope helpers.TeacherScope, queries *db_queries.Queries, req ChangeAbsencesRequest) (int64, error) {
			for _, absence_change := range req.AbsenceChanges {
				switch absence_change.Type {
				case ChangeDelete:
					row_count, err := queries.TeacherDeleteAbsence(ctx, db_queries.TeacherDeleteAbsenceParams{
						ID:        *absence_change.ID,
						TeacherID: teacher_scope.TeacherID,
						SchoolID:  teacher_scope.SchoolID,
					})

					if err != nil || row_count != 1 {
						return row_count, err
					}
				case ChangeAdd:
					row_count, err := queries.TeacherAddAbsence(ctx, db_queries.TeacherAddAbsenceParams{
						Date:            pgtype.Date{Time: time.Unix(req.Date, 0), Valid: true},
						Type:            string(*absence_change.AbsenceType),
						Note:            pgtype.Text{String: "", Valid: false},
						TeacherID:       teacher_scope.TeacherID,
						SchoolID:        teacher_scope.SchoolID,
						StudentID:       absence_change.StudentID,
						ClassSubjectsID: req.ClassSubjectsID,
					})

					if err != nil || row_count != 1 {
						return row_count, err
					}
				case ChangeUpdate:
					row_count, err := queries.TeacherEditAbsence(ctx, db_queries.TeacherEditAbsenceParams{
						ID:              *absence_change.ID,
						Date:            pgtype.Date{Time: time.Unix(req.Date, 0), Valid: true},
						Type:            string(*absence_change.AbsenceType),
						Note:            pgtype.Text{String: "", Valid: false},
						TeacherID:       teacher_scope.TeacherID,
						SchoolID:        teacher_scope.SchoolID,
						ClassSubjectsID: pgtype.Int4{Int32: req.ClassSubjectsID, Valid: true},
					})

					if err != nil || row_count != 1 {
						return row_count, err
					}
				default:
					return 0, errors.New("How is this possible?")
				}
			}

			return 1, nil
		})
}

package teacherapi

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

type AddGradeRequest struct {
	TermID        int32  `json:"term_id"`
	LessonGradeID int64  `json:"lesson_grade_id"`
	Value         int16  `json:"value"`
	Date          int64  `json:"date"`
	HasNote       bool   `json:"has_note"`
	Note          string `json:"note"`
	StudentID     int32  `json:"student_id"`
}

type AddFinalGradeRequest struct {
	TermID          int32 `json:"term_id"`
	Value           int16 `json:"value"`
	StudentID       int32 `json:"student_id"`
	ClassSubjectsID int32 `json:"class_subjects_id"`
}

type EditGradeRequest struct {
	GradeID       int64  `json:"grade_id"`
	TermID        int32  `json:"term_id"`
	Value         int16  `json:"value"`
	Date          int64  `json:"date"`
	HasNote       bool   `json:"has_note"`
	Note          string `json:"note"`
	LessonGradeID int64  `json:"lesson_grade_id"`
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
	TermID         int32  `json:"term_id"`
	Subject        string `json:"subject"`
	HasSubjectCode bool   `json:"has_subject_code"`
	SubjectCode    string `json:"subject_code"`
	Term           string `json:"term"`
	Type           string `json:"type"`
	Value          int16  `json:"value"`
	Date           int64  `json:"date"`
	HasNote        bool   `json:"has_note"`
	Note           string `json:"note"`
	StudentID      int32  `json:"student_id"`
	LessonID       int64  `json:"lesson_id"`
}

type FinalGradeSummary struct {
	ID             int64  `json:"id"`
	TermID         int32  `json:"term_id"`
	Subject        string `json:"subject"`
	HasSubjectCode bool   `json:"has_subject_code"`
	SubjectCode    string `json:"subject_code"`
	Term           string `json:"term"`
	Value          int16  `json:"value"`
}

type LessonGradeSummary struct {
	ID              int64  `json:"id"`
	ActualDate      int64  `json:"actual_date"`
	Title           string `json:"title"`
	Description     string `json:"description"`
	Type            int32  `json:"type"`
	ClassSubjectsID int32  `json:"class_subjects_id"`
	GroupID         int32  `json:"group_id"`
}

type AddLessonGradeRequest struct {
	ActualDate      int64  `json:"actual_date"`
	Title           string `json:"title"`
	Description     string `json:"description"`
	Type            int32  `json:"type"`
	ClassSubjectsID int32  `json:"class_subjects_id"`
	GroupID         int32  `json:"group_id"`
}

type EditLessonGradeRequest struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Type        int32  `json:"type"`
}

type RemoveLessonGradeRequest struct {
	ID int64 `json:"lesson_grade_id"`
}

type GradeType struct {
	Name   string `json:"name"`
	Weight int32  `json:"weight"`
}

func convertGrade(row db_queries.GetTeacherGradesRow) GradeSummary {
	return GradeSummary{
		ID:             row.ID,
		TermID:         row.TermID,
		Subject:        row.Subject,
		HasSubjectCode: row.SubjectCode.Valid,
		SubjectCode:    row.SubjectCode.String,
		Term:           row.Term,
		Type:           row.Type,
		Value:          row.Value,
		Date:           row.Date.Time.Unix(),
		HasNote:        row.Note.Valid,
		Note:           row.Note.String,
		StudentID:      row.StudentID,
		LessonID:       row.LessonGradeID,
	}
}

func convertGradeByLesson(row db_queries.GetTeacherGradesByLessonRow) GradeSummary {
	return GradeSummary{
		ID:             row.ID,
		TermID:         row.TermID,
		Subject:        row.Subject,
		HasSubjectCode: row.SubjectCode.Valid,
		SubjectCode:    row.SubjectCode.String,
		Term:           row.Term,
		Type:           row.Type,
		Value:          row.Value,
		Date:           row.Date.Time.Unix(),
		HasNote:        row.Note.Valid,
		Note:           row.Note.String,
		StudentID:      row.StudentID,
		LessonID:       row.LessonGradeID,
	}
}

func convertFinalGrade(row db_queries.GetTeacherFinalGradesRow) FinalGradeSummary {
	return FinalGradeSummary{
		ID:             row.ID,
		TermID:         row.TermID,
		Subject:        row.Subject,
		HasSubjectCode: row.SubjectCode.Valid,
		SubjectCode:    row.SubjectCode.String,
		Term:           row.Term,
		Value:          row.Value,
	}
}

func convertGradeTypes(row db_queries.GetTeacherGradeTypesRow) GradeType {
	return GradeType{
		Name:   row.Name,
		Weight: row.Weight,
	}
}

func convertLessonGrade(row db_queries.GetTeacherLessonGradesRow) LessonGradeSummary {
	return LessonGradeSummary{
		ID:              row.ID,
		ActualDate:      row.ActualDate,
		Title:           row.Title,
		Description:     row.Description,
		Type:            row.GradeTypeID,
		ClassSubjectsID: row.ClassSubjectsID,
		GroupID:         row.GroupID,
	}
}

func convertLessonGradeBySubject(row db_queries.GetTeacherLessonGradesBySubjectRow) LessonGradeSummary {
	return LessonGradeSummary{
		ID:              row.ID,
		ActualDate:      row.ActualDate,
		Title:           row.Title,
		Description:     row.Description,
		Type:            row.GradeTypeID,
		ClassSubjectsID: row.ClassSubjectsID,
		GroupID:         row.GroupID,
	}
}

func ListGrades(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	id, err := strconv.ParseInt(c.Query("lesson_id"), 10, 64)
	if err == nil {
		return helpers.TeacherSummary(c, pool, rdb, "GRADES_BY_LESSON_CACHE_TTL", "", nil, func(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, req bool, scope helpers.TeacherScope, ttl int32) ([]db_queries.GetTeacherGradesByLessonRow, error) {
			return helpers.CacheOrGetTeacherGradesByLesson(ctx, rdb, queries, scope, ttl, id)
		}, convertGradeByLesson)
	}

	return helpers.TeacherSummary(c, pool, rdb, "GRADES_CACHE_TTL", "", nil, func(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, req bool, scope helpers.TeacherScope, ttl int32) ([]db_queries.GetTeacherGradesRow, error) {
		return helpers.CacheOrGetTeacherGrades(ctx, rdb, queries, scope, ttl)
	}, convertGrade)
}

func ListFinalGrades(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return helpers.TeacherSummary(c, pool, rdb, "FINAL_GRADES_CACHE_TTL", "", nil, func(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, req bool, scope helpers.TeacherScope, ttl int32) ([]db_queries.GetTeacherFinalGradesRow, error) {
		return helpers.CacheOrGetTeacherFinalGrades(ctx, rdb, queries, scope, ttl)
	}, convertFinalGrade)
}

func GetGradeTypes(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return helpers.TeacherSummary(c, pool, rdb, "GRADE_TYPES_TTL", "", nil, func(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, req bool, scope helpers.TeacherScope, ttl int32) ([]db_queries.GetTeacherGradeTypesRow, error) {
		return helpers.CacheOrGetTeacherGradeTypes(ctx, rdb, queries, scope, ttl)
	}, convertGradeTypes)
}

func ListLessonGrades(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	id, err := strconv.ParseInt(c.Query("subject_id"), 10, 32)
	date, err2 := strconv.ParseInt(c.Query("actual_date"), 10, 64)
	if err == nil && err2 == nil {
		return helpers.TeacherSummary(c, pool, rdb, "LESSON_GRADES_BY_SUBJECT_CACHE_TTL", "", nil, func(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, req bool, scope helpers.TeacherScope, ttl int32) ([]db_queries.GetTeacherLessonGradesBySubjectRow, error) {
			return helpers.CacheOrGetTeacherLessonGradesBySubject(ctx, rdb, queries, scope, ttl, int32(id), date) // #nosec G115 -- subject_id parsed with bitSize 32, fits int32
		}, convertLessonGradeBySubject)
	}

	return helpers.TeacherSummary(c, pool, rdb, "LESSON_GRADES_CACHE_TTL", "", nil, func(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, req bool, scope helpers.TeacherScope, ttl int32) ([]db_queries.GetTeacherLessonGradesRow, error) {
		return helpers.CacheOrGetTeacherLessonGrades(ctx, rdb, queries, scope, ttl)
	}, convertLessonGrade)
}

func AddGrade(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return helpers.TeacherModify(c, pool, rdb,
		"",
		func(req AddGradeRequest) bool {
			return req.StudentID <= 0 || req.TermID <= 0 || req.LessonGradeID <= 0 || req.Value <= 0 || req.Date <= 0 || (req.HasNote && req.Note == "")
		},
		func(ctx context.Context, teacher_scope helpers.TeacherScope, queries *db_queries.Queries, req AddGradeRequest) (int64, error) {
			return queries.TeacherAddGrade(ctx, db_queries.TeacherAddGradeParams{
				TeacherID:     teacher_scope.TeacherID,
				TermID:        req.TermID,
				LessonGradeID: req.LessonGradeID,
				Value:         req.Value,
				Date:          pgtype.Date{Time: time.Unix(req.Date, 0), Valid: true},
				Note:          pgtype.Text{String: req.Note, Valid: req.HasNote},
				SchoolID:      teacher_scope.SchoolID,
				StudentID:     req.StudentID,
			})
		})
}

func AddFinalGrade(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return helpers.TeacherModify(c, pool, rdb,
		"",
		func(req AddFinalGradeRequest) bool {
			return req.TermID <= 0 || req.Value <= 0 || req.ClassSubjectsID <= 0
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
	return helpers.TeacherModify(c, pool, rdb,
		"",
		func(req RemoveGradeRequest) bool { return req.GradeID <= 0 },
		func(ctx context.Context, teacher_scope helpers.TeacherScope, queries *db_queries.Queries, req RemoveGradeRequest) (int64, error) {
			return queries.TeacherDeleteGrade(ctx, db_queries.TeacherDeleteGradeParams{
				ID:        req.GradeID,
				TeacherID: teacher_scope.TeacherID,
				SchoolID:  teacher_scope.SchoolID,
			})
		})
}

func EditGrade(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return helpers.TeacherModify(c, pool, rdb,
		"",
		func(req EditGradeRequest) bool {
			return req.GradeID <= 0 || req.TermID <= 0 || req.LessonGradeID <= 0 || req.Value <= 0 || req.Date <= 0 || (req.HasNote && req.Note == "")
		},
		func(ctx context.Context, teacher_scope helpers.TeacherScope, queries *db_queries.Queries, req EditGradeRequest) (int64, error) {
			return queries.TeacherEditGrade(ctx, db_queries.TeacherEditGradeParams{
				ID:            req.GradeID,
				TermID:        req.TermID,
				LessonGradeID: req.LessonGradeID,
				Value:         req.Value,
				Date:          pgtype.Date{Time: time.Unix(req.Date, 0), Valid: true},
				Note:          pgtype.Text{String: req.Note, Valid: req.HasNote},
				TeacherID:     teacher_scope.TeacherID,
				SchoolID:      teacher_scope.SchoolID,
			})
		})
}

func EditFinalGrade(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return helpers.TeacherModify(c, pool, rdb,
		"",
		func(req EditFinalGradeRequest) bool {
			return req.FinalGradeID <= 0 || req.TermID <= 0 || req.Value <= 0 || req.ClassSubjectsID <= 0
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

func AddLessonGrade(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return helpers.TeacherModify(c, pool, rdb,
		"",
		func(req AddLessonGradeRequest) bool {
			return req.ActualDate <= 0 || req.Title == "" || req.Type <= 0 || req.ClassSubjectsID <= 0 || req.GroupID <= 0
		},
		func(ctx context.Context, teacher_scope helpers.TeacherScope, queries *db_queries.Queries, req AddLessonGradeRequest) (int64, error) {
			return queries.TeacherAddLessonGrade(ctx, db_queries.TeacherAddLessonGradeParams{
				ActualDate:      req.ActualDate,
				Title:           req.Title,
				Description:     req.Description,
				Type:            req.Type,
				ClassSubjectsID: req.ClassSubjectsID,
				GroupID:         pgtype.Int4{Int32: req.GroupID, Valid: true},
				SchoolID:        teacher_scope.SchoolID,
				TeacherID:       teacher_scope.TeacherID,
			})
		},
	)
}

func EditLessonGrade(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return helpers.TeacherModify(c, pool, rdb,
		"",
		func(req EditLessonGradeRequest) bool {
			return req.Title == "" || req.Type <= 0
		},
		func(ctx context.Context, teacher_scope helpers.TeacherScope, queries *db_queries.Queries, req EditLessonGradeRequest) (int64, error) {
			return queries.TeacherEditLessonGrade(ctx, db_queries.TeacherEditLessonGradeParams{
				ID:          req.ID,
				Title:       req.Title,
				Description: req.Description,
				Type:        req.Type,
				SchoolID:    teacher_scope.SchoolID,
				TeacherID:   teacher_scope.TeacherID,
			})
		},
	)
}

func RemoveLessonGrade(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return helpers.TeacherModify(c, pool, rdb,
		"",
		func(req RemoveLessonGradeRequest) bool {
			return req.ID <= 0
		},
		func(ctx context.Context, teacher_scope helpers.TeacherScope, queries *db_queries.Queries, req RemoveLessonGradeRequest) (int64, error) {
			return queries.TeacherDeleteLessonGrade(ctx, db_queries.TeacherDeleteLessonGradeParams{
				ID:        req.ID,
				SchoolID:  teacher_scope.SchoolID,
				TeacherID: teacher_scope.TeacherID,
			})
		},
	)
}

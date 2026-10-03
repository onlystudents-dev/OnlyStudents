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

type ListHomeworkSubmissionsRequest struct {
	HomeworkID int32 `json:"homework_id"`
}

type AddHomeworkRequest struct {
	Title           string `json:"title"`
	HasDescription  bool   `json:"has_description"`
	Description     string `json:"description"`
	DueDate         int64  `json:"due_date"`
	ClassSubjectsID int32  `json:"class_subjects_id"`
}

type EditHomeworkRequest struct {
	HomeworkID      int64  `json:"homework_id"`
	Title           string `json:"title"`
	HasDescription  bool   `json:"has_description"`
	Description     string `json:"description"`
	DueDate         int64  `json:"due_date"`
	ClassSubjectsID int32  `json:"class_subjects_id"`
}

type UpdateHomeworkSubmissionRequest struct {
	HomeworkID      int32  `json:"homework_id"`
	StudentID       int32  `json:"student_id"`
	HasContent      bool   `json:"has_content"`
	Content         string `json:"content"`
	HasGradedValue  bool   `json:"has_graded_value"`
	GradedValue     int16  `json:"graded_value"`
	ClassSubjectsID int32  `json:"class_subjects_id"`
}

type RemoveHomeworkRequest struct {
	HomeworkID int64 `json:"homework_id"`
}

type HomeworkSummary struct {
	ID                    int64  `json:"id"`
	Subject               string `json:"subject"`
	HasSubjectCode        bool   `json:"has_subject_code"`
	SubjectCode           string `json:"subject_code"`
	Title                 string `json:"title"`
	HasDescription        bool   `json:"has_description"`
	Description           string `json:"description"`
	DueDate               int64  `json:"due_date"`
	CreatedAt             int64  `json:"created_at"`
	SubmissionCount       int64  `json:"submission_count"`
	OnTimeSubmissionCount int64  `json:"on_timesubmission_count"`
}

type HomeworkSubmissionSummary struct {
	ID             int64  `json:"id"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	HasContent     bool   `json:"has_content"`
	Content        string `json:"content"`
	SubmittedAt    int64  `json:"submitted_at"`
	HasGradedValue bool   `json:"has_graded_value"`
	GradedValue    int16  `json:"graded_value"`
}

func convertHomework(row db_queries.GetTeacherHomeworkRow) HomeworkSummary {
	return HomeworkSummary{
		ID:                    row.ID,
		Subject:               row.Subject,
		HasSubjectCode:        row.SubjectCode.Valid,
		SubjectCode:           row.SubjectCode.String,
		Title:                 row.Title,
		HasDescription:        row.Description.Valid,
		Description:           row.Description.String,
		DueDate:               row.DueDate.Time.Unix(),
		CreatedAt:             row.CreatedAt.Time.Unix(),
		SubmissionCount:       row.SubmissionCount,
		OnTimeSubmissionCount: row.OnTimeSubmissionCount,
	}
}

func convertHomeworkSubmission(row db_queries.GetTeacherHomeworkSubmissionsRow) HomeworkSubmissionSummary {
	return HomeworkSubmissionSummary{
		ID:        row.ID,
		FirstName: row.FirstName,
		LastName:  row.LastName,
	}
}

func ListHomeworkSubmissions(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return TeacherSummaryByID(c, pool, rdb, "HOMEWORK_SUBMISSIONS_CACHE_TTL",
		func(req ListHomeworkSubmissionsRequest) bool { return req.HomeworkID == 0 },
		func(req ListHomeworkSubmissionsRequest) int32 { return req.HomeworkID },
		helpers.CacheOrGetTeacherHomeworkSubmissions, convertHomeworkSubmission)
}

func ListHomework(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return TeacherSummary(c, pool, rdb, "HOMEWORK_CACHE_TTL", helpers.CacheOrGetTeacherHomework, convertHomework)
}

func AddHomework(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return TeacherModify(c, pool, rdb, true,
		func(req AddHomeworkRequest) int32 { return req.ClassSubjectsID },
		func(req AddHomeworkRequest) bool {
			return req.Title == "" || (req.HasDescription && req.Description == "") || req.DueDate == 0 || req.ClassSubjectsID == 0
		},
		func(ctx context.Context, teacher_scope helpers.TeacherScope, queries *db_queries.Queries, req AddHomeworkRequest) (int64, error) {
			return queries.TeacherAddHomework(ctx, db_queries.TeacherAddHomeworkParams{
				TeacherID:       teacher_scope.TeacherID,
				Title:           req.Title,
				Description:     pgtype.Text{String: req.Description, Valid: req.HasDescription},
				DueDate:         pgtype.Date{Time: time.Unix(req.DueDate, 0), Valid: true},
				ClassSubjectsID: req.ClassSubjectsID,
				SchoolID:        teacher_scope.SchoolID,
			})
		})
}

func RemoveHomework(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return TeacherModify(c, pool, rdb, false,
		func(req RemoveHomeworkRequest) int32 { return 0 },
		func(req RemoveHomeworkRequest) bool { return req.HomeworkID == 0 },
		func(ctx context.Context, teacher_scope helpers.TeacherScope, queries *db_queries.Queries, req RemoveHomeworkRequest) (int64, error) {
			return queries.TeacherDeleteHomework(ctx, db_queries.TeacherDeleteHomeworkParams{
				ID:        req.HomeworkID,
				TeacherID: teacher_scope.TeacherID,
				SchoolID:  teacher_scope.SchoolID,
			})
		})
}

func EditHomework(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return TeacherModify(c, pool, rdb, true,
		func(req EditHomeworkRequest) int32 { return req.ClassSubjectsID },
		func(req EditHomeworkRequest) bool {
			return req.HomeworkID == 0 || req.Title == "" || (req.HasDescription && req.Description == "") || req.DueDate == 0 || req.ClassSubjectsID == 0
		},
		func(ctx context.Context, teacher_scope helpers.TeacherScope, queries *db_queries.Queries, req EditHomeworkRequest) (int64, error) {
			return queries.TeacherEditHomework(ctx, db_queries.TeacherEditHomeworkParams{
				ID:              req.HomeworkID,
				ClassSubjectsID: req.ClassSubjectsID,
				Title:           req.Title,
				Description:     pgtype.Text{String: req.Description, Valid: req.HasDescription},
				DueDate:         pgtype.Date{Time: time.Unix(req.DueDate, 0), Valid: true},
				SchoolID:        teacher_scope.SchoolID,
				TeacherID:       teacher_scope.TeacherID,
			})
		})
}

func UpdateHomeworkSubmission(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return TeacherModify(c, pool, rdb, true,
		func(req UpdateHomeworkSubmissionRequest) int32 { return req.ClassSubjectsID },
		func(req UpdateHomeworkSubmissionRequest) bool {
			return req.HomeworkID == 0 || req.StudentID == 0 || (req.HasContent && req.Content == "") || (req.HasGradedValue && req.GradedValue == 0) || req.ClassSubjectsID == 0
		},
		func(ctx context.Context, teacher_scope helpers.TeacherScope, queries *db_queries.Queries, req UpdateHomeworkSubmissionRequest) (int64, error) {
			return queries.TeacherUpdateHomeworkSubmission(ctx, db_queries.TeacherUpdateHomeworkSubmissionParams{
				Content:     pgtype.Text{String: req.Content, Valid: req.HasContent},
				GradedValue: pgtype.Int2{Int16: req.GradedValue, Valid: req.HasGradedValue},
				HomeworkID:  req.HomeworkID,
				StudentID:   req.StudentID,
				SchoolID:    teacher_scope.SchoolID,
				TeacherID:   teacher_scope.TeacherID,
			})
		})
}

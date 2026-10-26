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
	HasAssignedDate bool   `json:"has_assigned_date"`
	AssignedDate    int64  `json:"assigned_date"`
	ClassSubjectsID int32  `json:"class_subjects_id"`
}

type EditHomeworkRequest struct {
	HomeworkID      int64  `json:"homework_id"`
	Title           string `json:"title"`
	HasDescription  bool   `json:"has_description"`
	Description     string `json:"description"`
	DueDate         int64  `json:"due_date"`
	HasAssignedDate bool   `json:"has_assigned_date"`
	AssignedDate    int64  `json:"assigned_date"`
	ClassSubjectsID int32  `json:"class_subjects_id"`
}

type GradeHomeworkSubmissionRequest struct {
	HomeworkID  int32 `json:"homework_id"`
	StudentID   int32 `json:"student_id"`
	GradedValue int16 `json:"graded_value"`
}

type RemoveHomeworkRequest struct {
	HomeworkID int64 `json:"homework_id"`
}

type GetHomeworkRequest struct {
	HomeworkID int32 `json:"homework_id" query:"homework_id"`
}

type HomeworkSummary struct {
	ID              int64  `json:"id"`
	ClassSubjectsID int32  `json:"class_subjects_id"`
	Subject         string `json:"subject"`
	Title           string `json:"title"`
	HasDescription  bool   `json:"has_description"`
	Description     string `json:"description"`
	DueDate         int64  `json:"due_date"`
	CreatedAt       int64  `json:"created_at"`
}

type UngradedHomeworkSubmissionSummary struct {
	ID              int64  `json:"id"`
	HomeworkID      int32  `json:"homework_id"`
	HomeworkTitle   string `json:"title"`
	Subject         string `json:"subject"`
	StudentID       int32  `json:"student_id"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	DueDate         int64  `json:"due_date"`
	FileName        string `json:"file_name"`
	FileType        string `json:"file_type"`
	FileSize        int64  `json:"file_size"`
	FileUUID        string `json:"file_uuid"`
	Note            string `json:"note"`
	SubmittedAt     int64  `json:"created_at"`
	ClassSubjectsID int32  `json:"class_subjects_id"`
}

type HomeworkSubmissionSummary struct {
	ID             int64  `json:"id"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	DueDate        int64  `json:"due_date"`
	Note           string `json:"note"`
	FileName       string `json:"file_name"`
	FileType       string `json:"file_type"`
	FileSize       int64  `json:"file_size"`
	FileUUID       string `json:"file_uuid"`
	SubmittedAt    int64  `json:"submitted_at"`
	HasGradedValue bool   `json:"has_graded_value"`
	GradedValue    int16  `json:"graded_value"`
}

func convertHomeworkSubmission(row db_queries.GetTeacherHomeworkSubmissionsRow) HomeworkSubmissionSummary {
	return HomeworkSubmissionSummary{
		ID:             row.ID,
		FirstName:      row.FirstName,
		LastName:       row.LastName,
		DueDate:        row.DueDate.Time.Unix(),
		FileName:       row.FileName.String,
		FileType:       row.FileType.String,
		FileSize:       row.FileSize.Int64,
		FileUUID:       row.FileUuid.String(),
		Note:           row.Note.String,
		SubmittedAt:    row.SubmittedAt.Time.Unix(),
		HasGradedValue: row.GradedValue.Valid,
		GradedValue:    row.GradedValue.Int16,
	}
}

func convertUngradedHomeworkSubmission(row db_queries.GetTeacherUngradedHomeworkSubmissionsRow) UngradedHomeworkSubmissionSummary {
	return UngradedHomeworkSubmissionSummary{
		ID:              row.ID,
		HomeworkID:      row.HomeworkID,
		HomeworkTitle:   row.HomeworkTitle,
		Subject:         row.Subject,
		StudentID:       row.StudentID,
		FirstName:       row.FirstName,
		LastName:        row.LastName,
		DueDate:         row.DueDate.Time.Unix(),
		FileName:        row.FileName.String,
		FileType:        row.FileType.String,
		FileSize:        row.FileSize.Int64,
		FileUUID:        row.FileUuid.String(),
		Note:            row.Note.String,
		SubmittedAt:     row.SubmittedAt.Time.Unix(),
		ClassSubjectsID: row.ClassSubjectsID,
	}
}

func convertHomework(row db_queries.GetTeacherHomeworkRow) HomeworkSummary {
	return HomeworkSummary{
		ID:              row.ID,
		ClassSubjectsID: row.ClassSubjectsID,
		Subject:         row.Subject,
		Title:           row.Title,
		HasDescription:  row.Description.Valid,
		Description:     row.Description.String,
		DueDate:         row.DueDate.Time.Unix(),
		CreatedAt:       row.CreatedAt.Time.Unix(),
	}
}

func ListHomeworkSubmissions(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return TeacherSummaryByID(c, pool, rdb, "HOMEWORK_SUBMISSIONS_CACHE_TTL", false,
		func(req ListHomeworkSubmissionsRequest) bool { return req.HomeworkID <= 0 },
		func(req ListHomeworkSubmissionsRequest) int32 { return req.HomeworkID },
		helpers.CacheOrGetTeacherHomeworkSubmissions, convertHomeworkSubmission)
}

func ListUngradedHomeworkSubmissions(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return TeacherSummary(c, pool, rdb, "HOMEWORK_SUBMISSIONS_CACHE_TTL", helpers.CacheOrGetTeacherUngradedHomeworkSubmissions, convertUngradedHomeworkSubmission)
}

func GetHomework(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return TeacherSummaryByID(c, pool, rdb, "HOMEWORK_CACHE_TTL", true,
		func(req GetHomeworkRequest) bool { return req.HomeworkID <= 0 },
		func(req GetHomeworkRequest) int32 { return req.HomeworkID },
		helpers.CacheOrGetTeacherHomework, convertHomework)
}

func AddHomework(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return TeacherModify(c, pool, rdb,
		func(req AddHomeworkRequest) bool {
			return req.Title == "" || (req.HasDescription && req.Description == "") || req.DueDate <= 0 || req.ClassSubjectsID <= 0
		},
		func(ctx context.Context, teacher_scope helpers.TeacherScope, queries *db_queries.Queries, req AddHomeworkRequest) (int64, error) {
			return queries.TeacherAddHomework(ctx, db_queries.TeacherAddHomeworkParams{
				TeacherID:       teacher_scope.TeacherID,
				Title:           req.Title,
				Description:     pgtype.Text{String: req.Description, Valid: req.HasDescription},
				DueDate:         pgtype.Date{Time: time.Unix(req.DueDate, 0), Valid: true},
				AssignedDate:    pgtype.Date{Time: time.Unix(req.AssignedDate, 0), Valid: req.HasAssignedDate},
				ClassSubjectsID: req.ClassSubjectsID,
				SchoolID:        teacher_scope.SchoolID,
			})
		})
}

func RemoveHomework(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return TeacherModify(c, pool, rdb,
		func(req RemoveHomeworkRequest) bool { return req.HomeworkID <= 0 },
		func(ctx context.Context, teacher_scope helpers.TeacherScope, queries *db_queries.Queries, req RemoveHomeworkRequest) (int64, error) {
			return queries.TeacherDeleteHomework(ctx, db_queries.TeacherDeleteHomeworkParams{
				ID:        req.HomeworkID,
				TeacherID: teacher_scope.TeacherID,
				SchoolID:  teacher_scope.SchoolID,
			})
		})
}

func EditHomework(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return TeacherModify(c, pool, rdb,
		func(req EditHomeworkRequest) bool {
			return req.HomeworkID <= 0 || req.Title == "" || (req.HasDescription && req.Description == "") || req.DueDate <= 0 || req.ClassSubjectsID <= 0
		},
		func(ctx context.Context, teacher_scope helpers.TeacherScope, queries *db_queries.Queries, req EditHomeworkRequest) (int64, error) {
			return queries.TeacherEditHomework(ctx, db_queries.TeacherEditHomeworkParams{
				ID:              req.HomeworkID,
				ClassSubjectsID: req.ClassSubjectsID,
				Title:           req.Title,
				Description:     pgtype.Text{String: req.Description, Valid: req.HasDescription},
				DueDate:         pgtype.Date{Time: time.Unix(req.DueDate, 0), Valid: true},
				AssignedDate:    pgtype.Date{Time: time.Unix(req.AssignedDate, 0), Valid: req.HasAssignedDate},
				SchoolID:        teacher_scope.SchoolID,
				TeacherID:       teacher_scope.TeacherID,
			})
		})
}

func GradeHomeworkSubmission(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return TeacherModify(c, pool, rdb,
		func(req GradeHomeworkSubmissionRequest) bool {
			return req.HomeworkID <= 0 || req.StudentID <= 0 || req.GradedValue == 0
		},
		func(ctx context.Context, teacher_scope helpers.TeacherScope, queries *db_queries.Queries, req GradeHomeworkSubmissionRequest) (int64, error) {
			return queries.TeacherGradeHomeworkSubmission(ctx, db_queries.TeacherGradeHomeworkSubmissionParams{
				GradedValue: pgtype.Int2{Int16: req.GradedValue, Valid: true},
				HomeworkID:  req.HomeworkID,
				StudentID:   req.StudentID,
				SchoolID:    teacher_scope.SchoolID,
				TeacherID:   teacher_scope.TeacherID,
			})
		})
}

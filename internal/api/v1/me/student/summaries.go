package studentapi

import (
	"context"
	db_queries "onlystudents/internal/db/store"
	"onlystudents/internal/helpers"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type RequestWithID struct {
	ID int64 `query:"id"`
}
type HomeworkSummary struct {
	ID               int64  `json:"id"`
	Subject          string `json:"subject"`
	HasSubjectCode   bool   `json:"has_subject_code"`
	SubjectCode      string `json:"subject_code"`
	TeacherLastName  string `json:"teacher_last_name"`
	TeacherFirstName string `json:"teacher_first_name"`
	Title            string `json:"title"`
	HasDescription   bool   `json:"has_description"`
	Description      string `json:"description"`
	DueDate          int64  `json:"due_date"`
	CreatedAt        int64  `json:"created_at"`
	SubmittedAt      int64  `json:"submitted_at"`
	HasGradedValue   bool   `json:"has_graded_value"`
	GradedValue      int16  `json:"graded_value"`
}

type AbsenceSummary struct {
	ID                  int64       `json:"id"`
	Subject             string      `json:"subject"`
	HasSubjectCode      bool        `json:"has_subject_code"`
	SubjectCode         string      `json:"subject_code"`
	HasTeacherLastName  bool        `json:"has_teacher_last_name"`
	HasTeacherFirstName bool        `json:"has_teacher_first_name"`
	TeacherLastName     string      `json:"teacher_last_name"`
	TeacherFirstName    string      `json:"teacher_first_name"`
	Date                int64       `json:"date"`
	ReasonType          string      `json:"reason_type"`
	FileName            string      `json:"file_name"`
	FileType            string      `json:"file_type"`
	FileSize            int64       `json:"file_size"`
	FileUUID            string      `json:"file_uuid"`
	ReasonNote          string      `json:"reason_note"`
	Note                string      `json:"note"`
	Justified           bool        `json:"justified"`
	VerifiedBy          interface{} `json:"verified_by"`
}

type ExamSummary struct {
	ID               int64  `json:"id"`
	Subject          string `json:"subject"`
	HasSubjectCode   bool   `json:"has_subject_code"`
	SubjectCode      string `json:"subject_code"`
	TeacherLastName  string `json:"teacher_last_name"`
	TeacherFirstName string `json:"teacher_first_name"`
	Title            string `json:"title"`
	HasDescription   bool   `json:"has_description"`
	Description      string `json:"description"`
	Date             int64  `json:"date"`
	StartTime        int64  `json:"start_time"`
	EndTime          int64  `json:"end_time"`
	HasRoom          bool   `json:"has_room"`
	Room             string `json:"room"`
}

type GradeSummary struct {
	ID               int64  `json:"id"`
	Subject          string `json:"subject"`
	HasSubjectCode   bool   `json:"has_subject_code"`
	SubjectCode      string `json:"subject_code"`
	TeacherLastName  string `json:"teacher_last_name"`
	TeacherFirstName string `json:"teacher_first_name"`
	Term             string `json:"term"`
	Type             string `json:"type"`
	Value            int16  `json:"value"`
	Date             int64  `json:"date"`
	HasNote          bool   `json:"has_note"`
	Note             string `json:"note"`
}

type FinalGradeSummary struct {
	ID               int64  `json:"id"`
	Subject          string `json:"subject"`
	HasSubjectCode   bool   `json:"has_subject_code"`
	SubjectCode      string `json:"subject_code"`
	TeacherLastName  string `json:"teacher_last_name"`
	TeacherFirstName string `json:"teacher_first_name"`
	Term             string `json:"term"`
	Value            int16  `json:"value"`
}

type SubjectSummary struct {
	ID               int32  `json:"id"`
	Name             string `json:"name"`
	Custom           bool   `json:"is_custom"`
	TeacherID        int32  `json:"teacher_id"`
	TeacherFirstName string `json:"teacher_first_name"`
	TeacherLastName  string `json:"teacher_last_name"`
	GroupID          int32  `json:"group_id"`
	GroupName        string `json:"group_name"`
}

func convertGrade(row db_queries.GetStudentGradesRow) GradeSummary {
	return GradeSummary{
		ID:               row.ID,
		Subject:          row.Subject,
		HasSubjectCode:   row.SubjectCode.Valid,
		SubjectCode:      row.SubjectCode.String,
		TeacherLastName:  row.TeacherLastName,
		TeacherFirstName: row.TeacherFirstName,
		Term:             row.Term,
		Type:             row.Type,
		Value:            row.Value,
		Date:             row.Date.Time.Unix(),
		HasNote:          row.Note.Valid,
		Note:             row.Note.String,
	}
}

func convertGradeBySubject(row db_queries.GetStudentGradesBySubjectRow) GradeSummary {
	return GradeSummary{
		ID:               row.ID,
		Subject:          row.Subject,
		HasSubjectCode:   row.SubjectCode.Valid,
		SubjectCode:      row.SubjectCode.String,
		TeacherLastName:  row.TeacherLastName,
		TeacherFirstName: row.TeacherFirstName,
		Term:             row.Term,
		Type:             row.Type,
		Value:            row.Value,
		Date:             row.Date.Time.Unix(),
		HasNote:          row.Note.Valid,
		Note:             row.Note.String,
	}
}

func convertFinalGrade(row db_queries.GetStudentFinalGradesRow) FinalGradeSummary {
	return FinalGradeSummary{
		ID:               row.ID,
		Subject:          row.Subject,
		HasSubjectCode:   row.SubjectCode.Valid,
		SubjectCode:      row.SubjectCode.String,
		TeacherLastName:  row.TeacherLastName,
		TeacherFirstName: row.TeacherFirstName,
		Term:             row.Term,
		Value:            row.Value,
	}
}

func convertAbsence(row db_queries.GetStudentAbsencesRow) AbsenceSummary {
	return AbsenceSummary{
		ID:                  row.ID,
		Subject:             row.Subject,
		HasSubjectCode:      row.SubjectCode.Valid,
		SubjectCode:         row.SubjectCode.String,
		HasTeacherFirstName: row.TeacherFirstName.Valid,
		HasTeacherLastName:  row.TeacherLastName.Valid,
		TeacherFirstName:    row.TeacherFirstName.String,
		TeacherLastName:     row.TeacherLastName.String,
		Date:                row.Date.Time.Unix(),
		ReasonType:          row.ReasonType.String,
		FileName:            row.FileName.String,
		FileType:            row.FileType.String,
		FileSize:            row.FileSize.Int64,
		FileUUID:            row.FileUuid.String(),
		ReasonNote:          row.ReasonNote.String,
		Note:                row.Note.String,
		Justified:           row.Justified,
		VerifiedBy:          row.VerifiedBy,
	}
}

func convertHomework(row db_queries.GetStudentHomeworkRow) HomeworkSummary {
	return HomeworkSummary{
		ID:               row.ID,
		Subject:          row.Subject,
		HasSubjectCode:   row.SubjectCode.Valid,
		SubjectCode:      row.SubjectCode.String,
		TeacherLastName:  row.TeacherLastName,
		TeacherFirstName: row.TeacherFirstName,
		Title:            row.Title,
		HasDescription:   row.Description.Valid,
		Description:      row.Description.String,
		DueDate:          row.DueDate.Time.Unix(),
		CreatedAt:        row.CreatedAt.Time.Unix(),
		SubmittedAt:      row.SubmittedAt.Time.Unix(),
		HasGradedValue:   row.GradedValue.Valid,
		GradedValue:      row.GradedValue.Int16,
	}
}

func convertHomeworkByID(row db_queries.GetStudentHomeworkByIDRow) HomeworkSummary {
	return HomeworkSummary{
		ID:               row.ID,
		Subject:          row.Subject,
		HasSubjectCode:   row.SubjectCode.Valid,
		SubjectCode:      row.SubjectCode.String,
		TeacherLastName:  row.TeacherLastName,
		TeacherFirstName: row.TeacherFirstName,
		Title:            row.Title,
		HasDescription:   row.Description.Valid,
		Description:      row.Description.String,
		DueDate:          row.DueDate.Time.Unix(),
		CreatedAt:        row.CreatedAt.Time.Unix(),
		SubmittedAt:      row.SubmittedAt.Time.Unix(),
		HasGradedValue:   row.GradedValue.Valid,
		GradedValue:      row.GradedValue.Int16,
	}
}

func convertExams(row db_queries.GetStudentExamsRow) ExamSummary {
	return ExamSummary{
		ID:               row.ID,
		Subject:          row.Subject,
		HasSubjectCode:   row.SubjectCode.Valid,
		SubjectCode:      row.SubjectCode.String,
		TeacherLastName:  row.TeacherLastName,
		TeacherFirstName: row.TeacherFirstName,
		Title:            row.Title,
		HasDescription:   row.Description.Valid,
		Description:      row.Description.String,
		Date:             row.Date.Time.Unix(),
		StartTime:        row.StartTime.Microseconds / 1000000,
		EndTime:          row.EndTime.Microseconds / 1000000,
		HasRoom:          row.Room.Valid,
		Room:             row.Room.String,
	}
}

func convertExam(row db_queries.GetStudentExamRow) ExamSummary {
	return ExamSummary{
		ID:               row.ID,
		Subject:          row.Subject,
		HasSubjectCode:   row.SubjectCode.Valid,
		SubjectCode:      row.SubjectCode.String,
		TeacherLastName:  row.TeacherLastName,
		TeacherFirstName: row.TeacherFirstName,
		Title:            row.Title,
		HasDescription:   row.Description.Valid,
		Description:      row.Description.String,
		Date:             row.Date.Time.Unix(),
		StartTime:        row.StartTime.Microseconds / 1000000,
		EndTime:          row.EndTime.Microseconds / 1000000,
		HasRoom:          row.Room.Valid,
		Room:             row.Room.String,
	}
}

func convertSubject(row db_queries.GetStudentSubjectsRow) SubjectSummary {
	return SubjectSummary{
		ID:               row.ID,
		Name:             row.Name,
		Custom:           row.Custom,
		TeacherID:        row.TeacherID,
		TeacherLastName:  row.TeacherLastName,
		TeacherFirstName: row.TeacherFirstName,
		GroupID:          row.GroupID,
		GroupName:        row.GroupName,
	}
}

type subjectArg struct {
	custom bool
	id     int32
}

func Grades(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	if subjects := c.Query("subject"); subjects != "" {
		var wanted []subjectArg

		for arg := range strings.SplitSeq(subjects, ",") {
			if len(arg) < 2 {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
			}

			n, err := strconv.ParseInt(arg[1:], 10, 32)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
			}

			wanted = append(wanted, subjectArg{custom: arg[:1] == "t", id: int32(n)})
		}

		return StudentSummary(c, pool, rdb, "GRADES_CACHE_TTL", false, nil, func(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, req bool, scope helpers.MeScope, ttl int32) ([]db_queries.GetStudentGradesBySubjectRow, error) {
			var grades []db_queries.GetStudentGradesBySubjectRow

			for _, w := range wanted {
				rows, err := helpers.CacheOrGetStudentGradesBySubject(ctx, rdb, queries, scope, ttl, w.custom, w.id)
				if err != nil {
					return nil, err
				}
				grades = append(grades, rows...)
			}

			return grades, nil
		}, convertGradeBySubject)
	}

	return StudentSummary(c, pool, rdb, "GRADES_CACHE_TTL", false, nil, func(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, req bool, scope helpers.MeScope, ttl int32) ([]db_queries.GetStudentGradesRow, error) {
		return helpers.CacheOrGetStudentGrades(ctx, rdb, queries, scope, ttl)
	}, convertGrade)
}

func FinalGrades(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return StudentSummary(c, pool, rdb, "FINAL_GRADES_CACHE_TTL", false, nil, func(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, req bool, scope helpers.MeScope, ttl int32) ([]db_queries.GetStudentFinalGradesRow, error) {
		return helpers.CacheOrGetStudentFinalGrades(ctx, rdb, queries, scope, ttl)
	}, convertFinalGrade)
}

func Absences(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return StudentSummary(c, pool, rdb, "ABSENCES_CACHE_TTL", false, nil, func(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, req bool, scope helpers.MeScope, ttl int32) ([]db_queries.GetStudentAbsencesRow, error) {
		return helpers.CacheOrGetStudentAbsences(ctx, rdb, queries, scope, ttl)
	}, convertAbsence)
}

func Homework(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return StudentSummary(c, pool, rdb, "HOMEWORK_CACHE_TTL", false, nil, func(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, req bool, scope helpers.MeScope, ttl int32) ([]db_queries.GetStudentHomeworkRow, error) {
		return helpers.CacheOrGetStudentHomework(ctx, rdb, queries, scope, ttl)
	}, convertHomework)
}

func HomeworkByID(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return StudentSummary(c, pool, rdb, "HOMEWORK_CACHE_TTL", true, func(req RequestWithID) bool { return req.ID <= 0 }, func(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, req RequestWithID, scope helpers.MeScope, ttl int32) ([]db_queries.GetStudentHomeworkByIDRow, error) {
		return helpers.CacheOrGetStudentHomeworkByID(ctx, rdb, queries, scope, ttl, req.ID)
	}, convertHomeworkByID)
}

func Exams(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return StudentSummary(c, pool, rdb, "EXAMS_CACHE_TTL", false, nil, func(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, req bool, scope helpers.MeScope, ttl int32) ([]db_queries.GetStudentExamsRow, error) {
		return helpers.CacheOrGetStudentExams(ctx, rdb, queries, scope, ttl)
	}, convertExams)
}

func Exam(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return StudentSummary(c, pool, rdb, "EXAMS_CACHE_TTL", true, func(req RequestWithID) bool { return req.ID <= 0 }, func(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, req RequestWithID, scope helpers.MeScope, ttl int32) ([]db_queries.GetStudentExamRow, error) {
		return helpers.CacheOrGetStudentExam(ctx, rdb, queries, scope, ttl, req.ID)
	}, convertExam)
}

func Subjects(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return StudentSummary(c, pool, rdb, "SUBJECTS_CACHE_TTL", false, nil, func(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, req bool, scope helpers.MeScope, ttl int32) ([]db_queries.GetStudentSubjectsRow, error) {
		return helpers.CacheOrGetStudentSubjects(ctx, rdb, queries, scope, ttl)
	}, convertSubject)
}

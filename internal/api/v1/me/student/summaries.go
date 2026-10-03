package studentapi

import (
	db_queries "onlystudents/internal/db/store"
	"onlystudents/internal/helpers"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

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
	Type                string      `json:"type"`
	Justified           bool        `json:"justified"`
	HasNote             bool        `json:"has_note"`
	Note                string      `json:"note"`
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
		Type:                row.Type,
		Justified:           row.Justified,
		HasNote:             row.Note.Valid,
		Note:                row.Note.String,
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

func convertExam(row db_queries.GetStudentExamsRow) ExamSummary {
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

func Grades(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return StudentSummary(c, pool, rdb, "grades", helpers.CacheOrGetStudentGrades, convertGrade)
}

func FinalGrades(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return StudentSummary(c, pool, rdb, "final_grades", helpers.CacheOrGetStudentFinalGrades, convertFinalGrade)
}

func Absences(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return StudentSummary(c, pool, rdb, "absences", helpers.CacheOrGetStudentAbsences, convertAbsence)
}

func Homework(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return StudentSummary(c, pool, rdb, "homework", helpers.CacheOrGetStudentHomework, convertHomework)
}

func Exams(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return StudentSummary(c, pool, rdb, "exams", helpers.CacheOrGetStudentExams, convertExam)
}

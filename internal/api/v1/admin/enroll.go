package adminapi

import (
	"context"
	"fmt"
	"log/slog"
	db_queries "onlystudents/internal/db/store"
	"onlystudents/internal/helpers"
	"onlystudents/internal/mail"
	opaquepkg "onlystudents/internal/opaque"
	"time"

	"github.com/bytemare/opaque"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type EnrollStudentRequest struct {
	// student params
	IDNumber             int32  `json:"id_number"`
	SchoolID             int32  `json:"school_id"`
	HasPhoneNumber       bool   `json:"has_phone_number"`
	PhoneNumber          string `json:"phone_number"`
	FirstName            string `json:"first_name"`
	LastName             string `json:"last_name"`
	BirthFirstName       string `json:"birth_first_name"`
	BirthLastName        string `json:"birth_last_name"`
	BirthDate            int64  `json:"birth_date"`
	BirthCity            string `json:"birth_city"`
	BirthCountry         string `json:"birth_country"`
	MotherBirthFirstName string `json:"mother_birth_first_name"`
	MotherBirthLastName  string `json:"mother_birth_last_name"`
	ClassesID            int32  `json:"classes_id"`
	PermamentAddress     string `json:"permanent_address"`
	TemporaryAddress     string `json:"temporary_address"`
	HasTaxNumber         bool   `json:"has_tax_number"`
	TaxNumber            int32  `json:"tax_number"`
	SsnNumber            int32  `json:"ssn_number"`
	BankName             string `json:"bank_name"`
	IbanOwner            string `json:"iban_owner"`
	IbanNumber           string `json:"iban_number"`
	DocumentType         string `json:"document_type"`
	DocumentNumber       string `json:"document_number"`

	// account params
	EmailAddress string `json:"email_address"`
}

type EnrollTeacherRequest struct {
	// teacher params
	PhoneNumber      string `json:"phone_number"`
	BirthFirstName   string `json:"birth_first_name"`
	BirthLastName    string `json:"birth_last_name"`
	BirthDate        int64  `json:"birth_date"`
	BirthCity        string `json:"birth_city"`
	BirthCountry     string `json:"birth_country"`
	PermamentAddress string `json:"permanent_address"`
	TemporaryAddress string `json:"temporary_address"`
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`

	// account params
	EmailAddress string `json:"email_address"`
}

type EnrollGuardianRequest struct {
	// guardian params
	PhoneNumber      string `json:"phone_number"`
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	BirthFirstName   string `json:"birth_first_name"`
	BirthLastName    string `json:"birth_last_name"`
	BirthDate        int64  `json:"birth_date"`
	BirthCity        string `json:"birth_city"`
	BirthCountry     string `json:"birth_country"`
	PermamentAddress string `json:"permanent_address"`
	TemporaryAddress string `json:"temporary_address"`

	// account params
	EmailAddress string `json:"email_address"`
}

type UserRegistrationEmailData struct {
	EnrollURL string `json:"enroll_url"`
	Name      string `json:"name"`
	AppName   string `json:"app_name"`
}

type MassEnrollRequest struct {
	Guardians []EnrollGuardianRequest `json:"guardians"`
	Students  []EnrollStudentRequest  `json:"students"`
	Teachers  []EnrollTeacherRequest  `json:"teachers"`
}

func SendEnrollToken(rdb *redis.Client, ctx context.Context, enroll_token string, email string, role string, account_uuid uuid.UUID) {
	err := rdb.Set(ctx, fmt.Sprintf("enroll:%s", enroll_token), opaquepkg.EnrollState{
		AccountUUID: account_uuid.String(),
	}, time.Duration(helpers.GetIntEnvFallback("ENROLL_TOKEN_TTL", 168, 30*24))*time.Hour).Err()

	if err != nil {
		return
	}

	// #nosec G118
	go func() {
		// do it in a separate thread to not block
		bg, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		mailer, e := mail.NewFromEnv()
		if e != nil {
			slog.Error("admin enroll mailer init failed", "err", e)
			return
		}
		enroll_url := fmt.Sprintf("%s/enroll?enroll_token=%s", helpers.GetEnvFallback("APP_URL", "http://localhost:8080"), enroll_token)
		if e := mailer.SendTemplateContext(bg, email, fmt.Sprintf("%s Registration", helpers.AppName), "registration.html", UserRegistrationEmailData{EnrollURL: enroll_url, Name: role, AppName: helpers.AppName}); e != nil {
			slog.Error("admin enroll email failed", "account", email, "err", e)
		}
	}()

}

func enroll_student(c fiber.Ctx, req EnrollStudentRequest, pool *pgxpool.Pool, rdb *redis.Client) int {
	enroll_token := opaquepkg.NewEnrollToken()
	account_uuid, err := uuid.NewRandom()

	if err != nil {
		return 500
	}

	queries := db_queries.New(pool)

	student_id, create_student_err := queries.CreateStudent(c.Context(), db_queries.CreateStudentParams{
		IDNumber:             req.IDNumber,
		SchoolID:             req.SchoolID,
		PhoneNumber:          pgtype.Text{String: req.PhoneNumber, Valid: req.HasPhoneNumber},
		FirstName:            req.FirstName,
		LastName:             req.LastName,
		BirthFirstName:       req.BirthFirstName,
		BirthLastName:        req.BirthLastName,
		BirthDate:            pgtype.Date{Time: time.Unix(req.BirthDate, 0), Valid: true},
		BirthCity:            req.BirthCity,
		BirthCountry:         req.BirthCountry,
		MotherBirthFirstName: req.MotherBirthFirstName,
		MotherBirthLastName:  req.MotherBirthLastName,
		ClassesID:            req.ClassesID,
		PermamentAddress:     req.PermamentAddress,
		TemporaryAddress:     req.TemporaryAddress,
		TaxNumber:            pgtype.Int4{Int32: req.TaxNumber, Valid: req.HasTaxNumber},
		SsnNumber:            req.SsnNumber,
		BankName:             req.BankName,
		IbanOwner:            req.IbanOwner,
		IbanNumber:           req.IbanNumber,
		DocumentType:         req.DocumentType,
		DocumentNumber:       req.DocumentNumber,
	})

	if create_student_err != nil {
		return 500
	}

	create_account_err := queries.CreateAccount(c.Context(), db_queries.CreateAccountParams{
		ID:            pgtype.UUID{Bytes: account_uuid, Valid: true},
		Role:          "student",
		StudentID:     pgtype.Int4{Int32: student_id, Valid: true},
		TeacherID:     pgtype.Int4{Valid: false},
		GuardianID:    pgtype.Int4{Valid: false},
		EmailAddress:  pgtype.Text{String: req.EmailAddress, Valid: true},
		EmailVerified: false,
	})

	if create_account_err != nil {
		return 500
	}

	SendEnrollToken(rdb, c.Context(), enroll_token, req.EmailAddress, "student", account_uuid)

	return 200
}

func EnrollStudent(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client, server *opaque.Server) error {
	_, ok := c.Locals("session").(helpers.SessionData)

	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "UNAUTHORIZED"})
	}

	var req EnrollStudentRequest

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	return c.SendStatus(enroll_student(c, req, pool, rdb))
}

func enroll_teacher(c fiber.Ctx, req EnrollTeacherRequest, pool *pgxpool.Pool, rdb *redis.Client) int {
	enroll_token := opaquepkg.NewEnrollToken()
	account_uuid, err := uuid.NewRandom()

	if err != nil {
		return 500
	}

	queries := db_queries.New(pool)

	teacher_id, create_teacher_err := queries.CreateTeacher(c.Context(), db_queries.CreateTeacherParams{
		PhoneNumber:      req.PhoneNumber,
		BirthFirstName:   req.BirthFirstName,
		BirthLastName:    req.BirthLastName,
		BirthDate:        pgtype.Date{Time: time.Unix(req.BirthDate, 0), Valid: true},
		BirthCity:        req.BirthCity,
		BirthCountry:     req.BirthCountry,
		PermamentAddress: req.PermamentAddress,
		TemporaryAddress: req.TemporaryAddress,
		FirstName:        req.FirstName,
		LastName:         req.LastName,
	})

	if create_teacher_err != nil {
		return 500
	}

	create_account_err := queries.CreateAccount(c.Context(), db_queries.CreateAccountParams{
		ID:            pgtype.UUID{Bytes: account_uuid, Valid: true},
		Role:          "teacher",
		StudentID:     pgtype.Int4{Valid: false},
		TeacherID:     pgtype.Int4{Int32: teacher_id, Valid: true},
		GuardianID:    pgtype.Int4{Valid: false},
		EmailAddress:  pgtype.Text{String: req.EmailAddress, Valid: true},
		EmailVerified: false,
	})

	if create_account_err != nil {
		return 500
	}

	SendEnrollToken(rdb, c.Context(), enroll_token, req.EmailAddress, "teacher", account_uuid)

	return 200
}

func EnrollTeacher(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client, server *opaque.Server) error {
	_, ok := c.Locals("session").(helpers.SessionData)

	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "UNAUTHORIZED"})
	}

	var req EnrollTeacherRequest

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	return c.SendStatus(enroll_teacher(c, req, pool, rdb))
}

func enroll_guardian(c fiber.Ctx, req EnrollGuardianRequest, pool *pgxpool.Pool, rdb *redis.Client) int {
	enroll_token := opaquepkg.NewEnrollToken()
	account_uuid, err := uuid.NewRandom()

	if err != nil {
		return 500
	}

	queries := db_queries.New(pool)

	guardian_id, create_guardian_err := queries.CreateGuardian(c.Context(), db_queries.CreateGuardianParams{
		PhoneNumber:      req.PhoneNumber,
		FirstName:        req.FirstName,
		LastName:         req.LastName,
		BirthFirstName:   req.BirthFirstName,
		BirthLastName:    req.BirthLastName,
		BirthDate:        pgtype.Date{Time: time.Unix(req.BirthDate, 0), Valid: true},
		BirthCity:        req.BirthCity,
		BirthCountry:     req.BirthCountry,
		PermamentAddress: req.PermamentAddress,
		TemporaryAddress: req.TemporaryAddress,
	})

	if create_guardian_err != nil {
		return 500
	}

	create_account_err := queries.CreateAccount(c.Context(), db_queries.CreateAccountParams{
		ID:            pgtype.UUID{Bytes: account_uuid, Valid: true},
		Role:          "guardian",
		StudentID:     pgtype.Int4{Valid: false},
		TeacherID:     pgtype.Int4{Valid: false},
		GuardianID:    pgtype.Int4{Int32: guardian_id, Valid: true},
		EmailAddress:  pgtype.Text{String: req.EmailAddress, Valid: true},
		EmailVerified: false,
	})

	if create_account_err != nil {
		return 500
	}

	SendEnrollToken(rdb, c.Context(), enroll_token, req.EmailAddress, "guardian", account_uuid)

	return 200
}

func EnrollGuardian(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client, server *opaque.Server) error {
	_, ok := c.Locals("session").(helpers.SessionData)

	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "UNAUTHORIZED"})
	}

	var req EnrollGuardianRequest

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	return c.SendStatus(enroll_guardian(c, req, pool, rdb))
}

// TODO: rework so this is a process in the background with jobs, this would be way too slow otherwise.
func MassEnroll(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client, server *opaque.Server) error {
	_, ok := c.Locals("session").(helpers.SessionData)

	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "UNAUTHORIZED"})
	}

	var req MassEnrollRequest

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	var student_errors int
	var teacher_errors int
	var guardian_errors int

	for _, student_request := range req.Students {
		code := enroll_student(c, student_request, pool, rdb)

		if code != fiber.StatusOK {
			student_errors += 1
		}
	}

	for _, teacher_request := range req.Teachers {
		code := enroll_teacher(c, teacher_request, pool, rdb)

		if code != fiber.StatusOK {
			teacher_errors += 1
		}
	}

	for _, guardian_request := range req.Guardians {
		code := enroll_guardian(c, guardian_request, pool, rdb)

		if code != fiber.StatusOK {
			guardian_errors += 1
		}
	}

	return c.SendStatus(fiber.StatusOK)
}

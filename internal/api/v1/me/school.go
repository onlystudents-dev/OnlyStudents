package meapi

import (
	"errors"
	"strconv"

	db_queries "onlystudents/internal/db/store"
	"onlystudents/internal/helpers"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type SubjectsData struct {
	ID               int32  `json:"id"`
	Name             string `json:"name"`
	Custom           bool   `json:"is_custom"`
	TeacherID        int32  `json:"teacher_id"`
	TeacherFirstName string `json:"teacher_first_name"`
	TeacherLastName  string `json:"teacher_last_name"`
	GroupID          int32  `json:"group_id"`
	GroupName        string `json:"group_name"`
}

func Subjects(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	sessionToken := c.Cookies("session_token", "")
	if sessionToken == "" {
		return errors.New("no session")
	}

	sessionData, err := helpers.SessionGet(c, rdb, sessionToken)
	if err != nil {
		return errors.New("no session")
	}

	schoolID, err := strconv.ParseInt(c.Get("X-School"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	queries := db_queries.New(pool)

	account, err := helpers.CacheOrGetAccount(c.Context(), rdb, *queries, sessionData.Role, sessionData.AccountID, helpers.GetInt32EnvFallback("ACCOUNT_CACHE_TTL", 5*60, 604800))
	if err != nil {
		return err
	}

	inSchool, err := queries.IsAccountInSchool(c.Context(), db_queries.IsAccountInSchoolParams{
		AccountID: account.ID,
		SchoolID:  int32(schoolID),
	})
	if err != nil {
		return err
	}
	if !inSchool {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "FORBIDDEN"})
	}

	rows, err := queries.GetAccountSubjects(c.Context(), db_queries.GetAccountSubjectsParams{
		AccountID: account.ID,
		SchoolID:  int32(schoolID),
	})
	if err != nil {
		return err
	}

	subjects := make([]SubjectsData, 0, len(rows))
	for _, r := range rows {
		subjects = append(subjects, SubjectsData{
			ID:               r.ID,
			Name:             r.SubjectName,
			Custom:           r.CustomSubject,
			TeacherID:        r.TeacherID,
			TeacherFirstName: r.TeacherFirstName,
			TeacherLastName:  r.TeacherLastName,
			GroupID:          r.GroupID,
			GroupName:        r.GroupName,
		})
	}

	return c.JSON(subjects)
}

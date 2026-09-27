package meapi

import (
	"errors"
	db_queries "onlystudents/internal/db/store"
	"onlystudents/internal/helpers"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type ChildrenData struct {
	AccountID int32  `json:"account_id"`
	SchoolID  int32  `json:"school_id"`
	ClassID   int32  `json:"class_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type StatusData struct {
	Role          string         `json:"role"`
	AccountID     int32          `json:"account_id"`
	SchoolID      int32          `json:"school_id"`
	ClassID       int32          `json:"class_id"`
	FirstName     string         `json:"first_name"`
	LastName      string         `json:"last_name"`
	PfpURL        string         `json:"pfp_url"`
	Nickname      string         `json:"nickname"`
	EmailAddress  string         `json:"email_address"`
	EmailVerified bool           `json:"email_verified"`
	Theme         string         `json:"theme"`
	Lang          string         `json:"lang"`
	Preferences   int32          `json:"preferences"`
	Children      []ChildrenData `json:"children"`
}

func GetStatusData(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) (*StatusData, error) {
	session_token := c.Cookies("session_token", "")

	if session_token == "" {
		return &StatusData{}, errors.New("no session")
	}

	session_data, err := helpers.SessionGet(c, rdb, session_token)

	if err != nil {
		return &StatusData{}, errors.New("no session")
	}

	queries := db_queries.New(pool)

	var status_data StatusData

	account, err := helpers.CacheOrGetAccount(c.Context(), rdb, *queries, session_data.Role, session_data.AccountID, helpers.GetInt32EnvFallback("ACCOUNT_CACHE_TTL", 5*60, 604800))

	if err != nil {
		return &StatusData{}, err
	}

	switch session_data.Role {
	case "student":
		student, err := helpers.CacheOrGetStudent(c.Context(), rdb, *queries, session_data.AccountID, helpers.GetInt32EnvFallback("PERSON_CACHE_TTL", 5*60, 604800))

		if err != nil {
			return &StatusData{}, err
		}

		status_data = StatusData{
			Role:          session_data.Role,
			AccountID:     session_data.AccountID,
			SchoolID:      student.SchoolID,
			FirstName:     student.FirstName,
			LastName:      student.LastName,
			PfpURL:        account.PfpUrl,
			Nickname:      account.Nickname,
			EmailAddress:  account.EmailAddress.String,
			EmailVerified: account.EmailVerified,
			Theme:         account.Theme,
			Lang:          account.Lang,
			Preferences:   account.Preferences,
			Children:      []ChildrenData{},
		}

	case "guardian":
		guardian, err := helpers.CacheOrGetGuardian(c.Context(), rdb, *queries, session_data.AccountID, helpers.GetInt32EnvFallback("PERSON_CACHE_TTL", 5*60, 604800))
		guardian_children, err := helpers.CacheOrGetGuardianChildren(c.Context(), rdb, *queries, session_data.AccountID, helpers.GetInt32EnvFallback("PERSON_CACHE_TTL", 5*60, 604800))

		guardian_children_data := []ChildrenData{}

		for _, children := range guardian_children {
			guardian_children_data = append(guardian_children_data, ChildrenData{
				AccountID: children.ID,
				SchoolID:  children.SchoolID,
				ClassID:   children.ClassID,
				FirstName: children.FirstName,
				LastName:  children.LastName,
			})
		}

		if err != nil {
			return &StatusData{}, err
		}

		status_data = StatusData{
			Role:          session_data.Role,
			AccountID:     session_data.AccountID,
			FirstName:     guardian.FirstName,
			LastName:      guardian.LastName,
			PfpURL:        account.PfpUrl,
			Nickname:      account.Nickname,
			EmailAddress:  account.EmailAddress.String,
			EmailVerified: account.EmailVerified,
			Theme:         account.Theme,
			Lang:          account.Lang,
			Preferences:   account.Preferences,
			Children:      guardian_children_data,
		}

	case "teacher":
		teacher, err := helpers.CacheOrGetTeacher(c.Context(), rdb, *queries, session_data.AccountID, helpers.GetInt32EnvFallback("PERSON_CACHE_TTL", 5*60, 604800))

		if err != nil {
			return &StatusData{}, err
		}

		status_data = StatusData{
			Role:          session_data.Role,
			AccountID:     session_data.AccountID,
			FirstName:     teacher.FirstName,
			LastName:      teacher.LastName,
			PfpURL:        account.PfpUrl,
			Nickname:      account.Nickname,
			EmailAddress:  account.EmailAddress.String,
			EmailVerified: account.EmailVerified,
			Theme:         account.Theme,
			Lang:          account.Lang,
			Preferences:   account.Preferences,
			Children:      []ChildrenData{},
		}

	default:
		return &StatusData{}, errors.New("invalid role")
	}

	return &status_data, nil
}

func Status(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	status_data, err := GetStatusData(c, pool, rdb)

	if err != nil {
		switch err.Error() {
		case "invalid role":
			return c.SendStatus(fiber.StatusBadRequest)
		default:
			return c.SendStatus(fiber.StatusUnauthorized)
		}
	}

	return c.JSON(status_data)
}

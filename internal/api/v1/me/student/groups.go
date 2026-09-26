package studentapi

import (
	"fmt"
	"log/slog"
	db_queries "onlystudents/internal/db/store"
	"onlystudents/internal/helpers"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type ClassData struct {
	ID             int32  `json:"id"`
	SchoolID       int32  `json:"school_id"`
	Name           string `json:"name"`
	TeacherID      int32  `json:"teacher_id"`
	HasCoTeacherID bool   `json:"has_co_teacher_id"`
	CoTeacherID    int32  `json:"co_teacher_id"`
	HasBellID      bool   `json:"has_bell_id"`
	BellID         int32  `json:"bell_id"`
}

type GroupData struct {
	ID        int32  `json:"id"`
	SchoolID  int32  `json:"school_id"`
	BellID    int32  `json:"bell_id"`
	GroupName string `json:"group_name"`
}

func ReadClass(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	scope, err := helpers.ResolveMeScope(c, pool, rdb)

	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	queries := db_queries.New(pool)

	class_data, err := helpers.CacheOrGetStudentClass(c.Context(), rdb, *queries, scope.StudentID, helpers.GetInt32EnvFallback("CLASS_CACHE_TTL", 5*60, 604800))

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	if class_data.SchoolID != scope.SchoolID {
		slog.Warn("Something went extremely wrong. Class school id does not match scope. Report to devs.")
		slog.Warn(fmt.Sprintf("Scope School ID: %d\nClass School ID: %d", scope.SchoolID, class_data.SchoolID))
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(ClassData{
		ID:             class_data.ID,
		SchoolID:       class_data.SchoolID,
		Name:           class_data.Name,
		TeacherID:      class_data.TeacherID,
		HasCoTeacherID: class_data.CoTeacherID.Valid,
		CoTeacherID:    class_data.CoTeacherID.Int32,
		HasBellID:      class_data.BellID.Valid,
		BellID:         class_data.BellID.Int32,
	})
}

func ReadGroups(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	scope, err := helpers.ResolveMeScope(c, pool, rdb)

	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	queries := db_queries.New(pool)

	groups_data, err := helpers.CacheOrGetStudentGroups(c.Context(), rdb, *queries, scope.StudentID, scope.SchoolID, helpers.GetInt32EnvFallback("GROUPS_CACHE_TTL", 5*60, 604800))

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	groups := []GroupData{}

	for _, group_data := range groups_data {
		if group_data.SchoolID != scope.SchoolID {
			slog.Warn("Something went extremely wrong. Group school id does not match scope. Report to devs.")
			slog.Warn(fmt.Sprintf("Scope School ID: %d\nGroup School ID: %d", scope.SchoolID, group_data.SchoolID))
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		group_data_go := GroupData{
			ID:        group_data.ID,
			SchoolID:  group_data.SchoolID,
			BellID:    group_data.BellID,
			GroupName: group_data.GroupName,
		}

		groups = append(groups, group_data_go)
	}

	return c.JSON(groups)
}

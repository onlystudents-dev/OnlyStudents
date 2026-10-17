package studentapi

import (
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

func convertClass(row db_queries.Class) ClassData {
	return ClassData{
		ID:             row.ID,
		SchoolID:       row.SchoolID,
		Name:           row.Name,
		TeacherID:      row.TeacherID,
		HasCoTeacherID: row.CoTeacherID.Valid,
		CoTeacherID:    row.CoTeacherID.Int32,
		HasBellID:      row.BellID.Valid,
		BellID:         row.BellID.Int32,
	}
}

func convertGroup(row db_queries.Group) GroupData {
	return GroupData{
		ID:        row.ID,
		SchoolID:  row.SchoolID,
		BellID:    row.BellID,
		GroupName: row.GroupName,
	}
}

func ReadClass(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return StudentSummarySingle(c, pool, rdb, "CLASS_CACHE_TTL", helpers.CacheOrGetStudentClass, convertClass)
}

func ReadGroups(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return StudentSummary(c, pool, rdb, "GROUPS_CACHE_TTL", helpers.CacheOrGetStudentGroups, convertGroup)
}

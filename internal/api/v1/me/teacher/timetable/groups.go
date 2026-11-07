package timetable

import (
	"context"
	db_queries "onlystudents/internal/db/store"
	"onlystudents/internal/helpers"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type CreateGroupRequest struct {
	BellId int32  `json:"bell_id"`
	Name   string `json:"name"`
}

type DeleteGroupRequest struct {
	Id int32 `json:"id"`
}

type EditGroupRequest struct {
	Id   int32  `json:"id"`
	Name string `json:"name"`
}

type ActionStudentGroupRequest struct {
	GroupID   int32 `json:"group_id"`
	StudentID int32 `json:"student_id"`
}

type ReadStudentGroupRequest struct {
	GroupID int32 `json:"group_id" query:"group_id"`
}

type GroupSummary struct {
	ID        int32  `json:"id"`
	SchoolID  int32  `json:"school_id"`
	BellID    int32  `json:"bell_id"`
	GroupName string `json:"group_name"`
}

type StudentInGroupSummary struct {
	ID        int32  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func convertGroup(row db_queries.Group) GroupSummary {
	return GroupSummary{
		ID:        row.ID,
		SchoolID:  row.SchoolID,
		BellID:    row.BellID,
		GroupName: row.GroupName,
	}
}

func convertStudentsInGroup(row db_queries.ReadListOfStudentsRow) StudentInGroupSummary {
	return StudentInGroupSummary{
		ID:        row.ID,
		FirstName: row.FirstName,
		LastName:  row.LastName,
	}
}

func CreateGroup(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return helpers.TeacherModify(c, pool, rdb, "MANAGE_GROUPS",
		func(req CreateGroupRequest) bool {
			return req.Name == "" || req.BellId <= 0
		},
		func(ctx context.Context, teacher_scope helpers.TeacherScope, queries *db_queries.Queries, req CreateGroupRequest) (int64, error, func()) {
			row_count, err := queries.CreateGroup(ctx, db_queries.CreateGroupParams{
				SchoolID:  teacher_scope.SchoolID,
				BellID:    req.BellId,
				GroupName: req.Name,
			})

			return row_count, err, func() { helpers.InvalidateCachedGroups(ctx, rdb, teacher_scope.SchoolID) }
		})
}

func DeleteGroup(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return helpers.TeacherModify(c, pool, rdb, "MANAGE_GROUPS",
		func(req DeleteGroupRequest) bool {
			return req.Id <= 0
		},
		func(ctx context.Context, teacher_scope helpers.TeacherScope, queries *db_queries.Queries, req DeleteGroupRequest) (int64, error, func()) {
			row_count, err := queries.DeleteGroup(ctx, db_queries.DeleteGroupParams{
				SchoolID: teacher_scope.SchoolID,
				ID:       req.Id,
			})

			return row_count, err, func() { helpers.InvalidateCachedGroups(ctx, rdb, teacher_scope.SchoolID) }
		})
}

func EditGroup(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return helpers.TeacherModify(c, pool, rdb, "MANAGE_GROUPS",
		func(req EditGroupRequest) bool {
			return req.Id <= 0 || req.Name == ""
		},
		func(ctx context.Context, teacher_scope helpers.TeacherScope, queries *db_queries.Queries, req EditGroupRequest) (int64, error, func()) {
			row_count, err := queries.EditGroup(ctx, db_queries.EditGroupParams{
				SchoolID:  teacher_scope.SchoolID,
				ID:        req.Id,
				GroupName: req.Name,
			})

			return row_count, err, func() { helpers.InvalidateCachedGroups(ctx, rdb, teacher_scope.SchoolID) }
		})
}

func ReadGroup(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return helpers.TeacherSummary(c, pool, rdb, "GROUPS_CACHE_TTL", "MANAGE_GROUPS", nil, func(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, req bool, scope helpers.TeacherScope, ttl int32) ([]db_queries.Group, error) {
		return helpers.CacheOrGetGroups(ctx, rdb, queries, scope, ttl)
	}, convertGroup)
}

func InsertStudentToGroup(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return helpers.TeacherModify(c, pool, rdb, "MANAGE_GROUPS",
		func(req ActionStudentGroupRequest) bool {
			return req.GroupID <= 0 || req.StudentID <= 0
		},
		func(ctx context.Context, teacher_scope helpers.TeacherScope, queries *db_queries.Queries, req ActionStudentGroupRequest) (int64, error, func()) {
			row_count, err := queries.InsertStudentToGroup(ctx, db_queries.InsertStudentToGroupParams{
				SchoolID:  teacher_scope.SchoolID,
				GroupID:   req.GroupID,
				StudentID: req.StudentID,
			})

			if row_count == 1 && err == nil {
				helpers.InvalidateCachedStudentGroups(ctx, rdb, helpers.MeScope{
					SchoolID:  teacher_scope.SchoolID,
					StudentID: req.StudentID,
				})
				helpers.InvalidateCachedStudentsInGroup(ctx, rdb, req.GroupID, teacher_scope.SchoolID)
			}

			return row_count, err, func() { helpers.InvalidateCachedGroups(ctx, rdb, teacher_scope.SchoolID) }
		})
}

func DeleteStudentFromGroup(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return helpers.TeacherModify(c, pool, rdb, "MANAGE_GROUPS",
		func(req ActionStudentGroupRequest) bool {
			return req.GroupID <= 0 || req.StudentID <= 0
		},
		func(ctx context.Context, teacher_scope helpers.TeacherScope, queries *db_queries.Queries, req ActionStudentGroupRequest) (int64, error, func()) {
			row_count, err := queries.DeleteStudentFromGroup(ctx, db_queries.DeleteStudentFromGroupParams{
				GroupID:   req.GroupID,
				SchoolID:  teacher_scope.SchoolID,
				StudentID: req.StudentID,
			})

			if row_count == 1 && err == nil {
				helpers.InvalidateCachedStudentGroups(ctx, rdb, helpers.MeScope{
					SchoolID:  teacher_scope.SchoolID,
					StudentID: req.StudentID,
				})
				helpers.InvalidateCachedStudentsInGroup(ctx, rdb, req.GroupID, teacher_scope.SchoolID)
			}

			return row_count, err, func() { helpers.InvalidateCachedGroups(ctx, rdb, teacher_scope.SchoolID) }
		})
}

func ReadStudentsInGroup(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return helpers.TeacherSummary(c, pool, rdb, "GROUP_STUDENTS_CACHE_TTL", "MANAGE_GROUPS", func(req ReadStudentGroupRequest) bool { return req.GroupID <= 0 }, func(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, req ReadStudentGroupRequest, scope helpers.TeacherScope, ttl int32) ([]db_queries.ReadListOfStudentsRow, error) {
		return helpers.CacheOrGetStudentsInGroup(ctx, rdb, queries, req.GroupID, scope, ttl)
	}, convertStudentsInGroup)
}

package timetable

import (
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

func CreateGroup(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req CreateGroupRequest

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	if req.BellId == 0 || req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	scope, status_code := helpers.ResolveTeacherCapabilityScope(c, pool, rdb, "MANAGE_GROUPS")

	if status_code != fiber.StatusOK {
		return c.SendStatus(status_code)
	}

	queries := db_queries.New(pool)

	params := db_queries.CreateGroupParams{
		SchoolID:  scope.SchoolID,
		BellID:    req.BellId,
		GroupName: req.Name,
	}

	err := queries.CreateGroup(c.Context(), params)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "SERVER_ERROR"})
	}

	return c.SendStatus(fiber.StatusOK)
}

func DeleteGroup(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req DeleteGroupRequest

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	if req.Id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	scope, status_code := helpers.ResolveTeacherCapabilityScope(c, pool, rdb, "MANAGE_GROUPS")

	if status_code != fiber.StatusOK {
		return c.SendStatus(status_code)
	}

	queries := db_queries.New(pool)

	params := db_queries.DeleteGroupParams{
		SchoolID: scope.SchoolID,
		ID:       req.Id,
	}

	err := queries.DeleteGroup(c.Context(), params)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "SERVER_ERROR"})
	}

	return c.SendStatus(fiber.StatusOK)
}

func EditGroup(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req EditGroupRequest

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	scope, status_code := helpers.ResolveTeacherCapabilityScope(c, pool, rdb, "MANAGE_GROUPS")

	if status_code != fiber.StatusOK {
		return c.SendStatus(status_code)
	}

	if req.Id == 0 || req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	queries := db_queries.New(pool)

	params := db_queries.EditGroupParams{
		SchoolID:  scope.SchoolID,
		ID:        req.Id,
		GroupName: req.Name,
	}

	err := queries.EditGroup(c.Context(), params)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "SERVER_ERROR"})
	}

	return c.SendStatus(fiber.StatusOK)
}

func ReadGroup(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {

	scope, status_code := helpers.ResolveTeacherCapabilityScope(c, pool, rdb, "MANAGE_GROUPS")

	if status_code != fiber.StatusOK {
		return c.SendStatus(status_code)
	}

	queries := db_queries.New(pool)

	data, err := queries.ReadGroup(c.Context(), scope.SchoolID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "SERVER_ERROR"})
	}

	return c.JSON(data)
}

func InsertStudentToGroup(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req ActionStudentGroupRequest

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	scope, status_code := helpers.ResolveTeacherCapabilityScope(c, pool, rdb, "MANAGE_GROUPS")

	if status_code != fiber.StatusOK {
		return c.SendStatus(status_code)
	}

	if req.GroupID == 0 || req.StudentID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	queries := db_queries.New(pool)

	params := db_queries.InsertStudentToGroupParams{
		SchoolID:  scope.SchoolID,
		GroupID:   req.GroupID,
		StudentID: req.StudentID,
	}

	err := queries.InsertStudentToGroup(c.Context(), params)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "SERVER_ERROR"})
	}

	return c.SendStatus(fiber.StatusOK)
}

func DeleteStudentFromGroup(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req ActionStudentGroupRequest

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	if req.GroupID == 0 || req.StudentID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	scope, status_code := helpers.ResolveTeacherCapabilityScope(c, pool, rdb, "MANAGE_GROUPS")

	if status_code != fiber.StatusOK {
		return c.SendStatus(status_code)
	}

	queries := db_queries.New(pool)

	params := db_queries.DeleteStudentFromGroupParams{
		GroupID:   req.GroupID,
		StudentID: req.StudentID,
		SchoolID:  scope.SchoolID,
	}

	err := queries.DeleteStudentFromGroup(c.Context(), params)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "SERVER_ERROR"})
	}

	return c.SendStatus(fiber.StatusOK)
}

func ReadStudentFromGroup(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req ReadStudentGroupRequest

	if err := c.Bind().Query(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	scope, status_code := helpers.ResolveTeacherCapabilityScope(c, pool, rdb, "MANAGE_GROUPS")

	if status_code != fiber.StatusOK {
		return c.SendStatus(status_code)
	}

	if req.GroupID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	queries := db_queries.New(pool)

	params := db_queries.ReadListOfStudentsParams{
		GroupID:  req.GroupID,
		SchoolID: scope.SchoolID,
	}

	data, err := queries.ReadListOfStudents(c.Context(), params)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "SERVER_ERROR"})
	}

	return c.JSON(data)
}

package timetable

import (
	db_queries "onlystudents/internal/db/store"
	"onlystudents/internal/helpers"
	"strconv"

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
		return c.SendStatus(400)
	}

	school_id, err := strconv.ParseInt(c.Get("X-School"), 10, 32)

	if err != nil {
		return c.SendStatus(400)
	}

	if school_id == 0 {
		return c.SendStatus(400)
	}

	if req.BellId == 0 || req.Name == "" {
		return c.SendStatus(400)
	}

	has_permission := helpers.CheckPermission(c, pool, rdb, "MANAGE_GROUPS", school_id)

	if !has_permission {
		return c.SendStatus(403)
	}

	queries := db_queries.New(pool)

	params := db_queries.CreateGroupParams{
		SchoolID:  int32(school_id),
		BellID:    req.BellId,
		GroupName: req.Name,
	}

	err = queries.CreateGroup(c.Context(), params)

	if err != nil {
		return c.SendStatus(500)
	}

	return c.SendStatus(200)
}

func DeleteGroup(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req DeleteGroupRequest

	if err := c.Bind().Body(&req); err != nil {
		return c.SendStatus(400)
	}

	school_id, err := strconv.ParseInt(c.Get("X-School"), 10, 32)

	if err != nil {
		return c.SendStatus(400)
	}

	if school_id == 0 {
		return c.SendStatus(400)
	}

	if req.Id == 0 {
		return c.SendStatus(400)
	}

	has_permission := helpers.CheckPermission(c, pool, rdb, "MANAGE_GROUPS", school_id)

	if !has_permission {
		return c.SendStatus(403)
	}

	queries := db_queries.New(pool)

	params := db_queries.DeleteGroupParams{
		SchoolID: int32(school_id),
		ID:       req.Id,
	}

	err = queries.DeleteGroup(c.Context(), params)

	if err != nil {
		return c.SendStatus(500)
	}

	return c.SendStatus(200)
}

func EditGroup(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req EditGroupRequest

	if err := c.Bind().Body(&req); err != nil {
		return c.SendStatus(400)
	}

	school_id, err := strconv.ParseInt(c.Get("X-School"), 10, 32)

	if err != nil {
		return c.SendStatus(400)
	}

	if school_id == 0 {
		return c.SendStatus(400)
	}

	if req.Id == 0 || req.Name == "" {
		return c.SendStatus(400)
	}

	has_permission := helpers.CheckPermission(c, pool, rdb, "MANAGE_GROUPS", school_id)

	if !has_permission {
		return c.SendStatus(403)
	}

	queries := db_queries.New(pool)

	params := db_queries.EditGroupParams{
		SchoolID:  int32(school_id),
		ID:        req.Id,
		GroupName: req.Name,
	}

	err = queries.EditGroup(c.Context(), params)

	if err != nil {
		return c.SendStatus(500)
	}

	return c.SendStatus(200)
}

func ReadGroup(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	school_id, err := strconv.ParseInt(c.Get("X-School"), 10, 32)

	if err != nil {
		return c.SendStatus(400)
	}

	if school_id == 0 {
		return c.SendStatus(400)
	}

	has_permission := helpers.CheckPermission(c, pool, rdb, "MANAGE_GROUPS", school_id)

	if !has_permission {
		return c.SendStatus(403)
	}

	queries := db_queries.New(pool)

	data, err := queries.ReadGroup(c.Context(), int32(school_id))

	if err != nil {
		return c.SendStatus(500)
	}

	return c.JSON(data)
}

func InsertStudentToGroup(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req ActionStudentGroupRequest

	if err := c.Bind().Body(&req); err != nil {
		return c.SendStatus(400)
	}

	school_id, err := strconv.ParseInt(c.Get("X-School"), 10, 32)

	if err != nil {
		return c.SendStatus(400)
	}

	if school_id == 0 {
		return c.SendStatus(400)
	}

	if req.GroupID == 0 || req.StudentID == 0 {
		return c.SendStatus(400)
	}

	has_permission := helpers.CheckPermission(c, pool, rdb, "MANAGE_GROUPS", school_id)

	if !has_permission {
		return c.SendStatus(403)
	}

	queries := db_queries.New(pool)

	params := db_queries.InsertStudentToGroupParams{
		SchoolID:  int32(school_id),
		GroupID:   req.GroupID,
		StudentID: req.StudentID,
	}

	err = queries.InsertStudentToGroup(c.Context(), params)

	if err != nil {
		return c.SendStatus(500)
	}

	return c.SendStatus(200)
}

func DeleteStudentFromGroup(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req ActionStudentGroupRequest

	if err := c.Bind().Body(&req); err != nil {
		return c.SendStatus(400)
	}

	school_id, err := strconv.ParseInt(c.Get("X-School"), 10, 32)

	if err != nil {
		return c.SendStatus(400)
	}

	if req.GroupID == 0 || req.StudentID == 0 || school_id == 0 {
		return c.SendStatus(400)
	}

	has_permission := helpers.CheckPermission(c, pool, rdb, "MANAGE_GROUPS", school_id)

	if !has_permission {
		return c.SendStatus(403)
	}

	queries := db_queries.New(pool)

	params := db_queries.DeleteStudentFromGroupParams{
		GroupID:   req.GroupID,
		StudentID: req.StudentID,
		SchoolID:  int32(school_id),
	}

	err = queries.DeleteStudentFromGroup(c.Context(), params)

	if err != nil {
		return c.SendStatus(500)
	}

	return c.SendStatus(200)
}

func ReadStudentFromGroup(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req ReadStudentGroupRequest

	if err := c.Bind().Query(&req); err != nil {
		return c.SendStatus(400)
	}

	school_id, err := strconv.ParseInt(c.Get("X-School"), 10, 32)

	if err != nil {
		return c.SendStatus(400)
	}

	if school_id == 0 {
		return c.SendStatus(400)
	}

	if req.GroupID == 0 {
		return c.SendStatus(400)
	}

	has_permission := helpers.CheckPermission(c, pool, rdb, "MANAGE_GROUPS", school_id)

	if !has_permission {
		return c.SendStatus(403)
	}

	queries := db_queries.New(pool)

	params := db_queries.ReadListOfStudentsParams{
		GroupID:  req.GroupID,
		SchoolID: int32(school_id),
	}

	data, err := queries.ReadListOfStudents(c.Context(), params)

	if err != nil {
		return c.SendStatus(500)
	}

	return c.JSON(data)
}

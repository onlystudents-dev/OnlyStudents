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
	Group_Id   int32 `json:"group_id"`
	Student_Id int32 `json:"student_id"`
}

type ReadStudentGroupRequest struct {
	Group_Id int32 `json:"group_id"`
}

func CreateGroup(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req CreateGroupRequest

	school_id, err := strconv.ParseInt(c.Get("X-School"), 10, 32)

	if err != nil {
		return c.SendStatus(500)
	}

	if school_id == 0 {
		return c.SendStatus(500)
	}

	if req.BellId == 0 || req.Name == "" {
		return c.SendStatus(400)
	}

	has_permission := helpers.CheckPermission(c, pool, rdb, "MANAGE_GROUPS", school_id)

	if !has_permission {
		return c.SendStatus(401)
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

	school_id, err := strconv.ParseInt(c.Get("X-School"), 10, 32)

	if err != nil {
		return c.SendStatus(500)
	}

	if school_id == 0 {
		return c.SendStatus(500)
	}

	if req.Id == 0 {
		return c.SendStatus(400)
	}

	has_permission := helpers.CheckPermission(c, pool, rdb, "MANAGE_GROUPS", school_id)

	if !has_permission {
		return c.SendStatus(401)
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

	school_id, err := strconv.ParseInt(c.Get("X-School"), 10, 32)

	if err != nil {
		return c.SendStatus(500)
	}

	if school_id == 0 {
		return c.SendStatus(500)
	}

	if req.Id == 0 || req.Name == "" {
		return c.SendStatus(400)
	}

	has_permission := helpers.CheckPermission(c, pool, rdb, "MANAGE_GROUPS", school_id)

	if !has_permission {
		return c.SendStatus(401)
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
		return c.SendStatus(500)
	}

	if school_id == 0 {
		return c.SendStatus(500)
	}

	has_permission := helpers.CheckPermission(c, pool, rdb, "MANAGE_GROUPS", school_id)

	if !has_permission {
		return c.SendStatus(401)
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

	school_id, err := strconv.ParseInt(c.Get("X-School"), 10, 32)

	if err != nil {
		return c.SendStatus(500)
	}

	if school_id == 0 {
		return c.SendStatus(500)
	}

	if req.Group_Id == 0 || req.Student_Id == 0 {
		return c.SendStatus(400)
	}

	has_permission := helpers.CheckPermission(c, pool, rdb, "MANAGE_GROUPS", school_id)

	if !has_permission {
		return c.SendStatus(401)
	}

	queries := db_queries.New(pool)

	params := db_queries.InsertStudentToGroupParams{
		SchoolID:  int32(school_id),
		GroupID:   req.Group_Id,
		StudentID: req.Student_Id,
	}

	err = queries.InsertStudentToGroup(c.Context(), params)

	if err != nil {
		return c.SendStatus(500)
	}

	return c.SendStatus(200)
}

func DeleteStudentFromGroup(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req ActionStudentGroupRequest

	school_id, err := strconv.ParseInt(c.Get("X-School"), 10, 32)

	if err != nil {
		return c.SendStatus(500)
	}

	if school_id == 0 {
		return c.SendStatus(500)
	}

	has_permission := helpers.CheckPermission(c, pool, rdb, "MANAGE_GROUPS", school_id)

	if !has_permission {
		return c.SendStatus(401)
	}

	queries := db_queries.New(pool)

	params := db_queries.DeleteStudentFromGroupParams{
		GroupID:   req.Group_Id,
		StudentID: req.Student_Id,
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

	school_id, err := strconv.ParseInt(c.Get("X-School"), 10, 32)

	if err != nil {
		return c.SendStatus(500)
	}

	if school_id == 0 {
		return c.SendStatus(500)
	}

	if req.Group_Id == 0 {
		return c.SendStatus(400)
	}

	has_permission := helpers.CheckPermission(c, pool, rdb, "MANAGE_GROUPS", school_id)

	if !has_permission {
		return c.SendStatus(401)
	}

	queries := db_queries.New(pool)

	params := db_queries.ReadListOfStudentsParams{
		GroupID:  req.Group_Id,
		SchoolID: int32(school_id),
	}

	data, err := queries.ReadListOfStudents(c.Context(), params)

	if err != nil {
		return c.SendStatus(500)
	}

	return c.JSON(data)
}

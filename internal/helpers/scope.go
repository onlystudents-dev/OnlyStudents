package helpers

import (
	"errors"
	"math"
	db_queries "onlystudents/internal/db/store"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type meScope struct {
	SchoolID  int32
	ClassID   int32
	StudentID int32
}

type teacherScope struct {
	SchoolID  int32
	ClassID   int32
	TeacherID int32
}

func CheckPermission(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client, permission string, school_id int64) bool {
	session_data, ok := c.Locals("session").(SessionData)

	if !ok {
		return false
	}

	if school_id < 0 || school_id > math.MaxInt32 {
		return false
	}

	queries := db_queries.New(pool)

	params := db_queries.CheckPermissionParams{
		Name:      permission,
		TeacherID: session_data.AccountID,
		SchoolID:  int32(school_id), // #nosec G115 -- bounds-checked above
	}

	has, err := CacheOrGetCheckPermission(c.Context(), rdb, *queries, params, GetInt32EnvFallback("CHECKPERMISSION_CACHE_TTL", 5*60, 604800))

	if err != nil {
		return false
	}

	return has
}

func ResolveMeScope(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) (meScope, error) {
	session, ok := c.Locals("session").(SessionData)
	if !ok {
		return meScope{}, errors.New("invalid scope")
	}

	queries := db_queries.New(pool)
	ttl := GetInt32EnvFallback("PERSON_CACHE_TTL", 5*60, 604800)

	var studentID int32
	switch session.Role {
	case "student":
		studentID = session.AccountID
	case "guardian":
		id, err := ResolvePerson(c, *queries, session)
		if err != nil {
			return meScope{}, errors.New("invalid scope")
		}
		studentID = id
	default:
		return meScope{}, errors.New("invalid scope")
	}

	student, err := CacheOrGetStudent(c.Context(), rdb, *queries, studentID, ttl)
	if err != nil {
		return meScope{}, errors.New("invalid scope")
	}

	return meScope{StudentID: studentID, SchoolID: student.SchoolID, ClassID: student.ClassesID}, nil
}

func ResolveTeacherCapabilityScope(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client, required_permission string) (teacherScope, int) {
	session, ok := c.Locals("session").(SessionData)
	if !ok {
		return teacherScope{}, fiber.StatusUnauthorized
	}

	var teacherID int32
	switch session.Role {
	case "teacher":
		teacherID = session.AccountID
	default:
		return teacherScope{}, fiber.StatusBadRequest
	}

	school_id, err := strconv.ParseInt(c.Get("X-School"), 10, 32)

	if err != nil {
		return teacherScope{}, fiber.StatusBadRequest
	}

	if school_id == 0 {
		return teacherScope{}, fiber.StatusBadRequest
	}

	has_permission := CheckPermission(c, pool, rdb, required_permission, school_id)

	if !has_permission {
		return teacherScope{}, fiber.StatusForbidden
	}

	return teacherScope{TeacherID: teacherID, SchoolID: int32(school_id)}, fiber.StatusOK
}

func ResolveTeacherScope(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) (teacherScope, int) {
	session, ok := c.Locals("session").(SessionData)
	if !ok {
		return teacherScope{}, fiber.StatusUnauthorized
	}

	var teacherID int32
	switch session.Role {
	case "teacher":
		teacherID = session.AccountID
	default:
		return teacherScope{}, fiber.StatusBadRequest
	}

	school_id, err := strconv.ParseInt(c.Get("X-School"), 10, 32)

	if err != nil {
		return teacherScope{}, fiber.StatusBadRequest
	}

	if school_id == 0 {
		return teacherScope{}, fiber.StatusBadRequest
	}

	queries := db_queries.New(pool)

	is_school_member, err := queries.IsTeacherSchoolMember(c.Context(), db_queries.IsTeacherSchoolMemberParams{
		TeacherID: teacherID,
		SchoolID:  int32(school_id),
	})

	if err != nil || !is_school_member {
		return teacherScope{}, fiber.StatusForbidden
	}

	return teacherScope{TeacherID: teacherID, SchoolID: int32(school_id)}, fiber.StatusOK
}

func ResolvePerson(c fiber.Ctx, queries db_queries.Queries, session_data SessionData) (int32, error) {
	switch session_data.Role {
	case "student":
		return session_data.AccountID, nil
	case "guardian":
		requested_student_id_str := c.Query("student_id")

		requested_student_id, err := strconv.ParseInt(requested_student_id_str, 10, 32)

		if err != nil {
			return 0, err
		}

		can_view_student, err := queries.CanViewStudent(c.Context(), db_queries.CanViewStudentParams{
			GuardianID: session_data.AccountID,
			StudentID:  int32(requested_student_id),
		})

		if err != nil {
			return 0, err
		}

		if can_view_student != 1 {
			return 0, errors.New("No access")
		}

		return int32(requested_student_id), nil

	case "teacher":
		return 0, errors.New("teacher cannot access this")
	default:
		return 0, errors.New("role doesn't exist")
	}
}

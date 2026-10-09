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

type MeScope struct {
	SchoolID  int32
	ClassID   int32
	StudentID int32
}

type TeacherScope struct {
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

func ResolveMeScope(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) (MeScope, error) {
	session, ok := c.Locals("session").(SessionData)
	if !ok {
		return MeScope{}, errors.New("invalid scope")
	}

	queries := db_queries.New(pool)

	var schoolID, classID, studentID int32
	var err error

	switch session.Role {
	case "student", "guardian":
		schoolID, classID, studentID, err = ResolvePerson(c, rdb, *queries, session)
		if err != nil {
			return MeScope{}, errors.New("invalid scope")
		}
	default:
		return MeScope{}, errors.New("invalid scope")
	}

	return MeScope{StudentID: studentID, SchoolID: schoolID, ClassID: classID}, nil
}

func ResolveTeacherCapabilityScope(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client, required_permission string) (TeacherScope, int) {
	session, ok := c.Locals("session").(SessionData)
	if !ok {
		return TeacherScope{}, fiber.StatusUnauthorized
	}

	var teacherID int32
	switch session.Role {
	case "teacher":
		teacherID = session.AccountID
	default:
		return TeacherScope{}, fiber.StatusBadRequest
	}

	school_id, err := strconv.ParseInt(c.Get("X-School"), 10, 32)

	if err != nil || school_id == 0 {
		return TeacherScope{}, fiber.StatusBadRequest
	}

	has_permission := CheckPermission(c, pool, rdb, required_permission, school_id)

	if !has_permission {
		return TeacherScope{}, fiber.StatusForbidden
	}

	return TeacherScope{TeacherID: teacherID, SchoolID: int32(school_id)}, fiber.StatusOK
}

func ResolveTeacherScope(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) (TeacherScope, int) {
	session, ok := c.Locals("session").(SessionData)
	if !ok {
		return TeacherScope{}, fiber.StatusUnauthorized
	}

	var teacherID int32
	switch session.Role {
	case "teacher":
		teacherID = session.AccountID
	default:
		return TeacherScope{}, fiber.StatusBadRequest
	}

	school_id, err := strconv.ParseInt(c.Get("X-School"), 10, 32)

	if err != nil || school_id == 0 {
		return TeacherScope{}, fiber.StatusBadRequest
	}

	queries := db_queries.New(pool)

	is_school_member, err := queries.IsTeacherSchoolMember(c.Context(), db_queries.IsTeacherSchoolMemberParams{
		TeacherID: teacherID,
		SchoolID:  int32(school_id),
	})

	if err != nil || !is_school_member {
		return TeacherScope{}, fiber.StatusForbidden
	}

	return TeacherScope{TeacherID: teacherID, SchoolID: int32(school_id)}, fiber.StatusOK
}

func ResolvePerson(c fiber.Ctx, rdb *redis.Client, queries db_queries.Queries, session_data SessionData) (int32, int32, int32, error) {
	school_id, err := strconv.ParseInt(c.Get("X-School"), 10, 32)

	if err != nil || school_id == 0 || school_id >= math.MaxInt32 {
		return 0, 0, 0, errors.New("invalid school")
	}

	var student_id int32

	switch session_data.Role {
	case "student":
		student_id = session_data.AccountID
	case "guardian":
		requested_student_id_str := c.Get("X-Child")

		requested_student_id, err := strconv.ParseInt(requested_student_id_str, 10, 32)

		if err != nil {
			return 0, 0, 0, err
		}

		can_view_student, err := queries.CanViewStudent(c.Context(), db_queries.CanViewStudentParams{
			GuardianID: session_data.AccountID,
			StudentID:  int32(requested_student_id),
		})

		if err != nil {
			return 0, 0, 0, err
		}

		if can_view_student != 1 {
			return 0, 0, 0, errors.New("No access")
		}

		student_id = int32(requested_student_id)
	case "teacher":
		return 0, 0, 0, errors.New("teacher cannot access this")
	default:
		return 0, 0, 0, errors.New("role doesn't exist")
	}

	ttl := GetInt32EnvFallback("SCHOOL_MEMBERSHIP_CACHE_TTL", 5*60, 604800)

	school_membership, err := CacheOrGetStudentMembership(c.Context(), rdb, queries, int32(school_id), int32(student_id), ttl)

	if err != nil {
		return 0, 0, 0, errors.New("student isnt part of school")
	}

	return school_membership.SchoolID, school_membership.ClassID, school_membership.StudentID, nil
}

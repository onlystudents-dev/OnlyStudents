package helpers

import (
	"context"
	"errors"
	"fmt"
	db_queries "onlystudents/internal/db/store"
	"time"

	"github.com/goccy/go-json"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/redis/go-redis/v9"
)

func CacheOrGet[T any](ctx context.Context, rdb *redis.Client, key string, ttl int32, load func() (T, error)) (T, error) {
	var zero T

	val, err := rdb.Get(ctx, key).Bytes()
	if err == nil {
		var t T
		if uerr := json.Unmarshal(val, &t); uerr == nil {
			return t, nil
		}
	} else if !errors.Is(err, redis.Nil) {
		return load()
	}

	t, err := load()
	if err != nil {
		return zero, err
	}

	b, err := json.Marshal(t)
	if err != nil {
		return zero, err
	}

	if err := rdb.Set(ctx, key, string(b), time.Duration(ttl)*time.Second).Err(); err != nil {
		return zero, err
	}

	return t, nil
}

func CacheOrGetAccount(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, role string, accountID int32, ttl int32) (db_queries.Account, error) {
	key := fmt.Sprintf("%s_account:%d", role, accountID)
	return CacheOrGet(ctx, rdb, key, ttl, func() (db_queries.Account, error) {
		var v pgtype.Int4 = pgtype.Int4{Int32: accountID, Valid: true}
		switch role {
		case "student":
			return queries.GetAccountByStudentID(ctx, v)
		case "guardian":
			return queries.GetAccountByGuardianID(ctx, v)
		case "teacher":
			return queries.GetAccountByTeacherID(ctx, v)
		default:
			return db_queries.Account{}, errors.New("role does not exist")
		}
	})
}

func CacheOrGetStudent(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, accountID int32, ttl int32) (db_queries.Student, error) {
	key := fmt.Sprintf("student:%d", accountID)
	return CacheOrGet(ctx, rdb, key, ttl, func() (db_queries.Student, error) {
		return queries.GetStudent(ctx, accountID)
	})
}

func CacheOrGetTeacher(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, accountID int32, ttl int32) (db_queries.Teacher, error) {
	key := fmt.Sprintf("teacher:%d", accountID)
	return CacheOrGet(ctx, rdb, key, ttl, func() (db_queries.Teacher, error) {
		return queries.GetTeacher(ctx, accountID)
	})
}

func CacheOrGetGuardian(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, accountID int32, ttl int32) (db_queries.Guardian, error) {
	key := fmt.Sprintf("guardian:%d", accountID)
	return CacheOrGet(ctx, rdb, key, ttl, func() (db_queries.Guardian, error) {
		return queries.GetGuardian(ctx, accountID)
	})
}

func CacheOrGetGuardianChildren(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, accountID int32, ttl int32) ([]db_queries.GetGuardianChildrenRow, error) {
	key := fmt.Sprintf("guardian_children:%d", accountID)
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.GetGuardianChildrenRow, error) {
		return queries.GetGuardianChildren(ctx, accountID)
	})
}

func CacheOrGetStudentGrades(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, accountID int32, ttl int32) ([]db_queries.GetStudentGradesRow, error) {
	key := fmt.Sprintf("student_grades:%d", accountID)
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.GetStudentGradesRow, error) {
		return queries.GetStudentGrades(ctx, accountID)
	})
}

func CacheOrGetTeacherGrades(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, accountID int32, schoolID int32, ttl int32) ([]db_queries.GetTeacherGradesRow, error) {
	key := fmt.Sprintf("teacher_grades:%d:%d", schoolID, accountID)
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.GetTeacherGradesRow, error) {
		return queries.GetTeacherGrades(ctx, db_queries.GetTeacherGradesParams{
			TeacherID: accountID,
			SchoolID:  schoolID,
		})
	})
}

func CacheOrGetStudentFinalGrades(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, accountID int32, ttl int32) ([]db_queries.GetStudentFinalGradesRow, error) {
	key := fmt.Sprintf("student_final_grades:%d", accountID)
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.GetStudentFinalGradesRow, error) {
		return queries.GetStudentFinalGrades(ctx, accountID)
	})
}

func CacheOrGetTeacherFinalGrades(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, accountID int32, schoolID int32, ttl int32) ([]db_queries.GetTeacherFinalGradesRow, error) {
	key := fmt.Sprintf("teacher_final_grades:%d:%d", schoolID, accountID)
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.GetTeacherFinalGradesRow, error) {
		return queries.GetTeacherFinalGrades(ctx, db_queries.GetTeacherFinalGradesParams{
			TeacherID: accountID,
			SchoolID:  schoolID,
		})
	})
}

func CacheOrGetStudentExams(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, accountID int32, ttl int32) ([]db_queries.GetStudentExamsRow, error) {
	key := fmt.Sprintf("student_exams:%d", accountID)
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.GetStudentExamsRow, error) {
		return queries.GetStudentExams(ctx, accountID)
	})
}

func CacheOrGetTeacherExams(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, accountID int32, schoolID int32, ttl int32) ([]db_queries.GetTeacherExamsRow, error) {
	key := fmt.Sprintf("teacher_exams:%d:%d", schoolID, accountID)
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.GetTeacherExamsRow, error) {
		return queries.GetTeacherExams(ctx, db_queries.GetTeacherExamsParams{
			TeacherID: accountID,
			SchoolID:  schoolID,
		})
	})
}

func CacheOrGetStudentHomework(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, accountID int32, ttl int32) ([]db_queries.GetStudentHomeworkRow, error) {
	key := fmt.Sprintf("student_homework:%d", accountID)
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.GetStudentHomeworkRow, error) {
		return queries.GetStudentHomework(ctx, accountID)
	})
}

func CacheOrGetTeacherHomework(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, accountID int32, schoolID int32, ttl int32) ([]db_queries.GetTeacherHomeworkRow, error) {
	key := fmt.Sprintf("teacher_homework:%d:%d", schoolID, accountID)
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.GetTeacherHomeworkRow, error) {
		return queries.GetTeacherHomework(ctx, db_queries.GetTeacherHomeworkParams{
			TeacherID: accountID,
			SchoolID:  schoolID,
		})
	})
}

func CacheOrGetTeacherHomeworkSubmissions(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, homeworkID int32, accountID int32, schoolID int32, ttl int32) ([]db_queries.GetTeacherHomeworkSubmissionsRow, error) {
	key := fmt.Sprintf("teacher_homework_submissions:%d:%d:%d", homeworkID, schoolID, accountID)
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.GetTeacherHomeworkSubmissionsRow, error) {
		return queries.GetTeacherHomeworkSubmissions(ctx, db_queries.GetTeacherHomeworkSubmissionsParams{
			HomeworkID: homeworkID,
			TeacherID:  accountID,
			SchoolID:   schoolID,
		})
	})
}

func CacheOrGetCheckPermission(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, check_permission_params db_queries.CheckPermissionParams, ttl int32) (bool, error) {
	key := fmt.Sprintf("check_permission:%s:%d:%d", check_permission_params.Name, check_permission_params.SchoolID, check_permission_params.TeacherID)
	return CacheOrGet(ctx, rdb, key, ttl, func() (bool, error) {
		return queries.CheckPermission(ctx, check_permission_params)
	})
}

func CacheOrGetStudentAbsences(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, accountID int32, ttl int32) ([]db_queries.GetStudentAbsencesRow, error) {
	key := fmt.Sprintf("student_absences:%d", accountID)
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.GetStudentAbsencesRow, error) {
		return queries.GetStudentAbsences(ctx, accountID)
	})
}

func CacheOrGetTeacherAbsences(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, accountID int32, schoolID int32, ttl int32) ([]db_queries.GetTeacherAbsencesRow, error) {
	key := fmt.Sprintf("teacher_absences:%d:%d", schoolID, accountID)
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.GetTeacherAbsencesRow, error) {
		return queries.GetTeacherAbsences(ctx, db_queries.GetTeacherAbsencesParams{
			TeacherID: accountID,
			SchoolID:  schoolID,
		})
	})
}

func CacheOrGetStudentGroups(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, StudentID int32, SchoolID int32, ttl int32) ([]db_queries.Group, error) {
	key := fmt.Sprintf("student_groups:%d:%d", SchoolID, StudentID)
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.Group, error) {
		return queries.GetGroupsByStudentID(ctx, db_queries.GetGroupsByStudentIDParams{
			StudentID: StudentID,
			SchoolID:  SchoolID,
		})
	})
}

func CacheOrGetStudentClass(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, classID int32, ttl int32) (db_queries.Class, error) {
	key := fmt.Sprintf("student_class:%d", classID)
	return CacheOrGet(ctx, rdb, key, ttl, func() (db_queries.Class, error) {
		return queries.GetClassByID(ctx, classID)
	})
}

func InvalidateCachedAccount(ctx context.Context, rdb *redis.Client, role string, account_id int32) {
	switch role {
	case "student":
		//
	case "guardian":
		//
	case "teacher":
		//
	default:
		return
	}

	rdb.Del(ctx, fmt.Sprintf("%s_account:%d", role, account_id))
}

func InvalidateCachedCheckPermission(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, check_permission_params db_queries.CheckPermissionParams, ttl int32) {
	rdb.Del(ctx, fmt.Sprintf("check_permission:%s:%d:%d", check_permission_params.Name, check_permission_params.SchoolID, check_permission_params.TeacherID))
}

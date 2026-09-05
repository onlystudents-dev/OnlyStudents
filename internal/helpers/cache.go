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

func CacheOrGetStudentGrades(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, accountID int32, ttl int32) ([]db_queries.GetStudentGradesRow, error) {
	key := fmt.Sprintf("student_grades:%d", accountID)
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.GetStudentGradesRow, error) {
		return queries.GetStudentGrades(ctx, accountID)
	})
}

func CacheOrGetTeacherGrades(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, accountID int32, ttl int32) ([]db_queries.GetTeacherGradesRow, error) {
	key := fmt.Sprintf("teacher_grades:%d", accountID)
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.GetTeacherGradesRow, error) {
		return queries.GetTeacherGrades(ctx, accountID)
	})
}

func CacheOrGetGuardianGrades(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, accountID int32, ttl int32) ([]db_queries.GetGuardianGradesRow, error) {
	key := fmt.Sprintf("guardian_grades:%d", accountID)
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.GetGuardianGradesRow, error) {
		return queries.GetGuardianGrades(ctx, accountID)
	})
}

func CacheOrGetStudentFinalGrades(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, accountID int32, ttl int32) ([]db_queries.GetStudentFinalGradesRow, error) {
	key := fmt.Sprintf("student_grades:%d", accountID)
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.GetStudentFinalGradesRow, error) {
		return queries.GetStudentFinalGrades(ctx, accountID)
	})
}

func CacheOrGetGuardianFinalGrades(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, accountID int32, ttl int32) ([]db_queries.GetGuardianFinalGradesRow, error) {
	key := fmt.Sprintf("guardian_grades:%d", accountID)
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.GetGuardianFinalGradesRow, error) {
		return queries.GetGuardianFinalGrades(ctx, accountID)
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

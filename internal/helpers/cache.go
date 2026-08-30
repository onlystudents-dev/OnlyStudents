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

type CacheStore struct {
	RedisDB *redis.Client
}

func (cs *CacheStore) CacheOrGetAccount(ctx context.Context, queries db_queries.Queries, role string, account_id int32, ttl int32) (db_queries.Account, error) {
	cached_account, cached := cs.GetAccount(ctx, role, account_id)

	if cached {
		return cached_account, nil
	} else {
		return cs.CacheAccount(ctx, queries, role, account_id, ttl)
	}
}

func (cs *CacheStore) CacheAccount(ctx context.Context, queries db_queries.Queries, role string, account_id int32, ttl int32) (db_queries.Account, error) {
	var account db_queries.Account
	var err error

	pg_user_id := pgtype.Int4{Int32: account_id, Valid: true}

	switch role {
	case "student":
		account, err = queries.GetAccountByStudentID(ctx, pg_user_id)
	case "guardian":
		account, err = queries.GetAccountByGuardianID(ctx, pg_user_id)
	case "teacher":
		account, err = queries.GetAccountByTeacherID(ctx, pg_user_id)
	default:
		return db_queries.Account{}, errors.New("role does not exist")
	}
	if err != nil {
		return db_queries.Account{}, err
	}

	account_value, err := json.Marshal(account)
	if err != nil {
		return db_queries.Account{}, err
	}

	if err := cs.RedisDB.Set(ctx, fmt.Sprintf("%s_account:%d", role, account_id), string(account_value), time.Duration(ttl)*time.Second).Err(); err != nil {
		return db_queries.Account{}, err
	}

	return account, nil
}

func (cs *CacheStore) GetAccount(ctx context.Context, role string, account_id int32) (db_queries.Account, bool) {
	switch role {
	case "student":
		//
	case "guardian":
		//
	case "teacher":
		//
	default:
		return db_queries.Account{}, false
	}

	val, err := cs.RedisDB.Get(ctx, fmt.Sprintf("%s_account:%d", role, account_id)).Bytes()
	if errors.Is(err, redis.Nil) {
		return db_queries.Account{}, false
	}

	if err != nil {
		return db_queries.Account{}, false
	}

	var data db_queries.Account
	if err := json.Unmarshal(val, &data); err != nil {
		return db_queries.Account{}, false
	}

	return data, true
}

func (cs *CacheStore) InvalidateCachedAccount(ctx context.Context, role string, account_id int32) {
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

	cs.RedisDB.Del(ctx, fmt.Sprintf("%s_account:%d", role, account_id))
}

// Student
// TODO: use a generic method
func (cs *CacheStore) CacheOrGetStudent(ctx context.Context, queries db_queries.Queries, student_id int32, ttl int32) (db_queries.Student, error) {
	cached, ok := cs.GetStudent(ctx, student_id)
	if ok {
		return cached, nil
	}

	return cs.CacheStudent(ctx, queries, student_id, ttl)
}

func (cs *CacheStore) CacheStudent(ctx context.Context, queries db_queries.Queries, student_id int32, ttl int32) (db_queries.Student, error) {
	student, err := queries.GetStudent(ctx, student_id)
	if err != nil {
		return db_queries.Student{}, err
	}

	student_value, err := json.Marshal(student)
	if err != nil {
		return db_queries.Student{}, err
	}

	if err := cs.RedisDB.Set(ctx, fmt.Sprintf("student:%d", student_id), string(student_value), time.Duration(ttl)*time.Second).Err(); err != nil {
		return db_queries.Student{}, err
	}

	return student, nil
}

func (cs *CacheStore) GetStudent(ctx context.Context, student_id int32) (db_queries.Student, bool) {
	val, err := cs.RedisDB.Get(ctx, fmt.Sprintf("student:%d", student_id)).Bytes()
	if errors.Is(err, redis.Nil) {
		return db_queries.Student{}, false
	}

	if err != nil {
		return db_queries.Student{}, false
	}

	var data db_queries.Student
	if err := json.Unmarshal(val, &data); err != nil {
		return db_queries.Student{}, false
	}

	return data, true
}

// Teacher
// TODO: use a generic method
func (cs *CacheStore) CacheOrGetTeacher(ctx context.Context, queries db_queries.Queries, teacher_id int32, ttl int32) (db_queries.Teacher, error) {
	cached, ok := cs.GetTeacher(ctx, teacher_id)
	if ok {
		return cached, nil
	}

	return cs.CacheTeacher(ctx, queries, teacher_id, ttl)
}

func (cs *CacheStore) CacheTeacher(ctx context.Context, queries db_queries.Queries, teacher_id int32, ttl int32) (db_queries.Teacher, error) {
	teacher, err := queries.GetTeacher(ctx, teacher_id)
	if err != nil {
		return db_queries.Teacher{}, err
	}

	teacher_value, err := json.Marshal(teacher)
	if err != nil {
		return db_queries.Teacher{}, err
	}

	if err := cs.RedisDB.Set(ctx, fmt.Sprintf("teacher:%d", teacher_id), string(teacher_value), time.Duration(ttl)*time.Second).Err(); err != nil {
		return db_queries.Teacher{}, err
	}

	return teacher, nil
}

func (cs *CacheStore) GetTeacher(ctx context.Context, teacher_id int32) (db_queries.Teacher, bool) {
	val, err := cs.RedisDB.Get(ctx, fmt.Sprintf("teacher:%d", teacher_id)).Bytes()
	if errors.Is(err, redis.Nil) {
		return db_queries.Teacher{}, false
	}

	if err != nil {
		return db_queries.Teacher{}, false
	}

	var data db_queries.Teacher
	if err := json.Unmarshal(val, &data); err != nil {
		return db_queries.Teacher{}, false
	}

	return data, true
}

// Guardian
// TODO: use a generic method
func (cs *CacheStore) CacheOrGetGuardian(ctx context.Context, queries db_queries.Queries, guardian_id int32, ttl int32) (db_queries.Guardian, error) {
	cached, ok := cs.GetGuardian(ctx, guardian_id)
	if ok {
		return cached, nil
	}

	return cs.CacheGuardian(ctx, queries, guardian_id, ttl)
}

func (cs *CacheStore) CacheGuardian(ctx context.Context, queries db_queries.Queries, guardian_id int32, ttl int32) (db_queries.Guardian, error) {
	guardian, err := queries.GetGuardian(ctx, guardian_id)
	if err != nil {
		return db_queries.Guardian{}, err
	}

	guardian_value, err := json.Marshal(guardian)
	if err != nil {
		return db_queries.Guardian{}, err
	}

	if err := cs.RedisDB.Set(ctx, fmt.Sprintf("guardian:%d", guardian_id), string(guardian_value), time.Duration(ttl)*time.Second).Err(); err != nil {
		return db_queries.Guardian{}, err
	}

	return guardian, nil
}

func (cs *CacheStore) GetGuardian(ctx context.Context, guardian_id int32) (db_queries.Guardian, bool) {
	val, err := cs.RedisDB.Get(ctx, fmt.Sprintf("guardian:%d", guardian_id)).Bytes()
	if errors.Is(err, redis.Nil) {
		return db_queries.Guardian{}, false
	}

	if err != nil {
		return db_queries.Guardian{}, false
	}

	var data db_queries.Guardian
	if err := json.Unmarshal(val, &data); err != nil {
		return db_queries.Guardian{}, false
	}

	return data, true
}

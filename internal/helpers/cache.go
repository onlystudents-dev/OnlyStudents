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

// guardian

func CacheOrGetGuardian(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, accountID, ttl int32) (db_queries.Guardian, error) {
	key := fmt.Sprintf("guardian:%d", accountID)
	return CacheOrGet(ctx, rdb, key, ttl, func() (db_queries.Guardian, error) {
		return queries.GetGuardian(ctx, accountID)
	})
}

func CacheOrGetGuardianChildren(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, accountID, ttl int32) ([]db_queries.GetGuardianChildrenRow, error) {
	key := fmt.Sprintf("guardian_children:%d", accountID)
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.GetGuardianChildrenRow, error) {
		return queries.GetGuardianChildren(ctx, accountID)
	})
}

// student

func CacheOrGetStudent(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, scope MeScope, ttl int32) (db_queries.Student, error) {
	key := fmt.Sprintf("student:%d", scope.StudentID)
	return CacheOrGet(ctx, rdb, key, ttl, func() (db_queries.Student, error) {
		return queries.GetStudent(ctx, scope.StudentID)
	})
}

func CacheOrGetStudentFinalGrades(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, scope MeScope, ttl int32) ([]db_queries.GetStudentFinalGradesRow, error) {
	key := fmt.Sprintf("student_final_grades:%d:%d", scope.SchoolID, scope.StudentID)
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.GetStudentFinalGradesRow, error) {
		return queries.GetStudentFinalGrades(ctx, db_queries.GetStudentFinalGradesParams{
			SchoolID:  scope.SchoolID,
			StudentID: scope.StudentID,
		})
	})
}

func CacheOrGetStudentHomework(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, scope MeScope, ttl int32) ([]db_queries.GetStudentHomeworkRow, error) {
	key := fmt.Sprintf("student_homework:%d:%d", scope.SchoolID, scope.StudentID)
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.GetStudentHomeworkRow, error) {
		return queries.GetStudentHomework(ctx, db_queries.GetStudentHomeworkParams{
			SchoolID:  scope.SchoolID,
			StudentID: scope.StudentID,
		})
	})
}

func CacheOrGetStudentMembership(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, scope MeScope, ttl int32) (db_queries.StudentSchool, error) {
	key := fmt.Sprintf("student_membership:%d:%d", scope.SchoolID, scope.StudentID)
	return CacheOrGet(ctx, rdb, key, ttl, func() (db_queries.StudentSchool, error) {
		return queries.GetStudentMembership(ctx, db_queries.GetStudentMembershipParams{
			SchoolID:  scope.SchoolID,
			StudentID: scope.StudentID,
		})
	})
}

func CacheOrListStudentMemberships(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, scope MeScope, ttl int32) ([]db_queries.ListStudentMembershipsRow, error) {
	key := fmt.Sprintf("student_memberships:%d", scope.StudentID)
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.ListStudentMembershipsRow, error) {
		return queries.ListStudentMemberships(ctx, scope.StudentID)
	})
}

func CacheOrGetStudentGrades(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, scope MeScope, ttl int32) ([]db_queries.GetStudentGradesRow, error) {
	key := fmt.Sprintf("student_grades:%d:%d", scope.SchoolID, scope.StudentID)
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.GetStudentGradesRow, error) {
		return queries.GetStudentGrades(ctx, db_queries.GetStudentGradesParams{
			SchoolID:  scope.SchoolID,
			StudentID: scope.StudentID,
		})
	})
}

func CacheOrGetStudentExams(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, scope MeScope, ttl int32) ([]db_queries.GetStudentExamsRow, error) {
	key := fmt.Sprintf("student_exams:%d:%d", scope.SchoolID, scope.StudentID)
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.GetStudentExamsRow, error) {
		return queries.GetStudentExams(ctx, db_queries.GetStudentExamsParams{
			SchoolID:  scope.SchoolID,
			StudentID: scope.StudentID,
		})
	})
}

func CacheOrGetStudentAbsences(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, scope MeScope, ttl int32) ([]db_queries.GetStudentAbsencesRow, error) {
	key := fmt.Sprintf("student_absences:%d", scope.StudentID)
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.GetStudentAbsencesRow, error) {
		return queries.GetStudentAbsences(ctx, db_queries.GetStudentAbsencesParams{
			StudentID: scope.StudentID,
			SchoolID:  scope.SchoolID,
		})
	})
}

func CacheOrGetStudentGroups(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, scope MeScope, ttl int32) ([]db_queries.Group, error) {
	key := fmt.Sprintf("student_groups:%d:%d", scope.SchoolID, scope.StudentID)
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.Group, error) {
		return queries.GetGroupsByStudentID(ctx, db_queries.GetGroupsByStudentIDParams{
			StudentID: scope.StudentID,
			SchoolID:  scope.SchoolID,
		})
	})
}

func CacheOrGetStudentClass(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, scope MeScope, ttl int32) ([]db_queries.Class, error) {
	key := fmt.Sprintf("student_class:%d:%d", scope.SchoolID, scope.StudentID)
	class, err := CacheOrGet(ctx, rdb, key, ttl, func() (db_queries.Class, error) {
		return queries.GetStudentClass(ctx, db_queries.GetStudentClassParams{
			StudentID: scope.StudentID,
			SchoolID:  scope.SchoolID,
		})
	})
	return []db_queries.Class{class}, err
}

func CacheOrGetStudentBaseScheduleClass(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, scope MeScope, ttl int32) ([]db_queries.ReadBaseScheduleClassRow, error) {
	key := fmt.Sprintf("base_schedule_classes:student:%d:%d", scope.SchoolID, scope.ClassID)
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.ReadBaseScheduleClassRow, error) {
		return queries.ReadBaseScheduleClass(ctx, db_queries.ReadBaseScheduleClassParams{
			SchoolID: scope.SchoolID,
			ClassID:  scope.ClassID,
		})
	})
}

func CacheOrGetStudentLessonTime(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, typeID int32, scope MeScope, ttl int32) ([]db_queries.BellSchedule, error) {
	key := fmt.Sprintf("lesson_time:student:%d:%d", scope.SchoolID, typeID)
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.BellSchedule, error) {
		return queries.ReadLessonTime(ctx, db_queries.ReadLessonTimeParams{
			SchoolID: scope.SchoolID,
			TypeID:   typeID,
		})
	})
}

func CacheOrGetStudentRooms(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, scope MeScope, ttl int32) ([]db_queries.ReadRoomRow, error) {
	key := fmt.Sprintf("rooms:%d", scope.SchoolID)
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.ReadRoomRow, error) {
		return queries.ReadRoom(ctx, scope.SchoolID)
	})
}

func CacheOrGetStudentCustomSubjects(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, scope MeScope, ttl int32) ([]db_queries.ReadCustomSubjectRow, error) {
	key := fmt.Sprintf("custom_subjects:%d", scope.SchoolID)
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.ReadCustomSubjectRow, error) {
		return queries.ReadCustomSubject(ctx, scope.SchoolID)
	})
}

func CacheOrGetStudentBellScheduleTypes(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, scope MeScope, ttl int32) ([]db_queries.ReadBellScheduleTypeRow, error) {
	key := fmt.Sprintf("bell_schedule_types:%d", scope.SchoolID)
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.ReadBellScheduleTypeRow, error) {
		return queries.ReadBellScheduleType(ctx, scope.SchoolID)
	})
}

func CacheOrGetTeacher(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, scope TeacherScope, ttl int32) (db_queries.Teacher, error) {
	key := fmt.Sprintf("teacher:%d", scope.TeacherID)
	return CacheOrGet(ctx, rdb, key, ttl, func() (db_queries.Teacher, error) {
		return queries.GetTeacher(ctx, scope.TeacherID)
	})
}

func CacheOrListTeacherMemberships(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, scope TeacherScope, ttl int32) ([]db_queries.ListTeacherSchoolMembershipsRow, error) {
	key := fmt.Sprintf("teacher_memberships:%d", scope.TeacherID)
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.ListTeacherSchoolMembershipsRow, error) {
		return queries.ListTeacherSchoolMemberships(ctx, scope.TeacherID)
	})
}

func CacheOrGetTeacherGrades(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, scope TeacherScope, ttl int32) ([]db_queries.GetTeacherGradesRow, error) {
	key := fmt.Sprintf("teacher_grades:%d:%d", scope.SchoolID, scope.TeacherID)
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.GetTeacherGradesRow, error) {
		return queries.GetTeacherGrades(ctx, db_queries.GetTeacherGradesParams{
			TeacherID: scope.TeacherID,
			SchoolID:  scope.SchoolID,
		})
	})
}

func CacheOrGetTeacherFinalGrades(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, scope TeacherScope, ttl int32) ([]db_queries.GetTeacherFinalGradesRow, error) {
	key := fmt.Sprintf("teacher_final_grades:%d:%d", scope.SchoolID, scope.TeacherID)
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.GetTeacherFinalGradesRow, error) {
		return queries.GetTeacherFinalGrades(ctx, db_queries.GetTeacherFinalGradesParams{
			TeacherID: scope.TeacherID,
			SchoolID:  scope.SchoolID,
		})
	})
}

func CacheOrGetTeacherExams(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, scope TeacherScope, ttl int32) ([]db_queries.GetTeacherExamsRow, error) {
	key := fmt.Sprintf("teacher_exams:%d:%d", scope.SchoolID, scope.TeacherID)
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.GetTeacherExamsRow, error) {
		return queries.GetTeacherExams(ctx, db_queries.GetTeacherExamsParams{
			TeacherID: scope.TeacherID,
			SchoolID:  scope.SchoolID,
		})
	})
}

func CacheOrGetTeacherUngradedHomeworkSubmissions(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, scope TeacherScope, ttl int32) ([]db_queries.GetTeacherUngradedHomeworkSubmissionsRow, error) {
	key := fmt.Sprintf("teacher_ungraded_homework_submissions:%d:%d", scope.SchoolID, scope.TeacherID)
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.GetTeacherUngradedHomeworkSubmissionsRow, error) {
		return queries.GetTeacherUngradedHomeworkSubmissions(ctx, db_queries.GetTeacherUngradedHomeworkSubmissionsParams{
			TeacherID: scope.TeacherID,
			SchoolID:  scope.SchoolID,
		})
	})
}

func CacheOrGetTeacherHomeworkSubmissions(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, homeworkID int32, scope TeacherScope, ttl int32) ([]db_queries.GetTeacherHomeworkSubmissionsRow, error) {
	key := fmt.Sprintf("teacher_homework_submissions:%d:%d:%d", homeworkID, scope.SchoolID, scope.TeacherID)
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.GetTeacherHomeworkSubmissionsRow, error) {
		return queries.GetTeacherHomeworkSubmissions(ctx, db_queries.GetTeacherHomeworkSubmissionsParams{
			HomeworkID: homeworkID,
			TeacherID:  scope.TeacherID,
			SchoolID:   scope.SchoolID,
		})
	})
}

func CacheOrGetTeacherAbsences(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, scope TeacherScope, ttl int32) ([]db_queries.GetTeacherAbsencesRow, error) {
	key := fmt.Sprintf("teacher_absences:%d:%d", scope.SchoolID, scope.TeacherID)
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.GetTeacherAbsencesRow, error) {
		return queries.GetTeacherAbsences(ctx, db_queries.GetTeacherAbsencesParams{
			TeacherID: scope.TeacherID,
			SchoolID:  scope.SchoolID,
		})
	})
}

func CacheOrGetCheckPermission(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, check_permission_params db_queries.CheckPermissionParams, ttl int32) (bool, error) {
	key := fmt.Sprintf("check_permission:%s:%d:%d", check_permission_params.Name, check_permission_params.SchoolID, check_permission_params.TeacherID)
	return CacheOrGet(ctx, rdb, key, ttl, func() (bool, error) {
		return queries.CheckPermission(ctx, check_permission_params)
	})
}

func CacheOrGetTeacherBellScheduleTypes(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, scope TeacherScope, ttl int32) ([]db_queries.ReadBellScheduleTypeRow, error) {
	key := fmt.Sprintf("bell_schedule_types:%d", scope.SchoolID)
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.ReadBellScheduleTypeRow, error) {
		return queries.ReadBellScheduleType(ctx, scope.SchoolID)
	})
}

func CacheOrGetTeacherCustomSubjects(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, scope TeacherScope, ttl int32) ([]db_queries.ReadCustomSubjectRow, error) {
	key := fmt.Sprintf("custom_subjects:%d", scope.SchoolID)
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.ReadCustomSubjectRow, error) {
		return queries.ReadCustomSubject(ctx, scope.SchoolID)
	})
}

func CacheOrGetGroups(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, scope TeacherScope, ttl int32) ([]db_queries.Group, error) {
	key := fmt.Sprintf("groups:%d", scope.SchoolID)
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.Group, error) {
		return queries.ReadGroup(ctx, scope.SchoolID)
	})
}

func CacheOrGetStudentsInGroup(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, groupID int32, scope TeacherScope, ttl int32) ([]db_queries.ReadListOfStudentsRow, error) {
	key := fmt.Sprintf("students_in_group:%d:%d", scope.SchoolID, groupID)
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.ReadListOfStudentsRow, error) {
		return queries.ReadListOfStudents(ctx, db_queries.ReadListOfStudentsParams{
			SchoolID: scope.SchoolID,
			GroupID:  groupID,
		})
	})
}

func CacheOrGetTeacherLessonTime(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, typeID int32, scope TeacherScope, ttl int32) ([]db_queries.BellSchedule, error) {
	key := fmt.Sprintf("lesson_time:teacher:%d:%d", scope.SchoolID, typeID)
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.BellSchedule, error) {
		return queries.ReadLessonTime(ctx, db_queries.ReadLessonTimeParams{
			SchoolID: scope.SchoolID,
			TypeID:   typeID,
		})
	})
}

func CacheOrGetTeacherRooms(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, scope TeacherScope, ttl int32) ([]db_queries.ReadRoomRow, error) {
	key := fmt.Sprintf("rooms:%d", scope.SchoolID)
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.ReadRoomRow, error) {
		return queries.ReadRoom(ctx, scope.SchoolID)
	})
}

func CacheOrGetTeacherBaseScheduleClass(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, classID int32, scope TeacherScope, ttl int32) ([]db_queries.ReadBaseScheduleClassRow, error) {
	key := fmt.Sprintf("base_schedule_classes:teacher:%d:%d", scope.SchoolID, classID)
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.ReadBaseScheduleClassRow, error) {
		return queries.ReadBaseScheduleClass(ctx, db_queries.ReadBaseScheduleClassParams{
			SchoolID: scope.SchoolID,
			ClassID:  classID,
		})
	})
}

func CacheOrGetBaseScheduleGroup(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, groupID int32, scope TeacherScope, ttl int32) ([]db_queries.ReadBaseScheduleGroupRow, error) {
	key := fmt.Sprintf("base_schedule_groups:%d:%d", scope.SchoolID, groupID)
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.ReadBaseScheduleGroupRow, error) {
		return queries.ReadBaseScheduleGroup(ctx, db_queries.ReadBaseScheduleGroupParams{
			SchoolID: scope.SchoolID,
			GroupID:  groupID,
		})
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

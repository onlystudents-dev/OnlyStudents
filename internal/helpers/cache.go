package helpers

import (
	"context"
	"errors"
	"fmt"
	db_queries "onlystudents/internal/db/store"
	"time"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/redis/go-redis/v9"
)

func CacheOrGet[T any](ctx context.Context, rdb *redis.Client, key string, ttl int32, load func() (T, error), tags ...string) (T, error) {
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

	for _, tag := range tags {
		if err := rdb.SAdd(ctx, tag, key).Err(); err != nil {
			continue // tagging is best-effort; never fail the request for it
		}
		rdb.Expire(ctx, tag, time.Duration(ttl)*time.Second)
	}

	return t, nil
}

func InvalidateCachedAccount(ctx context.Context, rdb *redis.Client, role string, accountID int32) {
	key := fmt.Sprintf("%s_account:%d", role, accountID)
	rdb.Del(ctx, key)
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

func InvalidateCachedGuardian(ctx context.Context, rdb *redis.Client, accountID int32) {
	key := fmt.Sprintf("guardian:%d", accountID)
	rdb.Del(ctx, key)
}

func CacheOrGetGuardian(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, accountID, ttl int32) (db_queries.Guardian, error) {
	key := fmt.Sprintf("guardian:%d", accountID)
	return CacheOrGet(ctx, rdb, key, ttl, func() (db_queries.Guardian, error) {
		return queries.GetGuardian(ctx, accountID)
	})
}

func InvalidateCachedGuardianChildren(ctx context.Context, rdb *redis.Client, accountID int32) {
	key := fmt.Sprintf("guardian_children:%d", accountID)
	rdb.Del(ctx, key)
}

func CacheOrGetGuardianChildren(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, accountID, ttl int32) ([]db_queries.GetGuardianChildrenRow, error) {
	key := fmt.Sprintf("guardian_children:%d", accountID)
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.GetGuardianChildrenRow, error) {
		return queries.GetGuardianChildren(ctx, accountID)
	})
}

// student

func InvalidateCachedStudent(ctx context.Context, rdb *redis.Client, scope MeScope) {
	key := fmt.Sprintf("student:%d", scope.StudentID)
	rdb.Del(ctx, key)
}

func CacheOrGetStudent(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, scope MeScope, ttl int32) (db_queries.Student, error) {
	key := fmt.Sprintf("student:%d", scope.StudentID)
	return CacheOrGet(ctx, rdb, key, ttl, func() (db_queries.Student, error) {
		return queries.GetStudent(ctx, scope.StudentID)
	})
}

func InvalidateCachedFinalGrades(ctx context.Context, rdb *redis.Client, schoolID, studentID int32) {
	members := rdb.SMembers(ctx, fmt.Sprintf("tag:final_grades:%d:%d", schoolID, studentID)).Val()
	if len(members) > 0 {
		rdb.Del(ctx, members...)
	}
}

func CacheOrGetStudentFinalGrades(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, scope MeScope, ttl int32) ([]db_queries.GetStudentFinalGradesRow, error) {
	key := fmt.Sprintf("student_final_grades:%d:%d", scope.SchoolID, scope.StudentID)
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.GetStudentFinalGradesRow, error) {
		return queries.GetStudentFinalGrades(ctx, db_queries.GetStudentFinalGradesParams{
			SchoolID:  scope.SchoolID,
			StudentID: scope.StudentID,
		})
	}, fmt.Sprintf("tag:final_grades:%d:%d", scope.SchoolID, scope.StudentID))
}

func InvalidateCachedHomework(ctx context.Context, rdb *redis.Client, schoolID int32) {
	members := rdb.SMembers(ctx, fmt.Sprintf("tag:homework:%d", schoolID)).Val()
	if len(members) > 0 {
		rdb.Del(ctx, members...)
	}
}

func CacheOrGetStudentHomework(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, scope MeScope, ttl int32) ([]db_queries.GetStudentHomeworkRow, error) {
	key := fmt.Sprintf("student_homework:%d:%d", scope.SchoolID, scope.StudentID)
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.GetStudentHomeworkRow, error) {
		return queries.GetStudentHomework(ctx, db_queries.GetStudentHomeworkParams{
			SchoolID:  scope.SchoolID,
			StudentID: scope.StudentID,
		})
	}, fmt.Sprintf("tag:homework:%d", scope.SchoolID))
}

func CacheOrGetStudentHomeworkByID(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, scope MeScope, ttl int32, hwid int64) ([]db_queries.GetStudentHomeworkByIDRow, error) {
	key := fmt.Sprintf("student_homework_byid:%d:%d:%d", scope.SchoolID, scope.StudentID, hwid)
	homework, err := CacheOrGet(ctx, rdb, key, ttl, func() (db_queries.GetStudentHomeworkByIDRow, error) {
		return queries.GetStudentHomeworkByID(ctx, db_queries.GetStudentHomeworkByIDParams{
			ID:        hwid,
			SchoolID:  scope.SchoolID,
			StudentID: scope.StudentID,
		})
	}, fmt.Sprintf("tag:homework:%d", scope.SchoolID))

	return []db_queries.GetStudentHomeworkByIDRow{homework}, err
}

func InvalidateCachedStudentHomework(ctx context.Context, rdb *redis.Client, schoolID, studentID int32, homeworkID int64) {
	rdb.Del(ctx,
		fmt.Sprintf("student_homework:%d:%d", schoolID, studentID),
		fmt.Sprintf("student_homework_byid:%d:%d:%d", schoolID, studentID, homeworkID),
	)
}

func InvalidateCachedStudentMembership(ctx context.Context, rdb *redis.Client, scope MeScope) {
	key := fmt.Sprintf("student_membership:%d:%d", scope.SchoolID, scope.StudentID)
	rdb.Del(ctx, key)
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

func InvalidateCachedGrades(ctx context.Context, rdb *redis.Client, schoolID, studentID int32) {
	members := rdb.SMembers(ctx, fmt.Sprintf("tag:grades:%d:%d", schoolID, studentID)).Val()
	if len(members) > 0 {
		rdb.Del(ctx, members...)
	}
}

func CacheOrGetStudentGrades(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, scope MeScope, ttl int32) ([]db_queries.GetStudentGradesRow, error) {
	key := fmt.Sprintf("student_grades:%d:%d", scope.SchoolID, scope.StudentID)
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.GetStudentGradesRow, error) {
		return queries.GetStudentGrades(ctx, db_queries.GetStudentGradesParams{
			SchoolID:  scope.SchoolID,
			StudentID: scope.StudentID,
		})
	}, fmt.Sprintf("tag:grades:%d:%d", scope.SchoolID, scope.StudentID))
}

func CacheOrGetStudentGradesBySubject(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, scope MeScope, ttl int32, custom bool, sid int32) ([]db_queries.GetStudentGradesBySubjectRow, error) {
	key := fmt.Sprintf("student_grades:%d:%d:%t:%d", scope.SchoolID, scope.StudentID, custom, sid)
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.GetStudentGradesBySubjectRow, error) {
		return queries.GetStudentGradesBySubject(ctx, db_queries.GetStudentGradesBySubjectParams{
			SchoolID:      scope.SchoolID,
			StudentID:     scope.StudentID,
			CustomSubject: custom,
			SubjectID:     sid,
		})
	}, fmt.Sprintf("tag:grades:%d:%d", scope.SchoolID, scope.StudentID))
}

func InvalidateCachedExams(ctx context.Context, rdb *redis.Client, schoolID int32) {
	members := rdb.SMembers(ctx, fmt.Sprintf("tag:exams:%d", schoolID)).Val()
	if len(members) > 0 {
		rdb.Del(ctx, members...)
	}
}

func CacheOrGetStudentExams(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, scope MeScope, ttl int32) ([]db_queries.GetStudentExamsRow, error) {
	key := fmt.Sprintf("exams:%d:%d", scope.SchoolID, scope.ClassID)
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.GetStudentExamsRow, error) {
		return queries.GetStudentExams(ctx, db_queries.GetStudentExamsParams{
			SchoolID:  scope.SchoolID,
			StudentID: scope.StudentID,
		})
	}, fmt.Sprintf("tag:exams:%d", scope.SchoolID))
}

func InvalidateCachedStudentExam(ctx context.Context, rdb *redis.Client, examID int64) {
	key := fmt.Sprintf("exams:%d", examID)
	rdb.Del(ctx, key)
}

func CacheOrGetStudentExam(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, scope MeScope, ttl int32, examid int64) ([]db_queries.GetStudentExamRow, error) {
	key := fmt.Sprintf("exams:%d", examid)
	exam, err := CacheOrGet(ctx, rdb, key, ttl, func() (db_queries.GetStudentExamRow, error) {
		return queries.GetStudentExam(ctx, db_queries.GetStudentExamParams{
			ID:        examid,
			SchoolID:  scope.SchoolID,
			StudentID: scope.StudentID,
		})
	})

	return []db_queries.GetStudentExamRow{exam}, err
}

func InvalidateCachedStudentAbsences(ctx context.Context, rdb *redis.Client, schoolID int32, studentID int32) {
	key := fmt.Sprintf("student_absences:%d:%d", schoolID, studentID)
	rdb.Del(ctx, key)
}

func CacheOrGetStudentAbsences(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, scope MeScope, ttl int32) ([]db_queries.GetStudentAbsencesRow, error) {
	key := fmt.Sprintf("student_absences:%d:%d", scope.SchoolID, scope.StudentID)
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.GetStudentAbsencesRow, error) {
		return queries.GetStudentAbsences(ctx, db_queries.GetStudentAbsencesParams{
			StudentID: scope.StudentID,
			SchoolID:  scope.SchoolID,
		})
	})
}

func InvalidateCachedStudentGroups(ctx context.Context, rdb *redis.Client, scope MeScope) {
	key := fmt.Sprintf("student_groups:%d:%d", scope.SchoolID, scope.StudentID)
	rdb.Del(ctx, key)
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

func InvalidateCachedStudentSubjects(ctx context.Context, rdb *redis.Client, scope MeScope) {
	key := fmt.Sprintf("student_subjects:%d:%d", scope.SchoolID, scope.StudentID)
	rdb.Del(ctx, key)
}

func CacheOrGetStudentSubjects(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, scope MeScope, ttl int32) ([]db_queries.GetStudentSubjectsRow, error) {
	key := fmt.Sprintf("student_subjects:%d:%d", scope.SchoolID, scope.StudentID)
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.GetStudentSubjectsRow, error) {
		return queries.GetStudentSubjects(ctx, db_queries.GetStudentSubjectsParams{
			StudentID: scope.StudentID,
			SchoolID:  scope.SchoolID,
		})
	})
}

func InvalidateCachedStudentClass(ctx context.Context, rdb *redis.Client, scope MeScope) {
	key := fmt.Sprintf("student_class:%d:%d", scope.SchoolID, scope.StudentID)
	rdb.Del(ctx, key)
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

func InvalidateCachedBaseScheduleClass(ctx context.Context, rdb *redis.Client, schoolID, classID int32) {
	key := fmt.Sprintf("base_schedule_classes:%d:%d", schoolID, classID)
	rdb.Del(ctx, key)
}

func CacheOrGetStudentBaseScheduleClass(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, scope MeScope, ttl int32) ([]db_queries.ReadBaseScheduleClassRow, error) {
	key := fmt.Sprintf("base_schedule_classes:%d:%d", scope.SchoolID, scope.ClassID)
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.ReadBaseScheduleClassRow, error) {
		return queries.ReadBaseScheduleClass(ctx, db_queries.ReadBaseScheduleClassParams{
			SchoolID: scope.SchoolID,
			ClassID:  scope.ClassID,
		})
	})
}

func InvalidateCachedLessonTime(ctx context.Context, rdb *redis.Client, typeID int32, schoolID int32) {
	key := fmt.Sprintf("lesson_time:%d:%d", schoolID, typeID)
	rdb.Del(ctx, key)
}

func CacheOrGetLessonTime(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, typeID int32, scope MeScope, ttl int32) ([]db_queries.BellSchedule, error) {
	key := fmt.Sprintf("lesson_time:%d:%d", scope.SchoolID, typeID)
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

func InvalidateCachedTeacher(ctx context.Context, rdb *redis.Client, scope TeacherScope) {
	key := fmt.Sprintf("teacher:%d", scope.TeacherID)
	rdb.Del(ctx, key)
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

func InvalidateCachedTeacherGradesByLesson(ctx context.Context, rdb *redis.Client, scope TeacherScope, id int64) {
	key := fmt.Sprintf("teacher_grades:%d:%d:%d", scope.SchoolID, scope.TeacherID, id)
	rdb.Del(ctx, key)
}

func CacheOrGetTeacherGradesByLesson(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, scope TeacherScope, ttl int32, id int64) ([]db_queries.GetTeacherGradesByLessonRow, error) {
	key := fmt.Sprintf("teacher_grades:%d:%d:%d", scope.SchoolID, scope.TeacherID, id)
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.GetTeacherGradesByLessonRow, error) {
		return queries.GetTeacherGradesByLesson(ctx, db_queries.GetTeacherGradesByLessonParams{
			TeacherID:     scope.TeacherID,
			SchoolID:      scope.SchoolID,
			LessonGradeID: id,
		})
	})
}

func InvalidateCachedTeacherGradeTypes(ctx context.Context, rdb *redis.Client, scope TeacherScope) {
	key := fmt.Sprintf("teacher_grade_types:%d", scope.SchoolID)
	rdb.Del(ctx, key)
}

func CacheOrGetTeacherGradeTypes(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, scope TeacherScope, ttl int32) ([]db_queries.GetTeacherGradeTypesRow, error) {
	key := fmt.Sprintf("teacher_grade_types:%d", scope.SchoolID)
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.GetTeacherGradeTypesRow, error) {
		return queries.GetTeacherGradeTypes(ctx, scope.SchoolID)
	})
}

func InvalidateCachedTeacherLessonGrades(ctx context.Context, rdb *redis.Client, scope TeacherScope) {
	tag := fmt.Sprintf("tag:teacher_lesson_grades:%d", scope.TeacherID)
	members := append(rdb.SMembers(ctx, tag).Val(), fmt.Sprintf("teacher_lesson_grades:%d", scope.TeacherID))
	rdb.Del(ctx, members...)
}

func CacheOrGetTeacherLessonGrades(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, scope TeacherScope, ttl int32) ([]db_queries.GetTeacherLessonGradesRow, error) {
	key := fmt.Sprintf("teacher_lesson_grades:%d", scope.TeacherID)
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.GetTeacherLessonGradesRow, error) {
		return queries.GetTeacherLessonGrades(ctx, db_queries.GetTeacherLessonGradesParams{
			TeacherID: scope.TeacherID,
			SchoolID:  scope.SchoolID,
		})
	}, fmt.Sprintf("tag:teacher_lesson_grades:%d", scope.TeacherID))
}

func InvalidateCachedTeacherLessonGradesBySubject(ctx context.Context, rdb *redis.Client, scope TeacherScope, id int32) {
	key := fmt.Sprintf("teacher_lesson_grades:%d:%d", scope.TeacherID, id)
	rdb.Del(ctx, key)
}

func CacheOrGetTeacherLessonGradesBySubject(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, scope TeacherScope, ttl int32, id int32, date int64) ([]db_queries.GetTeacherLessonGradesBySubjectRow, error) {
	key := fmt.Sprintf("teacher_lesson_grades:%d:%d", scope.TeacherID, id)
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.GetTeacherLessonGradesBySubjectRow, error) {
		return queries.GetTeacherLessonGradesBySubject(ctx, db_queries.GetTeacherLessonGradesBySubjectParams{
			TeacherID:  scope.TeacherID,
			SchoolID:   scope.SchoolID,
			ID:         id,
			ActualDate: date,
		})
	}, fmt.Sprintf("tag:teacher_lesson_grades:%d", scope.TeacherID))
}

func CacheOrGetTeacherExams(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, scope TeacherScope, ttl int32) ([]db_queries.GetTeacherExamsRow, error) {
	key := fmt.Sprintf("teacher_exams:%d:%d", scope.SchoolID, scope.TeacherID)
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.GetTeacherExamsRow, error) {
		return queries.GetTeacherExams(ctx, db_queries.GetTeacherExamsParams{
			TeacherID: scope.TeacherID,
			SchoolID:  scope.SchoolID,
		})
	}, fmt.Sprintf("tag:exams:%d", scope.SchoolID))
}

func InvalidateCachedTeacherUngradedHomeworkSubmissions(ctx context.Context, rdb *redis.Client, scope TeacherScope) {
	key := fmt.Sprintf("teacher_ungraded_homework_submissions:%d:%d", scope.SchoolID, scope.TeacherID)
	rdb.Del(ctx, key)
}

func InvalidateCachedTeacherUngradedHomeworkSubmissionsForSchool(ctx context.Context, rdb *redis.Client, schoolID int32) {
	members := rdb.SMembers(ctx, fmt.Sprintf("tag:ungraded_homework:%d", schoolID)).Val()
	if len(members) > 0 {
		rdb.Del(ctx, members...)
	}
}

func CacheOrGetTeacherUngradedHomeworkSubmissions(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, scope TeacherScope, ttl int32) ([]db_queries.GetTeacherUngradedHomeworkSubmissionsRow, error) {
	key := fmt.Sprintf("teacher_ungraded_homework_submissions:%d:%d", scope.SchoolID, scope.TeacherID)
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.GetTeacherUngradedHomeworkSubmissionsRow, error) {
		return queries.GetTeacherUngradedHomeworkSubmissions(ctx, db_queries.GetTeacherUngradedHomeworkSubmissionsParams{
			TeacherID: scope.TeacherID,
			SchoolID:  scope.SchoolID,
		})
	}, fmt.Sprintf("tag:ungraded_homework:%d", scope.SchoolID))
}

func InvalidateCachedTeacherHomeworkSubmissions(ctx context.Context, rdb *redis.Client, homeworkID int32, schoolID int32) {
	key := fmt.Sprintf("homework_submissions:%d:%d", schoolID, homeworkID)
	rdb.Del(ctx, key)
}

func CacheOrGetTeacherHomeworkSubmissions(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, homeworkID int32, scope TeacherScope, ttl int32) ([]db_queries.GetTeacherHomeworkSubmissionsRow, error) {
	key := fmt.Sprintf("homework_submissions:%d:%d", scope.SchoolID, homeworkID)
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.GetTeacherHomeworkSubmissionsRow, error) {
		return queries.GetTeacherHomeworkSubmissions(ctx, db_queries.GetTeacherHomeworkSubmissionsParams{
			HomeworkID: homeworkID,
			TeacherID:  scope.TeacherID,
			SchoolID:   scope.SchoolID,
		})
	})
}

func InvalidateCachedTeacherSubjects(ctx context.Context, rdb *redis.Client, scope TeacherScope) {
	key := fmt.Sprintf("teacher_subjects:%d:%d", scope.SchoolID, scope.TeacherID)
	rdb.Del(ctx, key)
}

func CacheOrGetTeacherSubjects(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, scope TeacherScope, ttl int32) ([]db_queries.GetTeacherSubjectsRow, error) {
	key := fmt.Sprintf("teacher_subjects:%d:%d", scope.SchoolID, scope.TeacherID)
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.GetTeacherSubjectsRow, error) {
		return queries.GetTeacherSubjects(ctx, db_queries.GetTeacherSubjectsParams{
			TeacherID: scope.TeacherID,
			SchoolID:  scope.SchoolID,
		})
	})
}

func InvalidateCachedTeacherHomework(ctx context.Context, rdb *redis.Client, homeworkID int32, scope TeacherScope) {
	key := fmt.Sprintf("teacher_homework:%d:%d:%d", homeworkID, scope.SchoolID, scope.TeacherID)
	rdb.Del(ctx, key)
}

func CacheOrGetTeacherHomework(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, homeworkID int32, scope TeacherScope, ttl int32) ([]db_queries.GetTeacherHomeworkRow, error) {
	key := fmt.Sprintf("teacher_homework:%d:%d:%d", homeworkID, scope.SchoolID, scope.TeacherID)
	homework, err := CacheOrGet(ctx, rdb, key, ttl, func() (db_queries.GetTeacherHomeworkRow, error) {
		return queries.GetTeacherHomework(ctx, db_queries.GetTeacherHomeworkParams{
			ID:        int64(homeworkID), //nolint:gosec // G115 -- homework_id is bounded to int32 by the request binding
			SchoolID:  scope.SchoolID,
			TeacherID: scope.TeacherID,
		})
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return []db_queries.GetTeacherHomeworkRow{homework}, nil
}

func InvalidateCachedTeacherMyLessonTimes(ctx context.Context, rdb *redis.Client, scope TeacherScope) {
	key := fmt.Sprintf("teacher_lesson_times:%d:%d", scope.SchoolID, scope.TeacherID)
	rdb.Del(ctx, key)
}

func CacheOrGetTeacherMyLessonTimes(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, scope TeacherScope, ttl int32) ([]db_queries.BellSchedule, error) {
	key := fmt.Sprintf("teacher_lesson_times:%d:%d", scope.SchoolID, scope.TeacherID)
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.BellSchedule, error) {
		return queries.ReadTeacherLessonTimes(ctx, db_queries.ReadTeacherLessonTimesParams{
			SchoolID:  scope.SchoolID,
			TeacherID: scope.TeacherID,
		})
	})
}

func InvalidateCachedTeacherAbsencesByLesson(ctx context.Context, rdb *redis.Client, schoolID int32) {
	members := rdb.SMembers(ctx, fmt.Sprintf("tag:teacher_absences:%d", schoolID)).Val()
	if len(members) > 0 {
		rdb.Del(ctx, members...)
	}
}

func CacheOrGetTeacherAbsencesByLesson(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, scope TeacherScope, ttl int32, classSubjectsID int32, date pgtype.Date) ([]db_queries.GetTeacherAbsencesByLessonRow, error) {
	key := fmt.Sprintf("teacher_absences:%d:%d:%d:%d", scope.SchoolID, scope.TeacherID, classSubjectsID, date.Time.Unix())
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.GetTeacherAbsencesByLessonRow, error) {
		return queries.GetTeacherAbsencesByLesson(ctx, db_queries.GetTeacherAbsencesByLessonParams{
			SchoolID:        scope.SchoolID,
			TeacherID:       scope.TeacherID,
			ClassSubjectsID: classSubjectsID,
			Date:            date,
		})
	}, fmt.Sprintf("tag:teacher_absences:%d", scope.SchoolID))
}

func InvalidateCachedTeacherMyGroups(ctx context.Context, rdb *redis.Client, scope TeacherScope) {
	key := fmt.Sprintf("teacher_groups:%d:%d", scope.SchoolID, scope.TeacherID)
	rdb.Del(ctx, key)
}

func CacheOrGetTeacherMyGroups(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, scope TeacherScope, ttl int32) ([]db_queries.Group, error) {
	key := fmt.Sprintf("teacher_groups:%d:%d", scope.SchoolID, scope.TeacherID)
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.Group, error) {
		return queries.ReadMyGroups(ctx, db_queries.ReadMyGroupsParams{
			TeacherID: scope.TeacherID,
			SchoolID:  scope.SchoolID,
		})
	})
}

func InvalidateCachedStudentsInMyGroup(ctx context.Context, rdb *redis.Client, groupID int32, scope TeacherScope) {
	key := fmt.Sprintf("students_in_my_group:%d:%d:%d", scope.SchoolID, scope.TeacherID, groupID)
	rdb.Del(ctx, key)
}

func CacheOrGetStudentsInMyGroup(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, groupID int32, scope TeacherScope, ttl int32) ([]db_queries.ReadStudentsInMyGroupRow, error) {
	key := fmt.Sprintf("students_in_my_group:%d:%d:%d", scope.SchoolID, scope.TeacherID, groupID)
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.ReadStudentsInMyGroupRow, error) {
		return queries.ReadStudentsInMyGroup(ctx, db_queries.ReadStudentsInMyGroupParams{
			SchoolID:  scope.SchoolID,
			TeacherID: scope.TeacherID,
			GroupID:   groupID,
		})
	})
}

func InvalidateCachedCheckPermission(ctx context.Context, rdb *redis.Client, check_permission_params db_queries.CheckPermissionParams) {
	key := fmt.Sprintf("check_permission:%s:%d:%d", check_permission_params.Name, check_permission_params.SchoolID, check_permission_params.TeacherID)
	rdb.Del(ctx, key)
}

func CacheOrGetCheckPermission(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, check_permission_params db_queries.CheckPermissionParams, ttl int32) (bool, error) {
	key := fmt.Sprintf("check_permission:%s:%d:%d", check_permission_params.Name, check_permission_params.SchoolID, check_permission_params.TeacherID)
	return CacheOrGet(ctx, rdb, key, ttl, func() (bool, error) {
		return queries.CheckPermission(ctx, check_permission_params)
	})
}

func InvalidateCachedBellScheduleTypes(ctx context.Context, rdb *redis.Client, schoolID int32) {
	key := fmt.Sprintf("bell_schedule_types:%d", schoolID)
	rdb.Del(ctx, key)
}

func CacheOrGetTeacherBellScheduleTypes(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, scope TeacherScope, ttl int32) ([]db_queries.ReadBellScheduleTypeRow, error) {
	key := fmt.Sprintf("bell_schedule_types:%d", scope.SchoolID)
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.ReadBellScheduleTypeRow, error) {
		return queries.ReadBellScheduleType(ctx, scope.SchoolID)
	})
}

func InvalidateCachedCustomSubjects(ctx context.Context, rdb *redis.Client, schoolID int32) {
	key := fmt.Sprintf("custom_subjects:%d", schoolID)
	rdb.Del(ctx, key)
}

func CacheOrGetTeacherCustomSubjects(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, scope TeacherScope, ttl int32) ([]db_queries.ReadCustomSubjectRow, error) {
	key := fmt.Sprintf("custom_subjects:%d", scope.SchoolID)
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.ReadCustomSubjectRow, error) {
		return queries.ReadCustomSubject(ctx, scope.SchoolID)
	})
}

func InvalidateCachedGroups(ctx context.Context, rdb *redis.Client, schoolID int32) {
	key := fmt.Sprintf("groups:%d", schoolID)
	rdb.Del(ctx, key)
}

func CacheOrGetGroups(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, scope TeacherScope, ttl int32) ([]db_queries.Group, error) {
	key := fmt.Sprintf("groups:%d", scope.SchoolID)
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.Group, error) {
		return queries.ReadGroup(ctx, scope.SchoolID)
	})
}

func InvalidateCachedStudentsInGroup(ctx context.Context, rdb *redis.Client, groupID int32, schoolID int32) {
	key := fmt.Sprintf("students_in_group:%d:%d", schoolID, groupID)
	rdb.Del(ctx, key)
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
	key := fmt.Sprintf("lesson_time:%d:%d", scope.SchoolID, typeID)
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.BellSchedule, error) {
		return queries.ReadLessonTime(ctx, db_queries.ReadLessonTimeParams{
			SchoolID: scope.SchoolID,
			TypeID:   typeID,
		})
	})
}

func InvalidateCachedRooms(ctx context.Context, rdb *redis.Client, schoolID int32) {
	key := fmt.Sprintf("rooms:%d", schoolID)
	rdb.Del(ctx, key)
}

func CacheOrGetTeacherRooms(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, scope TeacherScope, ttl int32) ([]db_queries.ReadRoomRow, error) {
	key := fmt.Sprintf("rooms:%d", scope.SchoolID)
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.ReadRoomRow, error) {
		return queries.ReadRoom(ctx, scope.SchoolID)
	})
}

func CacheOrGetTeacherBaseScheduleClass(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, classID int32, scope TeacherScope, ttl int32) ([]db_queries.ReadBaseScheduleClassRow, error) {
	key := fmt.Sprintf("base_schedule_classes:%d:%d", scope.SchoolID, classID)
	return CacheOrGet(ctx, rdb, key, ttl, func() ([]db_queries.ReadBaseScheduleClassRow, error) {
		return queries.ReadBaseScheduleClass(ctx, db_queries.ReadBaseScheduleClassParams{
			SchoolID: scope.SchoolID,
			ClassID:  classID,
		})
	})
}

func InvalidateCachedBaseScheduleGroup(ctx context.Context, rdb *redis.Client, schoolID int32, groupID int32) {
	key := fmt.Sprintf("base_schedule_groups:%d:%d", schoolID, groupID)
	rdb.Del(ctx, key)
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

func InvalidateCachedFile(ctx context.Context, rdb *redis.Client, file_uuid uuid.UUID) {
	key := fmt.Sprintf("file:%s", file_uuid)
	rdb.Del(ctx, key)
}

func CacheOrGetFile(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, file_uuid uuid.UUID, ttl int32) (db_queries.GetFileRow, error) {
	key := fmt.Sprintf("file:%s", file_uuid)
	return CacheOrGet(ctx, rdb, key, ttl, func() (db_queries.GetFileRow, error) {
		return queries.GetFile(ctx, pgtype.UUID{Bytes: file_uuid, Valid: true})
	})
}

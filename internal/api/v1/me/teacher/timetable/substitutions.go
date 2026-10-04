package timetable

import (
	"context"
	db_queries "onlystudents/internal/db/store"
	"onlystudents/internal/helpers"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type CanceledLessonRequest struct {
	Date            int64 `json:"date"`
	LessonNumber    int32 `json:"lesson_number"`
	TeacherId       int32 `json:"teacher_id"`
	RoomId          int32 `json:"room_id"`
	DayOfWeek       int32 `json:"day_of_week"`
	GroupId         int32 `json:"group_id"`
	IsCustomSubject bool  `json:"is_custom_subject"`
	CustomSubjectId int32 `json:"custom_subject_id"`
	SubjectId       int32 `json:"subject_id"`
}

type SubstitutionsLessonRequest struct {
	ID                    int32 `json:"id"`
	Date                  int64 `json:"date"`
	LessonNumber          int32 `json:"lesson_number"`
	TeacherId             int32 `json:"teacher_id"`
	RoomId                int32 `json:"room_id"`
	DayOfWeek             int32 `json:"day_of_week"`
	GroupId               int32 `json:"group_id"`
	IsCustomSubject       bool  `json:"is_custom_subject"`
	CustomSubjectId       int32 `json:"custom_subject_id"`
	SubjectId             int32 `json:"subject_id"`
	IsSubstitution        bool  `json:"is_substitution"`
	SubstitutionTeacherId int32 `json:"substitution_teacher_id"`
}

func AddCanceledLesson(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return TeacherTimeTableModify(c, pool, rdb, "MANAGE_SUBSTITUTIONS",
		func(req CanceledLessonRequest) bool {
			return req.Date <= 0 || req.LessonNumber <= 0 || req.TeacherId <= 0 || req.RoomId <= 0 || req.DayOfWeek <= 0 || req.GroupId <= 0 || (req.IsCustomSubject == false && req.SubjectId <= 0) || (req.IsCustomSubject == true && req.CustomSubjectId <= 0)
		},
		func(ctx context.Context, teacher_scope helpers.TeacherScope, queries *db_queries.Queries, req CanceledLessonRequest) (int64, error) {
			var custom_subject_id pgtype.Int4
			var subject_id pgtype.Int4

			if req.IsCustomSubject {
				custom_subject_id = pgtype.Int4{Int32: req.CustomSubjectId, Valid: true}
				subject_id = pgtype.Int4{Valid: false}
			} else {
				custom_subject_id = pgtype.Int4{Valid: false}
				subject_id = pgtype.Int4{Int32: req.SubjectId, Valid: true}
			}

			return queries.AddCanceledLesson(c.Context(), db_queries.AddCanceledLessonParams{
				SchoolID:        teacher_scope.SchoolID,
				TeacherID:       req.TeacherId,
				RoomID:          req.RoomId,
				DayOfWeek:       req.DayOfWeek,
				GroupID:         req.GroupId,
				CustomSubject:   req.IsCustomSubject,
				CustomSubjectID: custom_subject_id,
				SubjectID:       subject_id,
				ActualDate:      pgtype.Date{Time: time.Unix(req.Date, 0), Valid: true},
				LessonNum:       req.LessonNumber,
			})
		})
}

func RemoveCanceledLesson(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return TeacherTimeTableModify(c, pool, rdb, "MANAGE_SUBSTITUTIONS",
		func(req CanceledLessonRequest) bool {
			return req.Date <= 0 || req.LessonNumber <= 0 || req.TeacherId <= 0 || req.RoomId <= 0 || req.DayOfWeek <= 0 || req.GroupId <= 0 || (req.IsCustomSubject == false && req.SubjectId <= 0) || (req.IsCustomSubject == true && req.CustomSubjectId <= 0)
		},
		func(ctx context.Context, teacher_scope helpers.TeacherScope, queries *db_queries.Queries, req CanceledLessonRequest) (int64, error) {
			var custom_subject_id pgtype.Int4
			var subject_id pgtype.Int4

			if req.IsCustomSubject {
				custom_subject_id = pgtype.Int4{Int32: req.CustomSubjectId, Valid: true}
				subject_id = pgtype.Int4{Valid: false}
			} else {
				custom_subject_id = pgtype.Int4{Valid: false}
				subject_id = pgtype.Int4{Int32: req.SubjectId, Valid: true}
			}

			return queries.RemoveCanceledLesson(c.Context(), db_queries.RemoveCanceledLessonParams{
				SchoolID:        teacher_scope.SchoolID,
				TeacherID:       req.TeacherId,
				RoomID:          req.RoomId,
				DayOfWeek:       req.DayOfWeek,
				GroupID:         req.GroupId,
				CustomSubject:   req.IsCustomSubject,
				CustomSubjectID: custom_subject_id,
				SubjectID:       subject_id,
				ActualDate:      pgtype.Date{Time: time.Unix(req.Date, 0), Valid: true},
				LessonNum:       req.LessonNumber,
			})
		})
}

func UpdateSubstitution(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return TeacherTimeTableModify(c, pool, rdb, "MANAGE_SUBSTITUTIONS",
		func(req SubstitutionsLessonRequest) bool {
			return req.ID <= 0 || req.Date <= 0 || req.LessonNumber <= 0 || req.TeacherId <= 0 || req.RoomId <= 0 || req.DayOfWeek <= 0 || req.GroupId <= 0 || (req.IsCustomSubject == false && req.SubjectId <= 0) || (req.IsCustomSubject == true && req.CustomSubjectId <= 0 || !req.IsSubstitution || req.SubstitutionTeacherId <= 0)
		},
		func(ctx context.Context, teacher_scope helpers.TeacherScope, queries *db_queries.Queries, req SubstitutionsLessonRequest) (int64, error) {
			var custom_subject_id pgtype.Int4
			var subject_id pgtype.Int4

			if req.IsCustomSubject {
				custom_subject_id = pgtype.Int4{Int32: req.CustomSubjectId, Valid: true}
				subject_id = pgtype.Int4{Valid: false}
			} else {
				custom_subject_id = pgtype.Int4{Valid: false}
				subject_id = pgtype.Int4{Int32: req.SubjectId, Valid: true}
			}

			var substitution_teacher_id pgtype.Int4

			if req.IsSubstitution {
				substitution_teacher_id = pgtype.Int4{Int32: req.SubstitutionTeacherId, Valid: true}
			} else {
				substitution_teacher_id = pgtype.Int4{Valid: false}
			}

			return queries.ManageSubsitutionLesson(c.Context(), db_queries.ManageSubsitutionLessonParams{
				SchoolID:              teacher_scope.SchoolID,
				TeacherID:             req.TeacherId,
				RoomID:                req.RoomId,
				DayOfWeek:             req.DayOfWeek,
				GroupID:               req.GroupId,
				CustomSubject:         req.IsCustomSubject,
				CustomSubjectID:       custom_subject_id,
				SubjectID:             subject_id,
				ActualDate:            pgtype.Date{Time: time.Unix(req.Date, 0), Valid: true},
				LessonNum:             req.LessonNumber,
				IsSubstitution:        req.IsSubstitution,
				SubstitutionTeacherID: substitution_teacher_id,
				ID:                    req.ID,
			})
		})
}

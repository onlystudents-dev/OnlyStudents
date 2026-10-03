package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"
	"time"

	v1 "onlystudents/internal/api/v1"
	adminapi "onlystudents/internal/api/v1/admin"
	meapi "onlystudents/internal/api/v1/me"
	studentapi "onlystudents/internal/api/v1/me/student"
	teacherapi "onlystudents/internal/api/v1/me/teacher"
	timetableapi "onlystudents/internal/api/v1/me/teacher/timetable"
	"onlystudents/internal/db"
	"onlystudents/internal/helpers"
	"onlystudents/internal/middlewares"
	opaquepkg "onlystudents/internal/opaque"

	"github.com/bytemare/opaque"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func buildRoutes(pool *pgxpool.Pool, rdb *redis.Client, opaque_server *opaque.Server) fiber.Map {
	authRateMax := helpers.GetInt64EnvFallback("AUTH_RATE_MAX", 5, 1000000)
	authRateWindow := helpers.GetInt64EnvFallback("AUTH_RATELIMIT_WINDOW", 60, 1000000)
	apiRateMax := helpers.GetInt64EnvFallback("API_RATE_MAX", 60, 1000000)
	apiRateWindow := helpers.GetInt64EnvFallback("API_RATELIMIT_WINDOW", 60, 1000000)
	forgetRateMax := helpers.GetInt64EnvFallback("FORGET_RATE_MAX", 1, 1000000)
	forgetRateWindow := helpers.GetInt64EnvFallback("FORGET_RATELIMIT_WINDOW", 120, 1000000)

	authLimit := func(c fiber.Ctx) error {
		return middlewares.RateLimitMiddleware(c, rdb, "auth", authRateMax, authRateWindow)
	}
	apiLimit := func(c fiber.Ctx) error {
		return middlewares.RateLimitMiddleware(c, rdb, "api", apiRateMax, apiRateWindow)
	}
	forgetLimit := func(c fiber.Ctx) error {
		return middlewares.RateLimitMiddleware(c, rdb, "forget", forgetRateMax, forgetRateWindow)
	}
	auth := func(c fiber.Ctx) error { return middlewares.AuthMiddleware(c, rdb) }
	role := func(roles ...string) fiber.Handler {
		return func(c fiber.Ctx) error { return middlewares.RequireRoleMiddleware(c, rdb, roles) }
	}

	return fiber.Map{
		"config": helpers.Get(func(c fiber.Ctx) error {
			c.Set("Cache-Control", "public, max-age=86400")
			return c.JSON(fiber.Map{
				"password": fiber.Map{
					"minLength":        helpers.GetIntEnvFallback("PW_MIN_LEN", 12, 128),
					"maxLength":        helpers.GetIntEnvFallback("PW_MAX_LEN", 256, 1024),
					"hibpCheckEnabled": helpers.GetEnvFallback("HIBP_CHECK_ENABLED", "false") == "true",
				},
			})
		}),
		"admin_token_login": helpers.Get(func(c fiber.Ctx) error { return adminapi.AdminTokenLogin(c, pool, rdb) }),
		"login": fiber.Map{
			helpers.RoutesGroupMWKey: []fiber.Handler{authLimit},
			"init":                   helpers.Post(func(c fiber.Ctx) error { return v1.LoginInit(c, pool, rdb, opaque_server) }),
			"finish":                 helpers.Post(func(c fiber.Ctx) error { return v1.LoginFinish(c, pool, rdb, opaque_server) }),
		},
		"enroll": fiber.Map{
			helpers.RoutesGroupMWKey: []fiber.Handler{authLimit},
			"init":                   helpers.Post(func(c fiber.Ctx) error { return v1.EnrollInit(c, pool, rdb, opaque_server) }),
			"finish":                 helpers.Post(func(c fiber.Ctx) error { return v1.EnrollFinish(c, pool, rdb, opaque_server) }),
		},
		"forget_password": helpers.Post([]fiber.Handler{
			forgetLimit,
			func(c fiber.Ctx) error { return v1.ForgetPassword(c, pool, rdb) },
		}),
		"forget_password_confirm": helpers.Post([]fiber.Handler{
			forgetLimit,
			func(c fiber.Ctx) error { return v1.ForgetPasswordConfirm(c, pool, rdb) },
		}),

		"v1": fiber.Map{
			"admin": fiber.Map{
				helpers.RoutesGroupMWKey: []fiber.Handler{apiLimit, auth, role("admin")},
				"logs":                   helpers.Get(func(c fiber.Ctx) error { return adminapi.AdminLogs(c, pool, rdb) }),
				"status":                 helpers.Get(func(c fiber.Ctx) error { return adminapi.AdminStatus(c, pool, rdb) }),
				"enroll": fiber.Map{
					"student":  helpers.Post(func(c fiber.Ctx) error { return adminapi.EnrollStudent(c, pool, rdb, opaque_server) }),
					"teacher":  helpers.Post(func(c fiber.Ctx) error { return adminapi.EnrollTeacher(c, pool, rdb, opaque_server) }),
					"guardian": helpers.Post(func(c fiber.Ctx) error { return adminapi.EnrollGuardian(c, pool, rdb, opaque_server) }),
				},
			},
			"me": fiber.Map{
				helpers.RoutesGroupMWKey: []fiber.Handler{apiLimit, auth},
				"logout":                 helpers.Post(func(c fiber.Ctx) error { return v1.Logout(c, pool, rdb) }),
				"update_preferences":     helpers.Post(func(c fiber.Ctx) error { return meapi.UpdatePreferences(c, pool, rdb) }),
				"change_email":           helpers.Post(func(c fiber.Ctx) error { return meapi.ChangeEmail(c, pool, rdb) }),
				"verify_email": helpers.Post([]fiber.Handler{authLimit,
					func(c fiber.Ctx) error { return meapi.VerifyEmailRequest(c, pool, rdb) }}),
				"verify_email_confirm": helpers.Post([]fiber.Handler{authLimit,
					func(c fiber.Ctx) error { return meapi.VerifyEmailConfirm(c, pool, rdb) }}),
				"status": helpers.Get(func(c fiber.Ctx) error { return meapi.Status(c, pool, rdb) }),

				"student": fiber.Map{
					helpers.RoutesGroupMWKey: []fiber.Handler{role("student", "guardian")},
					"grades":                 helpers.Get(func(c fiber.Ctx) error { return studentapi.Grades(c, pool, rdb) }),
					"absences":               helpers.Get(func(c fiber.Ctx) error { return studentapi.Absences(c, pool, rdb) }),
					"final_grades":           helpers.Get(func(c fiber.Ctx) error { return studentapi.FinalGrades(c, pool, rdb) }),
					"exams":                  helpers.Get(func(c fiber.Ctx) error { return studentapi.Exams(c, pool, rdb) }),
					"homework":               helpers.Get(func(c fiber.Ctx) error { return studentapi.Homework(c, pool, rdb) }),
					"class":                  helpers.Get(func(c fiber.Ctx) error { return studentapi.ReadClass(c, pool, rdb) }),
					"groups":                 helpers.Get(func(c fiber.Ctx) error { return studentapi.ReadGroups(c, pool, rdb) }),
					"timetable": fiber.Map{
						fiber.MethodGet:      func(c fiber.Ctx) error { return studentapi.ReadMyRealTimeTable(c, pool, rdb) },
						"base":               helpers.Get(func(c fiber.Ctx) error { return studentapi.ReadBaseSchedule(c, pool, rdb) }),
						"lesson_time":        helpers.Get(func(c fiber.Ctx) error { return studentapi.ReadLessonTime(c, pool, rdb) }),
						"room":               helpers.Get(func(c fiber.Ctx) error { return studentapi.ReadRoom(c, pool, rdb) }),
						"custom_subject":     helpers.Get(func(c fiber.Ctx) error { return studentapi.ReadCustomSubject(c, pool, rdb) }),
						"bell_schedule_type": helpers.Get(func(c fiber.Ctx) error { return studentapi.ReadBellScheduleType(c, pool, rdb) }),
					},
				},

				"teacher": fiber.Map{
					helpers.RoutesGroupMWKey: []fiber.Handler{role("teacher")},

					"grades": fiber.Map{
						fiber.MethodGet:    func(c fiber.Ctx) error { return teacherapi.ListGrades(c, pool, rdb) },
						fiber.MethodPost:   func(c fiber.Ctx) error { return teacherapi.AddGrade(c, pool, rdb) },
						fiber.MethodPatch:  func(c fiber.Ctx) error { return teacherapi.EditGrade(c, pool, rdb) },
						fiber.MethodDelete: func(c fiber.Ctx) error { return teacherapi.RemoveGrade(c, pool, rdb) },
					},

					"final_grades": fiber.Map{
						fiber.MethodGet:   func(c fiber.Ctx) error { return teacherapi.ListFinalGrades(c, pool, rdb) },
						fiber.MethodPost:  func(c fiber.Ctx) error { return teacherapi.AddFinalGrade(c, pool, rdb) },
						fiber.MethodPatch: func(c fiber.Ctx) error { return teacherapi.EditFinalGrade(c, pool, rdb) },
					},

					"exams": fiber.Map{
						fiber.MethodGet:    func(c fiber.Ctx) error { return teacherapi.ListExams(c, pool, rdb) },
						fiber.MethodPost:   func(c fiber.Ctx) error { return teacherapi.AddExam(c, pool, rdb) },
						fiber.MethodPatch:  func(c fiber.Ctx) error { return teacherapi.EditExam(c, pool, rdb) },
						fiber.MethodDelete: func(c fiber.Ctx) error { return teacherapi.RemoveExam(c, pool, rdb) },
					},

					"absences": fiber.Map{
						fiber.MethodGet:    func(c fiber.Ctx) error { return teacherapi.ListAbsences(c, pool, rdb) },
						fiber.MethodPost:   func(c fiber.Ctx) error { return teacherapi.AddAbsence(c, pool, rdb) },
						fiber.MethodPatch:  func(c fiber.Ctx) error { return teacherapi.EditAbsence(c, pool, rdb) },
						fiber.MethodDelete: func(c fiber.Ctx) error { return teacherapi.RemoveAbsence(c, pool, rdb) },
					},

					"homework": fiber.Map{
						fiber.MethodGet:    func(c fiber.Ctx) error { return teacherapi.ListHomework(c, pool, rdb) },
						fiber.MethodPost:   func(c fiber.Ctx) error { return teacherapi.AddHomework(c, pool, rdb) },
						fiber.MethodPatch:  func(c fiber.Ctx) error { return teacherapi.EditHomework(c, pool, rdb) },
						fiber.MethodDelete: func(c fiber.Ctx) error { return teacherapi.RemoveHomework(c, pool, rdb) },
					},

					"homework_submissions": fiber.Map{
						fiber.MethodGet:   func(c fiber.Ctx) error { return teacherapi.ListHomeworkSubmissions(c, pool, rdb) },
						fiber.MethodPatch: func(c fiber.Ctx) error { return teacherapi.UpdateHomeworkSubmission(c, pool, rdb) },
					},

					"timetable": fiber.Map{
						"bell_schedule_type": fiber.Map{
							fiber.MethodPost:   func(c fiber.Ctx) error { return timetableapi.CreateBellScheduleType(c, pool, rdb) },
							fiber.MethodPatch:  func(c fiber.Ctx) error { return timetableapi.EditBellScheduleType(c, pool, rdb) },
							fiber.MethodDelete: func(c fiber.Ctx) error { return timetableapi.DeleteBellScheduleType(c, pool, rdb) },
							fiber.MethodGet:    func(c fiber.Ctx) error { return timetableapi.ReadBellScheduleType(c, pool, rdb) },
						},
						"lesson_time": fiber.Map{
							fiber.MethodPost:   func(c fiber.Ctx) error { return timetableapi.CreateLessonTime(c, pool, rdb) },
							fiber.MethodPatch:  func(c fiber.Ctx) error { return timetableapi.EditLessonTime(c, pool, rdb) },
							fiber.MethodDelete: func(c fiber.Ctx) error { return timetableapi.DeleteLessonTime(c, pool, rdb) },
							fiber.MethodGet:    func(c fiber.Ctx) error { return timetableapi.ReadLessonTime(c, pool, rdb) },
						},
						"custom_subject": fiber.Map{
							fiber.MethodPost:   func(c fiber.Ctx) error { return timetableapi.CreateCustomSubject(c, pool, rdb) },
							fiber.MethodPatch:  func(c fiber.Ctx) error { return timetableapi.EditCustomSubject(c, pool, rdb) },
							fiber.MethodDelete: func(c fiber.Ctx) error { return timetableapi.DeleteCustomSubject(c, pool, rdb) },
							fiber.MethodGet:    func(c fiber.Ctx) error { return timetableapi.ReadCustomSubject(c, pool, rdb) },
						},
						"room": fiber.Map{
							fiber.MethodPost:   func(c fiber.Ctx) error { return timetableapi.CreateRoom(c, pool, rdb) },
							fiber.MethodPatch:  func(c fiber.Ctx) error { return timetableapi.UpdateRoom(c, pool, rdb) },
							fiber.MethodDelete: func(c fiber.Ctx) error { return timetableapi.DeleteRoom(c, pool, rdb) },
							fiber.MethodGet:    func(c fiber.Ctx) error { return timetableapi.ReadRoom(c, pool, rdb) },
						},
						"group": fiber.Map{
							fiber.MethodPost:   func(c fiber.Ctx) error { return timetableapi.CreateGroup(c, pool, rdb) },
							fiber.MethodPatch:  func(c fiber.Ctx) error { return timetableapi.EditGroup(c, pool, rdb) },
							fiber.MethodDelete: func(c fiber.Ctx) error { return timetableapi.DeleteGroup(c, pool, rdb) },
							fiber.MethodGet:    func(c fiber.Ctx) error { return timetableapi.ReadGroup(c, pool, rdb) },
							"student": fiber.Map{
								fiber.MethodPost:   func(c fiber.Ctx) error { return timetableapi.InsertStudentToGroup(c, pool, rdb) },
								fiber.MethodDelete: func(c fiber.Ctx) error { return timetableapi.DeleteStudentFromGroup(c, pool, rdb) },
								fiber.MethodGet:    func(c fiber.Ctx) error { return timetableapi.ReadStudentFromGroup(c, pool, rdb) },
							},
						},
						"base_schedule": fiber.Map{
							fiber.MethodPost:   func(c fiber.Ctx) error { return timetableapi.CreateBaseSchedule(c, pool, rdb) },
							fiber.MethodPatch:  func(c fiber.Ctx) error { return timetableapi.UpdateBaseSchedule(c, pool, rdb) },
							fiber.MethodDelete: func(c fiber.Ctx) error { return timetableapi.DeleteBaseSchedule(c, pool, rdb) },
							"class":            helpers.Get(func(c fiber.Ctx) error { return timetableapi.ReadBaseScheduleClass(c, pool, rdb) }),
							"group":            helpers.Get(func(c fiber.Ctx) error { return timetableapi.ReadBaseScheduleGroup(c, pool, rdb) }),
						},
						"realtime": fiber.Map{
							fiber.MethodPost:   func(c fiber.Ctx) error { return timetableapi.CreateRealTimeLesson(c, pool, rdb) },
							fiber.MethodPatch:  func(c fiber.Ctx) error { return timetableapi.UpdateRealTimeLesson(c, pool, rdb) },
							fiber.MethodDelete: func(c fiber.Ctx) error { return timetableapi.DeleteRealTimeLesson(c, pool, rdb) },
							fiber.MethodGet:    func(c fiber.Ctx) error { return timetableapi.ReadRealTimeTable(c, pool, rdb) },
						},
						"canceled_lesson": fiber.Map{
							fiber.MethodPost:   func(c fiber.Ctx) error { return timetableapi.AddCanceledLesson(c, pool, rdb) },
							fiber.MethodDelete: func(c fiber.Ctx) error { return timetableapi.RemoveCanceledLesson(c, pool, rdb) },
						},
						"substitution": helpers.Patch(func(c fiber.Ctx) error { return timetableapi.UpdateSubstitution(c, pool, rdb) }),
					},
				},
			},
		},
	}
}

func main() {
	if helpers.GetEnvFallback("DEMO_MODE", "false") == "true" && helpers.GetEnvFallback("APP_ENV", "development") == "production" {
		slog.Warn("DEMO_MODE = true is not recommended in production!")
	}

	if !fiber.IsChild() {
		f, err := helpers.SetupLogging()

		if err != nil {
			panic(err)
		}

		defer f.Close()
	}

	pool := db.Connect()

	if !fiber.IsChild() {
		err := db.RunMigrations(pool)

		if err != nil {
			panic(fmt.Sprintf("Failed to apply migrations: %s", err))
		}
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", helpers.GetEnvFallback("REDIS_HOST", "localhost"), helpers.GetEnvFallback("REDIS_PORT", "6379")),
		Password: "",
		DB:       0,
	})
	defer rdb.Close()

	opaque_server, err := opaquepkg.CreateServerFromEnv()

	if err != nil {
		panic(fmt.Sprintf("Failed to create OPAQUE server: %s", err))
	}

	if err := opaquepkg.InitFakeRecord(opaquepkg.Conf); err != nil {
		panic(fmt.Sprintf("Failed to initialize OPAQUE fake record: %s", err))
	}

	if !fiber.IsChild() && helpers.GetEnvFallback("DEMO_MODE", "false") == "true" {
		if err := opaquepkg.SeedDemo(context.Background(), pool, opaque_server); err != nil {
			panic(fmt.Sprintf("Failed to seed demo accounts: %s", err))
		}
	}

	app := fiber.New(fiber.Config{
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
		AppName:      helpers.AppName,
		BodyLimit:    1 << 18,
		ReadTimeout:  time.Duration(helpers.GetInt64EnvFallback("SERVER_READ_TIMEOUT", 10, 3600)) * time.Second,
		WriteTimeout: time.Duration(helpers.GetInt64EnvFallback("SERVER_WRITE_TIMEOUT", 30, 3600)) * time.Second,
		IdleTimeout:  time.Duration(helpers.GetInt64EnvFallback("SERVER_IDLE_TIMEOUT", 120, 86400)) * time.Second,
		TrustProxy:   helpers.GetEnvFallback("TRUST_PROXY", "false") == "true",
		ProxyHeader:  fiber.HeaderXForwardedFor,
		TrustProxyConfig: fiber.TrustProxyConfig{
			Proxies: strings.Split(helpers.GetEnvFallback("TRUSTED_PROXIES", "127.0.0.1/32"), ","),
		},
	})

	app.Use(func(c fiber.Ctx) error {
		return middlewares.SecurityHeadersMiddleware(c)
	})

	apiRateMax := helpers.GetInt64EnvFallback("API_RATE_MAX", 60, 1000000)
	apiRateWindow := helpers.GetInt64EnvFallback("API_RATELIMIT_WINDOW", 60, 1000000)

	// frontend
	app.Use("/assets/fonts", static.New("frontend/dist/assets/fonts", static.Config{MaxAge: 31536000}))
	app.Use("/assets", static.New("frontend/dist/assets", static.Config{MaxAge: 3600, Compress: true}))
	indexHTML, err := os.ReadFile("frontend/dist/index.html")
	if err != nil {
		panic(err)
	}
	for _, path := range helpers.FrontendPaths {
		app.Get(path, func(c fiber.Ctx) error {
			return middlewares.FrontendMiddleware(c, pool, rdb, indexHTML, "api", apiRateMax, apiRateWindow)
		})
	}

	api := app.Group("/api")
	helpers.RegisterRoutes(api, buildRoutes(pool, rdb, opaque_server))

	log.Fatal(app.Listen(":8080",
		fiber.ListenConfig{
			EnablePrefork:         true,
			DisableStartupMessage: helpers.GetEnvFallback("APP_ENV", "development") == "production",
		}))
}

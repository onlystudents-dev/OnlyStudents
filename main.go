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
	env "onlystudents/internal/helpers"
	"onlystudents/internal/middlewares"
	opaquepkg "onlystudents/internal/opaque"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/redis/go-redis/v9"
)

func main() {
	if helpers.GetEnvFallback("DEMO_MODE", "false") == "true" && env.GetEnvFallback("APP_ENV", "development") == "production" {
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
		Addr:     fmt.Sprintf("%s:%s", env.GetEnvFallback("REDIS_HOST", "localhost"), env.GetEnvFallback("REDIS_PORT", "6379")),
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
		AppName:      "OnlyStudents",
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

	authRateMax := helpers.GetInt64EnvFallback("AUTH_RATE_MAX", 5, 1000000)
	authRateWindow := helpers.GetInt64EnvFallback("AUTH_RATELIMIT_WINDOW", 60, 1000000)
	apiRateMax := helpers.GetInt64EnvFallback("API_RATE_MAX", 60, 1000000)
	apiRateWindow := helpers.GetInt64EnvFallback("API_RATELIMIT_WINDOW", 60, 1000000)
	forgetRateMax := helpers.GetInt64EnvFallback("FORGET_RATE_MAX", 1, 1000000)
	forgetRateWindow := helpers.GetInt64EnvFallback("FORGET_RATELIMIT_WINDOW", 120, 1000000)

	// frontend
	app.Use("/assets/fonts", static.New("frontend/dist/assets/fonts", static.Config{MaxAge: 31536000}))
	app.Use("/assets", static.New("frontend/dist/assets", static.Config{MaxAge: 3600, Compress: true}))
	paths := []string{
		"/",
		"/me",
		"/timetable",
		"/grades",
		"/homeworks",
		"/absences",
	}
	indexHTML, err := os.ReadFile("frontend/dist/index.html")
	if err != nil {
		panic(err)
	}
	for _, path := range paths {
		app.Get(path, func(c fiber.Ctx) error {
			return middlewares.FrontendMiddleware(c, pool, rdb, indexHTML, "api", apiRateMax, apiRateWindow)
		})
	}

	app.Get("/ping", func(c fiber.Ctx) error {
		err := pool.Ping(c.Context())

		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to ping database: %s\n", err)
		}

		return c.SendString("Pong!")
	})

	api := app.Group("/api")

	api.Get("/config", func(c fiber.Ctx) error {
		c.Set("Cache-Control", "public, max-age=86400")
		return c.JSON(fiber.Map{
			"password": fiber.Map{
				"minLength":        env.GetIntEnvFallback("PW_MIN_LEN", 12, 128),
				"maxLength":        env.GetIntEnvFallback("PW_MAX_LEN", 256, 1024),
				"hibpCheckEnabled": env.GetEnvFallback("HIBP_CHECK_ENABLED", "false") == "true",
			},
		})
	})

	authLimit := func(c fiber.Ctx) error {
		return middlewares.RateLimitMiddleware(c, rdb, "auth", authRateMax, authRateWindow)
	}
	apiLimit := func(c fiber.Ctx) error {
		return middlewares.RateLimitMiddleware(c, rdb, "api", apiRateMax, apiRateWindow)
	}
	forgetLimit := func(c fiber.Ctx) error {
		return middlewares.RateLimitMiddleware(c, rdb, "forget", forgetRateMax, forgetRateWindow)
	}

	api.Post("/login/init", authLimit, func(c fiber.Ctx) error {
		return v1.LoginInit(c, pool, rdb, opaque_server)
	})

	api.Post("/login/finish", authLimit, func(c fiber.Ctx) error {
		return v1.LoginFinish(c, pool, rdb, opaque_server)
	})

	api.Post("/enroll/init", authLimit, func(c fiber.Ctx) error {
		return v1.EnrollInit(c, pool, rdb, opaque_server)
	})

	api.Post("/enroll/finish", authLimit, func(c fiber.Ctx) error {
		return v1.EnrollFinish(c, pool, rdb, opaque_server)
	})

	api.Post("/forget_password", forgetLimit, func(c fiber.Ctx) error {
		return v1.ForgetPassword(c, pool, rdb)
	})

	api.Post("/forget_password_confirm", forgetLimit, func(c fiber.Ctx) error {
		return v1.ForgetPasswordConfirm(c, pool, rdb)
	})

	api_v1 := api.Group("/v1")

	admin_group := api_v1.Group("/admin", apiLimit,
		func(c fiber.Ctx) error {
			return middlewares.AuthMiddleware(c, rdb)
		},
		func(c fiber.Ctx) error {
			return middlewares.RequireRoleMiddleware(c, rdb, []string{"admin"})
		})

	admin_group.Get("/logs", func(c fiber.Ctx) error {
		return adminapi.AdminLogs(c, pool, rdb)
	})

	admin_group.Get("/status", func(c fiber.Ctx) error {
		return adminapi.AdminStatus(c, pool, rdb)
	})

	admin_group.Post("/enroll/student", func(c fiber.Ctx) error {
		return adminapi.EnrollStudent(c, pool, rdb, opaque_server)
	})

	admin_group.Post("/enroll/teacher", func(c fiber.Ctx) error {
		return adminapi.EnrollTeacher(c, pool, rdb, opaque_server)
	})

	admin_group.Post("/enroll/guardian", func(c fiber.Ctx) error {
		return adminapi.EnrollGuardian(c, pool, rdb, opaque_server)
	})

	me_group := api_v1.Group("/me", apiLimit, func(c fiber.Ctx) error {
		return middlewares.AuthMiddleware(c, rdb)
	})

	me_group.Post("/logout", func(c fiber.Ctx) error {
		return v1.Logout(c, pool, rdb)
	})

	me_group.Post("/update_preferences", func(c fiber.Ctx) error {
		return meapi.UpdatePreferences(c, pool, rdb)
	})

	me_group.Post("/change_email", func(c fiber.Ctx) error {
		return meapi.ChangeEmail(c, pool, rdb)
	})

	me_group.Post("/verify_email", authLimit, func(c fiber.Ctx) error {
		return meapi.VerifyEmailRequest(c, pool, rdb)
	})

	me_group.Post("/verify_email_confirm", authLimit, func(c fiber.Ctx) error {
		return meapi.VerifyEmailConfirm(c, pool, rdb)
	})

	me_group.Get("/status", func(c fiber.Ctx) error {
		return meapi.Status(c, pool, rdb)
	})

	student_group := me_group.Group("/student",
		func(c fiber.Ctx) error {
			return middlewares.RequireRoleMiddleware(c, rdb, []string{"student", "guardian"})
		})

	student_group.Get("/grades", func(c fiber.Ctx) error {
		return studentapi.Grades(c, pool, rdb)
	})

	student_group.Get("/absences", func(c fiber.Ctx) error {
		return studentapi.Absences(c, pool, rdb)
	})

	student_group.Get("/final_grades", func(c fiber.Ctx) error {
		return studentapi.FinalGrades(c, pool, rdb)
	})

	student_group.Get("/exams", func(c fiber.Ctx) error {
		return studentapi.Exams(c, pool, rdb)
	})

	student_group.Get("/homework", func(c fiber.Ctx) error {
		return studentapi.Homework(c, pool, rdb)
	})

	student_group.Get("/timetable", func(c fiber.Ctx) error {
		return studentapi.ReadMyRealTimeTable(c, pool, rdb)
	})

	student_group.Get("/timetable/base", func(c fiber.Ctx) error {
		return studentapi.ReadBaseSchedule(c, pool, rdb)
	})

	student_group.Get("/timetable/lesson_time", func(c fiber.Ctx) error {
		return studentapi.ReadLessonTime(c, pool, rdb)
	})

	student_group.Get("/timetable/room", func(c fiber.Ctx) error {
		return studentapi.ReadRoom(c, pool, rdb)
	})

	student_group.Get("/timetable/custom_subject", func(c fiber.Ctx) error {
		return studentapi.ReadCustomSubject(c, pool, rdb)
	})

	student_group.Get("/class", func(c fiber.Ctx) error {
		return studentapi.ReadClass(c, pool, rdb)
	})

	student_group.Get("/groups", func(c fiber.Ctx) error {
		return studentapi.ReadGroups(c, pool, rdb)
	})

	student_group.Get("/timetable/bell_schedule_type", func(c fiber.Ctx) error {
		return studentapi.ReadBellScheduleType(c, pool, rdb)
	})

	teacher_group := me_group.Group("/teacher", func(c fiber.Ctx) error {
		return middlewares.RequireRoleMiddleware(c, rdb, []string{"teacher"})
	})

	timetable_group := teacher_group.Group("/timetable")

	timetable_group.Post("/bell_schedule_type", func(c fiber.Ctx) error {
		return timetableapi.CreateBellScheduleType(c, pool, rdb)
	})
	timetable_group.Patch("/bell_schedule_type", func(c fiber.Ctx) error {
		return timetableapi.EditBellScheduleType(c, pool, rdb)
	})
	timetable_group.Delete("/bell_schedule_type", func(c fiber.Ctx) error {
		return timetableapi.DeleteBellScheduleType(c, pool, rdb)
	})
	timetable_group.Get("/bell_schedule_type", func(c fiber.Ctx) error {
		return timetableapi.ReadBellScheduleType(c, pool, rdb)
	})

	timetable_group.Post("/lesson_time", func(c fiber.Ctx) error {
		return timetableapi.CreateLessonTime(c, pool, rdb)
	})
	timetable_group.Patch("/lesson_time", func(c fiber.Ctx) error {
		return timetableapi.EditLessonTime(c, pool, rdb)
	})
	timetable_group.Delete("/lesson_time", func(c fiber.Ctx) error {
		return timetableapi.DeleteLessonTime(c, pool, rdb)
	})
	timetable_group.Get("/lesson_time", func(c fiber.Ctx) error {
		return timetableapi.ReadLessonTime(c, pool, rdb)
	})

	timetable_group.Post("/custom_subject", func(c fiber.Ctx) error {
		return timetableapi.CreateCustomSubject(c, pool, rdb)
	})
	timetable_group.Patch("/custom_subject", func(c fiber.Ctx) error {
		return timetableapi.EditCustomSubject(c, pool, rdb)
	})
	timetable_group.Delete("/custom_subject", func(c fiber.Ctx) error {
		return timetableapi.DeleteCustomSubject(c, pool, rdb)
	})
	timetable_group.Get("/custom_subject", func(c fiber.Ctx) error {
		return timetableapi.ReadCustomSubject(c, pool, rdb)
	})

	timetable_group.Post("/room", func(c fiber.Ctx) error {
		return timetableapi.CreateRoom(c, pool, rdb)
	})
	timetable_group.Patch("/room", func(c fiber.Ctx) error {
		return timetableapi.UpdateRoom(c, pool, rdb)
	})
	timetable_group.Delete("/room", func(c fiber.Ctx) error {
		return timetableapi.DeleteRoom(c, pool, rdb)
	})
	timetable_group.Get("/room", func(c fiber.Ctx) error {
		return timetableapi.ReadRoom(c, pool, rdb)
	})

	timetable_group.Post("/group", func(c fiber.Ctx) error {
		return timetableapi.CreateGroup(c, pool, rdb)
	})
	timetable_group.Patch("/group", func(c fiber.Ctx) error {
		return timetableapi.EditGroup(c, pool, rdb)
	})
	timetable_group.Delete("/group", func(c fiber.Ctx) error {
		return timetableapi.DeleteGroup(c, pool, rdb)
	})
	timetable_group.Get("/group", func(c fiber.Ctx) error {
		return timetableapi.ReadGroup(c, pool, rdb)
	})
	timetable_group.Post("/group/student", func(c fiber.Ctx) error {
		return timetableapi.InsertStudentToGroup(c, pool, rdb)
	})
	timetable_group.Delete("/group/student", func(c fiber.Ctx) error {
		return timetableapi.DeleteStudentFromGroup(c, pool, rdb)
	})
	timetable_group.Get("/group/student", func(c fiber.Ctx) error {
		return timetableapi.ReadStudentFromGroup(c, pool, rdb)
	})

	timetable_group.Post("/base_schedule", func(c fiber.Ctx) error {
		return timetableapi.CreateBaseSchedule(c, pool, rdb)
	})
	timetable_group.Patch("/base_schedule", func(c fiber.Ctx) error {
		return timetableapi.UpdateBaseSchedule(c, pool, rdb)
	})
	timetable_group.Delete("/base_schedule", func(c fiber.Ctx) error {
		return timetableapi.DeleteBaseSchedule(c, pool, rdb)
	})
	timetable_group.Get("/base_schedule/class", func(c fiber.Ctx) error {
		return timetableapi.ReadBaseScheduleClass(c, pool, rdb)
	})
	timetable_group.Get("/base_schedule/group", func(c fiber.Ctx) error {
		return timetableapi.ReadBaseScheduleGroup(c, pool, rdb)
	})

	timetable_group.Post("/realtime", func(c fiber.Ctx) error {
		return timetableapi.CreateRealTimeLesson(c, pool, rdb)
	})
	timetable_group.Patch("/realtime", func(c fiber.Ctx) error {
		return timetableapi.UpdateRealTimeLesson(c, pool, rdb)
	})
	timetable_group.Delete("/realtime", func(c fiber.Ctx) error {
		return timetableapi.DeleteRealTimeLesson(c, pool, rdb)
	})
	timetable_group.Get("/realtime", func(c fiber.Ctx) error {
		return timetableapi.ReadRealTimeTable(c, pool, rdb)
	})

	timetable_group.Post("/canceled_lesson", func(c fiber.Ctx) error {
		return timetableapi.AddCanceledLesson(c, pool, rdb)
	})
	timetable_group.Delete("/canceled_lesson", func(c fiber.Ctx) error {
		return timetableapi.RemoveCanceledLesson(c, pool, rdb)
	})

	timetable_group.Patch("/substitution", func(c fiber.Ctx) error {
		return timetableapi.UpdateSubstitution(c, pool, rdb)
	})

	teacher_group.Get("/add_exam", func(c fiber.Ctx) error {
		return teacherapi.AddExam(c, pool, rdb)
	})

	teacher_group.Get("/remove_exam", func(c fiber.Ctx) error {
		return teacherapi.RemoveExam(c, pool, rdb)
	})

	teacher_group.Get("/edit_exam", func(c fiber.Ctx) error {
		return teacherapi.EditExam(c, pool, rdb)
	})

	teacher_group.Get("/add_grades", func(c fiber.Ctx) error {
		return teacherapi.AddGrade(c, pool, rdb)
	})

	teacher_group.Get("/remove_grades", func(c fiber.Ctx) error {
		return teacherapi.RemoveGrade(c, pool, rdb)
	})

	teacher_group.Get("/edit_grades", func(c fiber.Ctx) error {
		return teacherapi.EditGrade(c, pool, rdb)
	})

	log.Fatal(app.Listen(":8080",
		fiber.ListenConfig{
			EnablePrefork:         true,
			DisableStartupMessage: env.GetEnvFallback("APP_ENV", "development") == "production",
		}))
}

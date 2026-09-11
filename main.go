package main

import (
	"fmt"
	"log"
	"os"

	v1 "onlystudents/internal/api/v1"
	adminapi "onlystudents/internal/api/v1/admin"
	meapi "onlystudents/internal/api/v1/me"
	manageapi "onlystudents/internal/api/v1/me/manage"
	"onlystudents/internal/db"
	"onlystudents/internal/helpers"
	env "onlystudents/internal/helpers"
	"onlystudents/internal/middlewares"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/redis/go-redis/v9"
)

func main() {
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

	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
		AppName:     "OnlyStudents",
		BodyLimit:   1 << 20,
	})

	// frontend
	app.Use("/assets", static.New("frontend/dist/assets"))
	paths := []string{
		"/",
		"/me",
		"/timetable",
		"/grades",
		"/homeworks",
		"/absences",
	}
	for _, path := range paths {
		app.Get(path, func(c fiber.Ctx) error {
			return c.SendFile("frontend/dist/index.html")
		})
	}

	app.Get("/ping", func(c fiber.Ctx) error {
		err := pool.Ping(c.Context())

		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to ping database: %s\n", err)
		}

		return c.SendString("Pong!")
	})

	authRateMax := helpers.GetInt64EnvFallback("AUTH_RATE_MAX", 5, 1000000)
	authRateWindow := helpers.GetInt64EnvFallback("AUTH_RATELIMIT_WINDOW", 60, 1000000)
	apiRateMax := helpers.GetInt64EnvFallback("API_RATE_MAX", 60, 1000000)
	apiRateWindow := helpers.GetInt64EnvFallback("API_RATELIMIT_WINDOW", 60, 1000000)

	api := app.Group("/api")

	authLimit := func(c fiber.Ctx) error {
		return middlewares.RateLimitMiddleware(c, rdb, "auth", authRateMax, authRateWindow)
	}
	apiLimit := func(c fiber.Ctx) error {
		return middlewares.RateLimitMiddleware(c, rdb, "api", apiRateMax, apiRateWindow)
	}

	api.Post("/login", authLimit, func(c fiber.Ctx) error {
		return v1.Login(c, pool, rdb)
	})

	api.Post("/forget_password", func(c fiber.Ctx) error {
		return v1.ForgetPassword(c, pool, rdb)
	})

	api.Post("/forget_password_confirm", authLimit, func(c fiber.Ctx) error {
		return v1.ForgetPasswordConfirm(c, pool, rdb)
	})

	api_v1 := api.Group("/v1")

	admin_group := api_v1.Group("/admin", apiLimit, func(c fiber.Ctx) error {
		return middlewares.RequireRoleMiddleware(c, rdb, []string{"admin"})
	})

	admin_group.Post("/login", authLimit, func(c fiber.Ctx) error {
		return adminapi.AdminLogin(c, pool, rdb)
	})

	admin_group.Post("/debug", authLimit, func(c fiber.Ctx) error {
		return adminapi.AdminDebug(c, pool, rdb)
	})

	me_group := api_v1.Group("/me", apiLimit, func(c fiber.Ctx) error {
		return middlewares.AuthMiddleware(c, rdb)
	})

	me_group.Post("/logout", func(c fiber.Ctx) error {
		return v1.Logout(c, pool, rdb)
	})

	me_group.Post("/change_password", func(c fiber.Ctx) error {
		return meapi.ChangePassword(c, pool, rdb)
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

	me_group.Get("/timetable", func(c fiber.Ctx) error {
		return c.SendStatus(501) // NOTE: Future codebase update for timetable should replace this
	})

	me_group.Get("/subjects", func(c fiber.Ctx) error {
		return meapi.Subjects(c, pool, rdb)
	})

	me_group.Get("/grades", func(c fiber.Ctx) error {
		return meapi.Grades(c, pool, rdb)
	})

	me_group.Get("/final_grades", func(c fiber.Ctx) error {
		return meapi.FinalGrades(c, pool, rdb)
	})

	me_group.Get("/absences", func(c fiber.Ctx) error {
		return meapi.Absences(c, pool, rdb)
	})

	me_group.Get("/exams", func(c fiber.Ctx) error {
		return meapi.Exams(c, pool, rdb)
	})

	me_group.Get("/homework", func(c fiber.Ctx) error {
		return meapi.Homework(c, pool, rdb)
	})

	manage_group := me_group.Group("/manage", func(c fiber.Ctx) error {
		return middlewares.RequireRoleMiddleware(c, rdb, []string{"teacher"})
	})

	manage_group.Get("/absences", func(c fiber.Ctx) error {
		return manageapi.ManageAbsences(c, pool, rdb)
	})

	manage_group.Get("/exams", func(c fiber.Ctx) error {
		return manageapi.ManageExams(c, pool, rdb)
	})

	manage_group.Get("/grades", func(c fiber.Ctx) error {
		return manageapi.ManageGrades(c, pool, rdb)
	})

	manage_group.Get("/timetable", func(c fiber.Ctx) error {
		return manageapi.ManageTimeTable(c, pool, rdb)
	})

	log.Fatal(app.Listen(":8080", fiber.ListenConfig{
		EnablePrefork:         true,
		DisableStartupMessage: env.GetEnvFallback("APP_ENV", "development") == "production",
	}))
}

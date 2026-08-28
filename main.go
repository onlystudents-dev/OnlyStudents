package main

import (
	"context"
	"fmt"
	"log"
	"os"

	v1 "onlystudents/internal/api/v1"
	meapi "onlystudents/internal/api/v1/me"
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

	session_store := helpers.SessionStore{
		RedisDB: rdb,
	}

	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
		AppName:     "OnlyStudents",
	})

	// frontend
	app.Use("/assets", static.New("frontend/dist/assets"))
	paths := []string{
		"/",
		"/homeworks",
	}
	for _, path := range paths {
		app.Get(path, func(c fiber.Ctx) error {
			return c.SendFile("frontend/dist/index.html")
		})
	}

	app.Get("/ping", func(c fiber.Ctx) error {
		err := pool.Ping(context.Background())

		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to ping database: %s", err)
		}

		return c.SendString("Pong!")
	})

	api := app.Group("/api")

	api.Post("/login", func(c fiber.Ctx) error {
		return v1.Login(c, pool, &session_store)
	})

	api.Get("/logout", func(c fiber.Ctx) error {
		return v1.Logout(c, pool, &session_store)
	})

	api_v1 := api.Group("/v1", func(c fiber.Ctx) error {
		return middlewares.AuthMiddleware(c, &session_store)
	})

	me := api_v1.Group("/me")

	me.Get("/status", func(c fiber.Ctx) error {
		return meapi.Status(c, pool, &session_store)
	})

	me.Get("/timetable", func(c fiber.Ctx) error {
		return meapi.TimeTable(c, pool, &session_store)
	})

	me.Get("/subjects", func(c fiber.Ctx) error {
		return meapi.Subjects(c, pool, &session_store)
	})

	me.Get("/grades", func(c fiber.Ctx) error {
		return meapi.Grades(c, pool, &session_store)
	})

	me.Get("/absences", func(c fiber.Ctx) error {
		return meapi.Absences(c, pool, &session_store)
	})

	me.Get("/exams", func(c fiber.Ctx) error {
		return meapi.Exams(c, pool, &session_store)
	})

	me.Get("/homework", func(c fiber.Ctx) error {
		return meapi.Homework(c, pool, &session_store)
	})

	log.Fatal(app.Listen(":8080", fiber.ListenConfig{
		EnablePrefork:         true,
		DisableStartupMessage: env.GetEnvFallback("APP_ENV", "development") == "production",
	}))
}

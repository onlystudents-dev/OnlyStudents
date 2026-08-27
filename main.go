package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"onlystudents/internal/auth"
	"onlystudents/internal/db"
	db_queries "onlystudents/internal/db/store"
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

	session_store := auth.SessionStore{
		RedisDB: rdb,
	}

	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
		AppName:     "OnlyStudents",
	})

	app.Use("/", static.New("frontend/dist"))

	app.Get("/ping", func(c fiber.Ctx) error {
		err := pool.Ping(context.Background())

		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to ping database: %s", err)
		}

		return c.SendString("Pong!")
	})

	api := app.Group("/api")

	api_v1 := api.Group("/v1", func(c fiber.Ctx) error {
		return middlewares.AuthMiddleware(c, &session_store)
	})

	api_v1.Get("/students", func(c fiber.Ctx) error {
		queries := db_queries.New(pool)
		students, err := queries.ListStudents(context.Background())

		if err != nil {
			return c.SendStatus(500)
		}

		return c.JSON(students)
	})
	api_v1.Get("/guardians", func(c fiber.Ctx) error {
		queries := db_queries.New(pool)
		guardians, err := queries.ListGuardians(context.Background())

		if err != nil {
			return c.SendStatus(500)
		}

		return c.JSON(guardians)
	})
	api_v1.Get("/teachers", func(c fiber.Ctx) error {
		queries := db_queries.New(pool)
		teachers, err := queries.ListTeachers(context.Background())

		if err != nil {
			return c.SendStatus(500)
		}

		return c.JSON(teachers)
	})

	log.Fatal(app.Listen(":8080", fiber.ListenConfig{
		EnablePrefork:         true,
		DisableStartupMessage: env.GetEnvFallback("APP_ENV", "development") == "production",
	}))
}

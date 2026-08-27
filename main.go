package main

import (
	"context"
	"fmt"
	"log"
	"os"

	v1 "onlystudents/internal/api/v1"
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

	app.Use("/", static.New("frontend/dist"))

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

	api_v1 := api.Group("/v1", func(c fiber.Ctx) error {
		return middlewares.AuthMiddleware(c, &session_store)
	})

	api_v1.Get("/me", func(c fiber.Ctx) error {
		return v1.Me(c, pool, &session_store)
	})

	api_v1.Get("/students", func(c fiber.Ctx) error {
		return v1.Students(c, pool)
	})

	api_v1.Get("/guardians", func(c fiber.Ctx) error {
		return v1.Guardians(c, pool)
	})

	api_v1.Get("/teachers", func(c fiber.Ctx) error {
		return v1.Teachers(c, pool)
	})

	api_v1.Get("/schools", func(c fiber.Ctx) error {
		return v1.Schools(c, pool)
	})

	log.Fatal(app.Listen(":8080", fiber.ListenConfig{
		EnablePrefork:         true,
		DisableStartupMessage: env.GetEnvFallback("APP_ENV", "development") == "production",
	}))
}

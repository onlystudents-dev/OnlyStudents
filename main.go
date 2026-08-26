package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"onlystudents/internal/db"
	db_queries "onlystudents/internal/db/store"
	env "onlystudents/internal/helpers"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
)

func main() {
	pool := db.Connect()

	if !fiber.IsChild() {
		err := db.RunMigrations(pool)

		if err != nil {
			panic(fmt.Sprintf("Failed to apply migrations: %s", err))
		}
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

	app.Get("/api/students", func(c fiber.Ctx) error {
		queries := db_queries.New(pool)
		students, err := queries.ListStudents(context.Background())

		if err != nil {
			return c.SendStatus(500)
		}

		return c.JSON(students)
	})

	log.Fatal(app.Listen(":8080", fiber.ListenConfig{
		EnablePrefork:         true,
		DisableStartupMessage: env.GetEnvFallback("APP_ENV", "development") == "production",
	}))
}

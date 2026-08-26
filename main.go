package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"onlystudents/internal/db"
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
		rows, err := pool.Query(context.Background(), "SELECT * FROM students")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to query students table: %s", err)
			return c.SendStatus(500)
		}
		defer rows.Close()

		var students []map[string]any
		for rows.Next() {
			values, err := rows.Values()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to iterate student row values: %s", err)
				return c.SendStatus(500)
			}
			student := make(map[string]any)
			for i, fd := range rows.FieldDescriptions() {
				student[string(fd.Name)] = values[i]
			}
			students = append(students, student)
		}
		return c.JSON(students)
	})

	log.Fatal(app.Listen(":8080", fiber.ListenConfig{
		EnablePrefork:         true,
		DisableStartupMessage: env.GetEnvFallback("APP_ENV", "development") == "production",
	}))
}

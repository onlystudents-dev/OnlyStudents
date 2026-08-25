package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
)

func main() {
	app := fiber.New()

	app.Use("/", static.New("frontend/dist"))

	app.Get("/ping", func(c fiber.Ctx) error {
		return c.SendString("Pong!")
	})

	app.Get("/api/hello", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "hello world"})
	})

	log.Fatal(app.Listen(":8080"))
}

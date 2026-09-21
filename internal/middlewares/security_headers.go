package middlewares

import (
	"github.com/gofiber/fiber/v3"
)

func SecurityHeadersMiddleware(c fiber.Ctx) error {
	c.Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'wasm-unsafe-eval'; worker-src 'self' blob:; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self' data:; connect-src 'self'; base-uri 'self'; frame-ancestors 'none'; form-action 'self'")
	c.Set("X-Content-Type-Options", "nosniff")
	c.Set("Referrer-Policy", "no-referrer")
	c.Set("Cross-Origin-Resource-Policy", "same-origin")
	c.Set("Permissions-Policy", "camera=(), microphone=(), geolocation=(), payment=(), usb=()")
	return c.Next()
}

package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/helmet"
)

// function being used by csrf
func MiddleHelmet() fiber.Handler {
	var helmetConfig = helmet.Config{
		XSSProtection: "1",
		HSTSMaxAge:    36000,
	}
	return helmet.New(helmetConfig)
}

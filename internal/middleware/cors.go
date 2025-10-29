package middleware

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func MiddleCors() fiber.Handler {
	var corsConfig = cors.Config{
		AllowOrigins:     os.Getenv("ALLOW_ORIGINS"),
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization,X-Csrf-Token",
		MaxAge:           300,
		AllowCredentials: true,
	}
	return cors.New(corsConfig)
}

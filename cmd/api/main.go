package api

import (
	"os"
	"time"

	"classplanner/internal/middleware"
	"classplanner/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

func Microservice() {
	utils.LoadEnv()

	// load initial configuration
	app := fiber.New(fiber.Config{
		ServerHeader: os.Getenv("SERVER_HEADER"),
		AppName:      os.Getenv("APP_NAME")})

	// load middlewares
	app.Use(middleware.MiddleCsrf())
	app.Use(middleware.LoggerStarter())

	// health check
	app.Get("/up", func(c *fiber.Ctx) error {
		return c.SendString("service is up!")
	})

	// load files
	app.Static(os.Getenv("UPLOADS_URL"), os.Getenv("UPLOADS_PATH"), fiber.Static{
		Compress:      true,
		ByteRange:     true,
		Browse:        true,
		CacheDuration: 1 * time.Hour,
		MaxAge:        36000,
	})

	// initialize the api
	app.Listen(":3000")
}

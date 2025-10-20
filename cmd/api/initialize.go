package api

import (
	"os"
	"time"

	"classplanner/internal/infrastructure/database"
	"classplanner/internal/middleware"
	"classplanner/internal/repository"
	"classplanner/internal/service"
	"classplanner/internal/transport/users"
	"classplanner/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// AppDependencies holds all the dependencies to be injected into the handlers
type AppDependencies struct {
	UserHandler *users.UserHandler
	DB          *database.DatabaseInstance
	App         *fiber.App
}

// Initialize loads configuration, middlewares, database, repositories, and services
func Initialize() *AppDependencies {
	utils.LoadEnv()

	// Create a new Fiber app
	app := fiber.New(fiber.Config{
		ServerHeader: os.Getenv("SERVER_HEADER"),
		AppName:      os.Getenv("APP_NAME"),
	})

	// Connect to the database
	dbInstance := database.Connect()

	// Register middlewares
	app.Use(middleware.MiddleCsrf())
	app.Use(middleware.LoggerStarter())
	app.Use(middleware.HealthCheck())
	app.Use(middleware.MiddleHelmet())

	// Serve static files
	app.Static(os.Getenv("UPLOADS_URL"), os.Getenv("UPLOADS_PATH"), fiber.Static{
		Compress:      true,
		ByteRange:     true,
		Browse:        true,
		CacheDuration: 1 * time.Hour,
		MaxAge:        36000,
	})

	// Initialize repositories and services
	userRepo := repository.NewUserRepository(dbInstance.DB)
	userService := service.NewUserService(userRepo)
	userHandler := users.NewUserHandler(userService)

	return &AppDependencies{
		UserHandler: userHandler,
		DB:          dbInstance,
		App:         app,
	}
}

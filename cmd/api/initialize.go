package api

import (
	"context"
	"log"
	"os"
	"time"

	"classplanner/internal/infrastructure/database"
	"classplanner/internal/infrastructure/logger"
	"classplanner/internal/middleware"
	"classplanner/internal/repository"
	"classplanner/internal/service"
	"classplanner/internal/transport/users"
	"classplanner/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// AppDependencies holds all the dependencies to be injected into the handlers
type AppDependencies struct {
	UserHandler *users.UserHandler
	DB          *database.DatabaseInstance
	App         *fiber.App
	Logger      *logger.Logger
}

// Initialize loads configuration, middlewares, database, repositories, and services
func Initialize() *AppDependencies {
	utils.LoadEnv()

	app := fiber.New(fiber.Config{
		ServerHeader: os.Getenv("SERVER_HEADER"),
		AppName:      os.Getenv("APP_NAME"),
	})

	dbInstance := database.Connect()

	// Register middlewares
	app.Use(middleware.MiddleCsrf())
	app.Use(middleware.LoggerStarter())
	app.Use(middleware.HealthCheck())

	// Serve static files
	app.Static(os.Getenv("UPLOADS_URL"), os.Getenv("UPLOADS_PATH"), fiber.Static{
		Compress:      true,
		ByteRange:     true,
		Browse:        true,
		CacheDuration: 1 * time.Hour,
		MaxAge:        36000,
	})

	// Initialize logger
	logFile := "./logs/app.log"
	l, err := logger.NewLogger(logFile, zap.DebugLevel)
	if err != nil {
		log.Fatalf("no se pudo inicializar logger: %v", err)
	}

	// Initialize repositories and services
	userRepo := repository.NewUserRepository(dbInstance.DB)
	userService := service.NewUserService(userRepo, l)
	userHandler := users.NewUserHandler(userService)

	ctx := context.Background()
	l.Info(logger.EnsureContext(ctx), "Aplicación iniciada")

	return &AppDependencies{
		UserHandler: userHandler,
		DB:          dbInstance,
		App:         app,
		Logger:      l,
	}
}

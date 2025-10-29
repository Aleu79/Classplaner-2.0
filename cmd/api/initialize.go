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
	app := fiber.New(fiber.Config{
		ServerHeader: os.Getenv("SERVER_HEADER"),
		AppName:      os.Getenv("APP_NAME"),
	})

	dbInstance := database.Connect()

	// Initialize logger
	logFile := os.Getenv("LOGGER_PATH")
	mwLogFile := os.Getenv("MW_LOGGER_PATH")
	l, err := logger.NewLogger(logFile, zap.DebugLevel)
	if err != nil {
		log.Fatalf("no se pudo inicializar logger: %v", err)
	}

	// Register middlewares
	app.Use(middleware.MiddleCors()) // cors should be always in top
	app.Use(middleware.LoggerStarter(mwLogFile, l))
	app.Use(middleware.HealthCheck())
	app.Use(middleware.MiddleHelmet())
	app.Use(middleware.MiddleCsrf(l))

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
	userService := service.NewUserService(userRepo, l)
	userHandler := users.NewUserHandler(userService, l)

	ctx := context.Background()
	l.Info(logger.EnsureContext(ctx), "Aplicación iniciada")

	return &AppDependencies{
		UserHandler: userHandler,
		DB:          dbInstance,
		App:         app,
		Logger:      l,
	}
}

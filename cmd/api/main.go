package api

import (
	"os"
	"time"

	"classplanner/internal/infrastructure/database"
	"classplanner/internal/middleware"
	"classplanner/internal/transport/users"
	"classplanner/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

func Microservice() {
	utils.LoadEnv()

	// load initial configuration
	app := fiber.New(fiber.Config{
		ServerHeader: os.Getenv("SERVER_HEADER"),
		AppName:      os.Getenv("APP_NAME")})

	// load database
	database.Connect()

	// load middlewares
	app.Use(middleware.MiddleCsrf())
	app.Use(middleware.LoggerStarter())
	app.Use(middleware.HealthCheck())

	// load files
	app.Static(os.Getenv("UPLOADS_URL"), os.Getenv("UPLOADS_PATH"), fiber.Static{
		Compress:      true,
		ByteRange:     true,
		Browse:        true,
		CacheDuration: 1 * time.Hour,
		MaxAge:        36000,
	})

	//handlers
	// app.Get("/calendar", getCalendar)
	// app.Get("/calendar/:class_token", GetCalendarWithToken)

	// //classes
	// app.Get("/classes", GetClasses)
	// app.Post("/classes", createClass)
	// app.Post("/joinClass", joinClass)

	// //tasks
	// app.Get("/tasks", GetTasks)
	// app.Post("/tasks", createTask)

	// //submissions
	// app.Get("/submission", GetSubmission)
	// app.Get("/submissions", GetSubs)
	// app.Post("/submission", createSubmission)
	// app.Put("/submission/:id_submission", updateSubmission)

	// app.Get("/califications", GetCalifications)
	// app.Get("/usersclass", GetUsersFromClass)

	// //comments
	// app.Get("/comments", GetComments)
	// app.Post("/comments", createComment)

	// //user
	app.Get("/user", users.GetUser)
	// app.Delete("/user", DeleteUser)
	// app.Put("/user", updateUser)
	// app.Post("/login", login)
	// app.Post("/register", register)

	// initialize the api
	app.Listen(":3000")
}

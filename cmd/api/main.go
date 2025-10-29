package api

import (
	"classplanner/cmd/routes"
	"log"
	"os"
)

// En Microservice() después de Initialize()
func Microservice() {
	deps := Initialize() // Load DB, handlers, app

	// Register all routes by module
	routes.RegisterUserRoutes(deps.App, deps.UserHandler)
	// routes.RegisterTaskRoutes(deps.App, deps.TaskHandler)
	// routes.RegisterClassRoutes(deps.App, deps.ClassHandler)
	// routes.RegisterSubmissionRoutes(deps.App, deps.SubmissionHandler)
	// routes.RegisterCommentRoutes(deps.App, deps.CommentHandler)
	// routes.RegisterCalendarRoutes(deps.App, deps.CalendarHandler)
	// routes.RegisterPremiumRoutes(deps.App, deps.PremiumHandler)

	// Leer puerto desde variable de entorno o usar por defecto
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Servidor Classplanner corriendo en http://localhost:%s", port)

	// Start server
	if err := deps.App.Listen(":" + port); err != nil {
		log.Fatalf("Error iniciando servidor: %v", err)
	}
}

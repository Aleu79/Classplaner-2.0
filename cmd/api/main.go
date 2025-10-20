package api

import "classplanner/cmd/routes"

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

	// Start server
	deps.App.Listen(":3000")
}

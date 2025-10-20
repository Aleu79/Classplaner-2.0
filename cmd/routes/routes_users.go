package routes

import (
	"classplanner/internal/transport/users"

	"github.com/gofiber/fiber/v2"
)

func RegisterUserRoutes(app *fiber.App, userHandler *users.UserHandler) {
	userGroup := app.Group("/users")
	userGroup.Post("/register", userHandler.Register)
	userGroup.Post("/login", userHandler.Login)
	userGroup.Post("/logout", userHandler.Logout)
	userGroup.Get("/", userHandler.GetAll)
	userGroup.Get("/:id/exists", userHandler.Exists)
	userGroup.Put("/:id", userHandler.Update)
}

package users

import (
	"classplanner/internal/infrastructure/database"

	"github.com/gofiber/fiber/v2"
)

func GetUser(c *fiber.Ctx) error {
	users, err := database.DBInstance.Repository.UserStorage.GetAll()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(users)
}

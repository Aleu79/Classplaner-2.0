package users

import (
	"github.com/gofiber/fiber/v2"
)

// Obtener todos los usuarios
func (h *UserHandler) GetAll(c *fiber.Ctx) error {
	users, err := h.service.GetAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(users)
}

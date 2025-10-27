package users

import (
	"context"

	"github.com/gofiber/fiber/v2"
)

// Obtener todos los usuarios
func (h *UserHandler) GetAll(c *fiber.Ctx) error {
	ctx := context.Background()
	users, err := h.service.GetAll(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(users)
}

package users

import (
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// Obtener usuario por ID
func (h *UserHandler) GetByID(ctx context.Context, c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	user, err := h.service.GetByID(ctx, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if user == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "usuario no encontrado"})
	}
	return c.JSON(user)
}

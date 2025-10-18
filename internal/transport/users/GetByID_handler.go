package users

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// Obtener usuario por ID
func (h *UserHandler) GetByID(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	user, err := h.service.GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if user == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "usuario no encontrado"})
	}
	return c.JSON(user)
}

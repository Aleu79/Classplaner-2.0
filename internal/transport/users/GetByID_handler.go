package users

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// Obtener usuario por ID
func (h *UserHandler) GetByID(c *fiber.Ctx) error {
	ctx := h.GetCtx(c)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		h.logger.Warn(ctx, "ID inválido en GetByID: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	user, err := h.service.GetByID(ctx, id)
	if err != nil {
		h.logger.Error(ctx, "Error obteniendo usuario ID=%d: %v", id, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if user == nil {
		h.logger.Info(ctx, "Usuario no encontrado ID=%d", id)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "usuario no encontrado"})
	}

	h.logger.Info(ctx, "Usuario obtenido ID=%d", id)
	return c.JSON(user)
}

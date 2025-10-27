package users

import (
	"github.com/gofiber/fiber/v2"
)

// Buscar usuarios por nombre o email
func (h *UserHandler) Search(c *fiber.Ctx) error {
	ctx := h.GetCtx(c)

	query := c.Query("q")
	if query == "" {
		h.logger.Warn(ctx, "Query vacía en Search")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "query requerida"})
	}

	users, err := h.service.SearchByUserOrEmail(ctx, query)
	if err != nil {
		h.logger.Error(ctx, "Error buscando usuarios con query=%s: %v", query, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	h.logger.Info(ctx, "Busqueda realizada con query=%s, resultados=%d", query, len(users))
	return c.JSON(users)
}

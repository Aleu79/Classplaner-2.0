package users

import "github.com/gofiber/fiber/v2"

// Buscar usuarios por nombre o email
func (h *UserHandler) Search(c *fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "query requerida"})
	}

	users, err := h.service.SearchByUserOrEmail(query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(users)
}

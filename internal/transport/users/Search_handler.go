package users

import (
	"classplanner/pkg/response"

	"github.com/gofiber/fiber/v2"
)

// Search busca usuarios por nombre o email
func (h *UserHandler) Search(c *fiber.Ctx) error {
	ctx := h.GetCtx(c)

	query := c.Query("q")
	if query == "" {
		h.logger.Warn(ctx, "Query vacía en Search")
		return response.BadRequest(c, "Query requerida", nil)
	}

	users, err := h.service.SearchByUserOrEmail(ctx, query)
	if err != nil {
		h.logger.Error(ctx, "Error buscando usuarios con query=%s: %v", query, err)
		return response.InternalError(c, "Error buscando usuarios", err.Error())
	}

	h.logger.Info(ctx, "Busqueda realizada con query=%s, resultados=%d", query, len(users))
	return response.Success(c, "Usuarios encontrados correctamente", users)
}

package users

import (
	"classplanner/pkg/response"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// GetByID obtiene un usuario por su ID
func (h *UserHandler) GetByID(c *fiber.Ctx) error {
	ctx := h.GetCtx(c)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		h.logger.Warn(ctx, "ID inválido en GetByID: %v", err)
		return response.BadRequest(c, "ID inválido", err.Error())
	}

	user, err := h.service.GetByID(ctx, id)
	if err != nil {
		h.logger.Error(ctx, "Error obteniendo usuario ID=%d: %v", id, err)
		return response.InternalError(c, "Error obteniendo usuario", err.Error())
	}
	if user == nil {
		h.logger.Info(ctx, "Usuario no encontrado ID=%d", id)
		return response.NotFound(c, "Usuario no encontrado")
	}

	h.logger.Info(ctx, "Usuario obtenido ID=%d", id)
	return response.Success(c, "Usuario obtenido correctamente", user)
}

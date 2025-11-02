package users

import (
	"classplanner/pkg/response"

	"github.com/gofiber/fiber/v2"
)

// GetAll obtiene todos los usuarios
func (h *UserHandler) GetAll(c *fiber.Ctx) error {
	ctx := h.GetCtx(c)

	h.logger.Info(ctx, "GetAll usuarios request recibido")

	users, err := h.service.GetAll(ctx)
	if err != nil {
		h.logger.Error(ctx, "Error obteniendo usuarios: %v", err)
		return response.InternalError(c, "Error obteniendo usuarios", err.Error())
	}

	h.logger.Info(ctx, "Usuarios obtenidos exitosamente")
	return response.Success(c, "Usuarios obtenidos correctamente", users)
}

package users

import (
	"context"

	"classplanner/internal/infrastructure/logger"

	"github.com/gofiber/fiber/v2"
)

// Obtener todos los usuarios
func (h *UserHandler) GetAll(c *fiber.Ctx) error {
	ctx := context.Background()

	// Inyectar requestID si existe
	if reqID := c.Locals("requestid"); reqID != nil {
		if rid, ok := reqID.(string); ok {
			ctx = logger.WithReqID(ctx, rid)
		}
	}

	// Inyectar userID si existe
	if uid := c.Locals("userID"); uid != nil {
		if id, ok := uid.(int); ok {
			ctx = logger.WithUserID(ctx, id)
		}
	}

	// Log de entrada
	h.logger.Info(ctx, "GetAll usuarios request recibido")

	users, err := h.service.GetAll(ctx)
	if err != nil {
		h.logger.Error(ctx, "Error obteniendo usuarios: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	h.logger.Info(ctx, "Usuarios obtenidos exitosamente")
	return c.JSON(users)
}

package middleware

import (
	"context"
	"fmt"

	"classplanner/internal/infrastructure/logger"

	"github.com/gofiber/fiber/v2"
)

// middleware para rutas no encontradas
func CathAll(l *logger.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := context.Background()
		if l != nil {
			l.Warn(ctx, "uta no encontrada: %s %s", c.Method(), c.Path())
		} else {
			fmt.Printf("ruta no encontrada: %s %s\n", c.Method(), c.Path())
		}
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":  "ruta no encontrada",
			"method": c.Method(),
			"path":   c.Path(),
		})
	}
}

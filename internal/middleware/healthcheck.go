package middleware

import (
	"classplanner/internal/infrastructure/database"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
)

// provides status & readiness for fiber services
func HealthCheck() fiber.Handler {
	return healthcheck.New(healthcheck.Config{
		LivenessProbe: func(c *fiber.Ctx) bool {
			return true
		},
		LivenessEndpoint: "/live",
		ReadinessProbe: func(c *fiber.Ctx) bool {
			if database.DBInstance == nil {
				return false
			}
			return database.DBInstance.Ready()
		},
		ReadinessEndpoint: "/ready",
	})
}

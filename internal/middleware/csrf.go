package middleware

import (
	"classplanner/internal/infrastructure/logger"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/utils"
	"go.uber.org/zap"
)

const (
	HeaderName = "X-Csrf-Token"
)

func MiddleCsrf(l *logger.Logger) fiber.Handler {
	enableCSRF := getEnvBool("ENABLE_CSRF", false)
	environment := getEnv("APP_ENV", "development")

	// disable if in development
	if environment == "development" && !enableCSRF {
		return func(c *fiber.Ctx) error {
			return c.Next()
		}
	}
	// determinate if using https
	isHTTPS := environment == "production"

	var csrfConfig = csrf.Config{
		KeyLookup:         "header:" + HeaderName,
		CookieName:        getCsrfCookieName(isHTTPS),
		CookieSameSite:    "Lax",
		CookieSecure:      isHTTPS,
		CookieSessionOnly: true,
		CookieHTTPOnly:    true,
		Expiration:        1 * time.Hour,
		KeyGenerator:      utils.UUIDv4,
		Session:           SessionStore(),
		ErrorHandler:      CsrfErrorHandler(l), // Usar el error handler personalizado
		Extractor:         csrf.CsrfFromHeader(HeaderName),
		SessionKey:        "fiber.csrf.token",
		HandlerContextKey: "fiber.csrf.handler",
		// ignore some methods during preflight
		Next: func(c *fiber.Ctx) bool {
			// allow options
			if c.Method() == "OPTIONS" {
				return true
			}
			// ignore some routes if necessary
			return shouldSkipCSRF(c.Path())
		},
	}
	return csrf.New(csrfConfig)
}

// CsrfErrorHandler crea un manejador de errores para CSRF que usa tu logger
func CsrfErrorHandler(l *logger.Logger) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		// Contexto para logging
		ctx := c.Context()

		// Log del error CSRF
		l.ErrorFields(ctx, "CSRF validation failed",
			zap.String("path", c.Path()),
			zap.String("method", c.Method()),
			zap.String("ip", c.IP()),
			zap.String("user_agent", c.Get("User-Agent")),
			zap.Error(err),
		)

		// Devolver error en formato JSON
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    fiber.StatusForbidden,
				"message": "Invalid CSRF token",
				"details": "CSRF validation failed",
			},
		})
	}
}

func getCsrfCookieName(isHTTPS bool) string {
	if isHTTPS {
		return "__Host-csrf"
	}
	return "csrf"
}

// shouldSkipCSRF determinates which routes could be skiped
func shouldSkipCSRF(path string) bool {
	skipPaths := []string{
		"/health",
		"/metrics",
		"/webhook/",
		"/api/public/",
	}

	for _, skipPath := range skipPaths {
		if len(skipPath) > 0 && skipPath[len(skipPath)-1] == '/' {
			if len(path) >= len(skipPath) && path[:len(skipPath)] == skipPath {
				return true
			}
		} else if path == skipPath {
			return true
		}
	}
	return false
}

// getEnv obtains a variable from env or return a default
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvBool obtains a variable and returns a boolean
func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

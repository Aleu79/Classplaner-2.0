package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/utils"
)

const (
	HeaderName = "X-Csrf-Token"
)

func MiddleCsrf() fiber.Handler {
	var csrfConfig = csrf.Config{
		KeyLookup:  "header:" + HeaderName,
		CookieName: "csrf",
		// CookieName:        "__Host-csrf_", this value __host must be with a secure url
		CookieSameSite:    "Lax",
		CookieSecure:      true,
		CookieSessionOnly: true,
		CookieHTTPOnly:    true,
		Expiration:        1 * time.Hour,
		KeyGenerator:      utils.UUIDv4,
		// ErrorHandler:      defaultErrorHandler,
		Session:           SessionStore(),
		Extractor:         csrf.CsrfFromHeader(HeaderName),
		SessionKey:        "fiber.csrf.token",
		HandlerContextKey: "fiber.csrf.handler",
	}
	return csrf.New(csrfConfig)
}

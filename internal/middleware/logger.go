package middleware

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"classplanner/internal/infrastructure/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/natefinch/lumberjack"
)

// Este middleware nos da un registro confiable, seguro y estructurado de cada request HTTP,
// con capacidad de correlacionar logs en toda la arquitectura, además de ser robusto ante errores y
// fácil de mantener/extender siguiendo buenas prácticas de diseño

// Context keys personalizados para evitar colisiones
type ctxKey string

const (
	CtxRequestID ctxKey = "requestID"
	CtxUserID    ctxKey = "userID"
)

// LoggerStarter devuelve un middleware Fiber para logging HTTP
func LoggerStarter(logPath string, l *logger.Logger) fiber.Handler {

	// fallback si no se pasa loPath
	if logPath == "" {
		logPath = "./logs/http.log"
	}

	// Rotacion automatica de logs
	fileWriter := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    10, // MB
		MaxBackups: 5,
		MaxAge:     30, // días
		Compress:   true,
	}

	// Salida dual: archivo + consola
	output := io.MultiWriter(os.Stdout, fileWriter)

	reqIDMiddleware := requestid.New()

	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Ejecuta middleware de requestID
		if err := reqIDMiddleware(c); err != nil {
			log.Printf("Error generando requestID: %v", err)
		}

		reqID := c.Locals("requestid")
		c.Locals(CtxRequestID, reqID)

		// Continuar con el request
		err := c.Next()
		duration := time.Since(start)

		// Campos para log
		userID := c.Locals(CtxUserID)
		ip := c.IP()
		status := c.Response().StatusCode()
		method := c.Method()
		path := c.Path()

		// Construir log legible
		logLine := fmt.Sprintf("[%s] %s %s status=%d latency=%s ip=%s",
			time.Now().Format("2006-01-02 15:04:05"),
			method, path,
			status,
			duration,
			ip,
		)

		if reqID != nil {
			logLine += fmt.Sprintf(" requestID=%s", reqID.(string))
		}
		if userID != nil {
			logLine += fmt.Sprintf(" userID=%v", userID)
		}

		// Escribir en archivo + consola
		if _, writeErr := output.Write([]byte(logLine + "\n")); writeErr != nil {
			log.Printf("Error escribiendo log: %v", writeErr)
		}

		// Pasar log también al logger de la capa de infraestructura (Zap)
		if l != nil {
			ctx := context.Background()

			// requestID
			if reqID != nil {
				ctx = logger.WithReqID(ctx, reqID.(string))
			}

			// userID como int
			if userID != nil {
				if uid, ok := userID.(int); ok {
					ctx = logger.WithUserID(ctx, uid)
				}
			}

			l.Info(ctx, "%s", logLine)
		}

		return err
	}
}

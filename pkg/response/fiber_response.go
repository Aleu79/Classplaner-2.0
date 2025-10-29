package response

import (
	"classplanner/internal/infrastructure/logger"
	"time"

	"github.com/gofiber/fiber/v2"
)

type JSONResponse struct {
	Timestamp  string      `json:"timestamp"`
	Path       string      `json:"path,omitempty"`
	Status     string      `json:"status"`
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	Error      interface{} `json:"error,omitempty"`
}

// New crea una respuesta genérica
func New(c *fiber.Ctx, status string, code int, message string, data interface{}, err interface{}) error {
	res := JSONResponse{
		Timestamp:  time.Now().Format(time.RFC3339),
		Path:       c.Path(),
		Status:     status,
		StatusCode: code,
		Message:    message,
		Data:       data,
		Error:      err,
	}
	return c.Status(code).JSON(res)
}

// Respuetas exitosas
func Success(c *fiber.Ctx, message string, data interface{}) error {
	return New(c, "success", fiber.StatusOK, message, data, nil)
}

func Created(c *fiber.Ctx, message string, data interface{}) error {
	return New(c, "success", fiber.StatusCreated, message, data, nil)
}

// Para endpoints sin contenido (ej. DELETE)
func NoContent(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}

// Paginadas
func Paginated(c *fiber.Ctx, message string, data interface{}, page, limit, total int) error {
	meta := fiber.Map{
		"page":  page,
		"limit": limit,
		"total": total,
	}
	return New(c, "success", fiber.StatusOK, message, fiber.Map{
		"items": data,
		"meta":  meta,
	}, nil)
}

// Respuestas de error

func Error(c *fiber.Ctx, code int, message string, err interface{}) error {
	return New(c, "error", code, message, nil, err)
}

func BadRequest(c *fiber.Ctx, message string, err interface{}) error {
	return Error(c, fiber.StatusBadRequest, message, err)
}

func Unauthorized(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusUnauthorized, message, nil)
}

func Forbidden(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusForbidden, message, nil)
}

func NotFound(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusNotFound, message, nil)
}

func Conflict(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusConflict, message, nil)
}

func InternalError(c *fiber.Ctx, message string, err interface{}) error {
	return Error(c, fiber.StatusInternalServerError, message, err)
}

// Con loggind integrado
func LogAndError(c *fiber.Ctx, log *logger.Logger, code int, message string, err error) error {
	if log != nil {
		log.Error(c.Context(), "API error: %s - %v", message, err)
	}
	return Error(c, code, message, err.Error())
}

// Respuesta de validacion
func ValidationError(c *fiber.Ctx, validationErrs interface{}) error {
	return Error(c, fiber.StatusUnprocessableEntity, "Validation failed", validationErrs)
}

// Recomendaciones de uso
// Siempre usar las respuestas Success o Error en tus handlers de Fiber.
// Para endpoints POST usamos Created, para DELETE usamos NoContent.
// Para colecciones usamos Paginated.
// Para errores internos usamos LogAndError para tener logging estructurado.
// Para validaciones usamos ValidationError.

package errorsper

import (
	"fmt"
)

// AppError es el tipo de error central de la aplicación.
// usar **en repository y service**, nunca directamente en handlers.
// Contiene información del error real, un mensaje amigable y contexto
type AppError struct {
	Message    string      // mensaje amigable para frontend o logs
	Err        error       // error original (DB, validación, etc)
	Context    string      // lugar donde ocurrió (Repo, Service)
	StatusCode int         // opcional: código HTTP asociado (informativo)
	Metadata   interface{} // opcional: datos extra que quieras pasar
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v [Context: %s]", e.Message, e.Err, e.Context)
	}
	return fmt.Sprintf("%s [Context: %s]", e.Message, e.Context)
}

// Helper para errores internos de la app
func ErrInternal(err error, context string) *AppError {
	return &AppError{
		Message:    "Internal server error",
		Err:        err,
		Context:    context,
		StatusCode: 500,
	}
}

// Helper para errores de validación
func ErrValidation(message string, context string) *AppError {
	return &AppError{
		Message:    message,
		Context:    context,
		StatusCode: 422,
	}
}

// Helper para errores de not found
func ErrNotFound(message string, context string) *AppError {
	return &AppError{
		Message:    message,
		Context:    context,
		StatusCode: 404,
	}
}

// Helper para conflictos
func ErrConflict(message string, context string) *AppError {
	return &AppError{
		Message:    message,
		Context:    context,
		StatusCode: 409,
	}
}

package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// Imprime todas las rutas registradas en la app Fiber
func PrintRoutes(app *fiber.App) {
	fmt.Println("rutas registradas en fiber:")

	routes := app.GetRoutes()
	for _, route := range routes {
		fmt.Printf(" - %s %s\n", route.Method, route.Path)
	}

	fmt.Println("fsin de rutas")
}

package middleware

import (
	"io"
	"log"
	"os"

	"classplanner/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func LoggerStarter() fiber.Handler {
	// load dotenv files
	utils.LoadEnv()
	// Custom File Writer
	file, err := os.OpenFile(os.Getenv("LOGGER_PATH"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	bifurcation := io.MultiWriter(os.Stdout, file)

	var loggerConfig = logger.Config{
		Format:     "${time},${status},${method},${path},${latency}\n",
		TimeFormat: "2006-01-02-15:04:05",
		TimeZone:   "Argentina/Buenos_Aires",
		Output:     bifurcation,
	}

	return logger.New(loggerConfig)
}

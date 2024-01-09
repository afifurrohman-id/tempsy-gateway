package main

import (
	"os"
	"path"

	"github.com/afifurrohman-id/tempsy-gateway/pkg/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func init() {
	if os.Getenv("APP_ENV") != "production" {
		if err := godotenv.Load(path.Join("configs", ".env")); err != nil {
			log.Panic(err)
		}
	}
}

func main() {
	var (
		app = fiber.New(fiber.Config{
			BodyLimit:    30 << 20, // 30MB
			ErrorHandler: middleware.CatchServerError,
		})
		port = os.Getenv("PORT")
	)

	app.Use(middleware.Cors, recover.New(), compress.New(), logger.New(), favicon.New(), middleware.ProxyGateway)

	if err := app.Listen(":" + port); err != nil {
		log.Panic(err)
	}
}

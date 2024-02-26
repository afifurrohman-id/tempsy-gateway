package main

import (
	"os"
	"path"

	"github.com/afifurrohman-id/tempsy-gateway/internal/gateway/utils"
	"github.com/afifurrohman-id/tempsy-gateway/pkg/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func init() {
	if os.Getenv("APP_ENV") != "production" {
		utils.Check(godotenv.Load(path.Join("configs", ".env")))
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

	utils.Check(app.Listen(":" + port))
}

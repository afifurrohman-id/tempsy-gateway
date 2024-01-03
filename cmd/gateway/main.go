package main

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/proxy"
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
		app  = fiber.New()
		port = os.Getenv("PORT")
	)

	app.Use(logger.New(), func(ctx *fiber.Ctx) error {
		hostname := fmt.Sprintf("%s:%s", os.Getenv("HOSTNAME"), port)

		if strings.Contains(ctx.Get(fiber.HeaderAccept), fiber.MIMEApplicationJSON) || strings.Contains(ctx.Get(fiber.HeaderAccept), fiber.MIMEApplicationJSONCharsetUTF8) {
			log.Info("OK")
			return proxy.DomainForward(hostname, os.Getenv("SERVER_URL"))(ctx)
		}

		return proxy.DomainForward(hostname, os.Getenv("CLIENT_URL"))(ctx)
	})

	if err := app.Listen(":" + port); err != nil {
		log.Panic(err)
	}
}

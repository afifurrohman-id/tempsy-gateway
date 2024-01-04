package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func CatchServerError(ctx *fiber.Ctx, err error) error {
	log.Error("Server Error - ", err)
	return ctx.Status(fiber.StatusInternalServerError).SendString("Sorry, something went wrong!, we working to fix it")
}

package middleware

import (
	"strings"

	"github.com/afifurrohman-id/tempsy-gateway/internal/gateway/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

var Cors = cors.New(cors.Config{
	AllowMethods: strings.Join(auth.AllowedHttpMethod, ","),
	AllowHeaders: strings.Join([]string{fiber.HeaderContentType, fiber.HeaderContentLength, fiber.HeaderAccept, fiber.HeaderUserAgent, fiber.HeaderAcceptEncoding, fiber.HeaderAcceptCharset, fiber.HeaderAuthorization, fiber.HeaderOrigin, fiber.HeaderLocation, fiber.HeaderKeepAlive}, ","),
})

package middleware

import (
	"os"
	"regexp"
	"strings"

	"github.com/afifurrohman-id/tempsy-gateway/internal/gateway/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
)

func ProxyGateway(ctx *fiber.Ctx) error {
	// exclude path: /files/{username}/public/{filename}
	match, err := regexp.MatchString(`^/files/[a-zA-Z0-9_-]+/public/[a-zA-Z0-9_-]+\.+[a-zA-Z0-9_-]+$`, ctx.OriginalURL())
	utils.Check(err)

	if strings.Contains(ctx.Get(fiber.HeaderAccept), fiber.MIMEApplicationJSON) || strings.Contains(ctx.Get(fiber.HeaderAccept), fiber.MIMEApplicationJSONCharsetUTF8) || match {
		return proxy.Forward(os.Getenv("SERVER_URL") + ctx.OriginalURL())(ctx)
	}

	return proxy.Forward(os.Getenv("CLIENT_URL") + ctx.OriginalURL())(ctx)
}

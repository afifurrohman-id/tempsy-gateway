package auth

import "github.com/gofiber/fiber/v2"

var AllowedHttpMethod = []string{fiber.MethodGet, fiber.MethodDelete, fiber.MethodOptions, fiber.MethodPut, fiber.MethodPost}

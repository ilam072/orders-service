package rest

import "github.com/gofiber/fiber/v2"

func errorResponse(message string) fiber.Map {
	return fiber.Map{
		"status":  "error",
		"message": message,
	}
}

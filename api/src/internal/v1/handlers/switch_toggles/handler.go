package switch_toggles

import (
	"github.com/HlapovErop/MarkBot/src/internal/utils/logger"
	"github.com/HlapovErop/MarkBot/src/internal/utils/toggles"
	"github.com/gofiber/fiber/v2"
)

const ROUTE = "/switch_toggles"

func Handler(ctx *fiber.Ctx) error {
	body := new(Body)
	if err := ctx.BodyParser(body); err != nil {
		logger.GetLogger().Error("Error parsing body: ")
		return ctx.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Required fields: toggles",
		})
	}

	toggleStorage := toggles.GetTogglesStorage()
	for toggle, value := range body.Toggles {
		toggleStorage.Set(toggle, value)
	}

	return ctx.JSON(fiber.Map{
		"status":  "success",
		"message": "User registered",
		"id":      "1",
	})
}

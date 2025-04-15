package middlewares

import (
	"github.com/HlapovErop/MarkBot/src/internal/models"
	"github.com/gofiber/fiber/v2"
	"slices"
)

func GetByRole(role int64) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		user := ctx.Locals("User").(*models.User)
		if slices.Contains(user.GetRoles(), role) {
			return ctx.Next()
		} else {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "You don't have permission to do this action",
			})
		}
	}
}

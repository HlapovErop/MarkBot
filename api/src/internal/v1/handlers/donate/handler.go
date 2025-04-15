package donate

import (
	"github.com/HlapovErop/MarkBot/src/internal/models"
	user_storage "github.com/HlapovErop/MarkBot/src/internal/storage/user"
	"github.com/HlapovErop/MarkBot/src/internal/utils/logger"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

const ROUTE = "/donate"

func Handler(ctx *fiber.Ctx) error {
	body := new(Body)
	if err := ctx.BodyParser(body); err != nil {
		logger.GetLogger().Error("Error parsing body: ")
		return ctx.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Required fields: points, user_id",
		})
	}

	user := ctx.Locals("User").(*models.User)

	if user.GetID() == body.UserID {
		return ctx.Status(400).JSON(fiber.Map{
			"status":  "success",
			"message": "You can't donate to yourself",
		})
	}

	user, err := user_storage.GetByID(user.GetID())
	if err != nil {
		logger.GetLogger().Error("Error logging in user: ", zap.Error(err))
		return ctx.Status(401).JSON(fiber.Map{
			"status":  "error",
			"message": "User not found",
		})
	}

	if user.Points <= 0 || user.Points < body.Points {
		return ctx.Status(400).JSON(fiber.Map{
			"status":  "success",
			"message": "You don't have enough points",
		})
	}

	if body.Points <= 0 {
		return ctx.Status(400).JSON(fiber.Map{
			"status":  "success",
			"message": "You need to donate at least 1 point",
		})
	}

	err = user_storage.DonatePoints(user.GetID(), body.UserID, body.Points)
	if err != nil {
		logger.GetLogger().Error("Error donating points: ", zap.Error(err))
		return ctx.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Error donating points",
		})
	}

	return ctx.JSON(fiber.Map{
		"status":  "success",
		"message": "Points donated successfully",
	})
}

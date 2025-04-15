package who_am_i

import (
	"github.com/HlapovErop/MarkBot/src/internal/models"
	"github.com/HlapovErop/MarkBot/src/internal/utils"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

const ROUTE = "/who_am_i"

func Handler(ctx *fiber.Ctx) error {
	user := ctx.Locals("User").(*models.User)
	user = user.GetByID(user.GetID())
	if err := user.GetById(); err != nil {
		utils.GetLogger().Error("Error logging in user: ", zap.Error(err))
		return ctx.Status(401).JSON(fiber.Map{
			"status":  "error",
			"message": "User not found",
		})
	}

	session := &models.Session{
		UserID:    u.GetID(),
		UserAgent: ctx.Get("User-Agent"),
		IP:        ctx.IP(),
		Roles:     u.GetRoles(),
	}

	sessionID, err := models.SetSession(ctx.Context(), session)
	if err != nil {
		utils.GetLogger().Error("Error setting session: ", zap.Error(err))
		return ctx.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Error setting session",
		})
	}

	return ctx.JSON(fiber.Map{
		"status":  "success",
		"message": "User logged in",
		"session": sessionID,
	})
}

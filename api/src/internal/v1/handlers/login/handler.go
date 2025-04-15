package login

import (
	"github.com/HlapovErop/MarkBot/src/internal/models"
	"github.com/HlapovErop/MarkBot/src/internal/storage/user"
	"github.com/HlapovErop/MarkBot/src/internal/utils/logger"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

const ROUTE = "/login"

type Body struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Handler(ctx *fiber.Ctx) error {
	body := new(Body)
	if err := ctx.BodyParser(body); err != nil {
		logger.GetLogger().Error("Error parsing body: ")
		return ctx.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Email or password was not transmitted",
		})
	}

	u := &models.User{Email: body.Email, Password: body.Password}
	if err := user.Login(u); err != nil {
		logger.GetLogger().Error("Error logging in user: ", zap.Error(err))
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
		logger.GetLogger().Error("Error setting session: ", zap.Error(err))
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

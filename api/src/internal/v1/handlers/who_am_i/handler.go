package who_am_i

import (
	"github.com/HlapovErop/MarkBot/src/internal/models"
	user_storage "github.com/HlapovErop/MarkBot/src/internal/storage/user"
	"github.com/HlapovErop/MarkBot/src/internal/utils"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

const ROUTE = "/who_am_i"

func Handler(ctx *fiber.Ctx) error {
	user := ctx.Locals("User").(*models.User)
	user, err := user_storage.GetByID(user.GetID())
	if err != nil {
		utils.GetLogger().Error("Error logging in user: ", zap.Error(err))
		return ctx.Status(401).JSON(fiber.Map{
			"status":  "error",
			"message": "User not found",
		})
	}

	// Фактически просто берем и переносим поля юзера в новую структуру, а уже она попадет в вывод. Фактически таким образом ограничили вывод полей
	safeUser := outUser{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Roles:     user.Roles,
		Points:    user.Points,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	return ctx.JSON(fiber.Map{
		"status": "success",
		"user":   safeUser,
	})
}

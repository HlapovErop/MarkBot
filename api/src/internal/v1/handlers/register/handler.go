package register

import (
	"github.com/HlapovErop/MarkBot/src/database/postgresql"
	"github.com/HlapovErop/MarkBot/src/internal/models"
	"github.com/HlapovErop/MarkBot/src/internal/utils/logger"
	"github.com/HlapovErop/MarkBot/src/internal/utils/toggles"
	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
	"go.uber.org/zap"
)

const ROUTE = "/register"

type Body struct {
	Name     string `json:"name" binding:"required"`
	Surname  string `json:"surname" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Handler(ctx *fiber.Ctx) error {
	canRegister, _ := toggles.GetTogglesStorage().Get("CanRegister")
	if !canRegister.(bool) {
		return ctx.Status(401).JSON(fiber.Map{
			"status":  "error",
			"message": "Registration not allowed",
		})
	}

	body := new(Body)
	if err := ctx.BodyParser(body); err != nil {
		logger.GetLogger().Error("Error parsing body: ")
		return ctx.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Required fields: name, surname, email, password",
		})
	}

	u := &models.User{
		Email:    body.Email,
		Password: body.Password,
		Name:     body.Name,
		Surname:  body.Surname,
		Roles:    make(pq.Int64Array, models.RoleStudent), // регаться могут только студенты, учителя из сидов или бд (иногда делают для этого отдельную админку, но в данном проекте, где увеличения кол-ва учителей вообще не предполагается (им буду только я), в ней нет необходимости
	}
	result := postgresql.GetDB().Create(u)
	if result.Error != nil {
		logger.GetLogger().Error("Error register in user: ", zap.Error(result.Error))
		return ctx.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "User not registered",
		})
	}

	return ctx.JSON(fiber.Map{
		"status":  "success",
		"message": "User registered",
		"id":      u.ID,
	})
}

package user

import (
	"errors"
	"github.com/HlapovErop/MarkBot/src/database/postgresql"
	"github.com/HlapovErop/MarkBot/src/internal/models"
	"github.com/HlapovErop/MarkBot/src/internal/utils"
	"go.uber.org/zap"
)

func Login(u *models.User) error {
	if u.Email == "" || u.Password == "" {
		err := errors.New("email and password are required")
		utils.GetLogger().Error("Login error", zap.Error(err))
		return err
	}

	err := postgresql.GetDB().Where("email = ? AND password = ?", u.Email, u.Password).First(u).Error
	if err != nil {
		utils.GetLogger().Error("Login error", zap.Error(err))
		return err
	}

	return nil
}

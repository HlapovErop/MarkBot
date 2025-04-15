package user

import (
	"errors"
	"fmt"
	"github.com/HlapovErop/MarkBot/src/database/postgresql"
	"github.com/HlapovErop/MarkBot/src/internal/models"
	"github.com/HlapovErop/MarkBot/src/internal/utils/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func Login(u *models.User) error {
	if u.Email == "" || u.Password == "" {
		err := errors.New("email and password are required")
		logger.GetLogger().Error("Login error", zap.Error(err))
		return err
	}

	err := postgresql.GetDB().Where("email = ? AND password = ?", u.Email, u.Password).First(u).Error
	if err != nil {
		logger.GetLogger().Error("Login error", zap.Error(err))
		return err
	}

	return nil
}

func GetByID(id uint) (*models.User, error) {
	if id == 0 {
		err := errors.New("user ID cannot be zero")

		logger.GetLogger().Error("invalid user ID", zap.Error(err))
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	var user models.User

	err := postgresql.GetDB().First(&user, id).Error
	if err != nil {
		// Тут обработка ошибки, если просто нет такого юзера
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.GetLogger().Warn("user not found", zap.Uint("id", id))
			return nil, fmt.Errorf("user with ID %d not found: %w", id, err)
		}

		// А тут все остальные ошибки
		logger.GetLogger().Error("database error", zap.Uint("user_id", id), zap.Error(err))
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

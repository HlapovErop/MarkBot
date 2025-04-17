package user

import (
	"errors"
	"fmt"
	"github.com/HlapovErop/MarkBot/src/database/postgresql"
	"github.com/HlapovErop/MarkBot/src/internal/models"
	"github.com/HlapovErop/MarkBot/src/internal/utils/logger"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"regexp"
)

func Login(u *models.User) error {
	if u.Email == "" || u.Password == "" {
		err := errors.New("email and password are required")
		logger.GetLogger().Error("Login error", zap.Error(err))
		return err
	}

	password := u.Password
	err := postgresql.GetDB().Where("email = ?", u.Email).First(u).Error
	if err != nil {
		logger.GetLogger().Error("Login error", zap.Error(err))
		return err
	}
	if !u.CheckPassword(password) {
		logger.GetLogger().Error("Login error", zap.Error(err))
		return errors.New("password not valid")
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

func DonatePoints(senderID, receiverID uint, points int64) error {
	if senderID == receiverID {
		return errors.New("cannot donate points to yourself")
	}
	if points <= 0 {
		return errors.New("points amount must be positive")
	}

	db := postgresql.GetDB()

	// Да начнется транзакция
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var senderPoints int64
	var receiverExists bool

	// Запрос 1: Проверяем обоих пользователей и блокируем отправителя, чтобы не возникло гонки ресурсов(race condition)
	row := tx.Raw(`
        SELECT 
            (SELECT points FROM users WHERE id = ? FOR UPDATE) as sender_points,
            EXISTS(SELECT 1 FROM users WHERE id = ?) as receiver_exists
    `, senderID, receiverID).Row()

	if err := row.Scan(&senderPoints, &receiverExists); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to check users: %w", err)
	}

	if !receiverExists {
		tx.Rollback()
		return errors.New("receiver not found")
	}
	if senderPoints < points {
		tx.Rollback()
		return fmt.Errorf("not enough points (available: %d, requested: %d)", senderPoints, points)
	}

	// Запрос 2: Списание у отправителя
	if err := tx.Exec(`
        UPDATE users SET points = points - ? 
        WHERE id = ?
    `, points, senderID).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to deduct points: %w", err)
	}

	// Запрос 3: Начисление получателю
	if err := tx.Exec(`
        UPDATE users SET points = points + ? 
        WHERE id = ?
    `, points, receiverID).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to add points: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("commit failed: %w", err)
	}

	return nil
}

func Create(u *models.User) (uint, error) {
	if !isValidEmail(u.Email) {
		return 0, errors.New("email format is invalid")
	}

	if len(u.Password) < 8 {
		return 0, errors.New("password must be at least 8 characters")
	}

	if !u.ValidateNameSurname() {
		return 0, errors.New("name and surname combination is not allowed")
	}

	if exists, err := checkEmailExists(u.Email); err != nil {
		return 0, err
	} else if exists {
		return 0, errors.New("email already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	u.Password = string(hashedPassword)

	result := postgresql.GetDB().Create(u)
	if result.Error != nil {
		return 0, result.Error
	}

	return u.ID, nil
}

// Вспомогательные методы

func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func checkEmailExists(email string) (bool, error) {
	var count int64
	db := postgresql.GetDB()
	err := db.Model(&models.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

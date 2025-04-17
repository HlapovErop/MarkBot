package postgresql

import (
	"fmt"
	"github.com/HlapovErop/MarkBot/src/internal/models"
	logger2 "github.com/HlapovErop/MarkBot/src/internal/utils/logger"
	"github.com/HlapovErop/MarkBot/src/internal/utils/toggles"
	"go.uber.org/zap"
	"gorm.io/gorm/logger"
	"os"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var instance *gorm.DB // singleton - паттерн, в котором экземпляр класса может быть только один. Здесь он для работы с БД. Во время работы GetConnection() мы создаем ЕДИНОЖДЫ создаем экземпляр gorm.DB и сохраняем его в переменную connection, доступную далее по тому же методу. Больше никаких точек входа в бд создано не будет, тк базы данных это не любят
var once sync.Once

func GetDB() *gorm.DB {
	once.Do(connectDB)
	return instance
}

func connectDB() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Asia/Shanghai",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // У БД свой простенький логгер. Это связано как с реализацией GORM, так и с моим желанием не засорять файл логов - консольного вывода в данном проекте будет достаточно
	})

	if err != nil {
		logger2.GetLogger().Fatal("Failed to connect to database. \n", zap.Error(err))
	}

	logger2.GetLogger().Info("DB connected")

	logger2.GetLogger().Info("running migrations")
	db.AutoMigrate(&models.User{})

	instance = db

	seedsInstalled, _ := toggles.GetTogglesStorage().Get("SeedsInstalled")
	if !seedsInstalled.(bool) {
		logger2.GetLogger().Info("running seeds")
		seeds()
		toggles.GetTogglesStorage().Set("seedsInstalled", true)
	}
}

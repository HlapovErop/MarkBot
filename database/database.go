package database

import (
	"fmt"
	"github.com/HlapovErop/MarkBot/internal/models"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Dbinstance struct { // singleton - паттерн, в котором экземпляр класса может быть только один. Здесь он для работы с БД. Во время работы ConnectDb() мы создаем экземпляр класса Dbinstance и сохраняем его в переменную connection, доступную по всему проекту. Больше никаких точек входа в бд создано не будет, тк базы данных это не любят
	connection *gorm.DB
}

func (db Dbinstance) GetConnection() *gorm.DB {
	return db.connection
}

var DB Dbinstance

func ConnectDb() {
	dsn := fmt.Sprintf(
		"host=db user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Asia/Shanghai",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
		os.Exit(2)
	}

	log.Println("DB connected")
	db.Logger = logger.Default.LogMode(logger.Info)

	log.Println("running migrations")
	db.AutoMigrate(&models.User{})

	DB = Dbinstance{
		connection: db,
	}
}

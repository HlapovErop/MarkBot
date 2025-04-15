package redis

import (
	"context"
	"fmt"
	"github.com/HlapovErop/MarkBot/src/internal/utils"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"os"
	"sync"
)

var instance *redis.Client // singleton - паттерн, в котором экземпляр класса может быть только один. Здесь он для работы с БД. Во время работы ConnectDb() мы создаем экземпляр класса Dbinstance и сохраняем его в переменную connection, доступную по всему проекту. Больше никаких точек входа в бд создано не будет, тк базы данных это не любят
var once sync.Once

func GetRedis() *redis.Client {
	once.Do(connectRedis)
	return instance
}

func connectRedis() {
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:6379", os.Getenv("REDIS_HOST")),
	})

	err := client.Ping(context.Background()).Err() // пингуем Redis для проверки подключения. Если не подключено, то выйдет ошибка
	if err != nil {
		utils.GetLogger().Fatal("Failed to connect to database. \n", zap.Error(err))
	}

	utils.GetLogger().Info("Redis connected")

	instance = client
}

package database

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"os"
)

type RedisInstance struct { // singleton - паттерн, в котором экземпляр класса может быть только один. Здесь он для работы с БД. Во время работы ConnectDb() мы создаем экземпляр класса Dbinstance и сохраняем его в переменную connection, доступную по всему проекту. Больше никаких точек входа в бд создано не будет, тк базы данных это не любят
	connection *redis.Client
}

func (db RedisInstance) GetConnection() *redis.Client {
	return db.connection
}

var Redis RedisInstance

func ConnectRedis() {
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:6379", os.Getenv("REDIS_HOST")),
	})

	err := client.Ping(context.Background()).Err() // пингуем Redis для проверки подключения. Если не подключено, то выйдет ошибка
	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
		os.Exit(2)
	}

	log.Println("Redis connected")

	Redis = RedisInstance{
		connection: client,
	}
}

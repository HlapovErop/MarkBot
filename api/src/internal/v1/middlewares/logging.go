package middlewares

import (
	"github.com/HlapovErop/MarkBot/src/internal/utils/logger"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// LoggingMiddleware логирует каждый входящий запрос
func LoggingMiddleware(ctx *fiber.Ctx) error {
	start := time.Now()

	// Выполнение следующего middleware или обработчика
	err := ctx.Next()

	// Вычисление времени выполнения
	duration := time.Since(start)

	// Получение статуса ответа
	status := ctx.Response().StatusCode()

	// Логирование информации о запросе
	logger.GetLogger().Log(logger.RequestLevel, "Request", // Эта запись про запрос, на нее выделил отдельный уровень логирования, чтобы можно было легко их отследить, такое тоже может быть важно при выносе метрик
		zap.String("method", ctx.Method()),
		zap.String("path", ctx.Path()),
		zap.Int("status", status),
		zap.Duration("duration", duration),
		zap.String("ip", ctx.IP()),
		zap.String("user_agent", ctx.Get("User-Agent")), // Пример, что можно накинуть. В большинстве случаев хранить логи, связанные с юзерами - моветон. Но такие вещи применяются глубже - например для защиты от ddos (если nginx под это не настроен) или анализа, с каких систем юзеры чаще сидят
	)

	return err
}

package main

import (
	"fmt"
	"github.com/HlapovErop/MarkBot/src/consts"
	"github.com/HlapovErop/MarkBot/src/database/postgresql"
	"github.com/HlapovErop/MarkBot/src/database/redis"
	"github.com/HlapovErop/MarkBot/src/internal/models"
	"github.com/HlapovErop/MarkBot/src/internal/utils/logger"
	"github.com/HlapovErop/MarkBot/src/internal/utils/toggles"
	"github.com/HlapovErop/MarkBot/src/internal/v1/handlers/login"
	"github.com/HlapovErop/MarkBot/src/internal/v1/handlers/register"
	"github.com/HlapovErop/MarkBot/src/internal/v1/handlers/switch_toggles"
	"github.com/HlapovErop/MarkBot/src/internal/v1/handlers/who_am_i"
	"github.com/HlapovErop/MarkBot/src/internal/v1/middlewares"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"os"
)

// Точка входа в проект. Именно с этого файла стоит начать читать, чтобы сложилась полная картина
// В одном проекте может быть несколько точек входа, например одна для веб-сервера, а другая для крон-задач (крон-задачи - функции, исполняемые по расписанию и не привязанные к веб-проекту, например рассылки на почту или таймеры, подробнее читай про фоновые задачи)
// Все точки входа обычно находятся в cmd/, дальше делятся на отдельные пакеты, например cmd/api и cmd/cron, внутри которых один единственный файл по типу main.go
// Но не забываем, что это гошка, а не Ruby, Java или PHP со своими закостенелыми фреймворками, где есть четкие инструкции оформления и привязка к исполняемым директориям
// Из-за своей простоты гошка позволяет играть в креатив, и все, что будет рабочее, читаемое и документированное, можно считать нормальным проектом
// У air есть одна беда - если он не смог сбилдить проект, то использует последний прошедший билд. Ловить это сложно
func main() {

	// Это называется прогрев. Сейчас приложение только запускается, ему нужно установить все коннекты, в некоторых случаях подгрузить данные в редис, выполнить сиды и миграции. Если прогрева не будет, то все коннекты произойдут при первом запросе, который может изрядно подвиснуть
	logger.GetLogger()
	postgresql.GetDB()
	redis.GetRedis()
	err := toggles.GetTogglesStorage().Set("CanRegister", consts.INIT_CAN_REGISTER)
	if err != nil {
		logger.GetLogger().Fatal("Can't set toggle CanRegister", zap.Error(err))
	}

	app := fiber.New(
	//fiber.Config{
	//	Prefork: true, // Хорошая вещь для увеличения производительности на проде, но только не при air и других лайф-релодах! Если air потеряет контроль над дочерними процессами, то у тебя утечет память, останутся зомби-процессы, а трафик может улететь в старую версию
	//},
	)

	// Подключаем мидлвару для логирования
	app.Use(middlewares.LoggingMiddleware)

	// Мидлвара для установки /v1 в маршруте, и заодно добавляет в контекст, что версия 1. Глубже в бизнес-логике это может пригодиться, например в общих сторэджах или утилитах, где версия немного влияет на работу и бизнес-логику
	v1 := app.Group("/v1", func(ctx *fiber.Ctx) error {
		ctx.Locals("Version", "v1")
		return ctx.Next()
	})

	// Здесь мы регистрируем маршруты (роутинг).

	// В гошке, если у вас совсем простой проект, или вы любитель антипаттернов по типу класса бога, можно использовать анонимные функции, как в этом хэндлере
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("Project has been launched. Well done")
	})

	// А если не простой, то лучше выносить логику в отдельные места
	v1.Post(login.ROUTE, login.Handler)
	v1.Post(register.ROUTE, register.Handler)

	// Вот так просто и ненавязчиво сказали, что плюс ко всему ты должен пройти мидлвару авторизации и другие прежде чем получить доступ к хэндлеру
	v1.Post(who_am_i.ROUTE, middlewares.AuthMiddleware, who_am_i.Handler)
	v1.Post(switch_toggles.ROUTE, middlewares.AuthMiddleware, middlewares.GetByRole(models.RoleTeacher), switch_toggles.Handler)

	apiHost := os.Getenv("API_HOST")
	if apiHost == "" {
		apiHost = consts.DEFAULT_HOST
	}
	logger.GetLogger().Info(fmt.Sprintf("Server started: http://%s", apiHost))

	if err := app.Listen(apiHost); err != nil {
		logger.GetLogger().Fatal("Error in fiber app.Listen: ", zap.Error(err))
	}
}

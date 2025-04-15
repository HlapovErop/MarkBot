package middlewares

import (
	"github.com/HlapovErop/MarkBot/src/internal/models"
	"github.com/HlapovErop/MarkBot/src/internal/utils"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// В первой версии сделаем авторизацию через сессии. Во всяких обучалках сессии меньше раскрыты, чем токены. Поэтому решил зайти с них.
// Не скажу, кто лучше: токены или сессии. И то, и другое можно взломать, если юзер криворукий.
// Основное отличие в том, что сессии легко отозвать, тк хранят информацию на стороне сервера - просто удали из БД. Но из-за этого кол-во запросов в БД возрастает. Для этого я храню их в Redis - он быстрее, а сессии потерять не страшно.
// В это время токены вообще не требуют запросов в БД - необходимую инфу они могут хранить внутри - просто дешифруй. Но в случае взлома отозвать их не получится. И криворукий юзер будет ждать, пока мошенник не наиграется с его аккаунтом.
// Конечно, в токенах можно делать BlackList, или двухфакторку на важных операциях. Но тогда теряется их преимущество, и вообще BlackList считаю костылем.
// PS - никто, вообще никто не может запретить пользоваться несколькими видами авторизации одновременно ;)
func AuthMiddleware(ctx *fiber.Ctx) error {
	sessionID := ctx.Get("sessionID")
	utils.GetLogger().Info(sessionID)
	if sessionID == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error1",
			"message": "Unauthorized",
		})
	}

	session, err := models.GetSession(ctx.Context(), sessionID)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  err.Error(),
			"message": "Unauthorized",
		})
	}

	// Здесь выполняем проверки метаданных. Это отличительная черта сессий - они имеют привязку к месту, откуда произошел логин. Влияет даже браузер. Можно придумать и больше метаданных, но в данном примере нет необходимости
	if session == nil || session.IP != ctx.IP() || session.UserAgent != ctx.Get("User-Agent") {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error3",
			"message": "Unauthorized",
		})
	}

	// Почему не используем postgresql для получения остальных полей юзера? Потому что в большинстве случаев нам не нужны эти самые поля. Нам нужно знать его роли для контроля доступа к ресурсам. И его ID для того же контроля доступа и возможности получить остальные поля. Незачем лишний раз нагружать БД, у нее несварение
	user := &models.User{
		Model: gorm.Model{ID: session.UserID},
		Roles: session.Roles,
	}

	// Здесь мы добавляем пользователя в локальные переменные, чтобы он был доступен по проекту через контекст. Удобно? Удобно!
	ctx.Locals("User", user)
	return ctx.Next()
}

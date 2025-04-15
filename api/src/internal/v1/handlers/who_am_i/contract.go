package who_am_i

// Структура для ограничения полей в респонсе запроса. Необходимо, чтобы в ответ непопали поля внутреннего устройства приложения и мета.
// Есть и другие способы ограничить поля, например использование zek, кастомные JSON-маршаллеры или просто незаполнение тэгов json в модельке
// Но этот способ более явный и понятный для новичков, а проект несет обучающий характер
// Другие возможности отметил? Отметил - лапки чисты)
type outUser struct {
	ID        uint    `json:"id"`
	Name      string  `json:"name"`
	Email     string  `json:"email"`
	CreatedAt string  `json:"created_at"`
	Roles     []int64 `json:"roles"`
	Points    int64   `json:"points"`
}

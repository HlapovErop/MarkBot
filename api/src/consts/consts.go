package consts

// Константы, которым не нашлось места в проекте))
// Они достаточно общие, и хорошим тоном их выносят в конфигурации. Например как здесь пакет consts, или concerns/config
// Таким образом тебе не нужно витать по всему проекту в поисках заветных магических значений, а вшивать их в код без выноса в константы(!) - моветон
const (
	DEFAULT_HOST              = "0.0.0.0:3000"
	TOGGLES_FILE_PATH         = "./src/database/toggles.json"
	INIT_CAN_REGISTER         = false
	POINTS_AFTER_REGISTRATION = 100
)

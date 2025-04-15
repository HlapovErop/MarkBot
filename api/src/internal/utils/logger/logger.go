package logger

import (
	"go.uber.org/zap" // разных-приразных логеров достаточно. Есть встроенный из библиотеки "log", есть мидлвара из "fiber", есть известный "logrus", я выбрал уберовский - он неплох
	"go.uber.org/zap/zapcore"
	"os"
	"sync"
)

var logger *zap.Logger // наш старый добрый singleton
var loggerOnce sync.Once

const logFile = "./tmp/api.log"
const RequestLevel = zapcore.Level(15)

func GetLogger() *zap.Logger {
	loggerOnce.Do(func() {
		logger = newLogger(logFile)
	})
	return logger
}

func newLogger(logFile string) *zap.Logger {
	// Создаем файл для логов
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	// Создаем консольный вывод
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

	// Создаем файловый вывод
	fileEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())

	// Устанавливаем уровень логирования
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zap.InfoLevel),
		zapcore.NewCore(fileEncoder, zapcore.AddSync(file), zap.WarnLevel), // Устанавливаем уровень логирования для файлового вывода. Незачем хранить инфо-логи в файле - никакой памяти не хватит
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))

	return logger
}

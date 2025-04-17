package toggles

import (
	"encoding/json"
	"github.com/HlapovErop/MarkBot/src/consts"
	"github.com/HlapovErop/MarkBot/src/internal/utils/logger"
	"go.uber.org/zap"
	"os"
	"sync"
)

// fileStorage реализует хранение тоглов в JSON-файле
// Это еще один прикольный способ хранить данные, которые будут редко перезаписываться, и им отлично пользуются в реальных проектах
// Файлы легче контролировать, к ним не нужно подключение для редактирования разрабам, хватит и vim/nano
// Здесь нет необходимости мониторить файл, ставить крону или еще чего:
// При первом обращении скачаем из него все тоглы в поле data, и по мере обновлений будем перезаписывать data и файл
// Таким образом файл всегда будет свежим, а между перезапусками инфа всегда будет сохраняться
// При этом тоглы могут представлять из себя что угодно, даже структуру - главное переведи в json
type fileStorage struct {
	filePath string
	data     map[string]interface{}
	mu       sync.RWMutex
}

var instance *fileStorage
var once sync.Once

func GetTogglesStorage() *fileStorage {
	once.Do(func() {
		var err error
		instance, err = newFileStorage(consts.TOGGLES_FILE_PATH)
		if err != nil {
			logger.GetLogger().Fatal("Can't init toggles storage", zap.Error(err))
		}
	})
	return instance
}

// newFileStorage создает новый экземпляр хранилища, а заодно скачает инфу из файла
func newFileStorage(filePath string) (*fileStorage, error) {
	storage := &fileStorage{
		filePath: filePath,
		data:     make(map[string]interface{}),
	}

	// Загружаем существующие данные при инициализации
	if err := storage.load(); err != nil {
		return nil, err
	}

	return storage, nil
}

// Set устанавливает значение для ключа
func (s *fileStorage) Set(key string, value interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if currentValue, ok := s.data[key]; value == currentValue && ok {
		return nil
	} else {
		s.data[key] = value
		return s.save()
	}
}

// Get возвращает значение по ключу
func (s *fileStorage) Get(key string) (interface{}, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	val, exists := s.data[key]
	return val, exists
}

// Delete удаляет ключ
func (s *fileStorage) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, key)
	return s.save()
}

// GetAll возвращает все тогглы
func (s *fileStorage) GetAll() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Возвращаем копию данных
	copy := make(map[string]interface{}, len(s.data))
	for k, v := range s.data {
		copy[k] = v
	}
	return copy
}

func (s *fileStorage) GetStringSlice(key string) ([]string, bool) {
	val, exists := s.Get(key)
	if !exists {
		return nil, false
	}

	if slice, ok := val.([]string); ok {
		return slice, true
	}

	if slice, ok := val.([]interface{}); ok {
		result := make([]string, 0, len(slice))
		for _, item := range slice {
			if s, ok := item.(string); ok {
				result = append(result, s)
			}
		}
		return result, true
	}

	return nil, false
}

// load загружает данные из файла
func (s *fileStorage) load() error {
	file, err := os.ReadFile(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // Файл не существует - это нормально при первом запуске
		}
		return err
	}

	if len(file) == 0 {
		return nil // Пустой файл
	}

	return json.Unmarshal(file, &s.data)
}

// save сохраняет данные в файл
func (s *fileStorage) save() error {
	data, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return err
	}
	logger.GetLogger().Info(s.filePath)
	return os.WriteFile(s.filePath, data, 0644)
}

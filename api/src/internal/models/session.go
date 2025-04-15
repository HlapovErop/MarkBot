package models

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/HlapovErop/MarkBot/src/database/redis"
	"github.com/google/uuid"
	"strconv"
	"time"
)

type Session struct {
	UserID      uint      `json:"user_id" redis:"user_id"`
	Roles       []int64   `json:"roles" redis:"roles"`
	UserAgent   string    `json:"user_agent" redis:"user_agent"`
	IP          string    `json:"ip" redis:"ip"`
	LastUseDate time.Time `json:"last_use_date" redis:"last_use_date"`
}

const sessionsKey = "sessions"

func GetSession(ctx context.Context, sessionID string) (*Session, error) {
	client := redis.GetRedis()

	// Получаем все поля сессии одним запросом
	sessionMap, err := client.HGetAll(ctx, sessionHashKey(sessionID)).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}
	if len(sessionMap) == 0 {
		return nil, errors.New("session not found")
	}

	session, err := parseSessionFromMap(sessionMap)
	if err != nil {
		return nil, fmt.Errorf("failed to parse session: %w", err)
	}

	// Еще я хочу обновлять дату последнего использования при каждом получении сессии, чтобы в дальнейшем их легко чистить. Почему не пошел через время жизни записи сессии? Чтобы не разлогинивать активных юзеров раз в сутки	session.LastUseDate = time.Now()
	err = client.HSet(ctx, sessionHashKey(sessionID),
		"last_use_date", session.LastUseDate.Format(time.RFC3339)).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to update last use date: %w", err)
	}

	return session, nil
}

func SetSession(ctx context.Context, session *Session) (string, error) {
	id, err := generateSessionID(ctx)
	if err != nil {
		return "", err
	}

	// Обновляем дату последнего использования. Нужно для удаления сессий, которыми давно не пользовались
	session.LastUseDate = time.Now()

	err = redis.GetRedis().HSet(ctx, sessionHashKey(id),
		"user_id", strconv.FormatUint(uint64(session.UserID), 10),
		"roles", serializeRoles(session.Roles),
		"user_agent", session.UserAgent,
		"ip", session.IP,
		"last_use_date", session.LastUseDate.Format(time.RFC3339),
	).Err()

	if err != nil {
		return "", fmt.Errorf("failed to save session: %w", err)
	}

	return id, nil
}

// Вспомогательные функции

func sessionHashKey(sessionID string) string {
	return fmt.Sprintf("%s:%s", sessionsKey, sessionID)
}

func parseSessionFromMap(data map[string]string) (*Session, error) {
	userID, err := strconv.ParseUint(data["user_id"], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid user_id: %w", err)
	}

	roles, err := deserializeRoles(data["roles"])
	if err != nil {
		return nil, fmt.Errorf("invalid roles: %w", err)
	}

	lastUseDate, err := time.Parse(time.RFC3339, data["last_use_date"])
	if err != nil {
		return nil, fmt.Errorf("invalid last_use_date: %w", err)
	}

	return &Session{
		UserID:      uint(userID),
		Roles:       roles,
		UserAgent:   data["user_agent"],
		IP:          data["ip"],
		LastUseDate: lastUseDate,
	}, nil
}

func serializeRoles(roles []int64) string {
	data, _ := json.Marshal(roles)
	return string(data)
}

func deserializeRoles(rolesStr string) ([]int64, error) {
	var roles []int64
	err := json.Unmarshal([]byte(rolesStr), &roles)
	return roles, err
}

func generateSessionID(ctx context.Context) (string, error) {
	client := redis.GetRedis()

	// Делаем несколько попыток на случай коллизии UUID, чтобы не перетереть чужую сессию. Если с пяти попыток не получилось сгенерить несуществующий UUID - это повод задуматься, что в жизни творится что-то не так
	for i := 0; i < 5; i++ {
		id := uuid.New().String()
		exists, err := client.Exists(ctx, sessionHashKey(id)).Result()
		if err != nil {
			return "", err
		}
		if exists == 0 {
			return id, nil
		}
	}
	return "", errors.New("failed to generate unique session ID")
}

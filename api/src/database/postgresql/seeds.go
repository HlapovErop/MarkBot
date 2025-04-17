package postgresql

import (
	"context"
	"fmt"
	"time"
)

func seeds() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := EnsureNameSurnameIndex(ctx)
	if err != nil {
		fmt.Printf("Failed to ensure index: %v\n", err)
	}
}

func EnsureNameSurnameIndex(ctx context.Context) error {
	// 1. Проверяем существование индекса
	exists, err := checkIndexExists(ctx, "users", "idx_users_name_surname")
	if err != nil {
		return fmt.Errorf("failed to check index existence: %w", err)
	}

	if exists {
		return nil // Индекс уже существует
	}

	// 2. Создаем индекс с контекстом
	err = GetDB().WithContext(ctx).Exec(`
        CREATE UNIQUE INDEX CONCURRENTLY 
        idx_users_name_surname 
        ON users(name, surname)
    `).Error

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("index creation timed out after 5 seconds")
		}
		return fmt.Errorf("failed to create index: %w", err)
	}

	return nil
}

func checkIndexExists(ctx context.Context, table, index string) (bool, error) {
	var indexExists bool
	err := GetDB().WithContext(ctx).Raw(`
        SELECT EXISTS (
            SELECT 1 FROM pg_indexes 
            WHERE tablename = $1 
            AND indexname = $2
        )
    `, table, index).Scan(&indexExists).Error

	if err != nil && ctx.Err() == context.DeadlineExceeded {
		return false, fmt.Errorf("index check timed out")
	}

	return indexExists, err
}

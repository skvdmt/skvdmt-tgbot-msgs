package usecase

import (
	"context"
	"time"
)

// Repository Интерфейс репозиторного слоя.
type Repository interface {
	// Остановка.
	Stop(ctx context.Context) error
	// Сохранение сообщения.
	SaveMessage(ctx context.Context, telegramUserId int,
		telegramUserName string, text string) error
	// Получение времени создания последнего сообщения пользователем.
	UserMessageCreatedAt(ctx context.Context, telegramUserId int) (*time.Time, error)
}

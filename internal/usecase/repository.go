package usecase

import (
	"context"
	"time"

	"github.com/skvdmt/skvdmt-tgbot-msgs/internal/entities"
)

// Repository Интерфейс репозиторного слоя.
type Repository interface {
	// Остановка.
	Stop(ctx context.Context) error
	// Сохранение сообщения.
	SaveMessage(ctx context.Context, telegramUserId int, message *entities.Message) error
	// Получение времени создания последнего сообщения пользователем.
	UserMessageCreatedAt(ctx context.Context, telegramUserId int) (*time.Time, error)
	// Получение сообщений из репозитория.
	UpdateMessages(ctx context.Context) ([]*entities.Message, error)
}

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
	SaveMessage(ctx context.Context, telegramUserId int, text string) error
	// Свалка пользователей.
	DumpUsers(ctx context.Context, users map[int]*entities.User) error
	// Очистка реестра пользователей в базе данных.
	CleanUsersRegistry(ctx context.Context) error
	// Получение времени создания последнего сообщения пользователем.
	UserMessageCreatedAt(ctx context.Context, telegramUserId int) (*time.Time, error)
}

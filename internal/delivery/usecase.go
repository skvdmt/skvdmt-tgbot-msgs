package delivery

import (
	"context"

	"github.com/skvdmt/skvdmt-tgbot-msgs/internal/entities"
	"github.com/skvdmt/skvdmt-tgbot-msgs/internal/entities/configs"
)

// Usecase Интерфейс сервисного слоя.
type Usecase interface {
	// Остановка.
	Stop(ctx context.Context) error
	// Обработчик команд.
	CommandHandle(ctx context.Context,
		user *entities.User,
		update *entities.Update) (configs.RequestConfig, error)
	// Обработчик каптчи.
	CaptchaHandle(ctx context.Context,
		user *entities.User,
		update *entities.Update) (configs.RequestConfig, error)
	// Обработчик сообщения.
	MessageHandle(ctx context.Context,
		user *entities.User,
		update *entities.Update) (configs.RequestConfig, error)
	// Оптиизация реестра пользователей.
	OptimizeUserRegistry(ctx context.Context) error
	// Получение пользователя.
	User(ctx context.Context, telegramUserId int) (*entities.User, error)
}

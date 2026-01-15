package delivery

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/skvdmt/skvdmt-tgbot-msgs/internal/entities"
	"github.com/skvdmt/skvdmt-tgbot-msgs/internal/entities/configs"
	"github.com/skvdmt/skvdmt-tgbot-msgs/internal/messages"
	"github.com/skvdmt/skvdmt-tgbot-msgs/internal/model"
	"github.com/skvdmt/skvdmt-tgbot-msgs/internal/usecase"
)

// App Транспортный слой.
type App struct {
	// Канал остановки ресурсов.
	exit chan struct{}
	// Ресурсы.
	sources *sync.WaitGroup
	// Сервисный слой.
	usecase Usecase
	// Клиент для запросов к Telegram Bot API.
	client *Client
	// Реестр пользователей.
	users *entities.UserRegistry
	// Раннер оптимизации реестра пользователей.
	optimizeUserRegistry *time.Ticker
}

// NewApp Конструктор.
func NewApp(ctx context.Context, users *entities.UserRegistry) (*App, error) {
	model.Logs.Info.Info("delivery layer creating")
	a := &App{
		exit:    make(chan struct{}),
		sources: &sync.WaitGroup{},
		users:   users,
		optimizeUserRegistry: time.NewTicker(time.Minute *
			time.Duration(model.Config.Timers.OptimizeUserRegistryInterval)),
	}
	var err error
	// Создание клиента для запросов к Telegram Bot API.
	a.client, err = NewClient()
	if err != nil {
		return nil, err
	}
	// Создание сервисного слоя.
	a.usecase, err = usecase.NewApp(ctx, a.users)
	if err != nil {
		return nil, err
	}
	return a, nil
}

// Start Запуск.
func (a *App) Start(ctx context.Context) error {
	model.Logs.Info.Info("delivery layer started")
	// Начало работы ресурса приложения.
	a.sources.Add(1)
	go a.cleaner(ctx)
	if err := a.client.Start(ctx); err != nil {
		return err
	}
	for u := range a.client.Updates {
		if err := a.updateHandle(ctx, u); err != nil {
			return err
		}
	}
	return nil
}

// Stop Остановка.
func (a *App) Stop(ctx context.Context) error {
	a.exit <- struct{}{}
	// Ожидание завершения работы ресурсов.
	a.sources.Wait()
	// Остановка клиента делающего запросы к Telegram Bot API.
	if err := a.client.Stop(ctx); err != nil {
		return err
	}
	// Закрытие канала остановки ресурсов.
	close(a.exit)
	// Вызов остановки сервисного слоя.
	if err := a.usecase.Stop(ctx); err != nil {
		return err
	}
	model.Logs.Info.Info("delivery layer stopped")
	return nil
}

// cleaner Обработка сигналов раннеров оптимизации реестра
// пользователей и очистки базы данных реестра пользователей.
func (a *App) cleaner(ctx context.Context) {
	model.Logs.Info.Info("user registry cleaner started")
	for {
		select {
		case <-a.exit:
			model.Logs.Info.Info("user registry cleaner stopped")
			a.sources.Done()
			return
		case <-a.optimizeUserRegistry.C:
			if err := a.usecase.OptimizeUserRegistry(ctx); err != nil {
				model.Errors <- err
			}
		}
	}
}

// updateHandle Обработка обновления.
func (a *App) updateHandle(ctx context.Context, update *entities.Update) error {
	// Получение пользователя.
	user, err := a.usecase.User(ctx, update.Message.From.Id)
	if err != nil {
		return err
	}
	user.SetTelegramUsername(update.Message.From.Username)

	// Получение конфигурации запроса к телеграм боту,
	// чтобы дать отвера на запрос пользователя.
	cfg, err := a.updateRequestConfig(ctx, user, update)
	if err != nil {
		return err
	}
	cfg.SetChatId(update.Message.Chat.Id)
	// Отправка ответа на запрос пользователя.
	_, err = a.client.Do(ctx, cfg)
	if err != nil {
		return err
	}
	return nil
}

// updateRequestConfig Получение конфигурация запроса
// к телеграм боту, чтобы дать ответ на запрос пользователя.
func (a *App) updateRequestConfig(ctx context.Context,
	user *entities.User,
	update *entities.Update) (configs.RequestConfig, error) {
	if !a.updateHaveMessage(update) {
		// Можно отправлять только текстовые сообщения.
		return configs.NewSendMessage(messages.NeedTextMessage), nil
	}
	var cfg configs.RequestConfig
	var err error
	// Что ожидается от пользователя.
	switch user.BotWant() {
	case entities.BotWantCommand: // Команда.
		if cfg, err = a.usecase.CommandHandle(ctx, user, update); err != nil {
			return nil, err
		}
	case entities.BotWantCaptcha: // Каптча.
		if cfg, err = a.usecase.CaptchaHandle(ctx, user, update); err != nil {
			return nil, err
		}
	case entities.BotWantMessage: // Сообщение.
		if cfg, err = a.usecase.MessageHandle(ctx, user, update); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknown is want %d", user.BotWant())
	}
	return cfg, nil
}

// updateHaveMessage Обновление имеет сообщение.
func (a *App) updateHaveMessage(update *entities.Update) bool {
	return update.Message != nil
}

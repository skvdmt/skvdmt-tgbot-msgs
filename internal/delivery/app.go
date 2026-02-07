package delivery

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path"
	"sync"
	"time"

	erw "github.com/skvdmt/skvdmt-back/pkg/errwrap"
	"github.com/skvdmt/skvdmt-tgbot-msgs/internal/entities"
	"github.com/skvdmt/skvdmt-tgbot-msgs/internal/entities/configs"
	"github.com/skvdmt/skvdmt-tgbot-msgs/internal/messages"
	"github.com/skvdmt/skvdmt-tgbot-msgs/internal/model"
	"github.com/skvdmt/skvdmt-tgbot-msgs/internal/usecase"
)

const (
	defaultTimeout        = 10
	defaultMaxHeaderBytes = 1 << 20 // 1Mb

	pkg = "delivery"
	app = "app"

	get          = "GET %s"
	url_messages = "/messages"
)

// App Транспортный слой.
type App struct {
	// Ресурсы.
	sources *sync.WaitGroup
	// Сервисный слой.
	usecase Usecase
	// Клиент для запросов к Telegram Bot API.
	client *Client
	// Реестр пользователей.
	users *entities.UserRegistry
	// Раннер оптимизации реестра пользователей.
	tickerOptimizeUserRegistry *time.Ticker
	// Канал остановки системы оптимизации пользователей.
	exitOptimizeUserRegistry chan struct{}
	// Роутер.
	router *http.ServeMux
	// HTTP сервер.
	APIServer *http.Server
}

// NewApp Конструктор.
func NewApp(ctx context.Context, users *entities.UserRegistry) (*App, error) {
	model.Logs.Info.Info("delivery layer creating")
	r := http.NewServeMux()
	a := &App{
		sources: &sync.WaitGroup{},
		users:   users,
		tickerOptimizeUserRegistry: time.NewTicker(time.Minute *
			time.Duration(model.Config.Timers.OptimizeUserRegistryInterval)),
		exitOptimizeUserRegistry: make(chan struct{}),
		APIServer: &http.Server{
			Addr:           fmt.Sprintf(":%d", model.Config.APIServer.Port),
			Handler:        r,
			ReadTimeout:    defaultTimeout * time.Second,
			WriteTimeout:   defaultTimeout * time.Second,
			MaxHeaderBytes: defaultMaxHeaderBytes,
		},
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
	// Запуск ресурса обработка оптимизации реестра пользователей.
	a.sources.Go(func() {
		a.handlerOptimizeUserRegistry(ctx)
	})
	// Запуск клиента для запросов к боту.
	a.sources.Go(func() {
		if err := a.client.Start(ctx); err != nil {
			model.Errors <- err
			return
		}
		for u := range a.client.Updates {
			if err := a.updateHandle(ctx, u); err != nil {
				model.Errors <- err
				return
			}
		}
		a.sources.Done()
	})
	// Запуск API сервера для получения сообщений.
	a.sources.Go(func() {
		// Настройка маршрутов.
		model.Logs.Info.Info("API server routes creating")
		a.routes()
		// Обновление сообщений.
		model.Logs.Info.Info("API server messages updating")
		if err := a.usecase.UpdateMessages(ctx); err != nil {
			model.Errors <- err
			return
		}
		// API server starting
		model.Logs.Info.Info(fmt.Sprintf("API server starting on %d port",
			model.Config.APIServer.Port))
		if err := a.APIServer.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			model.Errors <- err
			return
		}
		a.sources.Done()
	})
	a.sources.Wait()
	return nil
}

// Stop Остановка.
func (a *App) Stop(ctx context.Context) error {
	// Остановка системы оптимизации пользователей.
	a.exitOptimizeUserRegistry <- struct{}{}
	// Остановка клиента делающего запросы к Telegram Bot API.
	if err := a.client.Stop(ctx); err != nil {
		return err
	}
	// Выключение API сервера.
	if err := a.APIServer.Shutdown(ctx); err != nil {
		return err
	}
	// Закрытие канала остановки ресурсов.
	close(a.exitOptimizeUserRegistry)
	// Вызов остановки сервисного слоя.
	if err := a.usecase.Stop(ctx); err != nil {
		return err
	}
	model.Logs.Info.Info("delivery layer stopped")
	return nil
}

// handlerOptimizeUserRegistry Обработка сигналов раннеров оптимизации реестра
// пользователей и очистки базы данных реестра пользователей.
func (a *App) handlerOptimizeUserRegistry(ctx context.Context) {
	model.Logs.Info.Info("handler optimize user registry started")
	for {
		select {
		case <-a.exitOptimizeUserRegistry:
			model.Logs.Info.Info("handler optimize user registry stopped")
			a.sources.Done()
			return
		case <-a.tickerOptimizeUserRegistry.C:
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

// routes Настройка маршрутов.
func (a *App) routes() error {
	bu := model.Config.APIServer.BaseUrl
	a.router.HandleFunc(fmt.Sprintf(get, path.Join(bu, url_messages)), a.messages)
	return nil
}

// messages Обработчик запроса сообщений.
func (a *App) messages(w http.ResponseWriter, r *http.Request) {
	mgs, err := a.usecase.Messages(r.Context())
	if err != nil {
		a.errorHandle(w, err)
		return
	}
	model.Logs.Info.Info("get messages")
	a.sendJSON(w, http.StatusOK, mgs)
}

// errorHandle Обработка HTTP ошибки.
func (a *App) errorHandle(w http.ResponseWriter, err error) {
	const m = "errorHandle"
	var e *erw.ErrorWrapper
	var ok bool
	if e, ok = err.(*erw.ErrorWrapper); !ok {
		e = erw.New(erw.Internal(
			erw.Location(pkg, app, m),
			erw.Error(fmt.Errorf(
				"%v; %v dosent match the type *errwrap.ErrorWrapper",
				fmt.Errorf("can't conversion error"), err),
			),
		))
	}
	switch {
	case 400 >= e.Code() && e.Code() <= 499:
		model.Logs.Info.Info(fmt.Sprintf("%v", e.Detailed()))
	case 500 >= e.Code() && e.Code() <= 599:
		model.Logs.Info.Info(fmt.Sprintf("%v", e.Detailed()))
	}
	a.sendJSON(w, e.Code(), map[string]string{"message": e.Message()})
}

// sendJSON Отправка ответа в JSON.
func (a *App) sendJSON(w http.ResponseWriter, code int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(value)
}

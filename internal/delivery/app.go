package delivery

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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
	pkg   = "delivery"
	app   = "app"
	limit = "limit"
	page  = "page"
)

// App Транспортный слой.
type App struct {
	// Сервисный слой.
	usecase Usecase
	// Клиент для запросов к Telegram Bot API.
	client *Client
	// Реестр пользователей.
	users *entities.UserRegistry
	// Раннер оптимизации реестра пользователей.
	tickerOptimizeUserRegistry *time.Ticker
	// Канал остановки системы оптимизации пользователей.
	stopOptimizeUserRegistry chan struct{}
	// Корректное завершение горутин.
	wg *sync.WaitGroup
}

// NewApp Конструктор.
func NewApp(ctx context.Context, users *entities.UserRegistry) (*App, error) {
	model.Logs.Info.Info("delivery layer creating")
	a := &App{
		users: users,
		wg:    &sync.WaitGroup{},
		tickerOptimizeUserRegistry: time.NewTicker(time.Minute *
			time.Duration(model.Config.Timers.OptimizeUserRegistryInterval)),
		stopOptimizeUserRegistry: make(chan struct{}, 1),
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
	model.Logs.Info.Info("delivery layer starting")
	// Запуск ресурса обработка оптимизации реестра пользователей.
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		a.handlerOptimizeUserRegistry(ctx)
	}()
	// Запуск клиента для запросов к боту и чтениее обновлений.
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
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
		model.Logs.Info.Info("update handle stopped")
	}()
	a.wg.Add(1)
	go func() {
		// Обновление сообщений.
		defer a.wg.Done()
		model.Logs.Info.Info("messages updating")
		if err := a.usecase.UpdateMessages(ctx); err != nil {
			model.Errors <- err
			return
		}
	}()
	return nil
}

// Stop Остановка.
func (a *App) Stop(ctx context.Context) error {
	// Остановка системы оптимизации пользователей.
	close(a.stopOptimizeUserRegistry)
	// Остановка клиента делающего запросы к Telegram Bot API.
	if err := a.client.Stop(ctx); err != nil {
		return err
	}
	// Вызов остановки сервисного слоя.
	if err := a.usecase.Stop(ctx); err != nil {
		return err
	}
	a.wg.Wait()
	model.Logs.Info.Info("delivery layer stopped")
	return nil
}

// messages Обработчик запроса сообщений.
func (a *App) Messages(w http.ResponseWriter, r *http.Request) {
	const m = "messages"
	ps := &entities.MessagesRequestParams{}
	var err error
	pm := r.URL.Query().Get(limit)
	if len(pm) > 0 {
		ps.Limit, err = strconv.Atoi(pm)
		if err != nil {
			a.errorHandle(w, erw.New(
				erw.CodeHTTP(http.StatusBadRequest),
				erw.Internal(
					erw.Location(pkg, app, m),
					erw.Error(fmt.Errorf("%v; %v can't convert %s to int",
						fmt.Errorf("conversion error"), err, pm)),
				)))
			return
		}
	}
	pm = r.URL.Query().Get(page)
	if len(pm) > 0 {
		ps.Page, err = strconv.Atoi(pm)
		if err != nil {
			a.errorHandle(w, erw.New(
				erw.CodeHTTP(http.StatusBadRequest),
				erw.Internal(
					erw.Location(pkg, app, m),
					erw.Error(fmt.Errorf("%v; %v can't convert %s to int",
						fmt.Errorf("conversion error"), err, pm)),
				)))
			return
		}
	}
	// По умолчанию первая страница
	if ps.Page == 0 {
		ps.Page = 1
	}
	// Валидация параметров.
	if ps.Limit < 0 {
		a.errorHandle(w, erw.New(
			erw.CodeHTTP(http.StatusBadRequest),
			erw.Internal(
				erw.Location(pkg, app, m),
				erw.Error(fmt.Errorf("%v the limit value must not be negative",
					fmt.Errorf("incorrect limit value: %d;", ps.Limit))),
			)))
		return
	}
	if ps.Page < 0 {
		a.errorHandle(w, erw.New(
			erw.CodeHTTP(http.StatusBadRequest),
			erw.Internal(
				erw.Location(pkg, app, m),
				erw.Error(fmt.Errorf("%v the page value must not be negative",
					fmt.Errorf("incorrect page value: %d;", ps.Page))),
			)))
		return
	}
	mgs, total, err := a.usecase.Messages(r.Context(), ps)
	if err != nil {
		a.errorHandle(w, err)
		return
	}
	model.Logs.Info.Info("get messages")
	a.sendJSON(w, http.StatusOK, map[string][]*entities.Message{"messages": mgs}, total)
}

// handlerOptimizeUserRegistry Обработка сигналов раннеров оптимизации реестра
// пользователей и очистки базы данных реестра пользователей.
func (a *App) handlerOptimizeUserRegistry(ctx context.Context) {
	model.Logs.Info.Info("handler optimize user registry starting")
	for {
		select {
		case <-a.stopOptimizeUserRegistry:
			model.Logs.Info.Info("handler optimize user registry stopped")
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
	a.sendJSON(w, e.Code(), map[string]string{"message": e.Message()}, 0)
}

// sendJSON Отправка ответа в JSON.
func (a *App) sendJSON(w http.ResponseWriter, code int, value any, total int) {
	w.Header().Set("Content-Type", "application/json")
	if total > 0 {
		w.Header().Set("X-Total-Count", fmt.Sprintf("%d", total))
	}
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(value)
}

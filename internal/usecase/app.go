package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/skvdmt/skvdmt-tgbot-msgs/internal/entities"
	"github.com/skvdmt/skvdmt-tgbot-msgs/internal/entities/configs"
	"github.com/skvdmt/skvdmt-tgbot-msgs/internal/messages"
	"github.com/skvdmt/skvdmt-tgbot-msgs/internal/model"
	"github.com/skvdmt/skvdmt-tgbot-msgs/internal/repository"
)

const (
	// Команда начала диалога с ботом.
	cmdStart = "/start"
	// Команда на запрос аутентификации.
	cmdAuth = "/auth"
	// Команда на запрос отправки сообщания.
	cmdMessage = "/message"
	// Текст ошибки пользователь уже авторизован.
	ErrAlreadyAuthorized = "user already authorized"
	// Текст ошибки поптки авторизации закончились.
	ErrAuthAttemptsLeft = "auth attempts left"
	// Текст ошибки пользователь не авторизован.
	ErrNotAuthorized = "user not authorized"
	// Текст ошибки пользователь уже отправил сообщение.
	ErrMessageAlreadySended = "message already sended"
)

// App Cервисный слой.
type App struct {
	repository Repository
	users      *entities.UserRegistry
	// Сообщения.
	messages   []*entities.Message
	muMessages *sync.RWMutex
}

// NewApp Конструктор.
func NewApp(ctx context.Context, users *entities.UserRegistry) (*App, error) {
	model.Logs.Info.Info("usecase layer creating")
	rep, err := repository.NewApp(ctx)
	if err != nil {
		return nil, err
	}
	return &App{
		users:      users,
		muMessages: &sync.RWMutex{},
		repository: rep,
	}, nil
}

// Stop Остановка.
func (a *App) Stop(ctx context.Context) error {
	if err := a.repository.Stop(ctx); err != nil {
		return err
	}
	model.Logs.Info.Info("usercase layer stopped")
	return nil
}

// CommandHandle Обработчик команд.
func (a *App) CommandHandle(ctx context.Context, user *entities.User, update *entities.Update) (configs.RequestConfig, error) {
	switch update.Message.Text {
	case cmdStart:
		return configs.NewSendMessage(fmt.Sprintf(messages.Start,
			update.Message.From.Username, model.Config.BotName,
			cmdAuth, cmdMessage, model.Config.MsgsUrl)), nil
	case cmdAuth:
		ok, err := user.CanAuth()
		if !ok {
			switch err.Error() {
			case ErrAlreadyAuthorized:
				user.SetBotWant(entities.BotWantCommand)
				return configs.NewSendMessage(fmt.Sprintf(messages.AlreadyAuth,
					cmdMessage)), nil
			case ErrAuthAttemptsLeft:
				user.SetBotWant(entities.BotWantCommand)
				return configs.NewSendMessage(fmt.Sprintf(messages.AuthOver,
					user.AuthCooldownLeft())), nil
			default:
				return nil, err
			}
		}
		i, v, err := model.Captcha()
		if err != nil {
			return nil, err
		}
		user.SetBotWant(entities.BotWantCaptcha)
		user.SetCaptcha(v)
		return configs.NewSendCaptcha(fmt.Sprintf(messages.AuthTitle,
			user.AuthAttemptsLeft()), i), nil
	case cmdMessage:
		ok, err := user.CanSendMessage()
		if !ok {
			switch err.Error() {
			case ErrNotAuthorized:
				user.SetBotWant(entities.BotWantCommand)
				return configs.NewSendMessage(fmt.Sprintf(messages.AuthIncomplete,
					cmdAuth)), nil
			case ErrMessageAlreadySended:
				user.SetBotWant(entities.BotWantCommand)
				return configs.NewSendMessage(fmt.Sprintf(messages.WillSaved,
					model.Config.MsgsUrl, user.SendMessageCooldownLeft())), nil
			default:
				return nil, err
			}
		}
		user.SetBotWant(entities.BotWantMessage)
		return configs.NewSendMessage(messages.SendMessageEnter), nil
	default:
		return configs.NewSendMessage(fmt.Sprintf(messages.UnknownCommand,
			update.Message.Text, cmdAuth, cmdMessage)), nil
	}
}

// CaptchaHanle Обработчик каптчи.
func (a *App) CaptchaHandle(ctx context.Context,
	user *entities.User,
	update *entities.Update) (configs.RequestConfig, error) {
	if !user.Auth(update.Message.Text) {
		user.SetBotWant(entities.BotWantCommand)
		user.SetCaptcha("")
		return configs.NewSendMessage(fmt.Sprintf(messages.AuthWrong,
			cmdAuth)), nil
	}
	user.SetBotWant(entities.BotWantCommand)
	return configs.NewSendMessage(fmt.Sprintf(messages.AuthComplete,
		cmdMessage)), nil
}

// MessageHandle Обработчик сообщения.
func (a *App) MessageHandle(ctx context.Context,
	user *entities.User,
	update *entities.Update) (configs.RequestConfig, error) {
	if !user.Authorized() {
		return nil, fmt.Errorf("user not authorized")
	}
	if !user.MessageCreatedAt().IsZero() {
		if time.Until(user.MessageCreatedAt()) <
			time.Minute*time.Duration(model.Config.Timers.SendMessageCooldown) {
			return nil, fmt.Errorf("message are saved")
		}
	}
	m := &entities.Message{
		Id:               uuid.New(),
		TelegramUserName: user.TelegramUsername(),
		Text:             update.Message.Text,
		CreatedAt:        time.Now(),
	}
	if err := a.repository.SaveMessage(ctx, user.TelegramUserId(), m); err != nil {
		return nil, err
	}

	a.muMessages.Lock()
	a.messages = append([]*entities.Message{m}, a.messages...)
	a.muMessages.Unlock()

	user.SetMessageCreatedAt(m.CreatedAt)
	user.SetBotWant(entities.BotWantCommand)
	return configs.NewSendMessage(fmt.Sprintf(messages.SendMessageSaved,
		model.Config.MsgsUrl, user.SendMessageCooldownLeft())), nil
}

// OptimizeUserRegistry Оптиизация реестра пользователей.
func (a *App) OptimizeUserRegistry(ctx context.Context) error {
	fresh := make(map[int]*entities.User)
	for k, v := range a.users.List() {
		if time.Until(v.LastUsedAt())+
			time.Minute*time.Duration(model.Config.Timers.UsedTimeout) > 0 {
			fresh[k] = v
		}
	}
	a.users.SetList(fresh)
	model.Logs.Info.Info("optimize user registry complete")
	return nil
}

// User Получение пользователя.
func (a *App) User(ctx context.Context, telegramUserId int) (*entities.User, error) {
	// Пользователь есть в карте.
	if a.users.Exists(telegramUserId) {
		u, err := a.users.Get(telegramUserId)
		if err != nil {
			return nil, err
		}
		u.SetLastUsedAt(time.Now())
		return u, nil
	}
	// Создает нового пользователя, помещает в карту и возвращает его.
	u, err := entities.NewUser(telegramUserId)
	if err != nil {
		return nil, err
	}
	mcat, err := a.repository.UserMessageCreatedAt(ctx, telegramUserId)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		mcat = &time.Time{}
	}
	u.SetMessageCreatedAt(*mcat)
	a.users.Set(u)
	return u, nil
}

// Messages Сервис сообщений.
func (a *App) Messages(ctx context.Context, params *entities.MessagesRequestParams) (mgs []*entities.Message, total int, err error) {
	a.muMessages.RLock()
	defer a.muMessages.RUnlock()
	// for _, m := range a.messages {
	// 	fmt.Println(m.Text)
	// }
	if params.Limit == 0 {
		// Лимит не указан.
		return a.messages, len(a.messages), nil
	}
	switch {
	case params.Limit*params.Page-params.Limit > len(a.messages):
		// Начало списка выходит за рамки.
		return []*entities.Message{}, len(a.messages), nil
	case params.Limit*params.Page > len(a.messages):
		// Конец списка выходит за рамки.
		return a.messages[params.Limit*params.Page-params.Limit:], len(a.messages), nil
	default:
		// Список в рамках.
		return a.messages[params.Limit*params.Page-params.Limit : params.Limit*params.Page], len(a.messages), nil
	}
}

// UpdateMessages Сервис обновления сообщений.
func (a *App) UpdateMessages(ctx context.Context) error {
	mgs, err := a.repository.UpdateMessages(ctx)
	if err != nil {
		return err
	}
	a.muMessages.Lock()
	a.messages = mgs
	a.muMessages.Unlock()
	return nil
}

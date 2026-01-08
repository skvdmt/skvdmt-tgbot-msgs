package entities

import (
	"fmt"
	"strings"
	"time"

	"github.com/skvdmt/skvdmt-tgbot-msgs/internal/model"
	"github.com/skvdmt/skvdmt-tgbot-msgs/pkg"
)

const (
	// Что ожидает бот от пользователя.
	BotWantCommand = iota + 1
	BotWantCaptcha
	BotWantMessage

	BotWantDefault = BotWantCommand

	// Описание ошибок.
	ErrAlreadyAuthorized    = "user already authorized"
	ErrAuthAttemptsLeft     = "auth attempts left"
	ErrNotAuthorized        = "user not authorized"
	ErrMessageAlreadySended = "message already sended"
)

// User Пользователь взаимодействующий с ботом.
type User struct {
	// Идентификатор пользователя в телеграм.
	telegramUserId int

	// Имя пользователя в телеграм.
	username string

	// Авторизован.
	authorized bool

	// Количество оставшихся попыток авторизации.
	authAttemptsLeft uint

	// Врямя начала восстановления попыток авторизации.
	authCooldownStartedAt time.Time

	// Информация, которую бот ожидает от пользователя.
	botWant int

	// Значение каптчи.
	captchaValue string

	// Время отправки последнего сообщения.
	messageCreatedAt time.Time

	// Время последнего использования.
	lastUsedAt time.Time
}

// NewUser Конструктор.
func NewUser(telegramUserId int) (*User, error) {
	return &User{
		telegramUserId:   telegramUserId,
		authAttemptsLeft: model.Config.DefaultMaxAuthAttempts,
		authorized:       false,
		botWant:          BotWantDefault,
		messageCreatedAt: time.Time{},
		lastUsedAt:       time.Now(),
	}, nil
}

// SetMessageCreatedAt Установка времени отправки последнего сообщения.
func (u *User) SetMessageCreatedAt(messageCreatedAt time.Time) {
	u.messageCreatedAt = messageCreatedAt
}

// SetUsername Установка имени пользователя.
func (u *User) SetUsername(username string) {
	u.username = username
}

// Authorized
func (u *User) Authorized() bool {
	return u.authorized
}

// Want Информация, которую бот ожидает от пользователя.
func (u *User) BotWant() int {
	return u.botWant
}

// TelegramUserId Идентификатор пользователя телеграм.
func (u *User) TelegramUserId() int {
	return u.telegramUserId
}

// MessageCreatedAt Время отправки последнего сообщения.
func (u *User) MessageCreatedAt() time.Time {
	return u.messageCreatedAt
}

// LastUsedAt Время последнего использования.
func (u *User) LastUsedAt() time.Time {
	return u.lastUsedAt
}

// SetLastUsedAt Устанавливает время последнего использования.
func (u *User) SetLastUsedAt(t time.Time) {
	u.lastUsedAt = t
}

// CanAuth если пользователь уже авторизован вернет false и причину(ошибку)
// если нет попыток авторизации вернет ошибку и причину(ошибку)
// в остальных случаях вернет true, также если время восстановления
// авторизации прошло, но не обнулено, то обнулить попытки авторизации
func (u *User) CanAuth() (ok bool, cause error) {
	if u.authorized {
		return false, fmt.Errorf(ErrAlreadyAuthorized)
	}
	if !u.authCooldownStartedAt.IsZero() &&
		time.Until(u.authCooldownStartedAt)+
			time.Minute*time.Duration(model.Config.Timers.AuthCooldown) <= 0 {
		// время восстановления попыток прошло
		// убрать восстановление и сбросить их
		// количество и значение по-умолчанию
		u.resetAuthAttempts()
	}
	if u.authAttemptsLeft == 0 {
		return false, fmt.Errorf(ErrAuthAttemptsLeft)
	}
	return true, nil
}

// resetAuthAttempts сбросить время восстановления,
// значение каптчи и количество попыток авторизации
func (u *User) resetAuthAttempts() {
	u.authCooldownStartedAt = time.Time{}
	u.captchaValue = ""
	u.authAttemptsLeft = model.Config.DefaultMaxAuthAttempts
}

// AuthCooldownLeft оставшееся время восстановления
// попыток авторизации в виде отформатированной строки
func (u *User) AuthCooldownLeft() string {
	return pkg.PrettyTime(time.Until(u.authCooldownStartedAt) +
		time.Minute*time.Duration(model.Config.Timers.AuthCooldown))
}

// SetCaptcha установить значение каптчи в значение value
func (u *User) SetCaptcha(value string) {
	u.captchaValue = value
}

// SetBotWant установить, какую информация бот будет ожидать от пользователя
func (u *User) SetBotWant(botWant int) error {
	switch botWant {
	case BotWantCommand,
		BotWantCaptcha,
		BotWantMessage:
		u.botWant = botWant
	default:
		return fmt.Errorf("unknown want value %d", botWant)
	}
	return nil
}

// AuthAttemptsLeft возвращает оставшееся
// количество попыток авторизации
func (u *User) AuthAttemptsLeft() int {
	return int(u.authAttemptsLeft)
}

// CanSendMessage если пользователь не авторизован вернет false и причину(ошибку)
// если нет
func (u *User) CanSendMessage() (ok bool, cause error) {
	if !u.authorized {
		return false, fmt.Errorf(ErrNotAuthorized)
	}
	if !u.messageCreatedAt.IsZero() &&
		time.Until(u.messageCreatedAt)+
			time.Minute*time.Duration(model.Config.Timers.SendMessageCooldown) <= 0 {
		// время восстановления отправки сообщения прошло
		u.messageCreatedAt = time.Time{}
	}
	if !u.messageCreatedAt.IsZero() {
		return false, fmt.Errorf(ErrMessageAlreadySended)
	}
	return true, nil
}

// SendMessageCooldownLeft оставшееся время восстановления
// до возможности отправки нового сообщения
func (u *User) SendMessageCooldownLeft() string {
	return pkg.PrettyTime(time.Until(u.messageCreatedAt) +
		time.Minute*time.Duration(model.Config.Timers.SendMessageCooldown))
}

// Auth попытка авторизации возвращает false
// если пользователь авторизован
// если у него нет попыток авторизации
// если значение каптчи в структуре пусто
// если не совпадает с переданым в параметре captcha
// иначе возвращает true
func (u *User) Auth(captcha string) bool {
	if u.authorized {
		return false
	}
	if len(u.captchaValue) == 0 {
		return false
	}
	if !strings.EqualFold(u.captchaValue, captcha) {
		u.authAttemptsLeft-- // reduces the number of available verification attempts
		u.captchaValue = ""  // force reset captcha value if bad auth attempt
		if u.authAttemptsLeft == 0 {
			u.authCooldownStartedAt = time.Now()
		}
		return false
	}
	u.authorized = true
	u.resetAuthAttempts()
	return true
}

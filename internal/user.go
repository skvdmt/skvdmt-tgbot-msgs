package internal

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/skvdmt/skvdmt-tgbot-msgs/internal/model"
	"github.com/skvdmt/skvdmt-tgbot-msgs/pkg"
)

const (
	// Описание ошибок
	ErrAlreadyAuthorized    = "user already authorized"
	ErrAuthAttemptsLeft     = "auth attempts left"
	ErrNotAuthorized        = "user not authorized"
	ErrMessageAlreadySended = "message already sended"

	// что ожидает бот от пользователя
	BotWantCommand = iota + 1
	BotWantCaptcha
	BotWantMessage

	BotWantDefault = BotWantCommand

	// Максимальное количество попыток аутентификации
	defaultMaxAuthAttempts = 3
	// Время восстановления отправки сообщения
	sendMessageCooldown = time.Hour * 24
	// Время восстановления новых попыток аутентификации
	authCooldown = time.Minute * 15
)

const ()

// User data
type User struct {
	// Unique identifier for this user or bot.
	// This number may have more than 32 significant bits and some programming languages
	// may have difficulty/silent defects in interpreting it.
	// But it has at most 52 significant bits, so a 64-bit integer or double-precision
	// float type are safe for storing this identifier.
	id int

	// Username telegram username who started a dialogue with the bot
	username string

	// Authorized trigger
	authorized bool

	// authLeft number of authorization attempts
	authAttemptsLeft uint

	// authCooldownStartedAt recovery time before the next attempt to pass the test
	authCooldownStartedAt time.Time

	// isWant information that the bot expects from the user
	isWant int

	// CaptchaValue value captcha
	captchaValue string

	// MessageCreatedAt date and time of publication of the last message
	messageCreatedAt time.Time

	// LastUsedAt время последнего использования
	lastUsedAt time.Time
}

// NewUser user constructor
func NewUser(ctx context.Context, id int, username string) (*User, error) {
	// db info check
	var mcat time.Time
	query := `SELECT message_created_at FROM users WHERE telegram_user_id = $1`
	if err := model.DB.QueryRowContext(ctx, query, id).Scan(&mcat); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}
	return &User{
		id:               id,
		username:         username,
		authAttemptsLeft: defaultMaxAuthAttempts,
		authorized:       false,
		isWant:           BotWantDefault,
		messageCreatedAt: mcat,
		lastUsedAt:       time.Now(),
	}, nil
}

// Authorized возвращает true если пользователь
// авторизован иначе возвращает false
func (u *User) Authorized() bool {
	return u.authorized
}

// IsWant информация, которую бот ожидает от пользователя
func (u *User) IsWant() int {
	return u.isWant
}

// SetWant установить, какую информация бот будет ожидать от пользователя
func (u *User) SetWant(want int) error {
	switch want {
	case BotWantCommand,
		BotWantCaptcha,
		BotWantMessage:
		u.isWant = want
	default:
		return fmt.Errorf("unknown want value %d", want)
	}
	return nil
}

// CanSendMessage если пользователь не авторизован вернет false и причину(ошибку)
// если нет
func (u *User) CanSendMessage() (ok bool, cause error) {
	if !u.authorized {
		return false, fmt.Errorf(ErrNotAuthorized)
	}
	if !u.messageCreatedAt.IsZero() &&
		time.Until(u.messageCreatedAt)+sendMessageCooldown <= 0 {
		// время восстановления отправки сообщения прошло
		u.messageCreatedAt = time.Time{}
	}
	if !u.messageCreatedAt.IsZero() {
		return false, fmt.Errorf(ErrMessageAlreadySended)
	}
	return true, nil
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
		time.Until(u.authCooldownStartedAt)+authCooldown <= 0 {
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

// SetCaptcha установить значение каптчи в значение value
func (u *User) SetCaptcha(value string) {
	u.captchaValue = value
}

// AuthAttemptsLeft возвращает оставшееся
// количество попыток авторизации
func (u *User) AuthAttemptsLeft() int {
	return int(u.authAttemptsLeft)
}

// BeforeCanTryAuth возвращает время, оставшееся до
// обновления колличества попыток авторизации
func (u *User) BeforeCanTryAuth() time.Duration {
	if u.authCooldownStartedAt.IsZero() {
		return 0
	}
	c := time.Until(u.authCooldownStartedAt)
	if c <= 0 {
		return 0
	}
	return c
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

// Message сохраняет отправленное пользователем сообщение
func (u *User) SaveMessage(ctx context.Context, message string) error {
	if !u.authorized {
		return fmt.Errorf("user not authorized")
	}
	if !u.messageCreatedAt.IsZero() {
		if time.Until(u.messageCreatedAt) < sendMessageCooldown {
			return fmt.Errorf("message are saved")
		}
	}
	_, err := model.DB.ExecContext(ctx, `INSERT INTO messages (user_id, message) VALUES ($1, $2);`, u.id, message)
	if err != nil {
		return err
	}
	fmt.Printf("MESSAGE %s SAVED\n", message)
	u.messageCreatedAt = time.Now()
	return nil
}

// resetAuthAttempts сбросить время восстановления,
// значение каптчи и количество попыток авторизации
func (u *User) resetAuthAttempts() {
	u.authCooldownStartedAt = time.Time{}
	u.captchaValue = ""
	u.authAttemptsLeft = defaultMaxAuthAttempts
}

// AuthCooldownLeft оставшееся время восстановления
// попыток авторизации в виде отформатированной строки
func (u *User) AuthCooldownLeft() string {
	return pkg.PrettyTime(time.Until(u.authCooldownStartedAt) + authCooldown)
}

// SendMessageCooldownLeft оставшееся время восстановления
// до возможности отправки нового сообщения
func (u *User) SendMessageCooldownLeft() string {
	return pkg.PrettyTime(time.Until(u.messageCreatedAt) + sendMessageCooldown)
}

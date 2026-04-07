package repository

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/skvdmt/skvdmt-tgbot-msgs/internal/entities"
	"github.com/skvdmt/skvdmt-tgbot-msgs/internal/model"
)

const (
	DB_PASSWORD     = "DB_PASSWORD"
	POSTGRES_DRIVER = "postgres"
	pkg             = "repository"
)

// App Репозиторный слой.
type App struct {
	// Название.
	name string
	// Соединение с базой данных.
	db *sql.DB
}

// NewApp Конструктор.
func NewApp(ctx context.Context) (*App, error) {
	model.Logs.Info.Info("repository layer creating")
	a := &App{
		name: "App",
	}
	var err error
	// Соединение с СУБД.
	a.db, err = a.openDB()
	if err != nil {
		return nil, err
	}
	return a, nil
}

// Stop Остановка.
func (a *App) Stop(ctx context.Context) error {
	if err := a.db.Close(); err != nil {
		return err
	}
	model.Logs.Info.Info("disconnect from database")
	model.Logs.Info.Info("repository layer stopped")
	return nil
}

// SaveMessage Сохранение сообщения.
func (a *App) SaveMessage(ctx context.Context, telegramUserId int, message *entities.Message) error {
	_, err := a.db.ExecContext(ctx,
		`INSERT INTO messages (telegram_user_id, telegram_user_name, text, created_at) VALUES ($1, $2, $3, $4);`,
		telegramUserId,
		a.nullString(message.TelegramUserName),
		message.Text,
		message.CreatedAt)
	if err != nil {
		return err
	}
	model.Logs.Info.Info(fmt.Sprintf(
		"message %s from telegram_user_id %d saved",
		message.Text, telegramUserId))
	return nil
}

// UserMessageCreatedAt получение времени создания последнего сообщения пользователем
func (a *App) UserMessageCreatedAt(ctx context.Context, telegramUserId int) (*time.Time, error) {
	mcat := time.Time{}
	if err := a.db.QueryRowContext(ctx,
		`SELECT created_at FROM messages WHERE telegram_user_id = $1 ORDER BY created_at DESC LIMIT 1;`,
		telegramUserId).Scan(&mcat); err != nil {
		return nil, err
	}
	return &mcat, nil
}

// UpdateMessages Репозиторий сообщений.
func (a *App) UpdateMessages(ctx context.Context) ([]*entities.Message, error) {
	query := "SELECT id, telegram_user_name, text, created_at FROM messages ORDER BY created_at DESC;"
	rows, err := a.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	var mgs []*entities.Message
	for rows.Next() {
		m := &entities.Message{}
		if err := rows.Scan(&m.Id,
			&m.TelegramUserName,
			&m.Text,
			&m.CreatedAt); err != nil {
			return nil, err
		}
		mgs = append(mgs, m)
	}
	return mgs, nil
}

// openDB Соединение с базой данных postgress.
func (a *App) openDB() (*sql.DB, error) {
	pwd, ok := os.LookupEnv(DB_PASSWORD)
	if !ok {
		return nil, fmt.Errorf("env %s not set", DB_PASSWORD)
	}
	pwd, err := url.QueryUnescape(pwd)
	if err != nil {
		return nil, err
	}
	db, err := sql.Open(POSTGRES_DRIVER, fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		model.Config.Postgres.Host,
		model.Config.Postgres.Port,
		model.Config.Postgres.User, pwd,
		model.Config.Postgres.Database,
	))
	if err != nil {
		return nil, err
	}
	model.Logs.Info.Info("connect to database success")
	return db, nil
}

// nullString Преобразует строку в тип sql.NullString.
func (a *App) nullString(value string) *sql.NullString {
	if len(value) > 0 {
		return &sql.NullString{
			String: value,
			Valid:  true,
		}
	}
	return &sql.NullString{}
}

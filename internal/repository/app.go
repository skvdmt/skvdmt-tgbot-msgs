package repository

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/skvdmt/skvdmt-back/pkg/errwrap"
	"github.com/skvdmt/skvdmt-tgbot-msgs/internal/entities"
	"github.com/skvdmt/skvdmt-tgbot-msgs/internal/model"
	"google.golang.org/grpc/codes"
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
func (a *App) SaveMessage(
	ctx context.Context,
	telegramUserId int,
	telegramUserName string,
	text string,
) error {
	_, err := a.db.ExecContext(ctx,
		`INSERT INTO messages (telegram_user_id, telegram_user_name, message) VALUES ($1, $2, $3);`,
		telegramUserId, a.nullString(telegramUserName), text)
	if err != nil {
		return err
	}
	model.Logs.Info.Info(fmt.Sprintf("message %s from telegram_user_id %d saved", text, telegramUserId))
	return nil
}

// UserMessageCreatedAt получение времени создания последнего сообщения пользователем
func (a *App) UserMessageCreatedAt(ctx context.Context, telegramUserId int) (*time.Time, error) {
	mcat := time.Time{}
	if err := a.db.QueryRowContext(ctx,
		`SELECT created_at FROM messages WHERE telegram_user_id = $1 ORDER BY created_at DESC LIMIT 1`,
		telegramUserId).Scan(&mcat); err != nil {
		return nil, err
	}
	return &mcat, nil
}

// Messages Репозиторий сообщений.
func (a *App) Messages(ctx context.Context) ([]*entities.Message, error) {
	mtd := "UpdateMessages"
	query := "SELECT id, message, created_at FROM messages ORDER BY created_at DESC"
	rows, err := a.db.QueryContext(ctx, query)
	if err != nil {
		return nil, errwrap.New(
			errwrap.CodegRPC(int(codes.Internal)),
			errwrap.Internal(
				errwrap.Location(pkg, a.name, mtd),
				errwrap.Error(err),
			),
		)
	}
	var mgs []*entities.Message
	for rows.Next() {
		m := &entities.Message{}
		if err := rows.Scan(m); err != nil {
			return nil, errwrap.New(
				errwrap.CodegRPC(int(codes.Internal)),
				errwrap.Internal(
					errwrap.Location(pkg, a.name, mtd),
					errwrap.Error(err),
				),
			)
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

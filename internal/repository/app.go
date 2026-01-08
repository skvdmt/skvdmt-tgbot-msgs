package repository

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/skvdmt/skvdmt-tgbot-msgs/internal/entities"
	"github.com/skvdmt/skvdmt-tgbot-msgs/internal/model"
)

const (
	DB_PASSWORD     = "DB_PASSWORD"
	uniqueViolation = "unique_violation"
)

// App Репозиторный слой.
type App struct {
	db    *sql.DB
	users *entities.UserRegistry
}

// NewApp Конструктор.
func NewApp(users *entities.UserRegistry) (*App, error) {
	model.Logs.Info.Info("repository layer creating")
	a := &App{
		users: users,
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
func (a *App) SaveMessage(ctx context.Context, telegramUserId int, text string) error {
	_, err := a.db.ExecContext(ctx,
		`INSERT INTO messages (user_id, message) VALUES ($1, $2);`,
		telegramUserId, text)
	if err != nil {
		return err
	}
	model.Logs.Info.Info(fmt.Sprintf("message %s from telegram_user_id %d saved", text, telegramUserId))
	return nil
}

// User получение времени создания последнего сообщения пользователем
func (a *App) UserMessageCreatedAt(ctx context.Context, telegramUserId int) (*time.Time, error) {
	mcat := time.Time{}
	query := `SELECT message_created_at FROM users WHERE telegram_user_id = $1`
	if err := a.db.QueryRowContext(ctx, query, telegramUserId).Scan(&mcat); err != nil {
		return nil, err
	}
	return &mcat, nil
}

// DumpUsers Свалка пользователей.
func (a *App) DumpUsers(ctx context.Context, users map[int]*entities.User) error {
	var wg sync.WaitGroup
	for _, v := range users {
		wg.Add(1)
		go func() {
			_, err := a.db.ExecContext(ctx,
				`INSERT INTO users (telegram_user_id, message_created_at) VALUES ($1, $2)`,
				v.TelegramUserId(), v.MessageCreatedAt(),
			)
			if err != nil {
				if err, ok := err.(*pq.Error); ok {
					// duplicate key value violates unique constraint
					if err.Code.Name() == uniqueViolation {
						_, err2 := a.db.ExecContext(ctx,
							`UPDATE users SET message_created_at = $1, updated_at = now() WHERE telegram_user_id = $2`,
							v.MessageCreatedAt(),
							v.TelegramUserId(),
						)
						if err2 != nil {
							model.Errors <- err2
						}
						wg.Done()
						return
					}
				}
				model.Errors <- err
			}
			wg.Done()
		}()
	}
	wg.Wait()
	model.Logs.Info.Info("dump user registry complete")
	return nil
}

// CleanUsersDB Очистка реестра пользователей в базе данных.
func (a *App) CleanUsersRegistry(ctx context.Context) error {
	_, err := a.db.ExecContext(ctx,
		fmt.Sprintf("DELETE FROM users WHERE updated_at < (now() - '%s'::interval)",
			model.Config.Timers.DbCleanInterval))
	if err != nil {
		return err
	}
	model.Logs.Info.Info("user registry cleaning complete")
	return nil
}

// GetMessageCreatedAt Получение времени создания последнего сообщения пользователем.
func (a *App) GetMessageCreatedAt(ctx context.Context, telegramUserId int) (*time.Time, error) {
	var mcat time.Time
	if err := a.db.QueryRowContext(ctx,
		`SELECT message_created_at FROM users WHERE telegram_user_id = $1`,
		telegramUserId).Scan(&mcat); err != nil {
		return nil, err
	}
	return &mcat, nil
}

// openDB Соединение с базой данных postgress.
func (a *App) openDB() (*sql.DB, error) {
	pwd, ok := os.LookupEnv(DB_PASSWORD)
	if !ok {
		return nil, fmt.Errorf("env %s not set", DB_PASSWORD)
	}
	db, err := sql.Open("postgres", fmt.Sprintf(
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

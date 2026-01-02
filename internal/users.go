package internal

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/lib/pq"
	"github.com/skvdmt/skvdmt-tgbot-msgs/internal/model"
)

const (
	mapCleanInterval      = time.Second * 30
	dbCleanInterval       = time.Second * 30
	usedTimeout           = time.Second * 30
	dbCleanIntervalString = "2 minutes"
	uniqueViolation       = "unique_violation"
)

// Users user map user id as key user as value
type Users struct {
	correct        sync.WaitGroup
	exit           chan struct{}
	mapCleanRunner *time.Ticker
	dbCleanRunner  *time.Ticker
	list           map[int]*User
	mu             sync.RWMutex
}

// NewUsers users map constructor
func NewUsers(ctx context.Context) *Users {
	u := &Users{
		mapCleanRunner: time.NewTicker(mapCleanInterval),
		dbCleanRunner:  time.NewTicker(dbCleanInterval),
		list:           make(map[int]*User),
		exit:           make(chan struct{}),
		correct:        sync.WaitGroup{},
	}
	u.cleaner(ctx)
	return u
}

// cleaner уборка карты пользователей
func (u *Users) cleaner(ctx context.Context) {
	u.correct.Add(1)
	go func() {
		for {
			select {
			case <-u.exit:
				if err := u.close(ctx); err != nil {
					model.Errors <- err
				}
				u.correct.Done()
				return
			case <-u.mapCleanRunner.C:
				if err := u.cleanMap(ctx, true); err != nil {
					model.Errors <- err
				}
			case <-u.dbCleanRunner.C:
				if err := u.cleanDB(ctx); err != nil {
					model.Errors <- err
				}
			}
		}
	}()
}

// save сохранение пользователей
func (u *Users) cleanMap(ctx context.Context, optimize bool) error {
	saving := make(map[int]*User)
	if optimize {
		// подготовить для сохранения не используемые записи
		fresh := make(map[int]*User)
		for k := range u.list {
			if time.Until(u.list[k].lastUsedAt)+usedTimeout <= 0 {
				saving[k] = u.list[k]
				delete(u.list, k)
				continue
			}
			fresh[k] = u.list[k]
		}
		u.list = fresh
	} else {
		// сохранить все
		saving = u.list
	}
	// save saving to database
	var wg sync.WaitGroup
	for k := range saving {
		wg.Add(1)
		go func() {
			_, err := model.DB.ExecContext(ctx,
				`INSERT INTO users (telegram_user_id, message_created_at) VALUES ($1, $2)`,
				saving[k].id, saving[k].messageCreatedAt,
			)
			if err != nil {
				if err, ok := err.(*pq.Error); ok {
					// duplicate key value violates unique constraint
					if err.Code.Name() == uniqueViolation {
						_, err2 := model.DB.ExecContext(ctx,
							`UPDATE users SET message_created_at = $1, updated_at = now() WHERE telegram_user_id = $2`,
							saving[k].messageCreatedAt,
							saving[k].id,
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
	return nil
}

// cleanDB очистка базы данных от устаревших записей
func (u *Users) cleanDB(ctx context.Context) error {
	_, err := model.DB.ExecContext(ctx,
		fmt.Sprintf("DELETE FROM users WHERE updated_at < (now() - '%s'::interval)", dbCleanIntervalString))
	if err != nil {
		return err
	}
	return nil
}

// Stop
func (u *Users) Stop() error {
	u.exit <- struct{}{}
	return nil
}

// Close закрытие обработчиков карты пользователей
func (u *Users) close(ctx context.Context) error {
	// save all users from map to database
	if err := u.cleanMap(ctx, false); err != nil {
		return err
	}
	// stop tickers
	u.mapCleanRunner.Stop()
	u.dbCleanRunner.Stop()
	// wait all process
	u.correct.Wait()
	// close exit channel
	close(u.exit)
	model.Logs.Info.Info("all users saved")
	return nil
}

// Exists user exist in map by user id.
func (u *Users) Exists(id int) bool {
	u.mu.RLock()
	defer u.mu.RUnlock()
	_, ok := (*u).list[id]
	return ok
}

// Add add user to users map.
func (u *Users) Add(user *User) {
	u.mu.Lock()
	defer u.mu.Unlock()
	(*u).list[user.id] = user
}

// GetOrCreateUserByIdAndName create user if not exists and return user from users map by user id and username.
func (u *Users) GetOrCreateUserByIdAndName(ctx context.Context, id int, username string) (*User, error) {
	if !u.Exists(id) {
		usr, err := NewUser(ctx, id, username)
		if err != nil {
			return nil, err
		}
		u.Add(usr)
	}
	u.mu.RLock()
	defer u.mu.RUnlock()
	usr, _ := (*u).list[id]
	usr.lastUsedAt = time.Now()
	return usr, nil
}

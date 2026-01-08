package entities

import (
	"fmt"
	"sync"

	"github.com/skvdmt/skvdmt-tgbot-msgs/internal/model"
)

// UserRegistry Реестр пользователей.
type UserRegistry struct {
	list map[int]*User
	mu   sync.RWMutex
}

// NewUserRegistry Конструктор.
func NewUserRegistry() (*UserRegistry, error) {
	model.Logs.Info.Info("user registry creating")
	u := &UserRegistry{
		list: make(map[int]*User),
	}
	return u, nil
}

// List Список пользователей.
func (u *UserRegistry) List() map[int]*User {
	return u.list
}

// SetList Устанавливает список пользователей.
func (u *UserRegistry) SetList(list map[int]*User) {
	u.list = list
}

// Exists Пользователь прусутствует в карте.
func (u *UserRegistry) Exists(telegramUserId int) bool {
	u.mu.RLock()
	defer u.mu.RUnlock()
	_, ok := (*u).list[telegramUserId]
	return ok
}

// Set Устанавливает пользователя в карту.
func (u *UserRegistry) Set(user *User) {
	u.mu.Lock()
	defer u.mu.Unlock()
	(*u).list[user.telegramUserId] = user
}

// Get Получет пользователя из карты.
func (u *UserRegistry) Get(telegramUserId int) (*User, error) {
	u.mu.Lock()
	defer u.mu.Unlock()
	usr, ok := (*u).list[telegramUserId]
	if !ok {
		return nil, fmt.Errorf("no user with telegram_user_id %d in user registry", telegramUserId)
	}
	return usr, nil
}

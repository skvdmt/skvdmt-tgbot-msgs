package entities

import (
	"time"

	"github.com/google/uuid"
)

// Message Сообщение.
type Message struct {
	Id               uuid.UUID `json:"id"`
	TelegramUserName string    `json:"telegram_user_name"`
	Message          string    `json:"message"`
	CreatedAt        time.Time `json:"created_at"`
}

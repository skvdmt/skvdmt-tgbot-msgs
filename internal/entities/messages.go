package entities

import (
	"time"

	"github.com/google/uuid"
)

// Message Сообщение.
type Message struct {
	Id                uuid.UUID `json:"id"`
	TelegramUserName  string    `json:"telegram_user_name"`
	Text              string    `json:"text"`
	CreatedAtFormated string    `json:"created_at"`
	CreatedAt         time.Time `json:"-"`
}

// MessagesRequestParams параметры запроса сообщений.
type MessagesRequestParams struct {
	Limit int
	Page  int
}

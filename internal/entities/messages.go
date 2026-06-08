package entities

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// Message Сообщение.
type Message struct {
	Id               uuid.UUID       `json:"id"`
	TelegramUserName *sql.NullString `json:"telegram_user_name"`
	Text             string          `json:"text"`
	CreatedAt        time.Time       `json:"created_at"`
}

// MessagesRequestParams параметры запроса сообщений.
type MessagesRequestParams struct {
	Limit int
	Page  int
}

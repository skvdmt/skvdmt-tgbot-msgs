package entities

import "github.com/skvdmt/skvdmt-tgbot-msgs/internal/dto"

// Update обновление из телеграма
type Update struct {
	Id      int          `json:"update_id"`
	Message *dto.Message `json:"message"`
	user    *User
}

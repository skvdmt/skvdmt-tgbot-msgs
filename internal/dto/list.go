package dto

import "encoding/json"

// Response dto
type Response struct {
	Ok          bool            `json:"ok"`
	Result      json.RawMessage `json:"result"`
	ErrorCode   int             `json:"error_code"`
	Description string          `json:"description"`
}

// Message dto
type Message struct {
	Id   int    `json:"message_id"`
	Text string `json:"text"`
	From *User  `json:"from"`
	Chat *Chat  `json:"chat"`
}

// User dto
type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
}

// Chat dto
type Chat struct {
	Id int `json:"id"`
}

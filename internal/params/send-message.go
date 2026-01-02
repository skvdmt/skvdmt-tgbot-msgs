package params

import (
	"io"
	"net/url"
	"strconv"
	"strings"
)

const (
	chatId            = "chat_id"
	replyToMessageID  = "reply_to_message_id"
	text              = "text"
	methodSendMessage = "sendMessage"
	urlValues         = "url_values"
)

// SendMessage params for send message request
type SendMessage struct {
	ChatId           int
	ReplyToMessageID int
	Text             string
}

// NewSendMessage конструктор
func NewSendMessage(chatId int, text string) *SendMessage {
	return &SendMessage{
		ChatId: chatId,
		Text:   text,
	}
}

// Method
func (s *SendMessage) Method() string {
	return methodSendMessage
}

// Body
func (s *SendMessage) Body() (io.Reader, error) {
	v := &url.Values{}
	v.Set(chatId, strconv.Itoa(s.ChatId))
	if s.ReplyToMessageID > 0 {
		v.Set(replyToMessageID, strconv.Itoa(s.ReplyToMessageID))
	}
	v.Set(text, s.Text)
	return strings.NewReader(v.Encode()), nil
}

// ContentType
func (s *SendMessage) ContentType() string {
	return urlEncoded
}

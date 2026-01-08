package configs

import "io"

// RequestConfig Интерфейс конфигурации запроса к телеграм боту.
type RequestConfig interface {
	Method() string
	Body() (io.Reader, error)
	ContentType() string
	SetChatId(id int)
}

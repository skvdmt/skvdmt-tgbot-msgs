package configs

import (
	"encoding/json"
	"io"
	"net/url"
	"strconv"
	"strings"
)

const (
	limit            = "limit"
	offset           = "offset"
	timeout          = "timeout"
	allowedUpdates   = "allowed_updates"
	methodGetUpdates = "getUpdates"

	defaultUpdateTimeout = 1
	defaultUpdateLimit   = 1
	urlEncoded           = "application/x-www-form-urlencoded"
)

// Getupdates Конфигурация запроса на получение обновлений.
type GetUpdates struct {
	Limit          int
	Offset         int
	Timeout        int
	AllowedUpdates []string
}

// NewGetUpdates Конструктор.
func NewGetUpdates() *GetUpdates {
	return &GetUpdates{
		Timeout: defaultUpdateTimeout,
		Limit:   defaultUpdateLimit,
	}
}

// Method
func (g *GetUpdates) Method() string {
	return methodGetUpdates
}

// Body
func (g *GetUpdates) Body() (io.Reader, error) {
	v := &url.Values{}
	if g.Limit > 0 {
		v.Set(limit, strconv.Itoa(g.Limit))
	}
	if g.Offset > 0 {
		v.Set(offset, strconv.Itoa(g.Offset))
	}
	if g.Timeout > 0 {
		v.Set(timeout, strconv.Itoa(g.Timeout))
	}
	if len(g.AllowedUpdates) > 0 {
		a, err := json.Marshal(g.AllowedUpdates)
		if err != nil {
			return nil, err
		}
		v.Set(allowedUpdates, string(a))
	}
	return strings.NewReader(v.Encode()), nil
}

// ContentType
func (g *GetUpdates) ContentType() string {
	return urlEncoded
}

// SetChatId empty
func (g *GetUpdates) SetChatId(chatId int) {
}

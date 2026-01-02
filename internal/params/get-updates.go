package params

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
	updateTimeout    = 10
	urlEncoded       = "application/x-www-form-urlencoded"
)

// Getupdates конфигурация тела запроса на получение обновления
type GetUpdates struct {
	Limit          int
	Offset         int
	Timeout        int
	AllowedUpdates []string
}

// NewGetUpdates конструктор
func NewGetUpdates() *GetUpdates {
	return &GetUpdates{Timeout: updateTimeout}
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

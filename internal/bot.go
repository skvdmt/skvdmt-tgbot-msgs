package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/skvdmt/skvdmt-tgbot-msgs/internal/dto"
	"github.com/skvdmt/skvdmt-tgbot-msgs/internal/model"
	"github.com/skvdmt/skvdmt-tgbot-msgs/internal/params"
)

const (
	name             = "skidanovdima_msgs_bot"
	TGBOT_TOKEN      = "TGBOT_TOKEN"
	requestMethod    = "POST"
	endpoint         = "https://api.telegram.org/bot%s/%s"
	msgsUrl          = "https://msgs.skvdmt.ru/"
	updateBufferSize = 100
	contentType      = "Content-Type"
	contextCanceled  = "context canceled"
)

// Params request configuration interface
type Params interface {
	Method() string
	Body() (io.Reader, error)
	ContentType() string
}

// Bot telegram bot structure for recive and save messages
type Bot struct {
	token            string
	client           *http.Client
	updates          chan *Update
	paramsGetUpdates *params.GetUpdates
	users            *Users
	exit             chan struct{}
}

// NewBot telegram bot constructor
func NewBot(ctx context.Context) (*Bot, error) {
	token, ok := os.LookupEnv(TGBOT_TOKEN)
	if !ok {
		return nil, fmt.Errorf("env %s not set", TGBOT_TOKEN)
	}
	return &Bot{
		token:            token,
		client:           &http.Client{},
		updates:          make(chan *Update, updateBufferSize),
		paramsGetUpdates: params.NewGetUpdates(),
		users:            NewUsers(ctx),
		exit:             make(chan struct{}),
	}, nil
}

// Start running telegram bot
func (b *Bot) Start(ctx context.Context) error {
	for {
		select {
		case <-b.exit:
			return b.close()
		default:
		}
		us, err := b.getUpdates(ctx)
		if err != nil {
			return err
		}
		for _, u := range us {
			u.bot = b
			if err := u.Handling(ctx); err != nil {
				return err
			}
		}
	}
}

// getUpdates get telegram bot updates
func (b *Bot) getUpdates(ctx context.Context) ([]Update, error) {
	// fmt.Println("REQUEST")
	resp, err := b.do(ctx, b.paramsGetUpdates)
	if err != nil {
		// handling and ignore context canceled error
		if strings.HasSuffix(err.Error(), contextCanceled) {
			return nil, nil
		}
		return nil, err
	}
	// fmt.Println("RESPONSE")
	us := []Update{}
	if err := json.Unmarshal(resp, &us); err != nil {
		return nil, err
	}
	return us, nil
}

// Stop
func (b *Bot) Stop() error {
	b.exit <- struct{}{}
	return nil
}

// stop correct telegram bot stopping
func (b *Bot) close() error {
	// close updates channel
	if err := b.users.Stop(); err != nil {
		return err
	}
	close(b.updates)
	close(b.exit)
	model.Logs.Info.Info("bot stopped")
	return nil
}

// do request with context and params
func (b *Bot) do(ctx context.Context, params Params) ([]byte, error) {
	body, err := params.Body()
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(
		ctx,
		requestMethod,
		fmt.Sprintf(endpoint, b.token, params.Method()),
		body,
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set(contentType, params.ContentType())

	resp, err := b.client.Do(req)
	if err != nil {
		return nil, err
	}

	rb, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err = resp.Body.Close(); err != nil {
		return nil, err
	}
	rsp := dto.Response{}
	if err := json.Unmarshal(rb, &rsp); err != nil {
		return nil, err
	}
	if !rsp.Ok {
		return nil, fmt.Errorf("code: %d, %s", rsp.ErrorCode, rsp.Description)
	}
	return rsp.Result, nil
}

package delivery

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/skvdmt/skvdmt-tgbot-msgs/internal/dto"
	"github.com/skvdmt/skvdmt-tgbot-msgs/internal/entities"
	"github.com/skvdmt/skvdmt-tgbot-msgs/internal/entities/configs"
	"github.com/skvdmt/skvdmt-tgbot-msgs/internal/model"
)

const (
	// HTTP метод, для запросов к телеграм боту.
	requestMethod = "POST"
	// Шаблон конечной точки запроса к телеграм боту.
	endpoint = "https://api.telegram.org/bot%s/%s"
	// Название параметра заголовка типа запроса.
	contentType = "Content-Type"
	// Название переменной окружения, которая должна
	// содержать токен авторизации с телеграм ботом.
	TGBOT_TOKEN = "TGBOT_TOKEN"
)

// Client Клиент для запросов к Telegram Bot API
// и получения на них ответов.
type Client struct {
	// Канал сигнала остановки получения обновлений.
	stopGetUpdates chan struct{}
	// HTTP клиент, через который делаются запросы.
	client *http.Client
	// Token для авторизации.
	token string
	// Канал обновлений.
	Updates chan *entities.Update
	// Конфигурация запроса на получение обновлений.
	ConfigGetUpdates *configs.GetUpdates
}

// NewClient Конструктор.
func NewClient() (*Client, error) {
	model.Logs.Info.Info("client creating")
	// Получение токена авторизации с телеграм ботом.
	tkn, ok := os.LookupEnv(TGBOT_TOKEN)
	if !ok {
		return nil, fmt.Errorf("env %s not set", TGBOT_TOKEN)
	}
	c := &Client{
		stopGetUpdates:   make(chan struct{}, 1),
		client:           &http.Client{},
		Updates:          make(chan *entities.Update),
		token:            tkn,
		ConfigGetUpdates: configs.NewGetUpdates(),
	}
	return c, nil
}

// Start Запуск.
func (c *Client) Start(ctx context.Context) error {
	model.Logs.Info.Info("client started")
	go c.getUpdates(ctx)
	return nil
}

// Stop Остановка.
func (c *Client) Stop(ctx context.Context) error {
	c.stopGetUpdates <- struct{}{}
	model.Logs.Info.Info("getting updates stopped")
	// Закрытие канала остановки получения обновлений.
	close(c.stopGetUpdates)
	// Закрытие канал обновлений.
	close(c.Updates)
	model.Logs.Info.Info("client stopped")
	return nil
}

// getUpdates Получение обновлений телеграм бота.
func (c *Client) getUpdates(ctx context.Context) {
	model.Logs.Info.Info("getting updates started")
	for {
		select {
		case <-c.stopGetUpdates:
			return
		default:
			res, err := c.Do(ctx, c.ConfigGetUpdates)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					<-c.stopGetUpdates
					return
				}
				model.Errors <- err
				return
			}
			uts := []entities.Update{}
			if err := json.Unmarshal(res, &uts); err != nil {
				model.Errors <- err
				return
			}
			for _, u := range uts {
				// Установка параметра params.Offset на получение
				// следующего update в конфигурации запросов обновлений.
				if u.Id >= c.ConfigGetUpdates.Offset {
					c.ConfigGetUpdates.Offset = u.Id + 1
				}
				c.Updates <- &u
			}
		}
	}
}

// Do Выполнение запроса к телеграм боту.
func (c *Client) Do(ctx context.Context, config configs.RequestConfig) ([]byte, error) {
	reqb, err := config.Body()
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(
		ctx,
		requestMethod,
		fmt.Sprintf(endpoint, c.token, config.Method()),
		reqb,
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set(contentType, config.ContentType())

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	resb, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if err = res.Body.Close(); err != nil {
		return nil, err
	}
	r := dto.Response{}
	if err := json.Unmarshal(resb, &r); err != nil {
		return nil, err
	}
	if !r.Ok {
		return nil, fmt.Errorf("code: %d, %s", r.ErrorCode, r.Description)
	}
	return r.Result, nil
}

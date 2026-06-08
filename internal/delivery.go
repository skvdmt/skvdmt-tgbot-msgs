package internal

import (
	"context"
	"net/http"
)

// Delivery Интерфейс транспортного слоя.
type Delivery interface {
	// Запуск.
	Start(ctx context.Context) error
	// Остановка.
	Stop(ctx context.Context) error
	// Сообщения.
	Messages(w http.ResponseWriter, r *http.Request)
}

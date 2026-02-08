package internal

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
	"github.com/skvdmt/skvdmt-tgbot-msgs/internal/delivery"
	"github.com/skvdmt/skvdmt-tgbot-msgs/internal/entities"
	"github.com/skvdmt/skvdmt-tgbot-msgs/internal/model"
)

// App Основная структура приложения.
type App struct {
	// Канал сигналов операционной системы
	// для остановки ресурсов.
	stopSources chan os.Signal
	// Контекст приложения.
	ctx context.Context
	// Функция отмены контекста всего приложения.
	cancel context.CancelFunc
	// Реестр пользователей.
	users *entities.UserRegistry
	// Транспортный слой.
	delivery Delivery
}

// NewApp Конструктор.
func NewApp() (*App, error) {
	model.Logs.Info.Info("telegram bot application creating")
	a := &App{
		stopSources: make(chan os.Signal),
	}
	var err error
	// Загрузка конфигурации.
	if err = model.LoadConfig(); err != nil {
		return nil, err
	}
	// Создание реестра пользователей.
	if a.users, err = entities.NewUserRegistry(); err != nil {
		return nil, err
	}
	// Создание контекста.
	a.ctx, a.cancel = context.WithCancel(context.Background())
	model.Logs.Info.Info("app context created")

	// Создание транспортного слоя из которого по
	// цепочки создаются остальные слои приложения.
	if a.delivery, err = delivery.NewApp(a.ctx, a.users); err != nil {
		return nil, err
	}
	return a, nil
}

// Start Запуск приложения.
func (a *App) Start() error {
	model.Logs.Info.Info("telegram bot application running")
	// Создане глобального канала ошибок для всего приложения.
	model.Errors = make(chan error)
	// Начало работы ресурсов приложения.
	go func() {
		var err error
		if err = a.delivery.Start(a.ctx); err != nil {
			model.Errors <- err
		}
	}()
	go a.signalHandling()
	// Обработка глобального канала ошибок.
	return a.errorHandling()
}

// signalHandling Отслеживание сигналов операционной системы.
func (a *App) signalHandling() {
	signal.Notify(a.stopSources, syscall.SIGTERM)
	<-a.stopSources
	model.Errors <- a.stop()
}

// errorHandling Обработка канала ошибок.
func (a *App) errorHandling() error {
	err := <-model.Errors
	close(model.Errors)
	return err
}

// stop Остановка.
func (a *App) stop() error {
	// Остановка транспортного слоя из которо по цепочке
	// останавливаются всех остальные слои.
	if err := a.delivery.Stop(a.ctx); err != nil {
		return err
	}
	// Закрытие канала остановки ресурсов.
	close(a.stopSources)
	// Отмена контекста.
	a.cancel()
	model.Logs.Info.Info("skidanovdima_msgs_bot stopped")
	return nil
}

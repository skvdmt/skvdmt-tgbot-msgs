package internal

import (
	"context"
	"os"
	"os/signal"
	"sync"
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
	exit chan os.Signal
	// Ресурсы.
	sources *sync.WaitGroup
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
		exit:    make(chan os.Signal),
		sources: &sync.WaitGroup{},
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
	// Начало работы ресурса приложения.
	a.sources.Go(func() {
		var err error
		if err = a.delivery.Start(a.ctx); err != nil {
			model.Errors <- err
		}
		// Завершение работы ресурса приложения.
	})

	go a.signalHandling()
	return a.errorHandling()
}

// signalHandling Отслеживание сигналов операционной системы.
func (a *App) signalHandling() {
	signal.Notify(a.exit, syscall.SIGTERM)
	<-a.exit
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
	// Ожидание завершения работы ресурсов приложения.
	a.sources.Wait()
	// Закрытие канала отслеживающего сигналы операционной системы.
	close(a.exit)
	// Отмена контекста.
	a.cancel()
	model.Logs.Info.Info("skidanovdima_msgs_bot stopped")
	return nil
}

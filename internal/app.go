package internal

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"github.com/skvdmt/skvdmt-tgbot-msgs/internal/delivery"
	"github.com/skvdmt/skvdmt-tgbot-msgs/internal/entities"
	"github.com/skvdmt/skvdmt-tgbot-msgs/internal/model"
)

// App Основная структура приложения.
type App struct {
	// Канал сигналов операционной системы для отслеживания
	// сигналов прерывания работы приложения.
	interrupt chan os.Signal
	// Контекст приложения.
	ctx context.Context
	// Функция отмены контекста всего приложения.
	cancel context.CancelFunc
	// Реестр пользователей.
	users *entities.UserRegistry
	// Транспортный слой.
	delivery Delivery
	// Роутер.
	router *http.ServeMux
	// Сервер.
	server *http.Server
	// Приложение запущено.
	started bool
	// Ошибки приложения.
	eg []error
	// Приложение уже закрывается.
	stopping bool
}

const (
	defaultTimeout        = 10
	defaultMaxHeaderBytes = 1 << 20 // 1Mb

	get          = "GET %s"
	url_messages = "/messages"
)

// NewApp Конструктор.
func NewApp() (*App, error) {
	model.Logs.Info.Info("telegram bot application creating")
	// Загрузка конфигурации.
	if err := model.LoadConfig(); err != nil {
		return nil, err
	}
	// Создане глобального канала ошибок для всего приложения.
	model.Errors = make(chan error)
	// Создание сервера.
	r := http.NewServeMux()
	a := &App{
		interrupt: make(chan os.Signal),
		router:    r,
		server: &http.Server{
			Addr:           fmt.Sprintf(":%d", model.Config.Server.Port),
			Handler:        r,
			ReadTimeout:    defaultTimeout * time.Second,
			WriteTimeout:   defaultTimeout * time.Second,
			MaxHeaderBytes: defaultMaxHeaderBytes,
		},
	}
	// Создание контекста.
	a.ctx, a.cancel = context.WithCancel(context.Background())
	model.Logs.Info.Info("app context created")
	// Создание реестра пользователей.
	var err error
	if a.users, err = entities.NewUserRegistry(); err != nil {
		return nil, err
	}
	// Создание транспортного слоя из которого по
	// цепочки создаются остальные слои приложения.
	if a.delivery, err = delivery.NewApp(a.ctx, a.users); err != nil {
		return nil, err
	}
	return a, nil
}

// Start Запуск приложения.
func (a *App) Start() error {
	go a.errorHandler()

	model.Logs.Info.Info("telegram bot application running")

	go func() {
		// Настройка и запуск сервера.
		a.routes()
		model.Logs.Info.Info(fmt.Sprintf("http server starting on %d port",
			model.Config.Server.Port))
		if err := a.server.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			model.Errors <- err
		}
	}()

	go func() {
		// Запуск слоев приложения по цепочке.
		if err := a.delivery.Start(a.ctx); err != nil {
			model.Errors <- err
		}
		a.started = true
	}()

	return a.interruptHandler()
}

// errorHandle Обработчик глобального канала ошибок.
func (a *App) errorHandler() {
	model.Logs.Info.Info("error handler starting")
	for {
		err := <-model.Errors
		if err == nil {
			return
		}
		if errors.Is(err, context.Canceled) {
			continue
		}
		a.eg = append(a.eg, err)
		if !a.stopping {
			a.interrupt <- syscall.SIGTERM
		}
	}
}

// interruptHandler Обработчик сигналов остановки приложения.
func (a *App) interruptHandler() error {
	model.Logs.Info.Info("error handler starting")
	signal.Notify(a.interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-a.interrupt
	// close sources
	if err := a.stop(); err != nil {
		model.Errors <- err
	}
	model.Errors <- nil
	close(model.Errors)
	if len(a.eg) > 0 {
		return fmt.Errorf("%v", a.eg)
	}
	return nil
}

// stop Остановка.
func (a *App) stop() error {
	a.stopping = true
	// Отмена контекста.
	a.cancel()
	model.Logs.Info.Info("context canceled")
	// Дождаться запуска приложения.
	for {
		if a.started {
			break
		}
	}
	// Остановка сервера.
	if err := a.server.Shutdown(context.Background()); err != nil {
		return err
	}
	model.Logs.Info.Info("http server shutdown")
	// Остановка транспортного слоя из которо по цепочке
	// останавливаются все остальные слои.
	if err := a.delivery.Stop(a.ctx); err != nil {
		return err
	}
	// Закрытие канала отслеживающего сигналы
	// прерывания операционной системы.
	close(a.interrupt)
	model.Logs.Info.Info(fmt.Sprintf("%s stopped", model.APP_NAME))
	return nil
}

// routes Настройка маршрутов.
func (a *App) routes() error {
	bu := model.Config.Server.BaseUrl
	a.router.HandleFunc(fmt.Sprintf(get, path.Join(bu, url_messages)), a.delivery.Messages)
	return nil
}

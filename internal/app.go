package internal

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	_ "github.com/lib/pq"
	"github.com/skvdmt/skvdmt-tgbot-msgs/internal/model"
)

// App appliocation struct
type App struct {
	exit      chan os.Signal
	bot       *Bot
	cancelBot context.CancelFunc
	correct   *sync.WaitGroup
}

// NewApp constructor
func NewApp() (*App, error) {
	model.Errors = make(chan error)
	return &App{
		exit:    make(chan os.Signal),
		correct: &sync.WaitGroup{},
	}, nil
}

// Start application run
func (a *App) Start() error {
	a.correct.Add(1)
	go func() {
		// async bot starting
		model.Logs.Info.Info("skidanovdima_msgs_bot starting")
		var err error
		if err = model.PostgressConnect(); err != nil {
			model.Errors <- err
		}
		var ctx context.Context
		ctx, a.cancelBot = context.WithCancel(context.Background())
		a.bot, err = NewBot(ctx)
		if err != nil {
			model.Errors <- err
		}
		if err := a.bot.Start(ctx); err != nil {
			model.Errors <- err
		}
		a.correct.Done()
	}()

	go a.signalHandling()
	return a.errorHandling()
}

// signalHandling os signal handling
func (a *App) signalHandling() {
	signal.Notify(a.exit, syscall.SIGTERM)
	<-a.exit
	model.Errors <- a.close()
}

// errorHandling error handling
func (a *App) errorHandling() error {
	err := <-model.Errors
	close(model.Errors)
	return err
}

// stop correct stop application
func (a *App) close() error {
	// stop bot
	if err := a.bot.Stop(); err != nil {
		return err
	}
	// close exit channel
	close(a.exit)
	// waiting correct stop
	a.correct.Wait()
	// cancel bot context
	a.cancelBot()
	// close database connection
	if err := model.DB.Close(); err != nil {
		return err
	}
	model.Logs.Info.Info("skidanovdima_msgs_bot stopped")
	return nil
}

package main

import (
	"os"

	"github.com/skvdmt/skvdmt-tgbot-msgs/internal"
	"github.com/skvdmt/skvdmt-tgbot-msgs/internal/model"
)

// Точка входа в приложение.
func main() {
	// Создание логгера.
	if err := model.LoadLogger(); err != nil {
		panic(err)
	}
	// Создание приложения.
	a, err := internal.NewApp()
	if err != nil {
		model.Logs.Error.Error(err.Error())
		os.Exit(1)
	}
	// Запуск приложения.
	if err := a.Start(); err != nil {
		model.Logs.Error.Error(err.Error())
		os.Exit(1)
	}
	// Закрытие логгера.
	if err := model.Logs.Close(); err != nil {
		panic(err)
	}
}

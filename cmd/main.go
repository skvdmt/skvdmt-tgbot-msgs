package main

import (
	"os"

	"github.com/skvdmt/skvdmt-tgbot-msgs/internal"
	"github.com/skvdmt/skvdmt-tgbot-msgs/internal/model"
)

func main() {
	// making logger
	if err := model.LoadLogger(); err != nil {
		panic(err)
	}
	// making app
	a, err := internal.NewApp()
	if err != nil {
		model.Logs.Error.Error(err.Error())
		os.Exit(1)
	}
	// starting app
	if err := a.Start(); err != nil {
		model.Logs.Error.Error(err.Error())
		os.Exit(1)
	}
	// close error log file
	if err := model.Logs.Close(); err != nil {
		panic(err)
	}
}

package model

import (
	"fmt"
	"log/slog"
	"os"
)

const (
	// env MODE
	MODE = "MODE"
	// MODE values
	DEV  = "dev"
	PROD = "prod" // default

	logErrorFile = "/var/log/skvdmt-tgbot-msgs/error.log"
	logFlag      = os.O_CREATE | os.O_APPEND | os.O_RDWR
	logPerm      = 0666
)

// Log
var Logs *log

// log simple journaling tool wrapper
type log struct {
	Info         *slog.Logger
	Error        *slog.Logger
	ErrorLogFile *os.File
}

// Close close open files
func (l *log) Close() error {
	if l.ErrorLogFile != nil {
		return l.ErrorLogFile.Close()
	}
	return nil
}

// Loadlogger construction logger and set to model.Log
func LoadLogger() error {
	n := "models.log.LoadLogger"
	m, ok := os.LookupEnv(MODE)
	if !ok {
		m = PROD
	}
	Logs = &log{
		Info: slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}
	switch m {
	case DEV:
		Logs.Error = slog.New(slog.NewTextHandler(os.Stderr, nil))
	case PROD:
		var err error
		Logs.ErrorLogFile, err = os.OpenFile(logErrorFile, logFlag, logPerm)
		if err != nil {
			return fmt.Errorf("%s %w", n, err)
		}
		Logs.Error = slog.New(slog.NewJSONHandler(Logs.ErrorLogFile, nil))
	default:
		return fmt.Errorf("%s unknown %s value %s", n, MODE, m)
	}
	return nil
}

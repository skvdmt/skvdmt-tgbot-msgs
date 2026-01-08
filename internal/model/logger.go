package model

import (
	"fmt"
	"log/slog"
	"os"
)

const (
	// Название переменной окружения режима.
	MODE = "MODE"
	// Значения режимов.
	DEV  = "dev"
	PROD = "prod"

	// Путь к файлу с журналом ошибок.
	logErrorFile = "/var/log/skvdmt-tgbot-msgs/error.log"
	logFlag      = os.O_CREATE | os.O_APPEND | os.O_RDWR
	logPerm      = 0666
)

// Logs Глобальная переменная логгера.
var Logs *log

// log Простая обертка журнала логов.
type log struct {
	Info         *slog.Logger
	Error        *slog.Logger
	ErrorLogFile *os.File
}

// Close Закрытие ресурсов логгера.
func (l *log) Close() error {
	if l.ErrorLogFile != nil {
		return l.ErrorLogFile.Close()
	}
	return nil
}

// Loadlogger Создать логгер и установить ссылку
// на него в глобальную переменную Logs.
func LoadLogger() error {
	n := "models.logger.LoadLogger"
	m, ok := os.LookupEnv(MODE)
	if !ok {
		return fmt.Errorf("env %s not set", MODE)
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

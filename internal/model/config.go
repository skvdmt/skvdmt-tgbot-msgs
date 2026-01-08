package model

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	// Путь к файлу конфигурации.
	configFilePath = "/etc/skvdmt-tgbot-msgs/config.yaml"
	// Директория с файлами шрифтов.
	fontsFolder = "/usr/local/share/fonts/"
)

// Timers Конфигурация временных интервалов.
type Timers struct {
	OptimizeUserRegistryInterval int    `yaml:"optimize-user-registry-interval"`
	SendMessageCooldown          int    `yaml:"send-message-cooldown"`
	AuthCooldown                 int    `yaml:"auth-cooldown"`
	DbCleanInterval              string `yaml:"db-clean-interval"`
	UsedTimeout                  int    `yaml:"used-timeout"`
}

// Postgres Конфигурация соединения с postgres.
type Postgres struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Database string `yaml:"database"`
}

// MainConfig Основная конфигурация.
type MainConfig struct {
	DefaultMaxAuthAttempts uint      `yaml:"default-max-auth-attempts"`
	BotName                string    `yaml:"bot-name"`
	MsgsUrl                string    `yaml:"msgs-url"`
	Timers                 Timers    `yaml:"timers"`
	Postgres               *Postgres `yaml:"postgres"`
	Fonts                  []string  `yaml:"fonts"`
}

// Cfg Глобальная конфигурация приложения.
var Config *MainConfig

// LoadConfig Загрузка основной конфигурации приложения и
// установки указателя на нее в глобальную переменную Config.
func LoadConfig() error {
	Logs.Info.Info("configuration loading")
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return err
	}
	Config = &MainConfig{}
	if err := yaml.Unmarshal(data, Config); err != nil {
		return err
	}
	for k := range Config.Fonts {
		Config.Fonts[k] = filepath.Join(fontsFolder, Config.Fonts[k])
	}
	return nil
}

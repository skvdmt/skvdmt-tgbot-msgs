package model

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	// Название приложения.
	APP_NAME = "skvdmt-tgbot-msgs"

	// Путь в директории конфигурации. (Добавляется директория с именем приложения).
	configDirectory = "/etc"
	// Имя файла конфигурации.
	configFileName = "config.yaml"
	// Путь в директории с файлами шрифтов. (Добавляется директория с файла шрифта).
	fontsDirectory = "/usr/local/share/fonts"
)

// Timers Конфигурация временных интервалов.
type TimersConfig struct {
	OptimizeUserRegistryInterval int    `yaml:"optimize-user-registry-interval"`
	SendMessageCooldown          int    `yaml:"send-message-cooldown"`
	AuthCooldown                 int    `yaml:"auth-cooldown"`
	DbCleanInterval              string `yaml:"db-clean-interval"`
	UsedTimeout                  int    `yaml:"used-timeout"`
}

// PostgresConfig Конфигурация соединения с postgres.
type PostgresConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Database string `yaml:"database"`
}

// APIServerConfig Конфигурация API сервера.
type APIServerConfig struct {
	BaseUrl string `yaml:"base_url"`
	Port    int    `yaml:"port"`
}

// MainConfig Основная конфигурация.
type MainConfig struct {
	DefaultMaxAuthAttempts uint            `yaml:"default-max-auth-attempts"`
	BotName                string          `yaml:"bot-name"`
	MsgsUrl                string          `yaml:"msgs-url"`
	Timers                 *TimersConfig   `yaml:"timers"`
	Postgres               *PostgresConfig `yaml:"postgres"`
	APIServer              APIServerConfig `yaml:"api_server"`
	Fonts                  []string        `yaml:"fonts"`
}

// Cfg Глобальная конфигурация приложения.
var Config *MainConfig

// LoadConfig Загрузка основной конфигурации приложения и
// установки указателя на нее в глобальную переменную Config.
func LoadConfig() error {
	Logs.Info.Info("configuration loading")
	d, err := os.ReadFile(filepath.Join(configDirectory, APP_NAME, configFileName))
	if err != nil {
		return err
	}
	cfg := &MainConfig{}
	if err := yaml.Unmarshal(d, cfg); err != nil {
		return err
	}
	Config = cfg
	for k := range Config.Fonts {
		Config.Fonts[k] = filepath.Join(fontsDirectory, Config.Fonts[k])
	}
	Config = cfg
	return nil
}

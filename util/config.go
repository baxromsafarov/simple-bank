package util

import (
	"github.com/spf13/viper"
)

// Config хранит конфигурацию приложения
type Config struct {
	DBSource      string `mapstructure:"DB_SOURCE"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
}

// LoadConfig загружает конфиг из app.env и переменных окружения
func LoadConfig(path string) (config Config, err error) {
	viper.SetConfigFile(path + "/app.env")

	// Переменные окружения перекрывают значения из файла
	viper.AutomaticEnv()

	// Читаем файл (если нет — не страшно, будет только из env)
	err = viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return
		}
	}

	err = viper.Unmarshal(&config)
	return
}

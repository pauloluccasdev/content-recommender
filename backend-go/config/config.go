package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	HTTPPort int    `mapstructure:"HTTP_PORT"`
	DBHost   string `mapstructure:"DB_HOST"`
	DBPort   int    `mapstructure:"DB_PORT"`
	DBUser   string `mapstructure:"DB_USER"`
	DBPass   string `mapstructure:"DB_PASSWORD"`
	DBName   string `mapstructure:"DB_NAME"`
	DBSSL    string `mapstructure:"DB_SSLMODE"`
	DBDriver string `mapstructure:"DB_DRIVER"`
	TimeZone string `mapstructure:"TZ"`
}

func LoadConfig(path string) (Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	viper.SetDefault("HTTP_PORT", 8080)
	viper.SetDefault("DB_PORT", 3306)
	viper.SetDefault("DB_SSLMODE", "disable")
	viper.SetDefault("DB_DRIVER", "mysql")
	viper.SetDefault("TZ", "UTC")

	var cfg Config
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("[config] usando apenas variáveis de ambiente (%v)\n", err)
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		return Config{}, err
	}

	if cfg.DBPass == "" {
		cfg.DBPass = viper.GetString("DB_PASS")
	}

	loc, err := time.LoadLocation(cfg.TimeZone)
	if err != nil {
		return Config{}, fmt.Errorf("timezone inválido: %w", err)
	}
	time.Local = loc

	return cfg, nil
}

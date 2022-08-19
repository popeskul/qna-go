package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/viper"
)

// Config represents application config
type Config struct {
	DB     Postgres
	Server struct {
		Port int `mapstructure:"port"`
	} `mapstructure:"server"`
}

// Postgres represents postgres config
type Postgres struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// New loads config from environment variables and viper config file
func New(folder, filename string) (*Config, error) {
	cfg := &Config{}

	viper.AddConfigPath(folder)
	viper.SetConfigName(filename)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}

	if err := envconfig.Process("db", cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

package db

import "fmt"

type ConfigDB struct {
	Host     string
	Port     int
	User     string
	DBName   string
	SSLMode  string
	Password string
}

func (cfg *ConfigDB) String() string {
	return fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=%s password=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.DBName, cfg.SSLMode, cfg.Password)
}

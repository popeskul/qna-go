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
	//return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)
	return fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=%s password=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.DBName, cfg.SSLMode, cfg.Password)
}

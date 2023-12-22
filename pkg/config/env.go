package config

import (
	"os"
	"strconv"
)

const defaultAppPort = 8080
const defaultDBPort = 3306

type Config struct {
	App            *App
	DatabaseConfig *DatabaseConfig
}

type App struct {
	Port int
}

type DatabaseConfig struct {
	Host     string
	User     string
	Port     int
	Password string
	Database string

	MaxIdleConn int
	MaxOpenConn int
}

func NewConfigFromEnv() *Config {
	conf := &Config{
		App: &App{
			Port: defaultAppPort,
		},
		DatabaseConfig: &DatabaseConfig{
			Host:     os.Getenv("DB_HOST"),
			User:     os.Getenv("DB_USER"),
			Port:     defaultDBPort,
			Password: os.Getenv("DB_PASSWORD"),
			Database: os.Getenv("DB_DATABASE"),
		},
	}

	if appPort := os.Getenv("APP_PORT"); appPort != "" {
		conf.App.Port, _ = strconv.Atoi(appPort)
	}

	if dbPort := os.Getenv("DB_PORT"); dbPort != "" {
		conf.DatabaseConfig.Port, _ = strconv.Atoi(dbPort)
	}

	return conf
}

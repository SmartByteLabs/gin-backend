package config

import (
	"os"
	"strconv"

	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/database"
)

const defaultAppPort = 8080
const defaultDBPort = 3306
const defaultSecret = "secret"

type Config struct {
	App            *App
	DatabaseConfig *database.Config
}

type App struct {
	JWTSecret string
	Port      int
}

func NewConfigFromEnv() *Config {
	conf := &Config{
		App: &App{
			Port:      defaultAppPort,
			JWTSecret: defaultSecret,
		},
		DatabaseConfig: &database.Config{
			Host:     os.Getenv("DB_HOST"),
			User:     os.Getenv("DB_USER"),
			Port:     defaultDBPort,
			Password: os.Getenv("DB_PASSWORD"),
			Database: os.Getenv("DB_DATABASE"),
		},
	}

	if jwtSecret := os.Getenv("JWT_SECRET"); jwtSecret != "" {
		conf.App.JWTSecret = jwtSecret
	}

	if appPort := os.Getenv("APP_PORT"); appPort != "" {
		conf.App.Port, _ = strconv.Atoi(appPort)
	}

	if dbPort := os.Getenv("DB_PORT"); dbPort != "" {
		conf.DatabaseConfig.Port, _ = strconv.Atoi(dbPort)
	}

	return conf
}

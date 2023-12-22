package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/logger"
)

type Config struct {
	Host     string
	User     string
	Port     int
	Password string
	Database string

	MaxIdleConn int
	MaxOpenConn int
}

/*
	Connect

This part will handle connection with database
read connection details from env and connect with database
if connection fails raise and panic
*/
func Connect(ctx context.Context, conf *Config) *sql.DB {
	log := logger.LoggerFromContext(ctx).WithField("func", "database.Connect")

	// connect to mysql database with
	url := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", conf.User, conf.Password, conf.Host, conf.Port, conf.Database)
	db, err := sql.Open("mysql", url)
	log.FatalIfError(err, "Failed to connect to database")

	err = db.Ping()
	log.FatalIfError(err, "Failed to ping database")

	if conf.MaxIdleConn > 0 {
		db.SetMaxIdleConns(conf.MaxIdleConn)
	}
	if conf.MaxOpenConn > 0 {
		db.SetMaxIdleConns(conf.MaxOpenConn)
	}

	return db
}

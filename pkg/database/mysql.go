package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/princeparmar/gin-backend.git/pkg/logger"
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

func QueryScanner[MODEL any](ctx context.Context, db *sql.DB, scanArr func(*MODEL) []interface{}, query string, args ...any) ([]MODEL, error) {
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requiredData []MODEL
	for rows.Next() {
		var model MODEL
		err := rows.Scan(scanArr(&model)...)
		if err != nil {
			return nil, err
		}
		requiredData = append(requiredData, model)
	}

	return nil, nil
}

type BaseDatabaseHelper[MODEL TableWithID[int]] struct {
	db        *sql.DB
	tableName string
	columns   []string
}

func NewBaseDatabaseHelper[MODEL TableWithID[int]](db *sql.DB, tableName string, columns []string) *BaseDatabaseHelper[MODEL] {
	return &BaseDatabaseHelper[MODEL]{db: db, tableName: tableName, columns: columns}
}

func (b *BaseDatabaseHelper[MODEL]) rowParser(m *MODEL) []interface{} {
	panic("implement me")
}

func (b *BaseDatabaseHelper[MODEL]) GetAll(ctx context.Context, where string, args ...any) ([]MODEL, error) {
	// Implement the GetAll method for RequiredData here
	if where != "" {
		where = "WHERE " + where
	}

	query := "SELECT " + strings.Join(b.columns, ", ") + " FROM " + b.tableName + " " + where

	return QueryScanner(ctx, b.db, b.rowParser, query, args...)
}

func (b *BaseDatabaseHelper[MODEL]) Get(ctx context.Context, id int) (*MODEL, error) {
	// Implement the Get method for RequiredData here

	ar, err := b.GetAll(ctx, "id = ?", fmt.Sprint(id))
	if err != nil {
		return nil, err
	}

	if len(ar) == 0 {
		return nil, errors.New("RequiredData not found")
	}

	return &ar[0], nil
}

func (b *BaseDatabaseHelper[MODEL]) Create(ctx context.Context, model *MODEL) (*MODEL, error) {
	// Implement the Create method for RequiredData here
	query := "INSERT INTO " + b.tableName + " (" + strings.Join(b.columns[1:], ", ") + ") VALUES (" + strings.Join(strings.Split(strings.Repeat("?", len(b.columns)-1), ""), ", ") + ")"

	res, err := b.db.ExecContext(ctx, query, b.rowParser(model)[1:]...)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, errors.New("RequiredData not created")
	}

	return b.Get(ctx, int(id))
}

func (b *BaseDatabaseHelper[MODEL]) Update(ctx context.Context, model *MODEL) (*MODEL, error) {
	// Implement the Update method for RequiredData here
	query := "UPDATE " + b.tableName + " SET " + strings.Join(b.columns[1:], " = ?, ") + " = ? WHERE id = ?"
	_, err := b.db.ExecContext(ctx, query, b.rowParser(model)[1:]...)
	if err != nil {
		return nil, err
	}

	return b.Get(ctx, (*model).GetID())
}

func (b *BaseDatabaseHelper[MODEL]) Delete(ctx context.Context, ID int) error {
	// Implement the Delete method for RequiredData here
	_, err := b.db.ExecContext(ctx, "DELETE FROM "+b.tableName+" WHERE id = ?", ID)

	return err
}

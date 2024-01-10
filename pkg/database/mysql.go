package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"

	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/constant"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/logger"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/utils"
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
func Connect(log logger.Logger, conf *Config) *sql.DB {
	log = log.WithField("func", "database.Connect")

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

func FromContext(ctx context.Context) *sql.DB {
	return ctx.Value(constant.CtxKey_DbConnection).(*sql.DB)
}

func QueryScanner[MODEL any](ctx context.Context, db *sql.DB, scanArr func(*MODEL) []interface{}, query string, args ...any) ([]MODEL, error) {
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := []MODEL{}
	for rows.Next() {
		var model MODEL
		err := rows.Scan(scanArr(&model)...)
		if err != nil {
			return nil, err
		}
		data = append(data, model)
	}

	return data, nil
}

type MysqlCurlHelper[MODEL TableWithID[int64]] CrudHelper[MysqlCondition, MODEL, int64]

type MysqlCondition struct {
	where string
	args  []any
}

type mysqlConditionHelper struct {
	MysqlCondition
}

func NewMysqlConditionHelper() Condition[MysqlCondition] {
	return &mysqlConditionHelper{}
}

func (c *mysqlConditionHelper) New() Condition[MysqlCondition] {
	return NewMysqlConditionHelper()
}

func (c *mysqlConditionHelper) Set(key string, opration ConditionOperation, value any) Condition[MysqlCondition] {
	c.where = key + " " + string(opration) + " ?"
	c.args = append(c.args, value)
	return c
}

func (c *mysqlConditionHelper) joinBy(str string, conditions ...Condition[MysqlCondition]) Condition[MysqlCondition] {
	where := []string{}

	if c.where != "" {
		where = append(where, c.where)
	}

	for _, condition := range conditions {
		where = append(where, condition.Final().where)
		c.args = append(c.args, condition.Final().args...)
	}

	c.where = "(" + strings.Join(where, ") AND (") + ")"

	return c
}

func (c *mysqlConditionHelper) And(conditions ...Condition[MysqlCondition]) Condition[MysqlCondition] {
	c.joinBy(" AND ", conditions...)
	return c
}

func (c *mysqlConditionHelper) Or(conditions ...Condition[MysqlCondition]) Condition[MysqlCondition] {
	c.joinBy(" OR ", conditions...)
	return c
}

func (c *mysqlConditionHelper) Final() *MysqlCondition {
	return &c.MysqlCondition
}

type BaseHelper[MODEL TableWithID[int64]] struct {
	db            *sql.DB
	tableName     string
	columnMapping func(*MODEL) map[string]interface{}
}

func NewBaseHelper[MODEL TableWithID[int64]](db *sql.DB, tableName string, columnMapping func(*MODEL) map[string]interface{}) *BaseHelper[MODEL] {
	return &BaseHelper[MODEL]{db: db, tableName: tableName, columnMapping: columnMapping}
}

func (b *BaseHelper[MODEL]) GetTableName() string {
	return b.tableName
}

func (b *BaseHelper[MODEL]) getParser(project []string) func(*MODEL) []interface{} {
	return func(m *MODEL) []interface{} {
		_map := b.columnMapping(m)
		var arr []interface{}
		for _, key := range project {
			arr = append(arr, _map[key])
		}

		return arr
	}
}

func (b *BaseHelper[MODEL]) GetColumns(project []string, withoutID bool) []string {

	out := utils.NewSet[string]()
	var m MODEL
	for key := range b.columnMapping(&m) {
		out.Add(key)
	}

	if withoutID {
		out.Remove("id")
	}

	if len(project) == 0 || project[0] == "*" {
		return out.ToSlice()
	}

	return out.GetCommonElements(project)
}

func (b *BaseHelper[MODEL]) Get(ctx context.Context, project []string, condition Condition[MysqlCondition]) ([]MODEL, error) {
	cond := condition.Final()
	where := cond.where
	if where != "" {
		where = "WHERE " + where
	}

	columns := b.GetColumns(project, false)
	query := "SELECT `" + strings.Join(columns, "`, `") + "` FROM " + b.tableName + " " + where

	return QueryScanner(ctx, b.db, b.getParser(columns), query, cond.args...)
}

func (b *BaseHelper[MODEL]) Create(ctx context.Context, model *MODEL, condition Condition[MysqlCondition]) (*MODEL, error) {
	args := []interface{}{}
	columns := []string{}

	for key, value := range b.columnMapping(model) {
		if key == "id" {
			continue
		}
		columns = append(columns, key)
		args = append(args, value)
	}

	query := "INSERT INTO " + b.tableName + " (`" + strings.Join(columns, "`, `") + "`) VALUES (" + strings.Join(strings.Split(strings.Repeat("?", len(columns)), ""), ", ") + ")"

	res, err := b.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, errors.New("requested entry not created")
	}

	d, err := b.Get(ctx, AllFields, NewMysqlConditionHelper().Set("id", ConditionOperationEqual, id))
	if err != nil {
		return nil, err
	}

	if len(d) == 0 {
		return nil, errors.New("RequiredData not found")
	}

	return &d[0], nil
}

func (b *BaseHelper[MODEL]) Update(ctx context.Context, model *MODEL, project []string, condition Condition[MysqlCondition]) error {
	cond := condition.Final()
	where := cond.where
	if where == "" {
		return errors.New("where condition required")
	}

	where = "WHERE " + where
	project = b.GetColumns(project, true)
	query := "UPDATE " + b.tableName + " SET " + strings.Join(project, " = ?, ") + " = ? " + where
	_, err := b.db.ExecContext(ctx, query, append(b.getParser(project)(model), cond.args...)...)
	return err
}

func (b *BaseHelper[MODEL]) Delete(ctx context.Context, condition Condition[MysqlCondition]) error {
	cond := condition.Final()
	where := cond.where
	if where == "" {
		return errors.New("where condition required")
	}

	where = " WHERE " + where

	_, err := b.db.ExecContext(ctx, "DELETE FROM "+b.tableName+where, cond.args...)

	return err
}

type MapWithID map[string]interface{} // map[string]interface{} with ID field

func (m MapWithID) GetID() int {
	return m["id"].(int)
}

func (m MapWithID) SetID(ID int) {
	m["id"] = ID
}

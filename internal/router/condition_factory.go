package router

import (
	"net/http"

	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/database"
)

// EmptyCondition will create simple MysqlConditionHelper with no condition
func EmptyCondition(req *http.Request) (database.Condition[database.MysqlCondition], error) {
	return database.NewMysqlConditionHelper(), nil
}

// templateConditionFactory returns helper with condition helper. it handles logic related to creating condition
// from http request.
func templateConditionFactory(req *http.Request) (database.Condition[database.MysqlCondition], error) {
	return database.NewMysqlConditionHelper(), nil
}

func templeConditionFactory(req *http.Request) (database.Condition[database.MysqlCondition], error) {
	return database.NewMysqlConditionHelper(), nil
}

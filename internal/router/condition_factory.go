package router

import (
	"net/http"

	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/database"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/rbac"
)

func templateConditionFactory(req *http.Request) (database.Condition[rbac.RbacCondition[database.MysqlCondition, int64]], error) {
	// get id => id on path
	// get all =>

	h := rbac.NewRbacCondition[database.MysqlCondition, int64](1, database.NewMysqlConditionHelper())
	h.SetUserLevelData("id", "id1")
	return h, nil
}

func templeConditionFactory(req *http.Request) (database.Condition[rbac.RbacCondition[database.MysqlCondition, int64]], error) {
	h := rbac.NewRbacCondition[database.MysqlCondition, int64](1, database.NewMysqlConditionHelper())
	h.SetUserLevelData("id", "id1")
	return h, nil
}

/*
	temple
		create => no condition required
		get =>
		getAll
		update
		delete =>
*/

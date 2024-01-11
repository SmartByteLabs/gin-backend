package mysql

import (
	"context"
	"database/sql"

	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/database"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/rbac"
)

func CreateAllTables(ctx context.Context, db *sql.DB) error {
	allCreateTableHelpers := []func(context.Context, *sql.DB) error{
		CreateAccessTable,
		CreateRoleTable,
		CreateRoleAccessMappingTable, // this table depends on access and role table
		CreateUserTable,
		CreateUserRoleMappingTable, // this table depends on user and role table
	}

	for _, createTableHelper := range allCreateTableHelpers {
		err := createTableHelper(ctx, db)
		if err != nil {
			return err
		}
	}

	return nil
}

// type accessJoinQueryWithConditionResponse struct {
// 	AccessName string `db:"name"`

// 	RoleKey   sql.NullString `db:"role_key"`
// 	RoleValue sql.NullString `db:"role_value"`

// 	UserKey   sql.NullString `db:"user_key"`
// 	UserValue sql.NullString `db:"user_value"`
// }

type mysqlRbacHelper struct {
	accessHelper            AccessHelper
	roleHelper              RoleHelper
	roleAccessMappingHelper RoleAccessMappingHelper
	userHelper              UserHelper
	userRoleMappingHelper   UserRoleMappingHelper

	db *sql.DB
}

func MysqlRbacHelper(db *sql.DB) rbac.DatabaseHelper[database.MysqlCondition, int64] {
	return &mysqlRbacHelper{
		accessHelper:            NewAccessHelper(db),
		roleHelper:              NewRoleHelper(db),
		roleAccessMappingHelper: NewRoleAccessMappingHelper(db),
		userHelper:              NewUserHelper(db),
		userRoleMappingHelper:   NewUserRoleMappingHelper(db),

		db: db,
	}
}

func (mrh *mysqlRbacHelper) GetAccessHelper() database.CrudHelper[database.MysqlCondition, rbac.Access[int64], int64] {
	return mrh.accessHelper
}

func (mrh *mysqlRbacHelper) GetRoleHelper() database.CrudHelper[database.MysqlCondition, rbac.Role[int64], int64] {
	return mrh.roleHelper
}

func (mrh *mysqlRbacHelper) GetRoleAccessMappingHelper() database.CrudHelper[database.MysqlCondition, rbac.RoleAccessMapping[int64], int64] {
	return mrh.roleAccessMappingHelper
}

func (mrh *mysqlRbacHelper) GetUserHelper() database.CrudHelper[database.MysqlCondition, rbac.User[int64], int64] {
	return mrh.userHelper
}

func (mrh *mysqlRbacHelper) GetUserRoleMappingHelper() database.CrudHelper[database.MysqlCondition, rbac.UserRoleMapping[int64], int64] {
	return mrh.userRoleMappingHelper
}

func (mrh *mysqlRbacHelper) GetAccessForUser(ctx context.Context, userID int64, access string) ([]rbac.AccessWithReferenceID[int64], error) {
	query := `SELECT access.name, role_access_mapping.project, user_role_mapping.reference_id
		FROM access
			INNER JOIN role_access_mapping ON access.id = role_access_mapping.access_id
			INNER JOIN role ON role_access_mapping.role_id = role.id
			INNER JOIN user_role_mapping ON role.id = user_role_mapping.role_id
		WHERE user_role_mapping.user_id = ? AND access.name = ?`

	return database.QueryScanner(ctx, mrh.db, func(m *rbac.AccessWithReferenceID[int64]) []interface{} {
		return []interface{}{
			&m.AccessName,
			&m.Project,
			&m.ReferenceID,
		}
	}, query, userID, access)
}

// TODO remove duplicate code
func (mrh *mysqlRbacHelper) GetAccessForUserWithReference(ctx context.Context, userID int64, access string, referenceID int64) ([]rbac.AccessWithReferenceID[int64], error) {
	query := `SELECT access.name, role_access_mapping.project
		FROM access
			INNER JOIN role_access_mapping ON access.id = role_access_mapping.access_id
			INNER JOIN role ON role_access_mapping.role_id = role.id
			INNER JOIN user_role_mapping ON role.id = user_role_mapping.role_id
		WHERE user_role_mapping.user_id = ? AND access.name = ? AND user_role_mapping.reference_id = ?`

	return database.QueryScanner(ctx, mrh.db, func(m *rbac.AccessWithReferenceID[int64]) []interface{} {
		return []interface{}{
			&m.AccessName,
			&m.Project,
			&m.ReferenceID,
		}
	}, query, userID, access, referenceID)

}

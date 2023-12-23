package mysql

import (
	"context"
	"database/sql"
	"strings"

	"github.com/princeparmar/gin-backend.git/pkg/database"
	"github.com/princeparmar/gin-backend.git/pkg/rbac"
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

type mysqlRbacHelper struct {
	accessHelper            AccessHelper
	roleHelper              RoleHelper
	roleAccessMappingHelper RoleAccessMappingHelper
	userHelper              UserHelper
	userRoleMappingHelper   UserRoleMappingHelper
	requiredDataHelper      RequiredDataHelper

	db *sql.DB
}

func MysqlRbacHelper(db *sql.DB) rbac.DatabaseHelper[int] {
	return &mysqlRbacHelper{
		accessHelper:            NewAccessHelper(db),
		roleHelper:              NewRoleHelper(db),
		roleAccessMappingHelper: NewRoleAccessMappingHelper(db),
		userHelper:              NewUserHelper(db),
		userRoleMappingHelper:   NewUserRoleMappingHelper(db),
		requiredDataHelper:      NewRequiredDataHelper(db),

		db: db,
	}
}

func (mrh *mysqlRbacHelper) GetAccessHelper() database.CRUDDatabaseHelper[rbac.Access[int], int] {
	return mrh.accessHelper
}

func (mrh *mysqlRbacHelper) GetRoleHelper() database.CRUDDatabaseHelper[rbac.Role[int], int] {
	return mrh.roleHelper
}

func (mrh *mysqlRbacHelper) GetRoleAccessMappingHelper() database.CRUDDatabaseHelper[rbac.RoleAccessMapping[int], int] {
	return mrh.roleAccessMappingHelper
}

func (mrh *mysqlRbacHelper) GetUserHelper() database.CRUDDatabaseHelper[rbac.User[int], int] {
	return mrh.userHelper
}

func (mrh *mysqlRbacHelper) GetUserRoleMappingHelper() database.CRUDDatabaseHelper[rbac.UserRoleMapping[int], int] {
	return mrh.userRoleMappingHelper
}

func (mrh *mysqlRbacHelper) GetRequiredDataHelper() database.CRUDDatabaseHelper[rbac.RequiredData[int], int] {
	return mrh.requiredDataHelper
}

func (mrh *mysqlRbacHelper) GetAccessMap(ctx context.Context, UserID, accessID int) ([]rbac.AccessMap, error) {
	// access > role_access_mapping > role > user_role_mapping > user > required_data (role) > required_data (user)

	query := `SELECT access.name as name, 
			role_rd.parent_id as ram_id, role_rd.key as role_key, role_rd.value as role_value, 
			user_rd.parent_id as urm_id, user_rd.key as user_key, user_rd.value as user_value
		FROM access
			INNER JOIN role_access_mapping as ram ON access.id = ram.access_id
			INNER JOIN role ON ram.role_id = role.id
			INNER JOIN user_role_mapping as urm ON role.id = urm.role_id
			INNER JOIN user ON urm.user_id = user.id
			LEFT JOIN required_data as role_rd ON role_rd.parent_id = ram.id AND role_rd.level = "role"
			LEFT JOIN required_data as user_rd ON user_rd.parent_id = urm.id AND user_rd.level = "user"
		WHERE user.id = ? AND access.id = ?`

	return database.QueryScanner(ctx, mrh.db, func(m *rbac.AccessMap) []interface{} {
		return nil
	}, query, UserID, accessID)

}

func (mrh *mysqlRbacHelper) GetAccessMapWithCondition(ctx context.Context, UserID, accessID int, userLevelCondition, roleLevelCondition map[string]string) ([]rbac.AccessMap, error) {
	query := `SELECT access.name as name, 
			role_rd.parent_id as ram_id, role_rd.key as role_key, role_rd.value as role_value, 
			user_rd.parent_id as urm_id, user_rd.key as user_key, user_rd.value as user_value
		FROM access
			INNER JOIN role_access_mapping as ram ON access.id = ram.access_id
			INNER JOIN role ON ram.role_id = role.id
			INNER JOIN user_role_mapping as urm ON role.id = urm.role_id
			INNER JOIN user ON urm.user_id = user.id
			LEFT JOIN required_data as role_rd ON role_rd.parent_id = ram.id AND role_rd.level = "role"
			LEFT JOIN required_data as user_rd ON user_rd.parent_id = urm.id AND user_rd.level = "user"
			`

	where := []string{"user.id = ?", "access.id = ?"}
	args := []interface{}{UserID, accessID}

	for key, value := range userLevelCondition {
		where = append(where, "user_rd.key = ? AND user_rd.value = ?")
		args = append(args, key, value)
	}

	for key, value := range roleLevelCondition {
		where = append(where, "role_rd.key = ? AND role_rd.value = ?")
		args = append(args, key, value)
	}

	query += " WHERE " + strings.Join(where, " AND ")

	return database.QueryScanner(ctx, mrh.db, func(m *rbac.AccessMap) []interface{} {
		return nil
	}, query, args...)
}

package mysql

import (
	"context"
	"database/sql"

	"github.com/princeparmar/gin-backend.git/pkg/database"
	"github.com/princeparmar/gin-backend.git/pkg/rbac"
)

func CreateUserRoleMappingTable(ctx context.Context, db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS user_role_mapping (
		id int(11) unsigned NOT NULL AUTO_INCREMENT,
		user_id int(11) unsigned NOT NULL,
		role_id int(11) unsigned NOT NULL,
		created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		PRIMARY KEY (id),
		FOREIGN KEY (user_id) REFERENCES user(id),
		FOREIGN KEY (role_id) REFERENCES role(id)
	) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;`

	_, err := db.ExecContext(ctx, query)
	return err
}

type UserRoleMappingHelper database.CRUDDatabaseHelper[rbac.UserRoleMapping[int], int]

type userRoleMappingHelper struct {
	*database.BaseDatabaseHelper[rbac.UserRoleMapping[int]]
}

func NewUserRoleMappingHelper(db *sql.DB) UserRoleMappingHelper {
	const tableName = "user_role_mapping"
	columns := []string{"id", "user_id", "role_id"}

	return &userRoleMappingHelper{
		BaseDatabaseHelper: database.NewBaseDatabaseHelper[rbac.UserRoleMapping[int]](db, tableName, columns),
	}
}

func (urmh *userRoleMappingHelper) rowParser(m *rbac.UserRoleMapping[int]) []interface{} {
	return []interface{}{&m.ID, &m.UserID, &m.RoleID}
}

package mysql

import (
	"context"
	"database/sql"

	"github.com/princeparmar/gin-backend.git/pkg/database"
	"github.com/princeparmar/gin-backend.git/pkg/rbac"
)

func CreateRoleTable(ctx context.Context, db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS role (
		id int(11) unsigned NOT NULL AUTO_INCREMENT,
		name varchar(255) NOT NULL DEFAULT '',
		created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		PRIMARY KEY (id)
	) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;`

	_, err := db.ExecContext(ctx, query)
	return err
}

type RoleHelper database.CRUDDatabaseHelper[rbac.Role[int], int]

type roleHelper struct {
	*database.BaseDatabaseHelper[rbac.Role[int]]
}

func NewRoleHelper(db *sql.DB) RoleHelper {
	const tableName = "role"
	columns := []string{"id", "name"}

	return &roleHelper{
		BaseDatabaseHelper: database.NewBaseDatabaseHelper[rbac.Role[int]](db, tableName, columns),
	}
}

func (rh *roleHelper) rowParser(m *rbac.Role[int]) []interface{} {
	return []interface{}{&m.ID, &m.Name}
}

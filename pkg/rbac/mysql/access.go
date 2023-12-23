package mysql

import (
	"context"
	"database/sql"

	"github.com/princeparmar/gin-backend.git/pkg/database"
	"github.com/princeparmar/gin-backend.git/pkg/rbac"
)

/*
	Query to create access table in database:
	CREATE TABLE `access` (
		`id` int(11) unsigned NOT NULL AUTO_INCREMENT,
		`name` varchar(255) NOT NULL DEFAULT '',
		`role_level_required_data` varchar(255),
		`created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
		`updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		PRIMARY KEY (`id`)
	) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;
*/

func CreateAccessTable(ctx context.Context, db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS access (
		id int(11) unsigned NOT NULL AUTO_INCREMENT,
		name varchar(255) NOT NULL DEFAULT '',
		role_level_required_data varchar(255),
		user_level_required_data varchar(255),
		created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		PRIMARY KEY (id)
	) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;`

	_, err := db.ExecContext(ctx, query)
	return err
}

type AccessHelper database.CRUDDatabaseHelper[rbac.Access[int], int]

type accessHelper struct {
	*database.BaseDatabaseHelper[rbac.Access[int]]
}

func NewAccessHelper(db *sql.DB) AccessHelper {
	const tableName = "access"
	columns := []string{"id", "name", "role_level_required_data", "user_level_required_data"}

	return &accessHelper{
		BaseDatabaseHelper: database.NewBaseDatabaseHelper[rbac.Access[int]](db, tableName, columns),
	}
}

// TODO: need to check if this is correct
func (ah *accessHelper) rowParser(m *rbac.Access[int]) []interface{} {
	return []interface{}{&m.ID, &m.Name, &m.RoleLevelRequiredData, &m.UserLevelRequiredData}
}

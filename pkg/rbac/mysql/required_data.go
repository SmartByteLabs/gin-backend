package mysql

import (
	"context"
	"database/sql"

	"github.com/princeparmar/gin-backend.git/pkg/database"
	"github.com/princeparmar/gin-backend.git/pkg/rbac"
)

/*
	Query to create required_data table in database:
	CREATE TABLE `required_data` (
		`id` int(11) unsigned NOT NULL AUTO_INCREMENT,
		`level` varchar(255) NOT NULL DEFAULT '',
		`parent_id` int(11) unsigned NOT NULL,
		`key` varchar(255) NOT NULL,
		`value` varchar(255) NOT NULL,
		`created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
		`updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		PRIMARY KEY (`id`)
	) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;
*/

func CreateRequiredDataTable(ctx context.Context, db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS required_data (
		id int(11) unsigned NOT NULL AUTO_INCREMENT,
		level varchar(255) NOT NULL DEFAULT '',
		parent_id int(11) unsigned NOT NULL,
		key varchar(255) NOT NULL,
		value varchar(255) NOT NULL,
		created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		PRIMARY KEY (id)
	) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;`

	_, err := db.ExecContext(ctx, query)
	return err
}

type RequiredDataHelper database.CRUDDatabaseHelper[rbac.RequiredData[int], int]

type requiredDataHelper struct {
	*database.BaseDatabaseHelper[rbac.RequiredData[int]]
}

func NewRequiredDataHelper(db *sql.DB) RequiredDataHelper {
	const tableName = "required_data"
	columns := []string{"id", "level", "parent_id", "key", "value"}

	return &requiredDataHelper{
		BaseDatabaseHelper: database.NewBaseDatabaseHelper[rbac.RequiredData[int]](db, tableName, columns),
	}
}

func (rdb *requiredDataHelper) rowParser(m *rbac.RequiredData[int]) []interface{} {
	return []interface{}{&m.ID, &m.Level, &m.ParentID, &m.Key, &m.Value}
}

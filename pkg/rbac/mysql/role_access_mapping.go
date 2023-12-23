package mysql

import (
	"context"
	"database/sql"

	"github.com/princeparmar/gin-backend.git/pkg/database"
	"github.com/princeparmar/gin-backend.git/pkg/rbac"
)

/*
	Query to create role_access_mapping table:
	CREATE TABLE `role_access_mapping` (
		`id` int(11) unsigned NOT NULL AUTO_INCREMENT,
		`role_id` int(11) unsigned NOT NULL,
		`access_id` int(11) unsigned NOT NULL,
		`created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
		`updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		PRIMARY KEY (`id`),
		FOREIGN KEY (`role_id`) REFERENCES `role`(`id`),
		FOREIGN KEY (`access_id`) REFERENCES `access`(`id`)
	) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;
*/

func CreateRoleAccessMappingTable(ctx context.Context, db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS role_access_mapping (
		id int(11) unsigned NOT NULL AUTO_INCREMENT,
		role_id int(11) unsigned NOT NULL,
		access_id int(11) unsigned NOT NULL,
		created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		PRIMARY KEY (id),
		FOREIGN KEY (role_id) REFERENCES role(id),
		FOREIGN KEY (access_id) REFERENCES access(id)
	) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;`

	_, err := db.ExecContext(ctx, query)
	return err
}

type RoleAccessMappingHelper database.CRUDDatabaseHelper[rbac.RoleAccessMapping[int], int]

type roleAccessMappingHelper struct {
	*database.BaseDatabaseHelper[rbac.RoleAccessMapping[int]]
}

func NewRoleAccessMappingHelper(db *sql.DB) RoleAccessMappingHelper {
	const tableName = "role_access_mapping"
	columns := []string{"id", "role_id", "access_id"}

	return &roleAccessMappingHelper{
		BaseDatabaseHelper: database.NewBaseDatabaseHelper[rbac.RoleAccessMapping[int]](db, tableName, columns),
	}
}

func (ramh *roleAccessMappingHelper) rowParser(m *rbac.RoleAccessMapping[int]) []interface{} {
	return []interface{}{&m.ID, &m.RoleID, &m.AccessID}
}

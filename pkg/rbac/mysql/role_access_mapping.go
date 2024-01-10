package mysql

import (
	"context"
	"database/sql"

	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/database"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/rbac"
)

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

type RoleAccessMappingHelper database.MysqlCurlHelper[rbac.RoleAccessMapping[int64]]

func NewRoleAccessMappingHelper(db *sql.DB) RoleAccessMappingHelper {
	return database.NewBaseHelper[rbac.RoleAccessMapping[int64]](db, "role_access_mapping", func(a *rbac.RoleAccessMapping[int64]) map[string]interface{} {
		return map[string]interface{}{
			"id":        &a.ID,
			"role_id":   &a.RoleID,
			"access_id": &a.AccessID,
		}
	})
}

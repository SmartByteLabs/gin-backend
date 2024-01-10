package mysql

import (
	"context"
	"database/sql"

	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/database"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/rbac"
)

func CreateRoleTable(ctx context.Context, db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS role (
		id int(11) unsigned NOT NULL AUTO_INCREMENT,
		name varchar(255) NOT NULL,
		created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		PRIMARY KEY (id),
		UNIQUE KEY (name)
	) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;`

	_, err := db.ExecContext(ctx, query)
	return err
}

type RoleHelper database.MysqlCurlHelper[rbac.Role[int64]]

func NewRoleHelper(db *sql.DB) RoleHelper {
	return database.NewBaseHelper[rbac.Role[int64]](db, "role", func(r *rbac.Role[int64]) map[string]interface{} {
		return map[string]interface{}{
			"id":   &r.ID,
			"name": &r.Name,
		}
	})
}

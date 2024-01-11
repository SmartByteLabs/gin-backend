package mysql

import (
	"context"
	"database/sql"

	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/database"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/rbac"
)

func CreateAccessTable(ctx context.Context, db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS access (
		id int(11) unsigned NOT NULL AUTO_INCREMENT,
		name varchar(255) NOT NULL,
		reference_required tinyint(1) NOT NULL DEFAULT 0,

		created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		PRIMARY KEY (id),
		UNIQUE KEY (name)
	) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;`

	_, err := db.ExecContext(ctx, query)
	return err
}

type AccessHelper database.MysqlCurlHelper[rbac.Access[int64]]

func NewAccessHelper(db *sql.DB) AccessHelper {
	return database.NewBaseHelper[rbac.Access[int64]](db, "access", func(a *rbac.Access[int64]) map[string]interface{} {
		return map[string]interface{}{
			"id":   &a.ID,
			"name": &a.Name,
		}
	})
}

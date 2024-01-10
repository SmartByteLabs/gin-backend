package mysql

import (
	"context"
	"database/sql"

	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/database"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/rbac"
)

func CreateUserTable(ctx context.Context, db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS user (
		id int(11) unsigned NOT NULL AUTO_INCREMENT,
		username varchar(255) NOT NULL,
		password varchar(255) NOT NULL,
		created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		PRIMARY KEY (id),
		UNIQUE KEY (username)
	) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;`

	_, err := db.ExecContext(ctx, query)
	return err
}

type UserHelper database.MysqlCurlHelper[rbac.User[int64]]

func NewUserHelper(db *sql.DB) UserHelper {
	return database.NewBaseHelper[rbac.User[int64]](db, "user", func(u *rbac.User[int64]) map[string]interface{} {
		return map[string]interface{}{
			"id":       &u.ID,
			"username": &u.Username,
			"password": &u.Password,
		}
	})
}

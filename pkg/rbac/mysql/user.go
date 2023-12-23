package mysql

import (
	"context"
	"database/sql"

	"github.com/princeparmar/gin-backend.git/pkg/database"
	"github.com/princeparmar/gin-backend.git/pkg/rbac"
)

func CreateUserTable(ctx context.Context, db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS user (
		id int(11) unsigned NOT NULL AUTO_INCREMENT,
		username varchar(255) NOT NULL DEFAULT '',
		password varchar(255) NOT NULL DEFAULT '',
		created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		PRIMARY KEY (id)
	) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;`

	_, err := db.ExecContext(ctx, query)
	return err
}

type UserHelper database.CRUDDatabaseHelper[rbac.User[int], int]

type userHelper struct {
	*database.BaseDatabaseHelper[rbac.User[int]]
}

func NewUserHelper(db *sql.DB) UserHelper {
	const tableName = "user"
	columns := []string{"id", "username", "password"}

	return &userHelper{
		BaseDatabaseHelper: database.NewBaseDatabaseHelper[rbac.User[int]](db, tableName, columns),
	}
}

func (uh *userHelper) rowParser(m *rbac.User[int]) []interface{} {
	return []interface{}{&m.ID, &m.Username, &m.Password}
}

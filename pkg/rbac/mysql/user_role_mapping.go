package mysql

import (
	"context"
	"database/sql"

	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/database"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/rbac"
)

func CreateUserRoleMappingTable(ctx context.Context, db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS user_role_mapping (
		id int(11) unsigned NOT NULL AUTO_INCREMENT,
		user_id int(11) unsigned NOT NULL,
		role_id int(11) unsigned NOT NULL,
		reference_id int(11) unsigned,
		created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		PRIMARY KEY (id),
		FOREIGN KEY (user_id) REFERENCES user(id),
		FOREIGN KEY (role_id) REFERENCES role(id)
	) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;`

	_, err := db.ExecContext(ctx, query)
	return err
}

type UserRoleMappingHelper database.MysqlCurlHelper[rbac.UserRoleMapping[int64]]

func NewUserRoleMappingHelper(db *sql.DB) UserRoleMappingHelper {
	return database.NewBaseHelper[rbac.UserRoleMapping[int64]](db, "user_role_mapping", func(a *rbac.UserRoleMapping[int64]) map[string]interface{} {
		return map[string]interface{}{
			"id":      &a.ID,
			"user_id": &a.UserID,
			"role_id": &a.RoleID,

			"reference_id": &a.ReferenceID,
		}
	})
}

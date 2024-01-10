package mysql

import (
	"context"
	"database/sql"

	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/database"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/rbac"
)

func CreateRequiredDataTable(ctx context.Context, db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS required_data (
		id int(11) unsigned NOT NULL AUTO_INCREMENT,
		level varchar(255) NOT NULL DEFAULT '',
		parent_id int(11) unsigned NOT NULL,
		` + "`key`" + ` varchar(255) NOT NULL,
		value varchar(255) NOT NULL,
		created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		PRIMARY KEY (id),
		UNIQUE KEY (level, parent_id, ` + "`key`" + `)
	) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;`

	_, err := db.ExecContext(ctx, query)
	return err
}

type RequiredDataHelper database.MysqlCurlHelper[rbac.RequiredData[int64]]

func NewRequiredDataHelper(db *sql.DB) RequiredDataHelper {

	return database.NewBaseHelper[rbac.RequiredData[int64]](db, "required_data", func(rd *rbac.RequiredData[int64]) map[string]interface{} {
		return map[string]interface{}{
			"id":        &rd.ID,
			"level":     &rd.Level,
			"parent_id": &rd.ParentID,
			"key":       &rd.Key,
			"value":     &rd.Value,
		}
	})
}

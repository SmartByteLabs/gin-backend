package tree

import (
	"context"
	"database/sql"
)

func CreateNodeTypeTable() func(ctx context.Context, db *sql.DB) error {
	return func(ctx context.Context, db *sql.DB) error {
		query := `CREATE TABLE IF NOT EXISTS tree_fields (
			id int(11) unsigned NOT NULL AUTO_INCREMENT,
			parent_id int(11) unsigned,
			config_id int(11) unsigned NOT NULL,
			name varchar(255) NOT NULL,
			type varchar(255) NOT NULL,
			required boolean NOT NULL,
			created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			PRIMARY KEY (id),
			UNIQUE KEY (name, parent_id, root_id)
		) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;`

		_, err := db.ExecContext(ctx, query)
		return err
	}
}

func CreateNodeValueTable(tableName string) func(ctx context.Context, db *sql.DB) error {
	return func(ctx context.Context, db *sql.DB) error {
		query := `CREATE TABLE IF NOT EXISTS ` + tableName + ` (
			id int(11) unsigned NOT NULL AUTO_INCREMENT,
			parent_id int(11) unsigned,
			root_id int(11) unsigned NOT NULL,
			node_type_id int(11) unsigned NOT NULL,
			value varchar(255) NOT NULL,
			created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			PRIMARY KEY (id),
			UNIQUE KEY (parent_id, root_id, node_type_id)
		) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;`

		_, err := db.ExecContext(ctx, query)
		return err
	}
}

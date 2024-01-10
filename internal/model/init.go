package model

import (
	"context"
	"database/sql"
)

func CreateAllTables(ctx context.Context, db *sql.DB) error {
	err := createTemplateTable(db)
	if err != nil {
		return err
	}

	err = createTempleTable(db)
	if err != nil {
		return err
	}

	return nil
}

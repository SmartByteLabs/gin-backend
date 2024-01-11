package model

import (
	"database/sql"

	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/database"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/rbac"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/rbac/mysql"
)

func createTemplateTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS template (
		id INT NOT NULL AUTO_INCREMENT,
		name VARCHAR(255) NOT NULL,
		datamapping VARCHAR(255) NOT NULL,
		filepath VARCHAR(255),
		PRIMARY KEY (id),
		UNIQUE KEY (name)
	);`)
	return err
}

type Template struct {
	database.TableID[int64]
	Name        string `json:"name,omitempty"`
	DataMapping string `json:"datamapping,omitempty"`
	FilePath    string `json:"filepath,omitempty"`
}

func GetTemplateHelper(db *sql.DB) database.CrudHelper[database.MysqlCondition, Template, int64] {
	h := database.NewBaseHelper(db, "template", func(t *Template) map[string]interface{} {
		return map[string]interface{}{
			"id":          &t.ID,
			"name":        &t.Name,
			"datamapping": &t.DataMapping,
			"filepath":    &t.FilePath,
		}
	})

	rbacHelper := rbac.NewRbacHelper(mysql.MysqlRbacHelper(db))
	rbacCrudHelper := rbac.NewCrudHelper(rbacHelper, h, rbac.UserFromCTX[int64])

	return rbacCrudHelper
}

package model

import (
	"database/sql"

	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/database"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/rbac"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/rbac/mysql"
)

func createTempleTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS temple (
		id INT NOT NULL AUTO_INCREMENT,
		name VARCHAR(255) NOT NULL,
		logo VARCHAR(255) NOT NULL,
		heroimage VARCHAR(255) NOT NULL,
		opentime VARCHAR(255) NOT NULL,
		closetime VARCHAR(255) NOT NULL,
		address_line1 VARCHAR(255) NOT NULL,
		address_line2 VARCHAR(255) NOT NULL, 
		city VARCHAR(255) NOT NULL,
		state VARCHAR(255) NOT NULL,
		PRIMARY KEY (id),
		UNIQUE KEY (name)
	);`)
	return err
}

type Temple struct {
	database.TableID[int64]
	Name      string `json:"name,omitempty"`
	Logo      string `json:"logo,omitempty"`
	HeroImage string `json:"heroimage,omitempty"`

	Timings Timings `json:"timings,omitempty"`

	Address Address `json:"address,omitempty"`
}

type Timings struct {
	OpenTime  string `json:"opentime,omitempty"`
	CloseTime string `json:"closetime,omitempty"`
}

type Address struct {
	Line1 string `json:"line1,omitempty"`
	Line2 string `json:"line2,omitempty"`
	City  string `json:"city,omitempty"`
	State string `json:"state,omitempty"`
}

func GetTempleHelper(db *sql.DB) database.CrudHelper[database.MysqlCondition, Temple, int64] {
	h := database.NewBaseHelper(db, "temple", func(t *Temple) map[string]interface{} {
		return map[string]interface{}{
			"id":            &t.ID,
			"name":          &t.Name,
			"logo":          &t.Logo,
			"heroimage":     &t.HeroImage,
			"opentime":      &t.Timings.OpenTime,
			"closetime":     &t.Timings.CloseTime,
			"address_line1": &t.Address.Line1,
			"address_line2": &t.Address.Line2,
			"city":          &t.Address.City,
			"state":         &t.Address.State,
		}
	})

	rbacHelper := rbac.NewRbacHelper(mysql.MysqlRbacHelper(db))
	rbacCrudHelper := rbac.NewCrudHelper(rbacHelper, h, rbac.UserFromCTX[int64]).ReferenceRequired()

	return rbacCrudHelper
}

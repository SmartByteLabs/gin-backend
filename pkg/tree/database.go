package tree

import (
	"context"
	"database/sql"

	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/database"
)

type CombinedNode struct {
	T NodeType
	V NodeValue
}

type NodeType struct {
	database.TableID[int64]
	ParentID sql.NullInt64 `json:"parent_id"`
	ConfigID int           `json:"config_id"`

	Name string    `json:"name"`
	Type ValueType `json:"type"`

	Required bool `json:"required"`
}

type NodeValue struct {
	database.TableID[int64]
	ParentID sql.NullInt64 `json:"parent_id"`
	RootID   int64         `json:"root_id"`

	NodeTypeID int64  `json:"node_type_id"`
	Value      string `json:"value"`
}

func NewNodeTypeDbHelper(db *sql.DB) *database.BaseHelper[NodeType] {
	return database.NewBaseHelper(db, "tree_fields", func(t *NodeType) map[string]interface{} {
		return map[string]interface{}{
			"id":        &t.ID,
			"parent_id": &t.ParentID,
			"config_id": &t.ConfigID,

			"name":     &t.Name,
			"type":     &t.Type,
			"required": &t.Required,
		}
	})
}

func NewNodeValueDbHelper(db *sql.DB, tableName string) *database.BaseHelper[NodeValue] {
	return database.NewBaseHelper(db, tableName, func(t *NodeValue) map[string]interface{} {
		return map[string]interface{}{
			"id":           &t.ID,
			"parent_id":    &t.ParentID,
			"root_id":      &t.RootID,
			"node_type_id": &t.NodeTypeID,
			"value":        &t.Value,
		}
	})
}

func GetFieldsWithValues(ctx context.Context, db *sql.DB, configID, rootID int64, tableName string) ([]CombinedNode, error) {
	query := `SELECT tree_fields.id as field_id, 
			tree_fields.parent_id as field_parent_id,
			name as field_name,
			type as field_type,
			required as field_required,

			` + tableName + `.parent_id as value_parent_id, 
			` + tableName + `.id as value_id, 
			` + tableName + `.value as value
		FROM tree_fields
		LEFT JOIN ` + tableName + ` on tree_fields.id = ` + tableName + `.node_type_id
		WHERE tree_fields.config_id = ? and ` + tableName + `.root_id = ?`

	return database.QueryScanner(ctx, db, func() func(*CombinedNode) []interface{} {
		return func(c *CombinedNode) []interface{} {
			return []interface{}{
				&c.T.ID,
				&c.T.ParentID,
				&c.T.Name,
				&c.T.Type,
				&c.T.Required,
				&c.V.ParentID,
				&c.V.ID,
				&c.V.Value,
			}
		}
	}(), query, configID, rootID)

}

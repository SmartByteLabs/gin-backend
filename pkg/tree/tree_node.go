package tree

import (
	"context"
	"database/sql"
	"errors"

	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/database"
)

type ValueType string

const (
	String  ValueType = "string"
	Number  ValueType = "number"
	Boolean ValueType = "boolean"
	Array   ValueType = "array"
	Object  ValueType = "object"
)

type Fields struct {
	ID       int64     `json:"id"`
	Name     string    `json:"name"`
	Type     ValueType `json:"type"`
	Required bool      `json:"required"`

	ChildField []*Fields `json:"fields_children,omitempty"`
}

func NewFieldsFromNodeType(nodeType NodeType) *Fields {
	return &Fields{
		ID:       nodeType.ID,
		Name:     nodeType.Name,
		Type:     nodeType.Type,
		Required: nodeType.Required,
	}
}

type Node struct {
	Fields

	Value interface{} `json:"value"`

	Children []*Node `json:"children,omitempty"`
}

/*
	SELECT tree_fields.{id, parent_id, name, type, required} tree_values.{id, parent_id, value} from tree_fields
		LEFT JOIN tree_values on tree_fields.id = tree_values.node_type_id
		WHERE tree_fields.config_id = {config_id} and tree_values.root_id = {temple_id}
*/

func CreateFieldTree(ctx context.Context, db *sql.DB, configID int) (*Fields, error) {
	h := NewNodeTypeDbHelper(db)
	nodeTypes, err := h.Get(ctx, []string{"id", "parent_id", "name", "type", "required"}, database.NewMysqlConditionHelper().Set("config_id", database.ConditionOperationEqual, configID))
	if err != nil {
		return nil, err
	}

	fieldMap := make(map[int64]*Fields)
	var root *Fields

	for _, nodeType := range nodeTypes {
		fieldMap[nodeType.ID] = NewFieldsFromNodeType(nodeType)
		if nodeType.ParentID.Valid {
			fieldMap[nodeType.ParentID.Int64].ChildField = append(fieldMap[nodeType.ParentID.Int64].ChildField, fieldMap[nodeType.ID])
		} else {
			root = fieldMap[nodeType.ID]
		}
	}

	if root == nil {
		return nil, errors.New("root field not found")
	}

	return root, nil
}

func CreateNodeTree(ctx context.Context, db *sql.DB, configID, templeID int64, tableName string) (*Node, error) {
	combinedNodes, err := GetFieldsWithValues(ctx, db, configID, templeID, tableName)
	if err != nil {
		return nil, err
	}

	nodeMap := make(map[int64]*Node)
	var root *Node

	for _, combinedNode := range combinedNodes {
		nodeMap[combinedNode.V.ID] = &Node{
			Fields: *NewFieldsFromNodeType(combinedNode.T),
			Value:  combinedNode.V.Value,
		}

		if combinedNode.V.ParentID.Valid {
			nodeMap[combinedNode.V.ParentID.Int64].Children = append(nodeMap[combinedNode.V.ParentID.Int64].Children, nodeMap[combinedNode.V.ID])
		} else {
			root = nodeMap[combinedNode.V.ID]
		}
	}

	return root, nil
}

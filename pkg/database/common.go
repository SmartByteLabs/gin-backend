package database

import (
	"context"
)

var AllFields = []string{"*"}

type TableWithID[IDTYPE int64 | string] interface {
	GetID() IDTYPE
}

type TableID[IDTYPE int64 | string] struct {
	ID IDTYPE `json:"id,omitempty"`
}

func (t TableID[IDTYPE]) GetID() IDTYPE {
	return t.ID
}

type CrudHelper[T any, MODEL TableWithID[IDTYPE], IDTYPE int64 | string] interface {
	GetTableName() string
	GetColumns(project []string, withoutID bool) []string
	Get(ctx context.Context, project []string, condition Condition[T]) ([]MODEL, error)
	Create(ctx context.Context, model *MODEL) (*MODEL, error)
	Update(ctx context.Context, model *MODEL, project []string, condition Condition[T]) error
	Delete(ctx context.Context, condition Condition[T]) error
}

type ConditionOperation string

const (
	ConditionOperationEqual              ConditionOperation = "="
	ConditionOperationNotEqual           ConditionOperation = "!="
	ConditionOperationGreaterThan        ConditionOperation = ">"
	ConditionOperationGreaterThanOrEqual ConditionOperation = ">="
	ConditionOperationLessThan           ConditionOperation = "<"
	ConditionOperationLessThanOrEqual    ConditionOperation = "<="
	ConditionOperationLike               ConditionOperation = "LIKE"
	ConditionOperationNotLike            ConditionOperation = "NOT LIKE"
	ConditionOperationIn                 ConditionOperation = "IN"
	ConditionOperationNotIn              ConditionOperation = "NOT IN"
)

type Condition[T any] interface {
	Final() *T
	New() Condition[T]
	Set(key string, operation ConditionOperation, value any) Condition[T]
	And(...Condition[T]) Condition[T]
	Or(...Condition[T]) Condition[T]
}

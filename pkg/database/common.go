package database

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"
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
	Create(ctx context.Context, model *MODEL, condition Condition[T]) (*MODEL, error)
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

type DbMap[K comparable, V any] struct {
	m map[K]V
}

func NewDbMap[K comparable, V any]() *DbMap[K, V] {
	return &DbMap[K, V]{
		m: make(map[K]V),
	}
}

func (d *DbMap[K, V]) Scan(value any) error {

	fmt.Println("scan", value)
	if s, ok := value.([]uint8); ok {
		return json.Unmarshal([]byte(string(s)), &d.m)
	}

	return nil
}

func (d *DbMap[K, V]) Value() (driver.Value, error) {
	fmt.Println("value", d.m)
	v, err := json.Marshal(d.m)
	return string(v), err
}

func (d *DbMap[K, V]) Map() map[K]V {
	return d.m
}

// MarshalJSON implements json.Marshaler.
func (d DbMap[K, V]) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.m)
}

// UnmarshalJSON implements json.Unmarshaler.
func (d DbMap[K, V]) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &d.m)
}

package database

import (
	"context"
)

type TableWithID[ID_TYPE int | string] interface {
	GetID() ID_TYPE
}

type TableID[ID_TYPE int | string] struct {
	ID ID_TYPE `json:"id"`
}

func (t TableID[ID_TYPE]) GetID() ID_TYPE {
	return t.ID
}

type CRUDDatabaseHelper[MODEL TableWithID[ID_TYPE], ID_TYPE int | string] interface {
	GetAll(ctx context.Context, where string, args ...any) ([]MODEL, error)
	Get(ctx context.Context, ID ID_TYPE) (*MODEL, error)
	Create(ctx context.Context, model *MODEL) (*MODEL, error)
	Update(ctx context.Context, model *MODEL) (*MODEL, error)
	Delete(ctx context.Context, ID ID_TYPE) error
}

type CreateTableHelper interface {
	CreateTable(ctx context.Context) error
}

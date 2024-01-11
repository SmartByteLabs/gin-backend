package database

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type DbMap[K comparable, V any] struct {
	m map[K]V
}

func NewDbMap[K comparable, V any]() *DbMap[K, V] {
	return &DbMap[K, V]{
		m: make(map[K]V),
	}
}

func (d *DbMap[K, V]) Scan(value any) error {

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

type DbSlice[V any] struct {
	s []V
}

func NewDbSlice[V any]() *DbSlice[V] {
	return &DbSlice[V]{
		s: make([]V, 0),
	}
}

func (d *DbSlice[V]) Scan(value any) error {
	if s, ok := value.([]uint8); ok {
		return json.Unmarshal([]byte(string(s)), &d.s)
	}

	return nil
}

func (d *DbSlice[V]) Value() (driver.Value, error) {
	v, err := json.Marshal(d.s)
	return string(v), err
}

func (d *DbSlice[V]) Slice() []V {
	return d.s
}

// MarshalJSON implements json.Marshaler.
func (d DbSlice[V]) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.s)
}

// UnmarshalJSON implements json.Unmarshaler.
func (d DbSlice[V]) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &d.s)
}

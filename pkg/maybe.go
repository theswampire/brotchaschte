package pkg

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type Maybe[T any] struct {
	v  T
	ok bool
}

func (m *Maybe[T]) Default(v T) T {
	if m.ok {
		return m.v
	}
	return v
}
func (m Maybe[T]) CopyDefault(v T) T {
	if m.ok {
		return m.v
	}
	return v
}
func (m *Maybe[T]) Unwrap() (T, error) {
	if !m.ok {
		return m.v, errors.New("no value")
	}
	return m.v, nil
}
func (m *Maybe[T]) Ok() bool {
	return m.ok
}
func (m Maybe[T]) MarshalJSON() ([]byte, error) {
	if m.Ok() {
		return json.Marshal(m.v)
	}
	return []byte("null"), nil
}

func Some[T any](v T) Maybe[T] {
	return Maybe[T]{v: v, ok: true}
}
func None[T any]() Maybe[T] {
	return Maybe[T]{ok: false}
}
func NewMaybe[T any](ok bool, v T) Maybe[T] {
	return Maybe[T]{ok: ok, v: v}
}

// all sql Nullable have Value
type nullable[T any] interface {
	Value() (driver.Value, error)
}

func MaybeFromSQL[T any](v nullable[T]) Maybe[T] {
	value, _ := v.Value()
	if value != nil {
		return Maybe[T]{
			ok: true,
			v:  value.(T),
		}
	}
	return Maybe[T]{}
}


package db

import (
	"context"
	"management/internal/shared/hook"
)

type Model interface {
	TableName() string
}

type ModelEvent[T Model] struct {
	hook.Event

	Tx      Transaction
	Context context.Context
	Model   *T
}

type ModelErrorEvent[T Model] struct {
	Error error
	ModelEvent[T]
}

func (e *ModelEvent[T]) Tags() []string {
	if e.Model == nil {
		return nil
	}
	return []string{(*e.Model).TableName()}
}

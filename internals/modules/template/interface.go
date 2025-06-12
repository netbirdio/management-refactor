package template

import (
	"context"

	"github.com/netbirdio/management-refactor/internals/shared/db"
)

type Manager interface {
	// Create
	CreateTemplate(ctx context.Context, tx db.Transaction, template *Template) (*Template, error)

	// Read
	GetAllTemplates(ctx context.Context, tx db.Transaction, strength db.LockingStrength, accountID string) ([]*Template, error)
	GetTemplateByID(ctx context.Context, tx db.Transaction, strength db.LockingStrength, id string) (*Template, error)

	// Update
	UpdateTemplate(ctx context.Context, tx db.Transaction, template *Template) (*Template, error)

	// Delete
	DeleteTemplate(ctx context.Context, tx db.Transaction, id string) error
}

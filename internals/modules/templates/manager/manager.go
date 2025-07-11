package manager

//go:generate go run github.com/golang/mock/mockgen -package manager -destination=manager_mock.go -source=../interface.go -build_flags=-mod=mod

import (
	"context"

	"github.com/netbirdio/management-refactor/internals/modules/templates"
	"github.com/netbirdio/management-refactor/internals/shared/db"
)

var _ templates.Manager = (*Manager)(nil)

type Manager struct {
	repo Repository
}

func NewManager(repo Repository) *Manager {
	return &Manager{repo: repo}
}

func (m *Manager) CreateTemplate(ctx context.Context, tx db.Transaction, template *templates.Template) (*templates.Template, error) {
	return m.repo.CreateTemplate(ctx, tx, template)
}

func (m *Manager) GetAllTemplates(ctx context.Context, tx db.Transaction, strength db.LockingStrength, accountID string) ([]*templates.Template, error) {
	return m.repo.GetAllTemplates(ctx, tx, strength, accountID)
}

func (m *Manager) GetTemplateByID(ctx context.Context, tx db.Transaction, strength db.LockingStrength, id string) (*templates.Template, error) {
	return m.repo.GetTemplateByID(ctx, tx, strength, id)
}

func (m *Manager) UpdateTemplate(ctx context.Context, tx db.Transaction, template *templates.Template) (*templates.Template, error) {
	return m.repo.UpdateTemplate(ctx, tx, template)
}

func (m *Manager) DeleteTemplate(ctx context.Context, tx db.Transaction, id string) error {
	template := &templates.Template{Id: id}
	return m.repo.DeleteTemplate(ctx, tx, template)
}

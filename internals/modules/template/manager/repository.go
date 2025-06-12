package manager

//go:generate go run github.com/golang/mock/mockgen -package manager -destination=repository_mock.go -source=./repository.go -build_flags=-mod=mod

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/netbirdio/management-refactor/internals/modules/template"
	"github.com/netbirdio/management-refactor/internals/shared/db"
)

var _ Repository = (*repository)(nil)

type Repository interface {
	RunInTx(fn func(tx db.Transaction) error) error
	CreateTemplate(ctx context.Context, tx db.Transaction, template *template.Template) (*template.Template, error)
	GetAllTemplates(ctx context.Context, tx db.Transaction, strength db.LockingStrength, accountId string) ([]*template.Template, error)
	GetTemplateByID(ctx context.Context, tx db.Transaction, strength db.LockingStrength, id string) (*template.Template, error)
	UpdateTemplate(ctx context.Context, tx db.Transaction, template *template.Template) (*template.Template, error)
	DeleteTemplate(ctx context.Context, tx db.Transaction, template *template.Template) error
}

type repository struct {
	store *db.Store
}

func NewRepository(s *db.Store) Repository {
	err := s.AutoMigrate(template.Template{})
	if err != nil {
		log.Fatalf("Failed to auto migrate: %v", err)
	}
	return &repository{store: s}
}

func (r *repository) RunInTx(fn func(tx db.Transaction) error) error {
	return r.store.RunInTx(fn)
}

func (r *repository) CreateTemplate(ctx context.Context, tx db.Transaction, template *template.Template) (*template.Template, error) {
	return template, r.store.Create(tx, template)
}

func (r *repository) GetAllTemplates(ctx context.Context, tx db.Transaction, strength db.LockingStrength, accountId string) (templates []*template.Template, err error) {
	err = r.store.GetMany(tx, strength, templates, "account_id = ?", accountId)
	return
}

func (r *repository) GetTemplateByID(ctx context.Context, tx db.Transaction, strength db.LockingStrength, id string) (template *template.Template, err error) {
	err = r.store.GetOne(tx, strength, &template, "id = ?", id)
	return
}

func (r *repository) UpdateTemplate(ctx context.Context, tx db.Transaction, template *template.Template) (*template.Template, error) {
	return template, r.store.Update(tx, template)
}

func (r *repository) DeleteTemplate(ctx context.Context, tx db.Transaction, template *template.Template) error {
	return r.store.Delete(tx, template)
}

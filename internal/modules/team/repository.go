package team

//go:generate go run github.com/golang/mock/mockgen -package team -destination=repository_mock.go -source=./repository.go -build_flags=-mod=mod

import "management/internal/shared/db"

type Repository interface {
	Store() *db.Store
}

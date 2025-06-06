package manager

import "github.com/netbirdio/management-refactor/pkg/logging"

var log = logging.LoggerForThisPackage()

type Manager struct {
	repo    Repository
	handler *handler
}

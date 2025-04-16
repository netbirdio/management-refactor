package accounts

import "management/pkg/logging"

var log = logging.LoggerForThisPackage()

type Manager struct {
	repo    Repository
	handler *handler
}

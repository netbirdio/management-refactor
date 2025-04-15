package errors

import (
	"errors"
	"fmt"
)

var (
	UnsupportedStoreEngineConfigError = errors.New("unsupported store engine")
	PermissionValidationError         = errors.New("permission validation failed")
	PermissionDeniedError             = errors.New("permission denied")
)

func NewUnsupportedStoreEngineConfigError(engine string) error {
	return fmt.Errorf("%w: %s", UnsupportedStoreEngineConfigError, engine)
}

func NewPermissionValidationError(err error) error {
	return fmt.Errorf("%w: %s", PermissionValidationError, err)
}

func NewPermissionDeniedError() error {
	return fmt.Errorf("%w", PermissionDeniedError)
}

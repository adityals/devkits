//go:build !darwin

package cert

import (
	"context"
	"errors"
)

type unsupported struct{}

func NewTrustStore() TrustStore {
	return unsupported{}
}

func (unsupported) Trusted(context.Context, string) (bool, error) {
	return false, errors.ErrUnsupported
}

func (unsupported) Trust(context.Context, string) error {
	return errors.ErrUnsupported
}

func (unsupported) Untrust(context.Context, string) error {
	return errors.ErrUnsupported
}

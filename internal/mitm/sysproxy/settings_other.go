//go:build !darwin

package sysproxy

import (
	"context"
	"errors"
)

type unsupported struct{}

func NewSettings() Settings {
	return unsupported{}
}

func (unsupported) Enabled(context.Context, string, string) (bool, error) {
	return false, errors.ErrUnsupported
}

func (unsupported) Enable(context.Context, string, string) error {
	return errors.ErrUnsupported
}

func (unsupported) Disable(context.Context) error {
	return errors.ErrUnsupported
}

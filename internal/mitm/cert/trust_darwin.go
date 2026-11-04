//go:build darwin

package cert

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"time"
)

const (
	securityBin = "/usr/bin/security"
	// Trust prompts for the user's password, so leave time to type it.
	securityTimeout = 2 * time.Minute
)

// keychain trusts the CA in the user's login keychain, so no admin rights are needed.
type keychain struct{}

func NewTrustStore() TrustStore {
	return keychain{}
}

func (keychain) Trusted(ctx context.Context, certPath string) (bool, error) {
	err := security(ctx, "verify-cert", "-c", certPath)

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return false, nil
	}

	return err == nil, err
}

func (keychain) Trust(ctx context.Context, certPath string) error {
	return security(ctx, "add-trusted-cert", "-r", "trustRoot", certPath)
}

// Untrust removes the trust settings by file, because add-trusted-cert may not store the cert in a keychain.
func (keychain) Untrust(ctx context.Context, certPath string) error {
	return security(ctx, "remove-trusted-cert", certPath)
}

func security(ctx context.Context, args ...string) error {
	ctx, cancel := context.WithTimeout(ctx, securityTimeout)
	defer cancel()

	out, err := exec.CommandContext(ctx, securityBin, args...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("security %s: %w: %s", args[0], err, bytes.TrimSpace(out))
	}

	return nil
}

//go:build darwin

package cert

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
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

// Trust stores the cert in the login keychain too: trust settings alone leave
// macOS unable to find the root when it builds a chain from a leaf.
func (keychain) Trust(ctx context.Context, certPath string) error {
	kc, err := loginKeychain(ctx)
	if err != nil {
		return err
	}

	_, err = security(ctx, "add-trusted-cert", "-r", "trustRoot", "-k", kc, certPath)

	return err
}

// Untrust deletes the cert by SHA-1, the keychain's identifier; -t also drops its trust settings.
func (keychain) Untrust(ctx context.Context, certPath string) error {
	kc, err := loginKeychain(ctx)
	if err != nil {
		return err
	}

	hash, err := sha1Hex(certPath)
	if err != nil {
		return err
	}

	_, err = security(ctx, "delete-certificate", "-Z", hash, "-t", kc)

	return err
}

func loginKeychain(ctx context.Context) (string, error) {
	out, err := security(ctx, "login-keychain")
	if err != nil {
		return "", err
	}

	return strings.Trim(strings.TrimSpace(out), `"`), nil
}

func sha1Hex(certPath string) (string, error) {
	data, err := os.ReadFile(certPath)
	if err != nil {
		return "", fmt.Errorf("read CA cert: %w", err)
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return "", errors.New("read CA cert: no PEM block")
	}

	return fmt.Sprintf("%X", sha1.Sum(block.Bytes)), nil
}

func security(ctx context.Context, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, securityTimeout)
	defer cancel()

	out, err := exec.CommandContext(ctx, securityBin, args...).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("security %s: %w: %s", args[0], err, bytes.TrimSpace(out))
	}

	return string(out), nil
}

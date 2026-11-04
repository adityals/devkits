package cert

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	mcert "github.com/lqqyt2423/go-mitmproxy/cert"
)

// certFile is the public PEM go-mitmproxy's SelfSignCA writes next to its key.
const certFile = "mitmproxy-ca-cert.pem"

type TrustStore interface {
	Trusted(ctx context.Context, certPath string) (bool, error)
	Trust(ctx context.Context, certPath string) error
	Untrust(ctx context.Context, certPath string) error
}

type Manager struct {
	dir   string
	trust TrustStore
}

func NewManager(
	dir string,
	trust TrustStore,
) *Manager {
	return &Manager{
		dir:   dir,
		trust: trust,
	}
}

// Load returns the root CA, creating it on first use.
func (m *Manager) Load() (mcert.CA, error) {
	ca, err := mcert.NewSelfSignCA(m.dir)
	if err != nil {
		return nil, fmt.Errorf("load CA from %s: %w", m.dir, err)
	}

	return ca, nil
}

func (m *Manager) CertPath() string {
	return filepath.Join(m.dir, certFile)
}

func (m *Manager) Trusted(ctx context.Context) (bool, error) {
	_, err := os.Stat(m.CertPath())
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}

	if err != nil {
		return false, fmt.Errorf("stat CA cert: %w", err)
	}

	return m.trust.Trusted(ctx, m.CertPath())
}

func (m *Manager) Use(ctx context.Context) error {
	if _, err := m.Load(); err != nil {
		return err
	}

	return m.trust.Trust(ctx, m.CertPath())
}

func (m *Manager) Remove(ctx context.Context) error {
	return m.trust.Untrust(ctx, m.CertPath())
}

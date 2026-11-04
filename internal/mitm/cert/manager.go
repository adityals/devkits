package cert

import (
	"context"
	"crypto/x509"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	mcert "github.com/lqqyt2423/go-mitmproxy/cert"
)

const (
	// certFile is the public PEM go-mitmproxy's SelfSignCA writes next to its key.
	certFile = "mitmproxy-ca-cert.pem"
	// probeHost names the test leaf Trusted verifies; any hostname works.
	probeHost = "devkits.localhost"
)

type TrustStore interface {
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

// Trusted verifies a leaf the CA issues against the system roots, as browsers do.
// Verifying the root alone is not enough: macOS accepts a self-signed root on its own.
func (m *Manager) Trusted(ctx context.Context) (bool, error) {
	_, err := os.Stat(m.CertPath())
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}

	if err != nil {
		return false, fmt.Errorf("stat CA cert: %w", err)
	}

	ca, err := m.Load()
	if err != nil {
		return false, err
	}

	probe, err := ca.GetCert(probeHost)
	if err != nil {
		return false, fmt.Errorf("issue probe cert: %w", err)
	}

	leaf, err := x509.ParseCertificate(probe.Certificate[0])
	if err != nil {
		return false, fmt.Errorf("parse probe cert: %w", err)
	}

	_, err = leaf.Verify(x509.VerifyOptions{DNSName: probeHost})

	return err == nil, nil
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

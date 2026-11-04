package sysproxy

import (
	"context"
	"sync/atomic"
)

type Settings interface {
	Enabled(ctx context.Context, host, port string) (bool, error)
	Enable(ctx context.Context, host, port string) error
	Disable(ctx context.Context) error
}

type Status struct {
	Enabled bool   `json:"enabled"`
	Address string `json:"address"`
}

type Manager struct {
	settings Settings
	host     string
	port     string
	// owned is set while this app has the system proxy pointed at itself.
	owned atomic.Bool
}

func NewManager(
	settings Settings,
	host string,
	port string,
) *Manager {
	return &Manager{
		settings: settings,
		host:     host,
		port:     port,
	}
}

func (m *Manager) Status(ctx context.Context) (Status, error) {
	enabled, err := m.settings.Enabled(ctx, m.host, m.port)

	return Status{
		Enabled: enabled,
		Address: m.host + ":" + m.port,
	}, err
}

func (m *Manager) Enable(ctx context.Context) error {
	if err := m.settings.Enable(ctx, m.host, m.port); err != nil {
		return err
	}

	m.owned.Store(true)

	return nil
}

func (m *Manager) Disable(ctx context.Context) error {
	if err := m.settings.Disable(ctx); err != nil {
		return err
	}

	m.owned.Store(false)

	return nil
}

// Release turns the system proxy off only if this app turned it on, leaving a user's own proxy alone.
func (m *Manager) Release(ctx context.Context) error {
	if !m.owned.CompareAndSwap(true, false) {
		return nil
	}

	if err := m.settings.Disable(ctx); err != nil {
		m.owned.Store(true)
		return err
	}

	return nil
}

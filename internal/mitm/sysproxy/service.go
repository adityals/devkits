package sysproxy

import "context"

type sysProxyService struct {
	m *Manager
}

func NewSysProxyService(m *Manager) *sysProxyService {
	return &sysProxyService{m: m}
}

func (s *sysProxyService) Status(ctx context.Context) (Status, error) {
	return s.m.Status(ctx)
}

func (s *sysProxyService) Enable(ctx context.Context) error {
	return s.m.Enable(ctx)
}

func (s *sysProxyService) Disable(ctx context.Context) error {
	return s.m.Disable(ctx)
}

// ServiceShutdown restores direct connections, so quitting never leaves the device pointed at a dead proxy.
func (s *sysProxyService) ServiceShutdown() error {
	return s.m.Release(context.Background())
}

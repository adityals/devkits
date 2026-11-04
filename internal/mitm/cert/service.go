package cert

import "context"

type Status struct {
	Path    string `json:"path"`
	Trusted bool   `json:"trusted"`
}

type certService struct {
	m *Manager
}

func NewCertService(m *Manager) *certService {
	return &certService{m: m}
}

func (cs *certService) Status(ctx context.Context) (Status, error) {
	trusted, err := cs.m.Trusted(ctx)

	return Status{
		Path:    cs.m.CertPath(),
		Trusted: trusted,
	}, err
}

func (cs *certService) Use(ctx context.Context) error {
	return cs.m.Use(ctx)
}

func (cs *certService) Remove(ctx context.Context) error {
	return cs.m.Remove(ctx)
}

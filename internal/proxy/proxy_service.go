package proxy

import (
	"context"
	"sync/atomic"

	prox "github.com/lqqyt2423/go-mitmproxy/proxy"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type proxyService struct {
	addr         string
	prox         *prox.Proxy
	eventManager *application.EventManager
	running      atomic.Bool
}

func NewProxyService(addr string) *proxyService {
	return &proxyService{
		addr: addr,
	}
}

func (ps *proxyService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	ps.eventManager = application.Get().Event
	return ps.Start()
}

func (ps *proxyService) ServiceShutdown() error {
	if ps.running.Load() {
		return ps.prox.Close()
	}
	return nil
}

func (ps *proxyService) Start() error {
	p, err := prox.NewProxy(&prox.Options{
		Addr:              ps.addr,
		StreamLargeBodies: 5 * 1024 * 1024, // 5MB
	})
	if err != nil {
		return err
	}

	ps.prox = p

	go func() {
		err := p.Start()
		ps.running.Store(true)
		if err != nil {
			ps.eventManager.Emit("proxy:stopped", err.Error())
		}
	}()

	return nil
}

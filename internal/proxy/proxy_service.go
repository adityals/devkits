package proxy

import (
	"context"
	"errors"
	"net/http"
	"sync"

	prox "github.com/lqqyt2423/go-mitmproxy/proxy"
	"github.com/wailsapp/wails/v3/pkg/application"
)

const (
	ProxyStatusEventName string = "proxy:status"
)

type ProxyStatus struct {
	Running bool   `json:"running"`
	Error   string `json:"error,omitempty"`
	Addr    string `json:"address"`
}

type proxyService struct {
	addr         string
	prox         *prox.Proxy
	eventManager *application.EventManager
	mu           sync.Mutex
}

func NewProxyService(addr string) *proxyService {
	return &proxyService{
		addr: addr,
	}
}

func (ps *proxyService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	ps.eventManager = application.Get().Event
	return nil
}

func (ps *proxyService) ServiceShutdown() error {
	return ps.Stop()
}

func (ps *proxyService) Stop() error {
	ps.mu.Lock()
	p := ps.prox
	ps.prox = nil
	ps.mu.Unlock()

	if p == nil {
		return nil
	}

	return p.Close()
}

func (ps *proxyService) Start() error {
	ps.mu.Lock()

	if ps.prox != nil {
		ps.mu.Unlock()
		return nil
	}

	p, err := prox.NewProxy(&prox.Options{
		Addr:              ps.addr,
		StreamLargeBodies: 5 * 1024 * 1024, // 5MB
	})
	if err != nil {
		ps.mu.Unlock()
		return err
	}

	ps.prox = p
	ps.eventManager.Emit(ProxyStatusEventName, ProxyStatus{
		Running: true,
		Addr:    ps.addr,
	})
	ps.mu.Unlock()

	go func() {
		err := ps.prox.Start()

		ps.mu.Lock()
		if ps.prox == p {
			ps.prox = nil
		}
		ps.mu.Unlock()

		if err != nil {
			errMsg := ""
			if !errors.Is(err, http.ErrServerClosed) {
				errMsg = err.Error()
			}
			ps.eventManager.Emit(ProxyStatusEventName, ProxyStatus{
				Running: false,
				Error:   errMsg,
			})
		}
	}()

	return nil
}

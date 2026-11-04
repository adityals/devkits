package proxy

import (
	"errors"
	"net/http"
	"sync"

	prox "github.com/lqqyt2423/go-mitmproxy/proxy"
)

const (
	ProxyStatusEventName string = "mitm:proxy:status"
)

type ProxyStatus struct {
	Running bool   `json:"running"`
	Error   string `json:"error,omitempty"`
	Addr    string `json:"address"`
}

type Emitter interface {
	Emit(name string, data ...any) bool
}

type proxyService struct {
	emitter Emitter
	addr    string
	caDir   string
	addOns  []prox.Addon
	prox    *prox.Proxy
	mu      sync.Mutex
}

func NewProxyService(
	emitter Emitter,
	addr string,
	caDir string,
	addOns ...prox.Addon,
) *proxyService {
	return &proxyService{
		emitter: emitter,
		addr:    addr,
		caDir:   caDir,
		addOns:  addOns,
	}
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
		// Re-emit so a remounted UI, which starts from "stopped", syncs to the running proxy.
		ps.emitter.Emit(ProxyStatusEventName, ProxyStatus{
			Running: true,
			Addr:    ps.addr,
		})
		return nil
	}

	p, err := prox.NewProxy(&prox.Options{
		Addr:              ps.addr,
		CaRootPath:        ps.caDir,
		StreamLargeBodies: 5 * 1024 * 1024, // 5MB
	})
	if err != nil {
		ps.mu.Unlock()
		return err
	}

	for _, a := range ps.addOns {
		p.AddAddon(a)
	}

	ps.prox = p
	ps.emitter.Emit(ProxyStatusEventName, ProxyStatus{
		Running: true,
		Addr:    ps.addr,
	})
	ps.mu.Unlock()

	go func() {
		err := p.Start()

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
			ps.emitter.Emit(ProxyStatusEventName, ProxyStatus{
				Running: false,
				Error:   errMsg,
			})
		}
	}()

	return nil
}

package mitm

import (
	"devkits/internal/mitm/cert"
	"devkits/internal/mitm/proxy"
	"devkits/internal/mitm/rule"
	"devkits/internal/mitm/sysproxy"
	"devkits/internal/mitm/traffic"
	"path/filepath"

	"github.com/wailsapp/wails/v3/pkg/application"
)

const (
	// proxyAddr listens on all interfaces so other devices can use it; the system proxy targets loopback.
	proxyAddr = ":" + proxyPort
	proxyHost = "127.0.0.1"
	proxyPort = "8000"
)

func init() {
	application.RegisterEvent[proxy.ProxyStatus](proxy.ProxyStatusEventName)
	application.RegisterEvent[traffic.Summary](traffic.FlowEventName)
}

// Services wires the MITM module. Order is startup order: rules load before the proxy can start.
func Services(
	emitter proxy.Emitter,
	configDir string,
) []application.Service {
	engine := rule.NewEngine()
	recorder := traffic.NewRecorder(emitter)
	caDir := filepath.Join(configDir, "ca")
	sysProxy := sysproxy.NewManager(
		sysproxy.NewSettings(),
		proxyHost,
		proxyPort,
	)

	return []application.Service{
		application.NewService(rule.NewRuleService(
			rule.NewFileStore(filepath.Join(configDir, "proxy_rules")),
			engine,
		)),
		application.NewService(traffic.NewTrafficService(recorder)),
		application.NewService(cert.NewCertService(cert.NewManager(
			caDir,
			cert.NewTrustStore(),
		))),
		application.NewService(proxy.NewProxyService(
			emitter,
			proxyAddr,
			caDir,
			sysProxy.Release,
			proxy.NewJsEngineAddOn(engine),
			proxy.NewRecordAddOn(recorder),
		)),
		application.NewService(sysproxy.NewSysProxyService(sysProxy)),
	}
}

package mitm

import (
	"devkits/internal/mitm/proxy"
	"devkits/internal/mitm/rule"
	"devkits/internal/mitm/traffic"
	"path/filepath"

	"github.com/wailsapp/wails/v3/pkg/application"
)

const proxyAddr = ":8000"

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

	return []application.Service{
		application.NewService(rule.NewRuleService(
			rule.NewFileStore(filepath.Join(configDir, "proxy_rules")),
			engine,
		)),
		application.NewService(traffic.NewTrafficService(recorder)),
		application.NewService(proxy.NewProxyService(
			emitter,
			proxyAddr,
			proxy.NewJsEngineAddOn(engine),
			proxy.NewRecordAddOn(recorder),
		)),
	}
}

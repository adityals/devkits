package proxy

import prox "github.com/lqqyt2423/go-mitmproxy/proxy"

type jsEngineAddOn struct {
	prox.BaseAddon
}

func NewJsEngineAddOn() *jsEngineAddOn {
	return &jsEngineAddOn{}
}

func (js *jsEngineAddOn) Requestheaders(f *prox.Flow) {}

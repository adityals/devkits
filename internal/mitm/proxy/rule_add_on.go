package proxy

import (
	"devkits/internal/mitm/rule"
	"log/slog"
	"strconv"

	prox "github.com/lqqyt2423/go-mitmproxy/proxy"
)

type RuleMatcher interface {
	Match(host, path string) *rule.Script
}

// jsEngineAddOn re-matches on every hook instead of caching per flow:
// go-mitmproxy skips Response for streamed bodies, so a per-flow cache would leak.
type jsEngineAddOn struct {
	prox.BaseAddon
	rules RuleMatcher
}

func NewJsEngineAddOn(rules RuleMatcher) *jsEngineAddOn {
	return &jsEngineAddOn{rules: rules}
}

func (js *jsEngineAddOn) Requestheaders(f *prox.Flow) {
	s := js.match(f, rule.HookRequestHeader)
	if s == nil {
		return
	}

	if err := s.OnRequestHeader(f.Request.Header); err != nil {
		slog.Warn("js rule request header", "url", f.Request.URL, "err", err)
	}
}

func (js *jsEngineAddOn) Request(f *prox.Flow) {
	s := js.match(f, rule.HookRequestBody)
	if s == nil {
		return
	}

	body, ok, err := s.OnRequestBody(f.Request.Body)
	if err != nil {
		slog.Warn("js rule request body", "url", f.Request.URL, "err", err)
		return
	}

	if ok {
		f.Request.Body = body
		f.Request.Header.Set("Content-Length", strconv.Itoa(len(body)))
	}
}

func (js *jsEngineAddOn) Responseheaders(f *prox.Flow) {
	s := js.match(f, rule.HookResponseHeader)
	if s == nil {
		return
	}

	if err := s.OnResponseHeader(f.Response.Header); err != nil {
		slog.Warn("js rule response header", "url", f.Request.URL, "err", err)
	}
}

// shortcut: bodies over StreamLargeBodies (5MB) are streamed and never reach Request/Response, so body hooks skip them; add StreamResponseModifier if large bodies matter.
func (js *jsEngineAddOn) Response(f *prox.Flow) {
	s := js.match(f, rule.HookResponseBody)
	if s == nil {
		return
	}

	f.Response.ReplaceToDecodedBody()

	body, ok, err := s.OnResponseBody(f.Response.Body)
	if err != nil {
		slog.Warn("js rule response body", "url", f.Request.URL, "err", err)
		return
	}

	if ok {
		f.Response.Body = body
		f.Response.Header.Set("Content-Length", strconv.Itoa(len(body)))
	}
}

func (js *jsEngineAddOn) match(f *prox.Flow, hook rule.Hook) *rule.Script {
	s := js.rules.Match(f.Request.URL.Host, f.Request.URL.Path)
	if s == nil || !s.Has(hook) {
		return nil
	}

	return s
}

package proxy

import (
	"devkits/internal/mitm/traffic"
	"sync"
	"time"

	prox "github.com/lqqyt2423/go-mitmproxy/proxy"
)

type FlowSink interface {
	Add(traffic.Flow)
	Recording() bool
}

// recordAddOn snapshots a flow once go-mitmproxy finishes it, so streamed and failed flows are recorded too.
type recordAddOn struct {
	prox.BaseAddon
	sink    FlowSink
	rules   RuleMatcher
	pending sync.Map
}

// pendingFlow holds what is known before the flow ends.
type pendingFlow struct {
	rule string
	err  string
}

func NewRecordAddOn(
	sink FlowSink,
	rules RuleMatcher,
) *recordAddOn {
	return &recordAddOn{
		sink:  sink,
		rules: rules,
	}
}

// Requestheaders decides once per flow, so toggling mid-flow never half-records it.
func (ra *recordAddOn) Requestheaders(f *prox.Flow) {
	if !ra.sink.Recording() {
		return
	}

	var p pendingFlow
	if s := ra.rules.Match(f.Request.URL.Host, f.Request.URL.Path); s != nil {
		p.rule = s.ID()
	}

	// Register the flow so RequestError only stores errors the snapshot goroutine will delete.
	ra.pending.Store(f.Id, p)
	go func() {
		<-f.Done()
		ra.sink.Add(ra.snapshot(f))
	}()
}

func (ra *recordAddOn) RequestError(f *prox.Flow, err error) {
	v, ok := ra.pending.Load(f.Id)
	if !ok {
		return
	}

	p := v.(pendingFlow)
	p.err = err.Error()
	ra.pending.Store(f.Id, p)
}

func (ra *recordAddOn) snapshot(f *prox.Flow) traffic.Flow {
	out := traffic.Flow{
		Summary: traffic.Summary{
			ID:         f.Id.String(),
			Method:     f.Request.Method,
			URL:        f.Request.URL.String(),
			StartTime:  f.StartTime,
			DurationMs: time.Since(f.StartTime).Milliseconds(),
		},
		Request: traffic.NewMessage(
			f.Request.Header,
			f.Request.Body,
			f.Stream,
		),
	}
	if v, ok := ra.pending.LoadAndDelete(f.Id); ok {
		p := v.(pendingFlow)
		out.Summary.Rule = p.rule
		out.Summary.Error = p.err
	}

	if f.Response == nil {
		return out
	}

	body := f.Response.Body
	if decoded, err := f.Response.DecodedBody(); err == nil {
		body = decoded
	}

	res := traffic.NewMessage(
		f.Response.Header,
		body,
		f.Stream,
	)
	out.Summary.Status = f.Response.StatusCode
	out.Summary.Size = len(body)
	out.Response = &res

	return out
}

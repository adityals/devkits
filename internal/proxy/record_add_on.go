package proxy

import (
	"devkits/internal/traffic"
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
	sink     FlowSink
	failures sync.Map
}

func NewRecordAddOn(sink FlowSink) *recordAddOn {
	return &recordAddOn{sink: sink}
}

// Requestheaders decides once per flow, so toggling mid-flow never half-records it.
func (ra *recordAddOn) Requestheaders(f *prox.Flow) {
	if !ra.sink.Recording() {
		return
	}

	// Register the flow so RequestError only stores errors the snapshot goroutine will delete.
	ra.failures.Store(f.Id, "")
	go func() {
		<-f.Done()
		ra.sink.Add(ra.snapshot(f))
	}()
}

func (ra *recordAddOn) RequestError(f *prox.Flow, err error) {
	if _, ok := ra.failures.Load(f.Id); ok {
		ra.failures.Store(f.Id, err.Error())
	}
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
	if msg, _ := ra.failures.LoadAndDelete(f.Id); msg != nil {
		out.Summary.Error = msg.(string)
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

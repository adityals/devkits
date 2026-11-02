package traffic

import "sync/atomic"

type Emitter interface {
	Emit(name string, data ...any) bool
}

// Recorder stores finished flows and announces each one to the UI.
type Recorder struct {
	flows     *ring
	emitter   Emitter
	recording atomic.Bool
}

func NewRecorder(emitter Emitter) *Recorder {
	r := &Recorder{
		flows:   newRing(DefaultCapacity),
		emitter: emitter,
	}
	r.recording.Store(true)

	return r
}

func (r *Recorder) Recording() bool {
	return r.recording.Load()
}

func (r *Recorder) SetRecording(on bool) {
	r.recording.Store(on)
}

func (r *Recorder) Add(f Flow) {
	r.flows.add(f)
	r.emitter.Emit(FlowEventName, f.Summary)
}

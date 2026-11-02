package traffic

import "fmt"

type trafficService struct {
	rec *Recorder
}

func NewTrafficService(rec *Recorder) *trafficService {
	return &trafficService{rec: rec}
}

func (ts *trafficService) List() []Summary {
	return ts.rec.flows.summaries()
}

func (ts *trafficService) Get(id string) (Flow, error) {
	f, ok := ts.rec.flows.get(id)
	if !ok {
		return Flow{}, fmt.Errorf("%w: %s", ErrNotFound, id)
	}

	return f, nil
}

func (ts *trafficService) Clear() {
	ts.rec.flows.clear()
}

func (ts *trafficService) Capacity() int {
	return DefaultCapacity
}

func (ts *trafficService) Recording() bool {
	return ts.rec.Recording()
}

func (ts *trafficService) SetRecording(on bool) {
	ts.rec.SetRecording(on)
}

package traffic

import "sync"

// ring keeps the newest flows up to its capacity, evicting the oldest first.
type ring struct {
	mu    sync.Mutex
	flows []Flow
	next  int
}

func newRing(capacity int) *ring {
	return &ring{flows: make([]Flow, 0, capacity)}
}

func (r *ring) add(f Flow) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(r.flows) < cap(r.flows) {
		r.flows = append(r.flows, f)
		return
	}

	r.flows[r.next] = f
	r.next = (r.next + 1) % len(r.flows)
}

// summaries returns flows oldest first.
func (r *ring) summaries() []Summary {
	r.mu.Lock()
	defer r.mu.Unlock()

	out := make([]Summary, 0, len(r.flows))
	for _, f := range r.flows[r.next:] {
		out = append(out, f.Summary)
	}

	for _, f := range r.flows[:r.next] {
		out = append(out, f.Summary)
	}

	return out
}

func (r *ring) get(id string) (Flow, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, f := range r.flows {
		if f.Summary.ID == id {
			return f, true
		}
	}

	return Flow{}, false
}

func (r *ring) clear() {
	r.mu.Lock()
	defer r.mu.Unlock()

	clear(r.flows)
	r.flows = r.flows[:0]
	r.next = 0
}

package traffic

import (
	"slices"
	"strconv"
	"testing"
)

func TestRingEvictsOldest(t *testing.T) {
	r := newRing(3)
	for i := range 5 {
		r.add(Flow{Summary: Summary{ID: strconv.Itoa(i)}})
	}

	var ids []string
	for _, s := range r.summaries() {
		ids = append(ids, s.ID)
	}

	if want := []string{"2", "3", "4"}; !slices.Equal(ids, want) {
		t.Fatalf("ids = %v, want %v", ids, want)
	}

	if _, ok := r.get("1"); ok {
		t.Fatal("evicted flow still found")
	}

	if _, ok := r.get("4"); !ok {
		t.Fatal("newest flow missing")
	}

	r.clear()
	r.add(Flow{Summary: Summary{ID: "x"}})
	if got := r.summaries(); len(got) != 1 || got[0].ID != "x" {
		t.Fatalf("after clear = %v", got)
	}
}

func TestNewMessageTruncates(t *testing.T) {
	m := NewMessage(
		nil,
		make([]byte, MaxBodySize+1),
		false,
	)
	if !m.Truncated || len(m.Body) != MaxBodySize {
		t.Fatalf("truncated = %v, len = %d", m.Truncated, len(m.Body))
	}
}

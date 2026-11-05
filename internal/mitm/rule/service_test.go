package rule

import (
	"errors"
	"fmt"
	"testing"
)

type memStore struct {
	rules []Rule
}

func (m *memStore) Load() ([]Rule, error) {
	return m.rules, nil
}

func (m *memStore) Save(Rule) error {
	return nil
}

func (m *memStore) Delete(string) error {
	return nil
}

func TestOverLimitTurnsAllRulesOff(t *testing.T) {
	store := &memStore{}
	for i := range MaxRules + 1 {
		store.rules = append(store.rules, Rule{
			ID:      fmt.Sprintf("r%d", i),
			Enabled: true,
			Script:  `const rule = { match: '/a', onRequestHeader() {} }`,
		})
	}

	engine := NewEngine()
	svc := NewRuleService(store, engine)

	if err := svc.Reload(); !errors.Is(err, ErrMaxRules) {
		t.Fatalf("Reload err = %v, want ErrMaxRules", err)
	}

	if got := len(svc.List()); got != MaxRules+1 {
		t.Fatalf("List len = %d, want %d", got, MaxRules+1)
	}

	if s := engine.Match("a.com", "/a"); s != nil {
		t.Fatal("rule ran while over the limit")
	}

	// Disabling one brings it back under the limit and the rest run again.
	off := svc.List()[0]
	off.Enabled = false
	if _, err := svc.Update(off); err != nil {
		t.Fatalf("Update: %v", err)
	}

	if s := engine.Match("a.com", "/a"); s == nil {
		t.Fatal("rules stayed off under the limit")
	}

	if _, err := svc.Update(Rule{ID: off.ID, Enabled: true, Script: off.Script}); !errors.Is(err, ErrMaxRules) {
		t.Fatalf("re-enable err = %v, want ErrMaxRules", err)
	}
}

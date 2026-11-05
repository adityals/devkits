package rule

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"sync"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type ruleService struct {
	store  Store
	engine *Engine
	mu     sync.Mutex
	rules  []Rule
}

func NewRuleService(store Store, engine *Engine) *ruleService {
	return &ruleService{
		store:  store,
		engine: engine,
	}
}

func (s *ruleService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	// a broken rule file must not stop the app; the user fixes it and calls Reload.
	if err := s.Reload(); err != nil {
		slog.Error("load proxy rules", "err", err)
	}

	return nil
}

// Capacity is the most rules that may be enabled at once.
func (s *ruleService) Capacity() int {
	return MaxRules
}

func (s *ruleService) List() []Rule {
	s.mu.Lock()
	defer s.mu.Unlock()

	return slices.Clone(s.rules)
}

func (s *ruleService) Create(r Rule) (Rule, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.index(r.ID) >= 0 {
		return Rule{}, fmt.Errorf("%w: %s", ErrExists, r.ID)
	}

	return s.put(append(slices.Clone(s.rules), r), r)
}

func (s *ruleService) Update(r Rule) (Rule, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	i := s.index(r.ID)
	if i < 0 {
		return Rule{}, fmt.Errorf("%w: %s", ErrNotFound, r.ID)
	}

	next := slices.Clone(s.rules)
	next[i] = r

	return s.put(next, r)
}

func (s *ruleService) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	i := s.index(id)
	if i < 0 {
		return fmt.Errorf("%w: %s", ErrNotFound, id)
	}

	next := slices.Delete(slices.Clone(s.rules), i, i+1)

	return s.apply(next, func() error {
		return s.store.Delete(id)
	})
}

// Reload re-reads rules from disk; on error the running rules stay active.
func (s *ruleService) Reload() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	rules, err := s.store.Load()
	if err != nil {
		return err
	}

	return s.apply(rules, func() error {
		return nil
	})
}

// put validates r before anything is written to disk.
func (s *ruleService) put(next []Rule, r Rule) (Rule, error) {
	if err := ValidateID(r.ID); err != nil {
		return Rule{}, err
	}

	// Refuse to go over the limit, but always allow disabling so the user can get back under it.
	if r.Enabled && enabledCount(next) > MaxRules {
		return Rule{}, fmt.Errorf("%w: disable a rule first", ErrMaxRules)
	}

	if _, err := s.engine.compile(r); err != nil {
		return Rule{}, err
	}

	err := s.apply(next, func() error {
		return s.store.Save(r)
	})
	if err != nil {
		return Rule{}, err
	}

	return r, nil
}

// apply compiles, persists, then swaps, so a failure leaves memory unchanged.
// Over the limit, every rule stays listed but none run: picking which ones would be a guess.
func (s *ruleService) apply(next []Rule, persist func() error) error {
	enabled := enabledCount(next)

	active := next
	if enabled > MaxRules {
		active = nil
	}

	set, err := s.engine.build(active)
	if err != nil {
		return err
	}

	if err := persist(); err != nil {
		return err
	}

	s.engine.swap(set)
	s.rules = next

	if enabled > MaxRules {
		return fmt.Errorf("%w: %d enabled, max %d; all rules are off", ErrMaxRules, enabled, MaxRules)
	}

	return nil
}

func enabledCount(rules []Rule) int {
	n := 0
	for _, r := range rules {
		if r.Enabled {
			n++
		}
	}

	return n
}

func (s *ruleService) index(id string) int {
	return slices.IndexFunc(s.rules, func(r Rule) bool {
		return r.ID == id
	})
}

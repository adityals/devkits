package rule

import (
	"sync/atomic"
	"time"
)

type ruleSet struct {
	scripts  []*Script
	matchers []matcher
}

type Engine struct {
	timeout time.Duration
	active  atomic.Pointer[ruleSet]
}

type Option func(*Engine)

func WithTimeout(d time.Duration) Option {
	return func(e *Engine) {
		e.timeout = d
	}
}

func NewEngine(opts ...Option) *Engine {
	e := &Engine{timeout: DefaultTimeout}
	for _, opt := range opts {
		opt(e)
	}

	e.active.Store(&ruleSet{})
	return e
}

// Match returns the best rule for path, or nil. Safe for concurrent use.
func (e *Engine) Match(path string) *Script {
	set := e.active.Load()
	if i := selectBest(set.matchers, path); i >= 0 {
		return set.scripts[i]
	}

	return nil
}

func (e *Engine) compile(r Rule) (*Script, error) {
	return compile(r, e.timeout)
}

// build compiles every enabled rule; it is all-or-nothing so a bad rule never half-applies.
func (e *Engine) build(rules []Rule) (*ruleSet, error) {
	set := &ruleSet{}
	for _, r := range rules {
		if !r.Enabled {
			continue
		}

		s, err := e.compile(r)
		if err != nil {
			return nil, err
		}

		set.scripts = append(set.scripts, s)
		set.matchers = append(set.matchers, s.matcher)
	}

	return set, nil
}

func (e *Engine) swap(set *ruleSet) {
	e.active.Store(set)
}

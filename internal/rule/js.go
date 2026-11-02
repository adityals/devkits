package rule

import (
	"errors"
	"fmt"
	"log/slog"
	"runtime"
	"time"

	"github.com/dop251/goja"
)

type Hook int

const (
	HookRequestHeader Hook = iota
	HookRequestBody
	HookResponseHeader
	HookResponseBody
	hookCount
)

var hookNames = [hookCount]string{
	"onRequestHeader",
	"onRequestBody",
	"onResponseHeader",
	"onResponseBody",
}

var ruleRef = goja.MustCompile("rule-ref", "rule", true)

type jsRuntime struct {
	vm    *goja.Runtime
	this  *goja.Object
	hooks [hookCount]goja.Callable
}

type Script struct {
	id      string
	host    hostMatcher
	matcher matcher
	program *goja.Program
	timeout time.Duration
	logger  *slog.Logger
	has     [hookCount]bool
	idle    chan *jsRuntime
}

func compile(r Rule, timeout time.Duration, logger *slog.Logger) (*Script, error) {
	prog, err := goja.Compile(r.ID+".js", r.Script, true)
	if err != nil {
		return nil, fmt.Errorf("%w: compile: %w", ErrInvalidRule, err)
	}

	s := &Script{
		id:      r.ID,
		program: prog,
		timeout: timeout,
		logger:  logger,
		idle:    make(chan *jsRuntime, runtime.GOMAXPROCS(0)),
	}
	rt, err := s.newRuntime()
	if err != nil {
		return nil, err
	}

	mod, okMod := stringProp(rt.this, "modifier")
	match, okMatch := stringProp(rt.this, "match")

	if !okMod || !okMatch {
		return nil, fmt.Errorf("%w: modifier and match must be strings", ErrInvalidRule)
	}

	if s.matcher, err = newMatcher(mod, match); err != nil {
		return nil, err
	}

	hosts, ok := stringsProp(rt.this, "host")
	if !ok {
		return nil, fmt.Errorf("%w: host must be a string or an array of strings", ErrInvalidRule)
	}

	if s.host, err = newHostMatcher(hosts); err != nil {
		return nil, err
	}

	for i, fn := range rt.hooks {
		s.has[i] = fn != nil
	}

	s.release(rt)
	return s, nil
}

func (s *Script) newRuntime() (*jsRuntime, error) {
	vm := goja.New()
	if err := installConsole(vm, s.logger, s.id); err != nil {
		return nil, fmt.Errorf("install console: %w", err)
	}

	v, err := s.guard(vm, func() (goja.Value, error) {
		if _, err := vm.RunProgram(s.program); err != nil {
			return nil, err
		}
		return vm.RunProgram(ruleRef)
	})
	if err != nil {
		return nil, fmt.Errorf("%w: script must define `const rule = {...}`: %w", ErrInvalidRule, err)
	}

	obj, ok := v.(*goja.Object)
	if !ok {
		return nil, fmt.Errorf("%w: rule must be an object", ErrInvalidRule)
	}

	rt := &jsRuntime{vm: vm, this: obj}
	found := false
	for i, name := range hookNames {
		hv := obj.Get(name)
		if hv == nil || goja.IsUndefined(hv) {
			continue
		}
		fn, ok := goja.AssertFunction(hv)
		if !ok {
			return nil, fmt.Errorf("%w: %s must be a function", ErrInvalidRule, name)
		}
		rt.hooks[i], found = fn, true
	}
	if !found {
		return nil, fmt.Errorf("%w: rule needs at least one of %v", ErrInvalidRule, hookNames)
	}
	return rt, nil
}

func (s *Script) acquire() (*jsRuntime, error) {
	select {
	case rt := <-s.idle:
		return rt, nil
	default:
		return s.newRuntime()
	}
}

func (s *Script) release(rt *jsRuntime) {
	select {
	case s.idle <- rt:
	default: // idle pool full; drop rt so GC reclaims it
	}
}

// guard interrupts the VM after s.timeout and leaves it reusable.
func (s *Script) guard(vm *goja.Runtime, fn func() (goja.Value, error)) (goja.Value, error) {
	fired := make(chan struct{})
	t := time.AfterFunc(s.timeout, func() {
		vm.Interrupt(ErrScriptTimeout)
		close(fired)
	})

	v, err := fn()
	if !t.Stop() {
		<-fired
	}

	vm.ClearInterrupt()

	if _, ok := errors.AsType[*goja.InterruptedError](err); ok {
		return nil, fmt.Errorf("%w after %s", ErrScriptTimeout, s.timeout)
	}
	return v, err
}

// stringProp reports false only when the property exists and is not a string.
func stringProp(o *goja.Object, name string) (string, bool) {
	v := o.Get(name)
	if v == nil || goja.IsUndefined(v) {
		return "", true
	}

	str, ok := v.Export().(string)
	return str, ok
}

// stringsProp accepts a string or an array of strings; false means any other type.
func stringsProp(o *goja.Object, name string) ([]string, bool) {
	v := o.Get(name)
	if v == nil || goja.IsUndefined(v) {
		return nil, true
	}

	switch x := v.Export().(type) {
	case string:
		return []string{x}, true
	case []any:
		out := make([]string, 0, len(x))
		for _, item := range x {
			str, ok := item.(string)
			if !ok {
				return nil, false
			}

			out = append(out, str)
		}

		return out, true
	}

	return nil, false
}

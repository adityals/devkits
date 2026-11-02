package rule

import (
	"fmt"
	"net/http"

	"github.com/dop251/goja"
)

func (s *Script) Has(h Hook) bool { return s.has[h] }

func (s *Script) OnRequestHeader(h http.Header) error {
	return s.callHeader(HookRequestHeader, h)
}

func (s *Script) OnResponseHeader(h http.Header) error {
	return s.callHeader(HookResponseHeader, h)
}

func (s *Script) OnRequestBody(b []byte) ([]byte, bool, error) {
	return s.callBody(HookRequestBody, b)
}

func (s *Script) OnResponseBody(b []byte) ([]byte, bool, error) {
	return s.callBody(HookResponseBody, b)
}

func (s *Script) callHeader(hook Hook, h http.Header) error {
	if !s.has[hook] {
		return nil
	}

	rt, err := s.acquire()
	if err != nil {
		return err
	}
	defer s.release(rt)

	_, err = s.guard(rt.vm, func() (goja.Value, error) {
		return rt.hooks[hook](rt.this, headerObject(rt.vm, h))
	})
	if err != nil {
		return fmt.Errorf("rule %s %s: %w", s.id, hookNames[hook], err)
	}
	return nil
}

// callBody returns the replacement body and true when the hook returned a string.
func (s *Script) callBody(hook Hook, b []byte) ([]byte, bool, error) {
	if !s.has[hook] {
		return nil, false, nil
	}

	rt, err := s.acquire()
	if err != nil {
		return nil, false, err
	}
	defer s.release(rt)

	v, err := s.guard(rt.vm, func() (goja.Value, error) {
		return rt.hooks[hook](rt.this, rt.vm.ToValue(string(b)))
	})
	if err != nil {
		return nil, false, fmt.Errorf("rule %s %s: %w", s.id, hookNames[hook], err)
	}

	if v == nil || goja.IsUndefined(v) || goja.IsNull(v) {
		return nil, false, nil
	}

	str, ok := v.Export().(string)
	if !ok {
		return nil, false, fmt.Errorf("rule %s %s: must return string or undefined", s.id, hookNames[hook])
	}

	return []byte(str), true, nil
}

func headerObject(vm *goja.Runtime, h http.Header) *goja.Object {
	o := vm.NewObject()
	_ = o.Set("get", h.Get)
	_ = o.Set("values", func(k string) []string {
		return append([]string(nil), h.Values(k)...)
	})
	_ = o.Set("set", h.Set)
	_ = o.Set("add", h.Add)
	_ = o.Set("del", h.Del)
	return o
}

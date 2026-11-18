package rule

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/dop251/goja"
)

type HooksFlow struct {
	URL *url.URL
}

func (s *Script) Has(h Hook) bool { return s.has[h] }

func (s *Script) OnRequestHeader(f HooksFlow, h http.Header) error {
	return s.callHeader(HookRequestHeader, f, h)
}

func (s *Script) OnResponseHeader(f HooksFlow, h http.Header) error {
	return s.callHeader(HookResponseHeader, f, h)
}

func (s *Script) OnRequestBody(f HooksFlow, b []byte) ([]byte, bool, error) {
	return s.callBody(HookRequestBody, f, b)
}

func (s *Script) OnResponseBody(f HooksFlow, b []byte) ([]byte, bool, error) {
	return s.callBody(HookResponseBody, f, b)
}

func (s *Script) callHeader(hook Hook, f HooksFlow, h http.Header) error {
	if !s.has[hook] {
		return nil
	}

	rt, err := s.acquire()
	if err != nil {
		return err
	}
	defer s.release(rt)

	_, err = s.guard(rt.vm, func() (goja.Value, error) {
		return rt.hooks[hook](rt.this, flowObject(rt.vm, f), headerObject(rt.vm, h))
	})
	if err != nil {
		return fmt.Errorf("rule %s %s: %w", s.id, hookNames[hook], err)
	}
	return nil
}

// callBody returns the replacement body and true when the hook returned a string.
func (s *Script) callBody(hook Hook, f HooksFlow, b []byte) ([]byte, bool, error) {
	if !s.has[hook] {
		return nil, false, nil
	}

	rt, err := s.acquire()
	if err != nil {
		return nil, false, err
	}
	defer s.release(rt)

	v, err := s.guard(rt.vm, func() (goja.Value, error) {
		return rt.hooks[hook](rt.this, flowObject(rt.vm, f), rt.vm.ToValue(string(b)))
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
	_ = o.Set("toObject", func() http.Header {
		return h.Clone()
	})
	_ = o.Set("set", h.Set)
	_ = o.Set("add", h.Add)
	_ = o.Set("del", h.Del)
	return o
}

func flowObject(vm *goja.Runtime, f HooksFlow) *goja.Object {
	u := f.URL
	if u == nil {
		u = &url.URL{}
	}

	o := vm.NewObject()
	_ = o.Set("host", u.Host)
	_ = o.Set("url", u.String())
	_ = o.Set("path", u.Path)
	_ = o.Set("query", map[string][]string(u.Query()))
	return o
}

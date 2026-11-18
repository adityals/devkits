package rule

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"sync"
	"testing"
	"time"
)

const fullScript = `
const rule = {
	modifier: '=',
	match: '/foo',
	onRequestHeader(f, h) { h.set('X-Req', 'yes'); h.del('X-Drop') },
	onRequestBody(f, body) { return body.toUpperCase() },
	onResponseHeader(f, h) { h.add('X-Res', h.get('X-Base') + '!') },
	onResponseBody(f, body) { return undefined },
}`

func TestScriptHooks(t *testing.T) {
	s, err := compile(Rule{ID: "full", Script: fullScript}, DefaultTimeout, slog.Default())
	if err != nil {
		t.Fatal(err)
	}

	if !s.matcher.matches("/foo") || s.matcher.matches("/foo/bar") {
		t.Fatal("matcher not built from script")
	}

	h := http.Header{
		"X-Drop": {"1"},
		"X-Base": {"b"},
	}

	if err := s.OnRequestHeader(HooksFlow{URL: &url.URL{Host: "testing.com"}}, h); err != nil {
		t.Fatal(err)
	}

	if err := s.OnResponseHeader(HooksFlow{URL: &url.URL{Host: "testing.com"}}, h); err != nil {
		t.Fatal(err)
	}

	if h.Get("X-Req") != "yes" || h.Get("X-Drop") != "" || h.Get("X-Res") != "b!" {
		t.Fatalf("headers not modified: %v", h)
	}

	body, ok, err := s.OnRequestBody(HooksFlow{}, []byte("abc"))
	if err != nil || !ok || string(body) != "ABC" {
		t.Fatalf("request body = %q, %v, %v", body, ok, err)
	}

	if _, ok, err := s.OnResponseBody(HooksFlow{}, []byte("abc")); err != nil || ok {
		t.Fatalf("undefined return must keep body: %v, %v", ok, err)
	}
}

func TestFlow(t *testing.T) {
	src := `const rule = {
		match: '/',
		onRequestHeader(f, h) {
			h.set('X-Host', f.host)
			h.set('X-Path', f.path)
			h.set('X-Url', f.url)
			h.set('X-Query', f.query.q[0])
		},
		onResponseHeader(f, h) { h.set('X-Res-Host', f.host) },
	}`
	s, err := compile(Rule{ID: "flow", Script: src}, DefaultTimeout, slog.Default())
	if err != nil {
		t.Fatal(err)
	}

	u, err := url.Parse("https://api.testing.com:8443/v1/items?q=go&q=js")
	if err != nil {
		t.Fatal(err)
	}

	h := http.Header{}
	if err := s.OnRequestHeader(HooksFlow{URL: u}, h); err != nil {
		t.Fatal(err)
	}

	if err := s.OnResponseHeader(HooksFlow{URL: u}, h); err != nil {
		t.Fatal(err)
	}

	want := map[string]string{
		"X-Host":     "api.testing.com:8443",
		"X-Path":     "/v1/items",
		"X-Url":      "https://api.testing.com:8443/v1/items?q=go&q=js",
		"X-Query":    "go",
		"X-Res-Host": "api.testing.com:8443",
	}
	for k, v := range want {
		if got := h.Get(k); got != v {
			t.Errorf("%s = %q, want %q", k, got, v)
		}
	}
}

func TestCompileRejectsInvalid(t *testing.T) {
	cases := map[string]string{
		"syntax":         `const rule = {`,
		"no rule":        `const other = {}`,
		"bad modifier":   `const rule = { modifier: '!', match: '/a', onRequestHeader() {} }`,
		"no match":       `const rule = { onRequestHeader() {} }`,
		"hook not fn":    `const rule = { match: '/a', onRequestBody: 1 }`,
		"no hooks":       `const rule = { match: '/a' }`,
		"top-level loop": `while (true) {}`,
		"host number":    `const rule = { host: 1, match: '/a', onRequestHeader() {} }`,
		"host mixed":     `const rule = { host: ['a.com', 2], match: '/a', onRequestHeader() {} }`,
		"host pattern":   `const rule = { host: 'a*.com', match: '/a', onRequestHeader() {} }`,
	}
	for name, src := range cases {
		t.Run(name, func(t *testing.T) {
			if _, err := compile(Rule{ID: "bad", Script: src}, 50*time.Millisecond, slog.Default()); err == nil {
				t.Fatal("expected error")
			}
		})
	}
}

func TestHookTimeoutKeepsRuntimeUsable(t *testing.T) {
	src := `const rule = {
		match: '/',
		onRequestBody(f, b) { if (b === 'loop') { while (true) {} } return b + '!' },
	}`
	s, err := compile(Rule{ID: "loop", Script: src}, 50*time.Millisecond, slog.Default())
	if err != nil {
		t.Fatal(err)
	}

	if _, _, err := s.OnRequestBody(HooksFlow{}, []byte("loop")); !errors.Is(err, ErrScriptTimeout) {
		t.Fatalf("want ErrScriptTimeout, got %v", err)
	}

	if body, _, err := s.OnRequestBody(HooksFlow{}, []byte("ok")); err != nil || string(body) != "ok!" {
		t.Fatalf("runtime unusable after timeout: %q, %v", body, err)
	}
}

func TestConcurrentHooks(t *testing.T) {
	s, err := compile(Rule{ID: "full", Script: fullScript}, DefaultTimeout, slog.Default())
	if err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup
	for i := range 50 {
		wg.Go(func() {
			want := fmt.Sprintf("BODY-%d", i)
			out, _, err := s.OnRequestBody(HooksFlow{}, fmt.Appendf(nil, "body-%d", i))
			if err != nil || string(out) != want {
				t.Errorf("got %q, %v, want %q", out, err, want)
			}

			if err := s.OnRequestHeader(HooksFlow{URL: &url.URL{Host: "testing.com"}}, http.Header{}); err != nil {
				t.Error(err)
			}
		})
	}

	wg.Wait()
}

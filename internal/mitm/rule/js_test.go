package rule

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"testing"
	"time"
)

const fullScript = `
const rule = {
	modifier: '=',
	match: '/foo',
	onRequestHeader(h) { h.set('X-Req', 'yes'); h.del('X-Drop') },
	onRequestBody(body) { return body.toUpperCase() },
	onResponseHeader(h) { h.add('X-Res', h.get('X-Base') + '!') },
	onResponseBody(body) { return undefined },
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
	if err := s.OnRequestHeader(h); err != nil {
		t.Fatal(err)
	}

	if err := s.OnResponseHeader(h); err != nil {
		t.Fatal(err)
	}

	if h.Get("X-Req") != "yes" || h.Get("X-Drop") != "" || h.Get("X-Res") != "b!" {
		t.Fatalf("headers not modified: %v", h)
	}

	body, ok, err := s.OnRequestBody([]byte("abc"))
	if err != nil || !ok || string(body) != "ABC" {
		t.Fatalf("request body = %q, %v, %v", body, ok, err)
	}

	if _, ok, err := s.OnResponseBody([]byte("abc")); err != nil || ok {
		t.Fatalf("undefined return must keep body: %v, %v", ok, err)
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
		onRequestBody(b) { if (b === 'loop') { while (true) {} } return b + '!' },
	}`
	s, err := compile(Rule{ID: "loop", Script: src}, 50*time.Millisecond, slog.Default())
	if err != nil {
		t.Fatal(err)
	}

	if _, _, err := s.OnRequestBody([]byte("loop")); !errors.Is(err, ErrScriptTimeout) {
		t.Fatalf("want ErrScriptTimeout, got %v", err)
	}

	if body, _, err := s.OnRequestBody([]byte("ok")); err != nil || string(body) != "ok!" {
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
			out, _, err := s.OnRequestBody(fmt.Appendf(nil, "body-%d", i))
			if err != nil || string(out) != want {
				t.Errorf("got %q, %v, want %q", out, err, want)
			}

			if err := s.OnRequestHeader(http.Header{}); err != nil {
				t.Error(err)
			}
		})
	}

	wg.Wait()
}

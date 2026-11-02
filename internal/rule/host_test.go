package rule

import "testing"

func TestSelectRuleByHost(t *testing.T) {
	var hosts []hostMatcher
	var paths []matcher
	for _, r := range []struct {
		hosts      []string
		mod, match string
	}{
		{[]string{"api.example.com"}, "", "/"},
		{[]string{"*.example.com"}, "", "/"},
		{[]string{"*.staging.example.com"}, "", "/"},
		{[]string{"api.*"}, "", "/"},
		{nil, "", "/"},
		{[]string{"localhost:3000"}, "", "/"},
		{[]string{"shop.com"}, "=", "/cart"},
	} {
		h, err := newHostMatcher(r.hosts)
		if err != nil {
			t.Fatal(err)
		}

		p, err := newMatcher(r.mod, r.match)
		if err != nil {
			t.Fatal(err)
		}

		hosts = append(hosts, h)
		paths = append(paths, p)
	}

	for _, tt := range []struct {
		name, host, path string
		want             int
	}{
		{"exact, case-insensitive", "API.Example.com", "/x", 0},
		{"longest leading wildcard", "a.staging.example.com", "/", 2},
		{"leading wildcard", "web.example.com", "/", 1},
		{"leading beats trailing", "api.foo.example.com", "/", 1},
		{"trailing wildcard", "api.local", "/", 3},
		{"trailing needs prefix match", "myapi.local", "/", 4},
		{"*. is subdomains only", "example.com", "/", 4},
		{"port-qualified name", "localhost:3000", "/", 5},
		{"other port falls to any host", "localhost:8080", "/", 4},
		{"host path hit", "shop.com", "/cart", 6},
		{"host path miss falls back", "shop.com", "/home", 4},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := selectRule(hosts, paths, tt.host, tt.path); got != tt.want {
				t.Errorf("selectRule(%q, %q) = %d, want %d", tt.host, tt.path, got, tt.want)
			}
		})
	}
}

func TestNewHostMatcherInvalid(t *testing.T) {
	for _, name := range []string{
		"",
		"*.",
		"a*.com",
		"*.a.*",
		"api.*:3000",
		"a.com/x",
	} {
		if _, err := newHostMatcher([]string{name}); err == nil {
			t.Errorf("newHostMatcher(%q) expected error", name)
		}
	}
}

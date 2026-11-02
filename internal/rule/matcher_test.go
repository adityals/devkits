package rule

import "testing"

func TestSelectBest(t *testing.T) {
	var ms []matcher
	for _, p := range [][2]string{
		{"", "/api"},
		{"", "/api/users"},
		{"~", `^/api/users/\d+$`},
		{"=", "/api/users/1"},
		{"^~", "/static"},
		{"~*", `\.PNG$`},
	} {
		m, err := newMatcher(p[0], p[1])
		if err != nil {
			t.Fatal(err)
		}
		ms = append(ms, m)
	}

	for path, want := range map[string]int{
		"/api/users/1":     3, // exact beats regex and prefix
		"/api/users/2":     2, // regex beats plain prefix
		"/api/users":       1, // longest plain prefix
		"/api/orders":      0,
		"/static/logo.png": 4, // ^~ beats regex
		"/img/logo.png":    5, // ~* case-insensitive
		"/other":           -1,
	} {
		if got := selectBest(ms, path); got != want {
			t.Errorf("selectBest(%q) = %d, want %d", path, got, want)
		}
	}

	for _, p := range [][2]string{
		{"!", "/a"},
		{"=", ""},
		{"~", "("},
	} {
		if _, err := newMatcher(p[0], p[1]); err == nil {
			t.Errorf("newMatcher(%q, %q) expected error", p[0], p[1])
		}
	}
}

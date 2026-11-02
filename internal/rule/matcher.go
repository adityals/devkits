package rule

import (
	"fmt"
	"regexp"
	"strings"
)

type modifier string

const (
	modExact       modifier = "="
	modPrefixFirst modifier = "^~"
	modRegex       modifier = "~"
	modRegexI      modifier = "~*"
	modPrefix      modifier = ""
)

type matcher struct {
	mod     modifier
	pattern string
	re      *regexp.Regexp
}

func newMatcher(mod, pattern string) (matcher, error) {
	if pattern == "" {
		return matcher{}, fmt.Errorf("%w: match must be a non-empty string", ErrInvalidRule)
	}

	m := matcher{mod: modifier(mod), pattern: pattern}
	switch m.mod {
	case modExact, modPrefixFirst, modPrefix:
	case modRegex, modRegexI:
		expr := pattern
		if m.mod == modRegexI {
			expr = "(?i)" + expr
		}
		re, err := regexp.Compile(expr)
		if err != nil {
			return matcher{}, fmt.Errorf("%w: match regex: %w", ErrInvalidRule, err)
		}
		m.re = re
	default:
		return matcher{}, fmt.Errorf("%w: unknown modifier %q", ErrInvalidRule, mod)
	}
	return m, nil
}

func (m matcher) matches(path string) bool {
	switch m.mod {
	case modExact:
		return path == m.pattern
	case modRegex, modRegexI:
		return m.re.MatchString(path)
	default:
		return strings.HasPrefix(path, m.pattern)
	}
}

// selectBest returns the index of the winning matcher using nginx location priority, or -1.
func selectBest(ms []matcher, path string) int {
	bestPrefix, bestPrefixFirst := -1, false
	for i, m := range ms {
		if m.mod == modExact && m.matches(path) {
			return i
		}
		if (m.mod == modPrefix || m.mod == modPrefixFirst) && m.matches(path) &&
			(bestPrefix < 0 || len(m.pattern) > len(ms[bestPrefix].pattern)) {
			bestPrefix, bestPrefixFirst = i, m.mod == modPrefixFirst
		}
	}

	if bestPrefixFirst {
		return bestPrefix
	}

	for i, m := range ms {
		if (m.mod == modRegex || m.mod == modRegexI) && m.matches(path) {
			return i
		}
	}

	return bestPrefix
}

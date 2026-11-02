package rule

import (
	"fmt"
	"net"
	"slices"
	"strings"
)

type hostKind int

const (
	hostExact hostKind = iota
	hostLeading
	hostTrailing
)

// tier scores keep exact > *.suffix > prefix.* apart; hostnames are at most 253 chars.
const (
	tierTrailing = 1000
	tierLeading  = 2000
	tierExact    = 3000
)

type hostName struct {
	kind     hostKind
	value    string
	withPort bool
}

// hostMatcher with no names matches any host, like nginx default_server.
type hostMatcher struct {
	names []hostName
}

func newHostMatcher(names []string) (hostMatcher, error) {
	var hm hostMatcher
	for _, raw := range names {
		n, err := parseHostName(strings.TrimSuffix(strings.ToLower(raw), "."))
		if err != nil {
			return hostMatcher{}, err
		}

		hm.names = append(hm.names, n)
	}

	return hm, nil
}

func parseHostName(name string) (hostName, error) {
	n := hostName{
		kind:  hostExact,
		value: name,
	}

	switch lead, trail := strings.HasPrefix(name, "*."), strings.HasSuffix(name, ".*"); {
	case lead && trail:
		return hostName{}, fmt.Errorf("%w: host %q may not start and end with *", ErrInvalidRule, name)
	case lead:
		n.kind, n.value = hostLeading, name[1:]
	case trail:
		n.kind, n.value = hostTrailing, name[:len(name)-1]
	}

	n.withPort = strings.Contains(n.value, ":")
	if len(n.value) < 2 || strings.ContainsAny(n.value, "*/ ") || (n.kind != hostExact && n.withPort) {
		return hostName{}, fmt.Errorf("%w: invalid host %q", ErrInvalidRule, name)
	}

	return n, nil
}

// score returns the best tier score for hostport, 0 when the matcher has no names, or -1.
func (hm hostMatcher) score(hostport string) int {
	if len(hm.names) == 0 {
		return 0
	}

	hostport = strings.ToLower(hostport)
	hostname := hostport
	if h, _, err := net.SplitHostPort(hostport); err == nil {
		hostname = h
	}

	hostname = strings.TrimSuffix(hostname, ".")

	best := -1
	for _, n := range hm.names {
		target := hostname
		if n.withPort {
			target = hostport
		}

		best = max(best, n.score(target))
	}

	return best
}

func (n hostName) score(target string) int {
	switch n.kind {
	case hostExact:
		if target == n.value {
			return tierExact
		}
	case hostLeading:
		if len(target) > len(n.value) && strings.HasSuffix(target, n.value) {
			return tierLeading + len(n.value)
		}
	case hostTrailing:
		if len(target) > len(n.value) && strings.HasPrefix(target, n.value) {
			return tierTrailing + len(n.value)
		}
	}

	return -1
}

// selectRule tries host tiers from best to worst and returns the first path hit, or -1.
func selectRule(hosts []hostMatcher, paths []matcher, host, path string) int {
	scores := make([]int, len(hosts))
	for i, h := range hosts {
		scores[i] = h.score(host)
	}

	tiers := slices.Compact(slices.Sorted(slices.Values(scores)))
	for _, tier := range slices.Backward(tiers) {
		if tier < 0 {
			break
		}

		var idx []int
		var ms []matcher
		for i, s := range scores {
			if s == tier {
				idx = append(idx, i)
				ms = append(ms, paths[i])
			}
		}

		if j := selectBest(ms, path); j >= 0 {
			return idx[j]
		}
	}

	return -1
}

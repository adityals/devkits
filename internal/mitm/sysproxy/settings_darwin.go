//go:build darwin

package sysproxy

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"slices"
	"strings"
	"time"
)

const (
	networksetupBin = "/usr/sbin/networksetup"
	// Changes may prompt for an admin password, so leave time to type it.
	networksetupTimeout = time.Minute
)

// networkSetup sets the HTTP and HTTPS proxy on every enabled network service.
// shortcut: Disable turns the proxy off instead of restoring prior settings, store and restore them if users already run a proxy.
type networkSetup struct{}

func NewSettings() Settings {
	return networkSetup{}
}

func (networkSetup) Enabled(ctx context.Context, host, port string) (bool, error) {
	services, err := listServices(ctx)
	if err != nil {
		return false, err
	}

	want := []string{
		"Enabled: Yes",
		"Server: " + host,
		"Port: " + port,
	}

	for _, svc := range services {
		out, err := networksetup(ctx, "-getwebproxy", svc)
		if err != nil {
			return false, err
		}

		lines := strings.Split(out, "\n")
		if !slices.ContainsFunc(want, func(w string) bool { return !slices.Contains(lines, w) }) {
			return true, nil
		}
	}

	return false, nil
}

func (networkSetup) Enable(ctx context.Context, host, port string) error {
	return eachService(ctx, func(svc string) error {
		if _, err := networksetup(ctx, "-setwebproxy", svc, host, port); err != nil {
			return err
		}

		_, err := networksetup(ctx, "-setsecurewebproxy", svc, host, port)

		return err
	})
}

func (networkSetup) Disable(ctx context.Context) error {
	return eachService(ctx, func(svc string) error {
		if _, err := networksetup(ctx, "-setwebproxystate", svc, "off"); err != nil {
			return err
		}

		_, err := networksetup(ctx, "-setsecurewebproxystate", svc, "off")

		return err
	})
}

func eachService(ctx context.Context, fn func(svc string) error) error {
	services, err := listServices(ctx)
	if err != nil {
		return err
	}

	for _, svc := range services {
		if err := fn(svc); err != nil {
			return err
		}
	}

	return nil
}

// listServices skips the header line and disabled services, which networksetup marks with "*".
func listServices(ctx context.Context) ([]string, error) {
	out, err := networksetup(ctx, "-listallnetworkservices")
	if err != nil {
		return nil, err
	}

	var services []string
	for _, line := range strings.Split(out, "\n")[1:] {
		if line != "" && !strings.HasPrefix(line, "*") {
			services = append(services, line)
		}
	}

	return services, nil
}

func networksetup(ctx context.Context, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, networksetupTimeout)
	defer cancel()

	out, err := exec.CommandContext(ctx, networksetupBin, args...).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("networksetup %s: %w: %s", args[0], err, bytes.TrimSpace(out))
	}

	return string(out), nil
}

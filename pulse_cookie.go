//go:build linux && !android

package antui

import (
	"os"
	"path/filepath"
)

// pulseCookie is the shared secret a session's programs authenticate with.
// A server that does not care about it — PipeWire does not — is given one
// anyway, because the field is not optional.
func pulseCookie() []byte {
	cookie := make([]byte, 256)
	paths := []string{os.Getenv("PULSE_COOKIE")}
	if config := os.Getenv("XDG_CONFIG_HOME"); config != "" {
		paths = append(paths, filepath.Join(config, "pulse", "cookie"))
	}
	if home, err := os.UserHomeDir(); err == nil {
		paths = append(paths,
			filepath.Join(home, ".config", "pulse", "cookie"),
			filepath.Join(home, ".pulse-cookie"))
	}
	for _, path := range paths {
		if path == "" {
			continue
		}
		if data, err := os.ReadFile(path); err == nil && len(data) >= len(cookie) {
			copy(cookie, data)
			return cookie
		}
	}
	return cookie
}

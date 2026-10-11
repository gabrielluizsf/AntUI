//go:build linux && !android

package antui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// pulseSocket is where the server is listening, in the order the ordinary
// tools look.
func pulseSocket() string {
	if server := os.Getenv("PULSE_SERVER"); server != "" {
		// A server address can name a transport and can list several; the
		// only one that matters here is a local socket.
		for part := range strings.SplitSeq(server, " ") {
			if path, ok := strings.CutPrefix(part, "unix:"); ok {
				return path
			}
			if strings.HasPrefix(part, "/") {
				return part
			}
		}
	}
	if dir := os.Getenv("XDG_RUNTIME_DIR"); dir != "" {
		return filepath.Join(dir, "pulse", "native")
	}
	return fmt.Sprintf("/run/user/%d/pulse/native", os.Getuid())
}

//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import (
	"os"
	"path/filepath"
)

// authPath is where the cookie file is: what XAUTHORITY says, or the default
// under the user's home directory.
func authPath() string {
	if path := os.Getenv("XAUTHORITY"); path != "" {
		return path
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".Xauthority")
}

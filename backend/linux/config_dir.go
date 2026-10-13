//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import "os"

// configDir is where a desktop keeps the settings it writes at login:
// XDG_CONFIG_HOME when the session named one, the home directory's .config
// when it did not, and nothing when there is neither — a session with no
// home directory has no settings file to read, and says so by having none.
func configDir() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		return ""
	}
	return dir
}

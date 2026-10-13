//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import (
	"os/exec"
	"strings"
	"time"
)

const settingsTimeout = 2 * time.Second

// askSettings runs the settings tool and returns what it printed, or nothing
// when it is not there, fails, or takes too long. The timeout mirrors ask in
// sysfont.go for the same reason: a session with a settings bus that will
// not answer can leave the call hanging, and no theme a desktop could name
// is worth a program that will not start.
func askSettings(args ...string) string {
	path := settingsHelper()
	if path == "" {
		return ""
	}
	done := make(chan string, 1)
	cmd := exec.Command(path, args...)
	go func() {
		out, err := cmd.Output()
		if err != nil {
			done <- ""
			return
		}
		done <- strings.TrimSpace(string(out))
	}()
	select {
	case answer := <-done:
		return answer
	case <-time.After(settingsTimeout):
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
		return ""
	}
}

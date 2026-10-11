package antui

import (
	"os/exec"
	"strings"
	"time"
)

// ask runs a program that answers with a path, and returns nothing at all if
// it is not there, fails, or takes too long.
//
// The timeout is not caution for its own sake: fc-match on a machine with a
// cold font cache can take seconds while it builds one, and a library that
// blocks a program's start-up on that has made things worse rather than
// better. Two seconds is far more than a warm cache needs.
func ask(name string, args ...string) string {
	path, err := exec.LookPath(name)
	if err != nil {
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
	case <-time.After(2 * time.Second):
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
		return ""
	}
}

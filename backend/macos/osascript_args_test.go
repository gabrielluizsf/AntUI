//go:build darwin

package macos

import (
	"strings"
	"testing"
)

// TestOsascriptArgsKeepTextOutOfTheScript: a title that is itself AppleScript
// stays an argument and never reaches the script's source.
func TestOsascriptArgsKeepTextOutOfTheScript(t *testing.T) {
	title := `x" & (do shell script "touch /tmp/pwned") & "`
	args := osascriptArgs(title, "body")
	if n := len(args); args[n-2] != title || args[n-1] != "body" {
		t.Fatalf("title and body should be the last two arguments, got %q", args)
	}
	for _, a := range args[:len(args)-2] {
		if strings.Contains(a, "pwned") {
			t.Fatalf("the title leaked into the script: %q", a)
		}
	}
}

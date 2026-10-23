package svg

import (
	"strings"
	"testing"
)

// TestResolveHrefKeepsADrivePathOnTheServer is the Windows half of
// CVE-2023-38633: a drive path in a drawing from a server resolves against the
// server, as any other bare path does, and never names this machine's disk.
func TestResolveHrefKeepsADrivePathOnTheServer(t *testing.T) {
	for _, href := range []string{`C:\Users\me\secret.png`, `c:/Users/me/secret.png`} {
		got := resolveHref("https://example.com/a.svg", href)
		if !strings.HasPrefix(got, "https://example.com/") {
			t.Errorf("resolveHref(server, %q) = %q, want it on the server", href, got)
		}
	}
	// A drawing on this machine, or one with no base, still reads the drive.
	for _, base := range []string{"", `D:\drawings\a.svg`} {
		if got := resolveHref(base, `C:\pic.png`); got != `C:\pic.png` {
			t.Errorf("resolveHref(%q, drive path) = %q, want it kept", base, got)
		}
	}
}

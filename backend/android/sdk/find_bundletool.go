package sdk

import (
	"os"
	"os/exec"
	"path/filepath"
)

// findBundletool looks where someone would have put it. It is never an
// error not to find it: it is not part of the SDK and nothing needs it to
// build.
func findBundletool(root string) string {
	if p := os.Getenv("BUNDLETOOL"); p != "" && exists(p) {
		return p
	}
	home, _ := os.UserHomeDir()
	for _, p := range []string{
		filepath.Join(root, "bundletool.jar"),
		filepath.Join(home, "bundletool.jar"),
		filepath.Join(home, ".local", "share", "bundletool.jar"),
		"/usr/local/lib/bundletool.jar",
	} {
		if exists(p) {
			return p
		}
	}
	if p, err := exec.LookPath("bundletool"); err == nil {
		return p
	}
	return ""
}

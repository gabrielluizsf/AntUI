package sdk

import (
	"os"
	"path/filepath"
)

// isSDK reports whether a directory looks like an SDK root rather than
// merely existing — an empty ANDROID_HOME pointing at a directory that was
// never populated is a common way to waste an afternoon.
func isSDK(dir string) bool {
	if dir == "" {
		return false
	}
	for _, sub := range []string{"platform-tools", "platforms", "build-tools", "cmdline-tools", "tools", "ndk"} {
		if st, err := os.Stat(filepath.Join(dir, sub)); err == nil && st.IsDir() {
			return true
		}
	}
	return false
}

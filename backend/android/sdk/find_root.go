package sdk

import (
	"os"
	"path/filepath"
	"runtime"
)

// findRoot looks where the environment says, then where each system puts it
// by default. ANDROID_HOME is the current name and ANDROID_SDK_ROOT the one
// before it; both are still set by real installs, so both are read.
func findRoot() string {
	for _, env := range []string{"ANDROID_HOME", "ANDROID_SDK_ROOT"} {
		if dir := os.Getenv(env); dir != "" && isSDK(dir) {
			return dir
		}
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	var candidates []string
	switch runtime.GOOS {
	case "darwin":
		candidates = []string{filepath.Join(home, "Library", "Android", "sdk")}
	case "windows":
		candidates = []string{filepath.Join(os.Getenv("LOCALAPPDATA"), "Android", "Sdk")}
	default:
		candidates = []string{
			filepath.Join(home, "Android", "Sdk"),
			filepath.Join(home, "android-sdk"),
			"/opt/android-sdk",
			"/usr/lib/android-sdk",
		}
	}
	for _, dir := range candidates {
		if isSDK(dir) {
			return dir
		}
	}
	return ""
}

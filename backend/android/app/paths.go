//go:build android

package app

import (
	"os"

	"github.com/gabrielluizsf/antui/backend/android/jni"
	"github.com/gabrielluizsf/antui/backend/android/ndk"
)

// setupPaths finds the directories the platform gives the app, and points
// the Go standard library at them.
//
// This matters more than it looks. Go's os.TempDir answers $TMPDIR or /tmp,
// and there is no /tmp on Android — so every package in the standard library
// that makes a temporary file, and every library that does, fails on a phone
// for a reason that has nothing to do with what it was doing. Setting three
// environment variables at startup fixes all of them at once.
func setupPaths(a *App) {
	cache := readCacheDir()
	if cache != "" {
		a.mu.Lock()
		a.info.Cache = cache
		a.mu.Unlock()
		// os.TempDir, and everything built on it.
		os.Setenv("TMPDIR", cache)
		os.Setenv("XDG_CACHE_HOME", cache)
	}
	files := a.Info().InternalData
	if files != "" {
		// os.UserHomeDir, os.UserConfigDir and os.UserCacheDir, which a
		// surprising amount of ordinary Go reaches for.
		os.Setenv("HOME", files)
		os.Setenv("XDG_CONFIG_HOME", files)
		os.Setenv("XDG_DATA_HOME", files)
	}
}

// readCacheDir asks the platform. It is not among the paths a native
// activity is handed — only the files, external and obb ones are — so it is
// the one directory that has to be fetched from Java.
func readCacheDir() string {
	var out string
	err := jni.Do(func(e *jni.Env) error {
		dir, err := e.Invoke(Context(), "getCacheDir",
			jni.Sig(jni.TClass("java/io/File")))
		if err != nil || dir.IsNil() {
			return err
		}
		out, err = e.InvokeString(dir, "getAbsolutePath", jni.Sig(jni.TString))
		return err
	})
	if err != nil {
		ndk.Warnf("cannot find the cache directory: %v", err)
	}
	return out
}

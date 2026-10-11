//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import (
	"os"
	"time"
)

// SystemDark is the color scheme the desktop paints its own interface in, and
// whether the desktop had an answer to give.
//
// It is asked the way a program asks without a library to link against: the
// environment the session exported, the files the desktop writes when it
// starts, and — where a desktop keeps its answer behind a tool that will
// speak to the program — that tool, spawned only when the file it reads has
// changed since the last time. What comes back stands for a second at a
// time, so no frame pays for the search, and a theme the user changes in the
// desktop's own settings still lands on the frame after it.
//
// A desktop that says nothing is reported as saying nothing: (false, false)
// is not "light", it is "the desktop did not say". See themeFrom for the
// order the sources are asked in.
//
// This is asked on every frame of a window that is up, so everything that
// costs anything — reading a file, working out which tool there is to ask —
// happens only when a reading is due, and the frame in between gets a lock
// and a comparison. TestSystemDarkCostsNothingBetweenReadings is that.
func SystemDark() (bool, bool) {
	themeMu.Lock()
	defer themeMu.Unlock()
	now := time.Now()
	if !now.Before(theme.next) {
		env := os.Getenv
		var ask settingsAsk
		if canAskSettings(env) {
			ask = askSettings
		}
		refreshTheme(now, env, configDir(), ask)
	}
	return theme.dark, theme.known
}

//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import (
	"path/filepath"
	"time"
)

// refreshTheme reads the desktop again, with themeMu held.
//
// The settings helper is the expensive source — it costs a process to spawn
// — so it runs only while there is a desktop to ask and only when the file
// the preference lives in has changed since the last time it ran: a desktop
// that has not written a new answer has not changed what it would say. Its
// answer, once it has one, stands over what the files say, because both are
// written by the same desktop and the tool is the one that reads back what
// it wrote.
func refreshTheme(now time.Time, env func(string) string, dir string, ask settingsAsk) {
	stamp := fileStampOf(filepath.Join(dir, "dconf", "user"))
	switch {
	case ask == nil:
		// There is no tool to ask, so nothing one of them said before
		// should go on standing over the files.
		theme.setting, theme.stamp = darkAnswer{}, fileStamp{}
	case stamp != theme.stamp:
		// A tool with nothing to say is recorded as having been asked
		// all the same: two processes a second, forever, would cost more
		// than any answer it is ever going to give.
		theme.setting.dark, theme.setting.ok = themeInSettings(ask)
		theme.stamp = stamp
	}
	theme.dark, theme.known = themeFrom(env, dir, theme.setting)
	theme.next = now.Add(themeTTL)
}

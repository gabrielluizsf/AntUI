//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
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

// themeTTL is how long a reading stands before the session's environment and
// the desktop's files are read again. A second is far inside the time it
// takes a person to change a theme and look at the result, and short enough
// that they see it land; reading files on every frame would spend more than
// the answer is worth on a state that changes twice a year.
const themeTTL = time.Second

// darkAnswer is one source's word: ok says whether it had any, dark what it
// said. A source that had none leaves ok false rather than claiming light,
// and the question moves on to the next one.
type darkAnswer struct {
	dark bool
	ok   bool
}

// fileStamp is a file's identity without opening it — whether it is there,
// how big it is, when it was last written — which is all it takes to notice
// that the desktop wrote a new answer into one.
type fileStamp struct {
	size int64
	mod  time.Time
	ok   bool
}

// themeState is what the desktop said the last time it was asked, and when
// it is due to be asked again.
type themeState struct {
	dark, known bool
	setting     darkAnswer // what the settings tool said, as of stamp
	stamp       fileStamp  // the file that answer came from
	next        time.Time  // when the cheap sources are read again
}

var (
	themeMu sync.Mutex
	theme   themeState
)

// themeNow is SystemDark with every input handed to it — the clock the cache
// runs on, the environment and directory to read, and the settings helper to
// ask (nil when there is none to ask) — so the caching and the precedence
// can be exercised without a desktop to take them from.
func themeNow(now time.Time, env func(string) string, dir string, ask settingsAsk) (bool, bool) {
	themeMu.Lock()
	defer themeMu.Unlock()
	if now.Before(theme.next) {
		return theme.dark, theme.known
	}
	refreshTheme(now, env, dir, ask)
	return theme.dark, theme.known
}

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

// themeFrom reads the color scheme from the sources that could know it,
// best first: the environment is set for this one program and is the most
// specific word about it, the desktop's settings are what the user chose for
// everything they run, the files a desktop writes at login are what is left
// when there is no tool to ask, and the colours a terminal exports describe
// one window rather than the desktop behind it. The first source that says
// something wins, and when none does the answer is that none did.
func themeFrom(env func(string) string, dir string, setting darkAnswer) (dark, ok bool) {
	if dark, said := themeInName(env("GTK_THEME")); said {
		return dark, true
	}
	if setting.ok {
		return setting.dark, true
	}
	if dark, said := themeInConfig(dir); said {
		return dark, true
	}
	return themeInColorFGBG(env("COLORFGBG"))
}

// themeInName reads a theme's name, which is where a desktop spells its side
// out: a theme that is dark says so in its name or in the variant after a
// colon (Adwaita:dark), and turning dark on is exactly what that marker is
// for — so a name carrying neither is a light theme. No name at all is not
// an answer.
func themeInName(value string) (dark, ok bool) {
	if strings.TrimSpace(value) == "" {
		return false, false
	}
	return strings.Contains(strings.ToLower(value), "dark"), true
}

// themeInColorFGBG reads the colour pair a terminal exports in COLORFGBG,
// "foreground;background" as palette indices. Black through cyan is a dark
// background and white or bright white is a light one; the indices in
// between are a pair someone found readable, not a statement about the
// desktop, and say nothing.
func themeInColorFGBG(value string) (dark, ok bool) {
	_, background, found := strings.Cut(value, ";")
	if !found {
		return false, false
	}
	index, err := strconv.Atoi(strings.TrimSpace(background))
	if err != nil {
		return false, false
	}
	switch {
	case index >= 0 && index <= 6:
		return true, true
	case index == 7 || index == 15:
		return false, true
	}
	return false, false
}

//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import (
	"time"
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

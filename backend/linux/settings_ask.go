//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

// settingsAsk runs a desktop's own settings tool with the given arguments
// and returns what it printed, or nothing when it has none to give.
type settingsAsk func(args ...string) string

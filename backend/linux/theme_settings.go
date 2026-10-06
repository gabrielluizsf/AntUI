//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import (
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

// settingsAsk runs a desktop's own settings tool with the given arguments
// and returns what it printed, or nothing when it has none to give.
type settingsAsk func(args ...string) string

// settingsHelper is the program a desktop's settings are read from, looked
// up once: a machine without one is not asked about it again.
var settingsHelper = sync.OnceValue(func() string {
	path, err := exec.LookPath("gsettings")
	if err != nil {
		return ""
	}
	return path
})

// canAskSettings reports whether there is a tool worth asking, and a desktop
// worth asking it about. The tool reads GNOME's settings, which is what a
// GNOME-family desktop answers through and what its desktops grew from; a
// desktop of another family keeps its answer in a file themeInConfig reads,
// and a session that names no desktop is a window manager with no theme of
// its own to report.
func canAskSettings(env func(string) string) bool {
	if settingsHelper() == "" {
		return false
	}
	names := strings.ToLower(env("XDG_CURRENT_DESKTOP") + ";" +
		env("XDG_SESSION_DESKTOP") + ";" + env("DESKTOP_SESSION"))
	for _, family := range []string{"gnome", "ubuntu", "cinnamon", "budgie", "unity", "pantheon"} {
		if strings.Contains(names, family) {
			return true
		}
	}
	return false
}

// themeInSettings asks the desktop's settings tool which scheme it is
// drawing in.
//
// The two keys are the ones GNOME grew and its desktops have kept writing
// into ever since — Cinnamon and MATE both write the theme back through
// them, which is why GTK applications pick a desktop's choice up at all. The
// scheme is asked first, because it is the answer for the newer desktops
// that have one; "default" there means the theme already is what the user
// asked for, and the theme's name is then where that desktop says which one
// that is.
func themeInSettings(ask settingsAsk) (dark, ok bool) {
	switch scheme := unquote(ask("get", "org.gnome.desktop.interface", "color-scheme")); scheme {
	case "prefer-dark":
		return true, true
	case "prefer-light":
		return false, true
	}
	return themeInName(unquote(ask("get", "org.gnome.desktop.interface", "gtk-theme")))
}

// askSettings runs the settings tool and returns what it printed, or nothing
// when it is not there, fails, or takes too long. The timeout mirrors ask in
// sysfont.go for the same reason: a session with a settings bus that will
// not answer can leave the call hanging, and no theme a desktop could name
// is worth a program that will not start.
func askSettings(args ...string) string {
	path := settingsHelper()
	if path == "" {
		return ""
	}
	done := make(chan string, 1)
	cmd := exec.Command(path, args...)
	go func() {
		out, err := cmd.Output()
		if err != nil {
			done <- ""
			return
		}
		done <- strings.TrimSpace(string(out))
	}()
	select {
	case answer := <-done:
		return answer
	case <-time.After(2 * time.Second):
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
		return ""
	}
}

// unquote takes the quotes the settings tool puts around the value it
// prints, and the whitespace it puts around that.
func unquote(value string) string {
	return strings.Trim(strings.TrimSpace(value), `'"`)
}

// fileStampOf reads a file's identity without opening it — whether it is
// there, how big it is, when it was last written — which is all it takes to
// notice that the desktop wrote a new answer into it. A file that is not
// there has a stamp of its own, so a desktop that stops writing settings
// is not asked about the ones it no longer keeps.
func fileStampOf(path string) fileStamp {
	info, err := os.Stat(path)
	if err != nil {
		return fileStamp{}
	}
	return fileStamp{size: info.Size(), mod: info.ModTime(), ok: true}
}

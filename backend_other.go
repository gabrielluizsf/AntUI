//go:build !linux && !freebsd && !openbsd && !netbsd && !dragonfly && !windows && !darwin

package antui

import (
	"errors"
	"runtime"

	"github.com/gabrielluizsf/antui/canvas"
)

// noWindows is the fallback backend while a platform has no port yet. It
// reports an error from Open and answers every question with a nothing.
// Drawing to a Canvas works everywhere, window or not, so the library builds
// and runs anywhere even before the platform files exist.
//
// A platform with a port (linux/BSD → backend/linux, windows →
// backend/windows, darwin → backend/macos), is taken off this build
// constraint by its own backend_*.go file in this package.
type noWindows struct{}

func newBackend() platform { return noWindows{} }

var errNoBackend = errors.New("antui: there is no window backend for " + runtime.GOOS +
	"; the Canvas API works, but opening a window does not")

func (noWindows) open(*Window, string, int, int) error { return errNoBackend }
func (noWindows) close()                               {}
func (noWindows) pump(*Window)                         {}
func (noWindows) present(*Window, canvas.Area)         {}
func (noWindows) setTitle(string)                      {}
func (noWindows) setFullscreen(bool) bool              { return false }
func (noWindows) displaySize() (int, int, bool)        { return 0, 0, false }
func (noWindows) displayRefresh() int                  { return 0 }

func (noWindows) setLimits(*Window, Limits)      {}
func (noWindows) setSize(*Window, int, int) bool { return false }
func (noWindows) contentScale() float64          { return 0 }

func (noWindows) clipboard() (string, []string) { return "", nil }

func (noWindows) setIcon([]*canvas.Canvas) bool { return false }

func (noWindows) setClipboard(string) bool { return false }
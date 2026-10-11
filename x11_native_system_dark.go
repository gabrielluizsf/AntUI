//go:build (linux && !android) || freebsd || openbsd || netbsd || dragonfly

package antui

import (
	x11 "github.com/gabrielluizsf/antui/backend/linux"
)

// systemDark takes the answer the desktop gives, which on X11 is the only
// answer there is: the server has no signal for the theme a desktop is
// drawing in, so it is read the way every toolkit reads it — out of the
// session's environment, the files the desktop writes at login and, where a
// desktop keeps it behind one, its own settings tool (see the linux
// package's SystemDark). A desktop with nothing to say still reports
// nothing, and [Window.SetSystemDark] stands in for it.
func (n *x11Native) systemDark() (bool, bool) { return x11.SystemDark() }

//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import (
	"os"
	"path/filepath"

	"github.com/gabrielluizsf/antui/backend/dbus"
)

// SendNotification shows a desktop notification, reporting whether the
// notification server took it. It does not go through X at all: the
// freedesktop notification server lives on the session bus, which is why it
// works the same under Wayland, and why it needs no window to be alive.
func (x *Driver) SendNotification(title, body string) bool {
	return dbus.Notify(filepath.Base(os.Args[0]), title, body) == nil
}

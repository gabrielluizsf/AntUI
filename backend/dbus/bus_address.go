package dbus

import (
	"errors"
	"os"
	"path"
)

// busAddress is the session bus's address: what DBUS_SESSION_BUS_ADDRESS
// says, or the socket systemd puts in the runtime directory when it says
// nothing.
func busAddress() (string, error) {
	if addr := os.Getenv("DBUS_SESSION_BUS_ADDRESS"); addr != "" {
		return addr, nil
	}
	if dir := os.Getenv("XDG_RUNTIME_DIR"); dir != "" {
		return "unix:path=" + path.Join(dir, "bus"), nil
	}
	return "", errors.New("dbus: DBUS_SESSION_BUS_ADDRESS is not set, so there is no session bus to talk to")
}

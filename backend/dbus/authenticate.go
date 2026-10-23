package dbus

import (
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// authenticate runs the SASL handshake the bus asks for before any message:
// a NUL byte, then EXTERNAL with our uid — the bus checks it against the
// socket's own credentials, so there is no secret to send — and BEGIN once
// the bus says OK.
func (c *conn) authenticate() error {
	uid := hex.EncodeToString([]byte(strconv.Itoa(os.Getuid())))
	if _, err := c.Write([]byte("\x00AUTH EXTERNAL " + uid + "\r\n")); err != nil {
		return fmt.Errorf("dbus: authenticating: %w", err)
	}
	line, err := c.r.ReadString('\n')
	if err != nil {
		return fmt.Errorf("dbus: authenticating: %w", err)
	}
	if !strings.HasPrefix(line, "OK ") {
		return fmt.Errorf("dbus: the bus refused to authenticate: %q", strings.TrimSpace(line))
	}
	if _, err := c.Write([]byte("BEGIN\r\n")); err != nil {
		return fmt.Errorf("dbus: authenticating: %w", err)
	}
	return nil
}

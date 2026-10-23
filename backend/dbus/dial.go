package dbus

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

// dial connects to the session bus and authenticates, leaving a connection
// that every read and write gives up on at deadline.
func dial(deadline time.Time) (*conn, error) {
	address, err := busAddress()
	if err != nil {
		return nil, err
	}
	paths := socketPaths(address)
	if len(paths) == 0 {
		return nil, fmt.Errorf("dbus: the session bus address %q names no unix socket", address)
	}
	var lastErr error
	for _, path := range paths {
		raw, err := net.DialTimeout("unix", path, time.Until(deadline))
		if err != nil {
			lastErr = err
			continue
		}
		raw.SetDeadline(deadline)
		c := &conn{Conn: raw, r: bufio.NewReader(raw)}
		if err := c.authenticate(); err != nil {
			raw.Close()
			return nil, err
		}
		return c, nil
	}
	return nil, fmt.Errorf("dbus: could not connect to the session bus: %w", lastErr)
}

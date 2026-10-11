//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

// dialDisplay opens the socket named by DISPLAY and reports which display
// number it turned out to be, which is what picks the auth cookie.
func dialDisplay() (net.Conn, int, error) {
	display := os.Getenv("DISPLAY")
	if display == "" {
		return nil, 0, errors.New("antui: DISPLAY is not set, so there is no X server to talk to")
	}
	colon := strings.LastIndex(display, ":")
	if colon < 0 {
		return nil, 0, fmt.Errorf("antui: DISPLAY %q is not a display address", display)
	}
	host := display[:colon]
	number := 0
	if tail := display[colon+1:]; tail != "" {
		// "0" and "0.1" both mean display 0; the screen after the dot is not
		// something this library picks between.
		if dot := strings.Index(tail, "."); dot >= 0 {
			tail = tail[:dot]
		}
		if n, err := strconv.Atoi(tail); err == nil {
			number = n
		}
	}

	if host == "" || host == "unix" {
		path := fmt.Sprintf("/tmp/.X11-unix/X%d", number)
		if conn, err := net.Dial("unix", path); err == nil {
			return conn, number, nil
		}
		// Some servers only listen on the Linux abstract namespace, which Go
		// spells with a leading NUL.
		if conn, err := net.Dial("unix", "\x00"+path); err == nil {
			return conn, number, nil
		}
		return nil, 0, fmt.Errorf("antui: could not connect to the X server on %s", path)
	}

	addr := net.JoinHostPort(host, strconv.Itoa(6000+number))
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, 0, fmt.Errorf("antui: could not connect to the X server at %s: %w", addr, err)
	}
	return conn, number, nil
}

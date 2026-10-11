//go:build linux && !android

package antui

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
)

// dialPulse connects to the PulseAudio daemon and identifies the client.
func dialPulse(name string, a *Audio) (*pulse, error) {
	conn, err := net.Dial("unix", pulseSocket())
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrNoAudio, err)
	}
	p := &pulse{conn: conn, audio: a}

	if _, err := p.ask(cmdAuth, func(t *tags) {
		t.u32(pulseVersion)
		t.bytes(pulseCookie())
	}); err != nil {
		conn.Close()
		return nil, err
	}

	if _, err := p.ask(cmdSetClientName, func(t *tags) {
		t.props([2]string{"application.name", name},
			[2]string{"application.process.binary", filepath.Base(os.Args[0])},
			[2]string{"application.icon_name", "applications-games"})
	}); err != nil {
		conn.Close()
		return nil, err
	}
	return p, nil
}

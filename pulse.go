//go:build linux && !android

package antui

import (
	"net"
	"time"
)

// A pulse is one playback stream and the goroutine feeding it.
type pulse struct {
	conn    net.Conn
	audio   *Audio
	channel uint32 // the stream, as the server numbers it
	serial  uint32
	delay   time.Duration

	block []float32 // handed to the game to fill
	out   []byte    // and the same thing as bytes for the socket
}

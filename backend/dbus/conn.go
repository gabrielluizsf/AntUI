package dbus

import (
	"bufio"
	"net"
)

// conn is a connection to the bus, buffered on the way in: the
// authentication lines and the messages after them are read off the same
// reader, so nothing it read ahead is lost between the two.
type conn struct {
	net.Conn
	r      *bufio.Reader
	serial uint32
}

// nextSerial numbers an outgoing message. Zero is not a serial, so the first
// one is 1.
func (c *conn) nextSerial() uint32 {
	c.serial++
	return c.serial
}

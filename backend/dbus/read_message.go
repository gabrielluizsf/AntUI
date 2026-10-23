package dbus

import (
	"encoding/binary"
	"fmt"
	"io"
)

// message is what read keeps of an incoming message: enough to match an
// answer to its call and to say why an error is one.
type message struct {
	typ         byte
	serial      uint32
	replySerial uint32
	member      string
	errorName   string
	signature   string
	order       binary.ByteOrder
	body        []byte
}

// maxMessage is the protocol's own ceiling on a message's size; anything
// claiming more is a broken stream, not a message to allocate for.
const maxMessage = 128 << 20

// readMessage reads one whole message off the connection.
func (c *conn) readMessage() (message, error) {
	var fixed [16]byte
	if _, err := io.ReadFull(c.r, fixed[:]); err != nil {
		return message{}, err
	}
	m := message{typ: fixed[1], order: binary.LittleEndian}
	switch fixed[0] {
	case 'l':
	case 'B':
		m.order = binary.BigEndian
	default:
		return message{}, fmt.Errorf("dbus: a message with endianness %q", fixed[0])
	}
	bodyLen := m.order.Uint32(fixed[4:])
	m.serial = m.order.Uint32(fixed[8:])
	fieldsLen := m.order.Uint32(fixed[12:])
	if uint64(bodyLen)+uint64(fieldsLen) > maxMessage {
		return message{}, fmt.Errorf("dbus: a message of %d bytes", uint64(bodyLen)+uint64(fieldsLen))
	}
	// The fields end padded to 8, counted from the start of the message.
	padded := (16 + int(fieldsLen) + 7) &^ 7
	rest := make([]byte, padded-16+int(bodyLen))
	if _, err := io.ReadFull(c.r, rest); err != nil {
		return message{}, err
	}
	m.parseFields(rest[:fieldsLen])
	m.body = rest[padded-16:]
	return m, nil
}

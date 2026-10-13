//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import (
	"encoding/binary"
	"errors"
	"fmt"
)

// handshake introduces the client and reads back the server's description of
// itself: the resource id range, the root window, the visual and the pixel
// format everything else depends on.
func (x *Driver) handshake(displayNumber int) error {
	auth := loadAuth(displayNumber)

	// Byte order is declared once, here, and everything after it is read and
	// written little-endian regardless of what the machine is.
	request := make([]byte, xHandshakeRequestBytes)
	request[0] = 'l'
	binary.LittleEndian.PutUint16(request[2:], 11) // protocol major
	binary.LittleEndian.PutUint16(request[4:], 0)  // protocol minor
	binary.LittleEndian.PutUint16(request[6:], uint16(len(auth.name)))
	binary.LittleEndian.PutUint16(request[8:], uint16(len(auth.data)))

	request = append(request, auth.name...)
	request = append(request, padding[:pad(len(auth.name))]...)
	request = append(request, auth.data...)
	request = append(request, padding[:pad(len(auth.data))]...)
	if !x.write(request) {
		return errors.New("antui: could not send the handshake to the X server")
	}

	if !x.need(xProtocolHeaderBytes, xReplyTimeout) {
		return errors.New("antui: the X server did not answer the handshake")
	}
	header := append([]byte(nil), x.rbuf[x.rpos:x.rpos+xProtocolHeaderBytes]...)
	x.consume(xProtocolHeaderBytes)

	bodyLen := int(binary.LittleEndian.Uint16(header[6:])) * 4
	var body []byte
	if bodyLen > 0 {
		if !x.need(bodyLen, xReplyTimeout) {
			return errors.New("antui: the X server's answer was cut short")
		}
		body = append([]byte(nil), x.rbuf[x.rpos:x.rpos+bodyLen]...)
		x.consume(bodyLen)
	}

	if header[0] != 1 {
		reason := ""
		if header[0] == 0 && len(body) > 0 {
			reason = ": " + string(body[:min(int(header[1]), len(body))])
		}
		return fmt.Errorf("antui: the X server refused the connection%s", reason)
	}
	return x.parseSetup(body)
}

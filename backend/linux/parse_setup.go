//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import (
	"encoding/binary"
	"errors"
	"fmt"
)

// parseSetup walks the success reply. Its layout is fixed by the protocol:
// a header, the vendor string, the pixel formats, then the screens.
func (x *Driver) parseSetup(body []byte) error {
	if len(body) < xSetupHeaderBytes {
		return errors.New("antui: the X server's answer was too short to read")
	}
	x.idBase = binary.LittleEndian.Uint32(body[4:])
	x.idMask = binary.LittleEndian.Uint32(body[8:])
	vendorLen := int(binary.LittleEndian.Uint16(body[16:]))
	maxRequest := binary.LittleEndian.Uint16(body[18:])
	screens := int(body[20])
	formats := int(body[21])
	x.imageByteOrder = body[22]
	x.minKeycode = body[26]
	x.maxKeycode = body[27]

	x.maxRequestBytes = max(uint32(maxRequest)*xProtocolWordBytes, xMinRequestBytes)

	at := xSetupHeaderBytes + vendorLen
	at += pad(at)
	formatsAt := at
	screenAt := at + formats*xFormatRecordBytes

	if screens == 0 || screenAt+xScreenMinimumBytes > len(body) {
		return errors.New("antui: the X server described no screen this library can draw on")
	}
	x.parseScreen(body, screenAt, formatsAt, formats)

	if x.bitsPerPixel != xPixel32Bits {
		return fmt.Errorf("antui: the X server uses %d bits per pixel; this library needs 32",
			x.bitsPerPixel)
	}
	x.idNext = 1
	x.seq = 0
	return nil
}

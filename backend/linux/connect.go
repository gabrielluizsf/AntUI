//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type xAuth struct {
	name string
	data []byte
}

// loadAuth digs the MIT-MAGIC-COOKIE-1 for this display out of ~/.Xauthority.
// A missing or unreadable file is not an error: a server on a local socket
// often needs no cookie at all, and the handshake will say so if it does.
func loadAuth(displayNumber int) xAuth {
	path := os.Getenv("XAUTHORITY")
	if path == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return xAuth{}
		}
		path = filepath.Join(home, ".Xauthority")
	}
	f, err := os.Open(path)
	if err != nil {
		return xAuth{}
	}
	defer f.Close()

	wanted := strconv.Itoa(displayNumber)
	readBlock := func() ([]byte, error) {
		var length uint16
		if err := binary.Read(f, binary.BigEndian, &length); err != nil {
			return nil, err
		}
		buf := make([]byte, length)
		_, err := io.ReadFull(f, buf)
		return buf, err
	}

	for {
		var family uint16
		if err := binary.Read(f, binary.BigEndian, &family); err != nil {
			return xAuth{}
		}
		if _, err := readBlock(); err != nil { // address
			return xAuth{}
		}
		number, err := readBlock()
		if err != nil {
			return xAuth{}
		}
		name, err := readBlock()
		if err != nil {
			return xAuth{}
		}
		data, err := readBlock()
		if err != nil {
			return xAuth{}
		}
		// An entry with no display number applies to every display.
		if string(name) == "MIT-MAGIC-COOKIE-1" &&
			(len(number) == 0 || string(number) == wanted) {
			return xAuth{name: string(name), data: data}
		}
	}
}

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

// pad is how many bytes of padding follow a field of this length, the
// protocol requiring everything to sit on a four-byte boundary.
func pad(length int) int { return (4 - length%4) % 4 }

var padding [3]byte

// handshake introduces the client and reads back the server's description of
// itself: the resource id range, the root window, the visual and the pixel
// format everything else depends on.
func (x *Driver) handshake(displayNumber int) error {
	auth := loadAuth(displayNumber)

	// Byte order is declared once, here, and everything after it is read and
	// written little-endian regardless of what the machine is.
	request := make([]byte, 12)
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

	if !x.need(8, xReplyTimeout) {
		return errors.New("antui: the X server did not answer the handshake")
	}
	header := append([]byte(nil), x.rbuf[x.rpos:x.rpos+8]...)
	x.consume(8)

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

// parseSetup walks the success reply. Its layout is fixed by the protocol:
// a header, the vendor string, the pixel formats, then the screens.
func (x *Driver) parseSetup(body []byte) error {
	if len(body) < 32 {
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

	x.maxRequestBytes = max(uint32(maxRequest)*4, 4096)

	at := 32 + vendorLen
	at += pad(at)
	formatsAt := at
	screenAt := at + formats*8

	if screens == 0 || screenAt+40 > len(body) {
		return errors.New("antui: the X server described no screen this library can draw on")
	}
	x.root = binary.LittleEndian.Uint32(body[screenAt:])
	x.screenWidth = binary.LittleEndian.Uint16(body[screenAt+20:])
	x.screenHeight = binary.LittleEndian.Uint16(body[screenAt+22:])
	x.screenWidthMM = binary.LittleEndian.Uint16(body[screenAt+24:])
	x.screenHeightMM = binary.LittleEndian.Uint16(body[screenAt+26:])
	rootVisual := binary.LittleEndian.Uint32(body[screenAt+32:])
	depth := body[screenAt+38]
	x.visual = rootVisual
	x.rootDepth = depth

	for i := range formats {
		f := body[formatsAt+i*8:]
		if f[0] == depth {
			x.bitsPerPixel = f[1]
			break
		}
	}

	// The root visual's channel masks say whether red and blue sit the other
	// way round from the order the canvas holds them in.
	var redMask, blueMask uint32
	at = screenAt + 40
	for range int(body[screenAt+39]) {
		if at+8 > len(body) {
			break
		}
		visuals := int(binary.LittleEndian.Uint16(body[at+2:]))
		at += 8
		for range visuals {
			if at+24 > len(body) {
				break
			}
			if binary.LittleEndian.Uint32(body[at:]) == rootVisual {
				redMask = binary.LittleEndian.Uint32(body[at+8:])
				blueMask = binary.LittleEndian.Uint32(body[at+16:])
			}
			at += 24
		}
	}
	x.swapRB = redMask == 0x000000FF && blueMask == 0x00FF0000
	x.samePixels = x.imageByteOrder == 0 && !x.swapRB &&
		binary.NativeEndian.Uint32([]byte{1, 0, 0, 0}) == 1

	if x.bitsPerPixel != 32 {
		return fmt.Errorf("antui: the X server uses %d bits per pixel; this library needs 32",
			x.bitsPerPixel)
	}
	x.idNext = 1
	x.seq = 0
	return nil
}
//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import (
	"encoding/binary"
	"errors"
	"net"
	"os"
	"syscall"
	"time"

	"github.com/gabrielluizsf/antui/backend"
)

// Request opcodes.
const (
	xCreateWindow      = 1
	xDestroyWindow     = 4
	xMapWindow         = 8
	xConfigureWindow   = 12
	xInternAtom        = 16
	xChangeProperty    = 18
	xGetProperty       = 20
	xSetSelectionOwner = 22
	xConvertSelection  = 24
	xSendEvent         = 25
	xCreateGC          = 55
	xFreeGC            = 60
	xPutImage          = 72
	xQueryExtension    = 98
	xGetKeyboardMap    = 101
)

// RandR, for the refresh rate. The core protocol has no idea what one is.
const (
	xRRQueryVersion  = 0
	xRRGetScreenInfo = 5
)

// Event and reply codes.
const (
	xError           = 0
	xReply           = 1
	xKeyPress        = 2
	xKeyRelease      = 3
	xButtonPress     = 4
	xButtonRelease   = 5
	xMotionNotify    = 6
	xFocusIn         = 9
	xFocusOut        = 10
	xExpose          = 12
	xMapNotify       = 19
	xConfigureNotify = 22
	xSelectionClear  = 29
	xSelectionRequest = 30
	xSelectionNotify = 31
	xClientMessage   = 33
	xMappingNotify   = 34
)

// The events the window asks the server for: key press/release, button
// press/release, pointer motion, exposure, structure notify and focus change.
const xEventMask = 0x1 | 0x2 | 0x4 | 0x8 | 0x40 | 0x8000 | 0x20000 | 0x200000

// Predefined atoms the protocol fixes, so they need no InternAtom round trip.
const (
	atomString   = 31
	atomAtom     = 4
	atomCardinal = 6
	atomWMName   = 39
	atomWMHints  = 35
	atomWMClass  = 67
)

// How long to wait for a reply while opening. A server that has not answered
// in five seconds is not going to.
const xReplyTimeout = 5 * time.Second

// Driver is one window on an X11 display. It owns the connection the way the
// window core owns the input state.
type Driver struct {
	conn net.Conn
	win  backend.Face

	idBase, idMask, idNext   uint32
	root, visual, window, gc uint32
	maxRequestBytes          uint32
	seq                      uint16

	rootDepth      byte
	imageByteOrder byte // 0 = LSB first
	bitsPerPixel   byte
	swapRB         bool // a visual with red and blue the other way round

	atomWMProtocols, atomWMDeleteWindow      uint32
	atomNetWMName, atomUTF8String            uint32
	atomNetWMState, atomNetWMStateFullscreen uint32

	// The atoms a drag and a paste need, looked up the first time one
	// happens, and the window that is dragging something over us.
	atomXdndAware, atomXdndEnter, atomXdndPosition uint32
	atomXdndStatus, atomXdndLeave, atomXdndDrop    uint32
	atomXdndFinished, atomXdndSelection            uint32
	atomXdndActionCopy, atomUriList                uint32
	atomClipboard, atomDropProperty                uint32
	atomNetWMIcon                                  uint32
	atomTargets, atomText                          uint32

	// copied is what this window has put on the clipboard. On X11 there is
	// nowhere else for it to live: the window that copied is the clipboard,
	// and every paste anywhere is a question put to it.
	copied   string
	dragFrom uint32

	// Where the window is on the screen, which a drag's coordinates are
	// given in and the pointer's are not.
	windowX, windowY int

	// Where the pointer is in the window, which a drop event reports.
	mouseX, mouseY int

	screenWidth, screenHeight     uint16
	screenWidthMM, screenHeightMM uint16 // what the server claims, often a lie
	randrOpcode                   byte
	refreshHz                     int // 0 until asked, -1 when unknown

	minKeycode, maxKeycode, keysymsPerCode byte
	keysyms                                []uint32

	deadKey   uint32 // a dead accent waiting for the next key
	scratch   []byte // the staging buffer PutImage sends from
	putBuffer []byte // and the request it is sent in, kept between frames

	// samePixels is whether a canvas's own bytes are already what the server
	// wants, which is every ordinary desktop and turns the conversion into a
	// copy. Worked out once, when the server has said what it wants.
	samePixels bool

	rbuf  []byte // bytes received and not consumed yet
	rpos  int
	alive bool

	// dead records why the connection was given up on, for diagnosis.
	dead string

	// raw is the socket underneath conn, for reading without waiting. A
	// deadline cannot express that: Go's poller checks the deadline before
	// it attempts the read, so one already in the past returns at once and
	// reads nothing even when the kernel has bytes waiting.
	raw syscall.RawConn
}

// ---------------------------------------------------------------------------
// Socket
// ---------------------------------------------------------------------------

// fill reads whatever the server has sent. A zero timeout polls: it takes
// what has already arrived and returns rather than waiting for more.
func (x *Driver) fill(timeout time.Duration) bool {
	if !x.alive {
		return false
	}
	x.makeRoom()
	if timeout <= 0 {
		return x.fillNow()
	}

	x.conn.SetReadDeadline(time.Now().Add(timeout))
	at := len(x.rbuf)
	n, err := x.conn.Read(x.rbuf[at:cap(x.rbuf)])
	x.rbuf = x.rbuf[:at+n]

	if err != nil {
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			return n > 0
		}
		// Anything else — EOF included — means the server is gone.
		x.die("blocking read: " + err.Error())
		return false
	}
	return n > 0
}

// fillNow takes whatever has already arrived, without waiting for more.
//
// It reads the socket directly rather than through a deadline, because a
// deadline in the past does not mean "do not wait" to Go's poller — it means
// "fail immediately", and the read never happens. Getting that wrong is
// silent: drawing still works, because writing is a separate path, and only
// input disappears.
func (x *Driver) fillNow() bool {
	at := len(x.rbuf)
	into := x.rbuf[at:cap(x.rbuf)]

	// Clear any deadline left over from a blocking read. Both the poller and
	// RawConn.Read honour it, so one set during the handshake and never
	// cleared makes every later read fail the moment it passes — and the
	// connection would look dead five seconds in for no reason at all.
	x.conn.SetReadDeadline(time.Time{})

	if x.raw == nil {
		// No raw access to the socket: wait for a moment instead. Slower to
		// notice input, but it does notice it.
		x.conn.SetReadDeadline(time.Now().Add(time.Millisecond))
		defer x.conn.SetReadDeadline(time.Time{})
		n, err := x.conn.Read(into)
		x.rbuf = x.rbuf[:at+n]
		if err != nil {
			var netErr net.Error
			if !errors.As(err, &netErr) || !netErr.Timeout() {
				x.die("fallback read: " + err.Error())
				return false
			}
		}
		return n > 0
	}

	var n int
	var readErr error
	controlErr := x.raw.Read(func(fd uintptr) bool {
		for {
			var err error
			n, err = syscall.Read(int(fd), into)
			if err == syscall.EINTR {
				continue
			}
			readErr = err
			return true // never wait for readability
		}
	})
	if controlErr != nil {
		// A deadline that slipped past is not a dead server, only a read
		// that did not happen.
		if errors.Is(controlErr, os.ErrDeadlineExceeded) {
			return false
		}
		x.die("raw control: " + controlErr.Error())
		return false
	}

	switch {
	case readErr == syscall.EAGAIN || readErr == syscall.EWOULDBLOCK:
		return false // nothing waiting, and nothing wrong
	case readErr != nil:
		x.die("raw read: " + readErr.Error())
		return false
	case n == 0:
		x.die("the server closed the connection")
		return false
	}
	x.rbuf = x.rbuf[:at+n]
	return true
}

// makeRoom compacts what has been consumed and leaves space to read into, so
// the buffer does not grow just because the read position walked forward.
func (x *Driver) makeRoom() {
	if x.rpos > 0 {
		x.rbuf = x.rbuf[:copy(x.rbuf, x.rbuf[x.rpos:])]
		x.rpos = 0
	}
	if at := len(x.rbuf); cap(x.rbuf)-at < 8192 {
		grown := make([]byte, at, max(cap(x.rbuf)*2, at+8192))
		copy(grown, x.rbuf)
		x.rbuf = grown
	}
}

func (x *Driver) available() int { return len(x.rbuf) - x.rpos }

// need waits until at least count bytes are buffered.
func (x *Driver) need(count int, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for x.available() < count {
		if !x.fill(time.Until(deadline)) && (!x.alive || !time.Now().Before(deadline)) {
			return false
		}
	}
	return true
}

func (x *Driver) consume(count int) { x.rpos += count }

// writeChunk bounds how much goes out between two reads. Large enough that a
// frame is a handful of writes, small enough that the server never waits long
// for us to listen.
const writeChunk = 64 << 10

// write sends a request, listening to the server as it goes.
//
// The listening is the point. A frame of pixels is over a megabyte, and while
// we push it the server is pushing events back; if neither side reads, both
// block on a full socket and the program hangs with a half-drawn window. So
// every chunk is preceded by a drain of whatever has arrived, which is what
// keeps the two of us from waiting on each other.
func (x *Driver) write(data []byte) bool {
	deadline := time.Now().Add(xReplyTimeout)

	for len(data) > 0 && x.alive {
		// Listen before pushing more at it.
		for x.fill(0) {
		}

		chunk := min(len(data), writeChunk)
		x.conn.SetWriteDeadline(time.Now().Add(xReplyTimeout))
		n, err := x.conn.Write(data[:chunk])
		data = data[n:]

		if n > 0 {
			deadline = time.Now().Add(xReplyTimeout)
		}
		if err == nil {
			continue
		}
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() && time.Now().Before(deadline) {
			// The server is not taking it yet. That is only fatal if it
			// never does.
			continue
		}
		x.die("write: " + err.Error())
		return false
	}
	return x.alive
}

// die gives up on the connection, remembering why.
func (x *Driver) die(reason string) {
	if x.alive {
		x.alive = false
		x.dead = reason
	}
}

// takePacket pulls one whole packet out of the buffer: the fixed 32 bytes,
// plus the variable tail a reply carries. It reports false when only part of
// one has arrived, leaving what it has for the next pass.
func (x *Driver) takePacket() (packet, extra []byte, ok bool) {
	if x.available() < 32 {
		return nil, nil, false
	}
	p := x.rbuf[x.rpos:]
	total := 32
	if p[0] == xReply {
		total += int(binary.LittleEndian.Uint32(p[4:])) * 4
	}
	if x.available() < total {
		return nil, nil, false
	}
	// The slices point into the read buffer, which the next fill will move,
	// so both are copied before anything else can touch it.
	packet = append([]byte(nil), x.rbuf[x.rpos:x.rpos+32]...)
	if total > 32 {
		extra = append([]byte(nil), x.rbuf[x.rpos+32:x.rpos+total]...)
	}
	x.consume(total)
	return packet, extra, true
}

func (x *Driver) newID() uint32 {
	id := x.idBase | (x.idNext & x.idMask)
	x.idNext++
	return id
}
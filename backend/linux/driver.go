//go:build linux || freebsd || openbsd || netbsd || dragonfly

// Package linux drives windows over the X11 protocol, spoken straight over
// the socket with nothing to link against.
//
//	X11 is asked for, and a compositor answers on the same protocols.
//
// The window core of AntUI lives elsewhere; this package is the platform
// half. It receives the window to drive as a [backend.Face], which is the
// half of Window a platform is allowed to touch, and drives the display with
// the exported calls below: Open starts the connection and the window,
// Pump folds whatever events are waiting into the window this frame, and
// Present ships the dirty rectangle back to the server.
//
// The wire protocol is spoken directly over the socket — there is no libX11
// here and nothing to link against — which is what lets an AntUI binary run
// on a machine that has an X server and nothing else.
package linux

import (
	"encoding/binary"
	"errors"
	"syscall"
	"unsafe"

	"github.com/gabrielluizsf/antui/backend"
	"github.com/gabrielluizsf/antui/canvas"
)

// Open connects to the display and creates the window that Options describes.
func Open(win backend.Face, opts backend.Options) (*Driver, error) {
	if win == nil {
		return nil, errors.New("antui: x11 backend with no window to drive")
	}
	d := &Driver{alive: true}
	if err := d.open(win, opts.Title, opts.Width, opts.Height); err != nil {
		return nil, err
	}
	return d, nil
}

// open connects to the server, creates the window and maps it.
func (x *Driver) open(win backend.Face, title string, width, height int) error {
	x.win = win

	conn, displayNumber, err := dialDisplay()
	if err != nil {
		return err
	}
	x.conn = conn
	// Raw access to the socket, for reading without waiting. Without it the
	// backend still works, a millisecond at a time.
	if sc, ok := conn.(interface {
		SyscallConn() (syscall.RawConn, error)
	}); ok {
		x.raw, _ = sc.SyscallConn()
	}
	if err := x.handshake(displayNumber); err != nil {
		conn.Close()
		x.alive = false
		return err
	}

	x.atomWMProtocols = x.internAtom("WM_PROTOCOLS")
	x.atomWMDeleteWindow = x.internAtom("WM_DELETE_WINDOW")
	x.atomNetWMName = x.internAtom("_NET_WM_NAME")
	x.atomUTF8String = x.internAtom("UTF8_STRING")
	x.atomNetWMState = x.internAtom("_NET_WM_STATE")
	x.atomNetWMStateFullscreen = x.internAtom("_NET_WM_STATE_FULLSCREEN")

	// Say the window will take a drag. Without this property no source
	// offers it anything, so a window that does not ask never finds out that
	// anybody was dragging.
	x.announceDrops()
	x.loadKeymap()

	// CreateWindow. The depth, visual and colormap are all inherited from the
	// parent, which is the root: this window wants exactly what the screen
	// already is.
	x.window = x.newID()
	body := make([]byte, 36)
	binary.LittleEndian.PutUint32(body[0:], x.window)
	binary.LittleEndian.PutUint32(body[4:], x.root)
	binary.LittleEndian.PutUint16(body[8:], 0)  // x
	binary.LittleEndian.PutUint16(body[10:], 0) // y
	binary.LittleEndian.PutUint16(body[12:], uint16(width))
	binary.LittleEndian.PutUint16(body[14:], uint16(height))
	binary.LittleEndian.PutUint16(body[16:], 0) // border width
	binary.LittleEndian.PutUint16(body[18:], 1) // class = InputOutput
	binary.LittleEndian.PutUint32(body[20:], 0) // visual = the parent's
	binary.LittleEndian.PutUint32(body[24:], 0x2|0x800)
	binary.LittleEndian.PutUint32(body[28:], 0xFF1F1F1F) // background pixel
	binary.LittleEndian.PutUint32(body[32:], xEventMask)
	x.send(xCreateWindow, 0, body) // detail 0 = inherit the parent's depth

	x.setTitleOn(title)
	x.changeProperty(atomWMClass, atomString, 8, []byte("antui\x00AntUI\x00"), 12)

	// WM_HINTS with the input bit: without it some window managers never hand
	// keyboard focus to the window at all.
	hints := le32(1|2, 1, 1, 0, 0, 0, 0, 0, 0) // InputHint|StateHint, input, NormalState
	x.changeProperty(atomWMHints, atomWMHints, 32, hints, 9)

	if x.atomWMProtocols != 0 && x.atomWMDeleteWindow != 0 {
		// Without this the window manager kills the connection to close the
		// window instead of asking, and the program never gets to save.
		x.changeProperty(x.atomWMProtocols, atomAtom, 32, le32(x.atomWMDeleteWindow), 1)
	}

	x.gc = x.newID()
	gcBody := make([]byte, 12)
	binary.LittleEndian.PutUint32(gcBody[0:], x.gc)
	binary.LittleEndian.PutUint32(gcBody[4:], x.window)
	x.send(xCreateGC, 0, gcBody)

	x.send(xMapWindow, 0, le32(x.window))

	if !x.alive {
		return errors.New("antui: the connection to the X server dropped while opening the window")
	}
	return nil
}

// Close destroys the window and lets go of the connection.
func (x *Driver) Close() {
	if x.conn == nil {
		return
	}
	if x.alive {
		if x.gc != 0 {
			x.send(xFreeGC, 0, le32(x.gc))
		}
		if x.window != 0 {
			x.send(xDestroyWindow, 0, le32(x.window))
		}
	}
	x.conn.Close()
	x.conn = nil
	x.alive = false
	x.keysyms = nil
	x.scratch = nil
	x.rbuf = nil
}

// Pump handles whatever is waiting on the connection this frame, folding it
// into the window through win.push.
func (x *Driver) Pump(win backend.Face) {
	// Take the whole of what has arrived first, so a frame's worth of input
	// is not spread over the next several frames.
	for x.fill(0) {
	}
	for {
		packet, _, ok := x.takePacket()
		if !ok {
			// Half a packet is buffered: try to complete it without blocking.
			if !x.fill(0) {
				break
			}
			continue
		}
		if packet[0] != xReply {
			x.handleEvent(packet)
		}
	}
	if !x.alive {
		win.SetShouldClose()
	}
}

// SetTitle names the window in the window manager's decoration.
func (x *Driver) SetTitle(title string) {
	if x.alive {
		x.setTitleOn(title)
	}
}

// SetFullscreen asks the window manager for a full-screen window, reporting
// whether the request could be made at all.
func (x *Driver) SetFullscreen(on bool) bool {
	if !x.alive || x.atomNetWMStateFullscreen == 0 {
		return false // no EWMH here
	}
	action := uint32(0)
	if on {
		action = 1
	}
	return x.wmState(action, x.atomNetWMStateFullscreen)
}

// DisplaySize is the size of the display the window is on.
func (x *Driver) DisplaySize() (w, h int, ok bool) {
	if x.screenWidth == 0 || x.screenHeight == 0 {
		return 0, 0, false
	}
	return int(x.screenWidth), int(x.screenHeight), true
}

// DisplayRefresh is how often the display repaints, in hertz. 0 when it
// cannot be known, which is any server without RandR.
func (x *Driver) DisplayRefresh() int {
	if !x.alive {
		return 0
	}
	if x.refreshHz != 0 {
		return max(x.refreshHz, 0)
	}
	x.refreshHz = -1 // asked, and unknown unless RandR says otherwise
	if x.randr() == 0 {
		return 0
	}
	seq := x.send(x.randrOpcode, xRRGetScreenInfo, le32(x.root))
	reply, _, ok := x.awaitReply(seq)
	if !ok {
		return 0
	}
	if rate := binary.LittleEndian.Uint16(reply[26:]); rate > 0 && rate < 1000 {
		x.refreshHz = int(rate)
	}
	return max(x.refreshHz, 0)
}

// Present copies the window's dirty rectangle to the drawable, in as many
// chunks as the server's maximum request size demands.
func (x *Driver) Present(win backend.Face, dirty canvas.Area) {
	if !x.alive || dirty.Width <= 0 || dirty.Height <= 0 {
		return
	}
	cv := win.Canvas()
	rowPixels := dirty.Width
	rowBytes := rowPixels * 4
	rowsPerChunk := max(int((x.maxRequestBytes-32)/uint32(rowBytes)), 1)

	if len(x.scratch) < rowBytes*rowsPerChunk {
		x.scratch = make([]byte, rowBytes*rowsPerChunk)
	}

	for y := 0; y < dirty.Height; y += rowsPerChunk {
		rows := min(dirty.Height-y, rowsPerChunk)

		for row := range rows {
			at := (dirty.Y+y+row)*cv.Stride + dirty.X
			x.packRow(x.scratch[row*rowBytes:], cv.Pixels[at:at+rowPixels])
		}

		header := make([]byte, 20)
		binary.LittleEndian.PutUint32(header[0:], x.window)
		binary.LittleEndian.PutUint32(header[4:], x.gc)
		binary.LittleEndian.PutUint16(header[8:], uint16(dirty.Width))
		binary.LittleEndian.PutUint16(header[10:], uint16(rows))
		binary.LittleEndian.PutUint16(header[12:], uint16(int16(dirty.X)))
		binary.LittleEndian.PutUint16(header[14:], uint16(int16(dirty.Y+y)))
		header[16] = 0 // left pad
		header[17] = x.rootDepth

		// PutImage carries the pixels as its body, so the request is built
		// here rather than through send: the data is already padded to four
		// bytes by being whole pixels.
		x.seq++
		total := 24 + rowBytes*rows

		// The request buffer is kept rather than made: at a megabyte and a
		// half a frame, allocating it is a garbage collection every couple of
		// seconds and a pause in the middle of somebody's game.
		if cap(x.putBuffer) < total {
			x.putBuffer = make([]byte, total)
		}
		request := x.putBuffer[:total]
		request[0] = xPutImage
		request[1] = 2 // ZPixmap
		binary.LittleEndian.PutUint16(request[2:], uint16(total/4))
		copy(request[4:], header)
		copy(request[24:], x.scratch[:rowBytes*rows])

		if !x.write(request) {
			return
		}
	}
}

// packRow writes one row of pixels in the byte order the server asked for.
//
// The ordinary case is that it asked for the order the pixels are already in
// — a little-endian server on a little-endian machine, with red where red
// goes — and then this is a memmove. It matters more than it looks: a frame
// is a third of a million pixels, and doing them one at a time through a
// function pointer was costing more per frame than drawing the world did.
//
// The slow path is still there and still right, because a server that wants
// the other order is not a server that should get the wrong picture.
func (x *Driver) packRow(dst []byte, src []canvas.Color) {
	if len(src) == 0 {
		return
	}
	if x.samePixels {
		copy(dst, unsafe.Slice((*byte)(unsafe.Pointer(&src[0])), len(src)*4))
		return
	}

	put := binary.LittleEndian.PutUint32
	if x.imageByteOrder != 0 {
		put = binary.BigEndian.PutUint32
	}
	for i, pixel := range src {
		c := uint32(pixel)
		if x.swapRB {
			c = c&0xFF00FF00 | (c&0xFF)<<16 | (c>>16)&0xFF
		}
		put(dst[i*4:], c)
	}
}
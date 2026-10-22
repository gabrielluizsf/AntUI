//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import (
	"net"
	"syscall"
	"time"

	"github.com/gabrielluizsf/antui/backend"
)

// Driver is one window on an X11 display. It owns the connection the way the
// window core owns the input state.
type Driver struct {
	conn net.Conn
	win  backend.Face

	idBase, idMask, idNext   uint32
	root, visual, window, gc uint32
	maxRequestBytes          uint32
	seq                      uint16

	// cursor is the custom cursor [Driver.SetCursor] last gave the window,
	// freed the moment it is replaced rather than left for the server to
	// notice the window no longer wants it.
	cursor uint32

	rootDepth      byte
	imageByteOrder byte // 0 = LSB first
	bitsPerPixel   byte
	swapRB         bool // a visual with red and blue the other way round

	// bitmapBitOrder and bitmapScanlinePad are how a 1-bit-deep image — the
	// source and mask [Driver.SetCursor] sends a cursor's shape as — has to
	// be packed: which end of each byte is the leftmost pixel, and how many
	// bits a row is padded out to. Nothing else on this connection carries
	// a bitmap, which is why nothing else reads them.
	bitmapBitOrder    byte
	bitmapScanlinePad byte

	atomWMProtocols, atomWMDeleteWindow      uint32
	atomNetWMName, atomUTF8String            uint32
	atomNetWMState, atomNetWMStateFullscreen uint32
	atomNetWMOpacity                         uint32

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

	// scale is the display's scale, read once and kept: Xft.dpi is a round
	// trip to the server and a frame is no place to spend one. scaleKnown
	// says whether the display said, and scaleDirty says the desktop wrote
	// the resource database since it was read, which is the one thing that
	// makes the next read ask the server again. See [Driver.DisplayScale].
	scale      float64
	scaleKnown bool
	scaleDirty bool

	minKeycode, maxKeycode, keysymsPerCode byte
	keysyms                                []uint32

	deadKey   uint32 // a dead accent waiting for the next key
	scratch   []byte // the staging buffer PutImage sends from
	putBuffer []byte // and the request it is sent in, kept between frames
	header    [20]byte

	// samePixels is whether a canvas's own bytes are already what the server
	// wants, which is every ordinary desktop and turns the conversion into a
	// copy. Worked out once, when the server has said what it wants.
	samePixels bool

	// due is when the program's next frame is wanted, which the window core
	// says; an idle Pump waits until then. See [Driver.SetFrameDeadline].
	due time.Time

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

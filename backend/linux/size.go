//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import (
	"encoding/binary"
	"strconv"
	"strings"

	"github.com/gabrielluizsf/antui/backend"
)

// The size hints, and how much they are worth.
//
// ICCCM gives a client one way to say how big its window may be:
// WM_NORMAL_HINTS, a property the window manager is asked to read. Asked. A
// reparenting window manager honours it, a tiling one does not read it at
// all, and neither is wrong — i3 exists to decide the layout itself, and a
// program insisting otherwise is the thing that is out of place. So the hints
// go on the window and nothing here believes them afterwards: the size that
// arrives in ConfigureNotify is the size the window is.

// WM_SIZE_HINTS flag bits, from ICCCM. Only four of the nine are set here:
// the window's position is the window manager's business.
const (
	xSizeHintMinSize  = 16
	xSizeHintMaxSize  = 32
	xSizeHintAspect   = 128
	xSizeHintBaseSize = 256
	xWMNormalHints    = 40 // XA_WM_NORMAL_HINTS
	xWMSizeHints      = 41 // XA_WM_SIZE_HINTS
	xResourceManager  = 23 // XA_RESOURCE_MANAGER, where Xft.dpi lives
)

// SetLimits passes the size constraints on to the window manager as
// WM_NORMAL_HINTS. A window with no bounds at all still gets the property
// written, because leaving a stale one behind would lock a window that had
// just been unlocked.
func (x *Driver) SetLimits(limits backend.Limits) {
	if !x.alive || x.window == 0 {
		return
	}

	// Eighteen 32-bit fields: flags, four obsolete position and size fields
	// no one has read since 1989, then the sizes, the increments, the two
	// aspect ratios, the base size and the gravity.
	hints := make([]uint32, 18)
	if limits.MinWidth > 0 || limits.MinHeight > 0 {
		hints[0] |= xSizeHintMinSize
		hints[5] = uint32(max(limits.MinWidth, 1))
		hints[6] = uint32(max(limits.MinHeight, 1))
	}
	if limits.MaxWidth > 0 || limits.MaxHeight > 0 {
		hints[0] |= xSizeHintMaxSize
		// A maximum on one axis only still has to name both, so the axis
		// with no maximum is given one no display will reach.
		hints[7] = uint32(orElse(limits.MaxWidth, 1<<15))
		hints[8] = uint32(orElse(limits.MaxHeight, 1<<15))
	}
	if limits.Aspect > 0 {
		// The ratio goes as a fraction of two integers, and both ends are
		// set to it: the window manager keeps the shape between them, so an
		// equal pair is how an exact ratio is said.
		num, den := ratio(limits.Aspect)
		hints[0] |= xSizeHintAspect
		hints[11], hints[12] = num, den // smallest allowed
		hints[13], hints[14] = num, den // largest allowed
	}
	if limits.MinWidth > 0 && limits.MinHeight > 0 {
		// Without a base size the increments and the aspect are measured
		// from zero, which is not what a minimum meant.
		hints[0] |= xSizeHintBaseSize
		hints[15] = uint32(limits.MinWidth)
		hints[16] = uint32(limits.MinHeight)
	}

	x.changeProperty(xWMNormalHints, xWMSizeHints, 32, le32(hints...), 18)
}

// orElse is v, or the fallback when v says nothing.
func orElse(v, fallback int) int {
	if v > 0 {
		return v
	}
	return fallback
}

// ratio turns a width-to-height ratio into the fraction X11 wants. Two
// decimal places of denominator is enough for every ratio a game uses — 16:9,
// 16:10, 4:3, 21:9 — and asking for more would only invite a rounding that
// makes the window drift by a pixel each drag.
func ratio(aspect float64) (num, den uint32) {
	const scale = 1000
	n := uint32(aspect*scale + 0.5)
	if n == 0 {
		return 1, 1
	}
	return n, scale
}

// SetSize asks the window manager for a new drawable size, reporting
// whether the request could be made. The window manager may refuse, and the
// canvas follows the ConfigureNotify rather than this.
func (x *Driver) SetSize(win backend.Face, width, height int) bool {
	if !x.alive || x.window == 0 {
		return false
	}
	// window, then a mask of which values follow, then those values: the
	// width and the height and nothing else, so the window manager keeps
	// deciding where the window goes.
	body := make([]byte, 16)
	binary.LittleEndian.PutUint32(body[0:], x.window)
	binary.LittleEndian.PutUint16(body[4:], 0x4|0x8)
	binary.LittleEndian.PutUint32(body[8:], uint32(max(width, 1)))
	binary.LittleEndian.PutUint32(body[12:], uint32(max(height, 1)))
	x.send(xConfigureWindow, 0, body)
	return true
}

// ContentScale is how much the display is scaled by, or 0 when it cannot be
// known.
//
// X11 has no answer to this, which is the honest summary — there is a
// physical size in the screen record and it is very often a lie, since a
// server with no EDID makes one up that lands on exactly 96 dpi. What there
// is instead is a convention: every toolkit reads Xft.dpi out of the X
// resource database, and that is what the user's desktop settings write. So
// that is read first and the physical size is the fallback, believed only
// when it lands somewhere a real display could be.
func (x *Driver) ContentScale() float64 {
	if !x.alive {
		return 0
	}
	if dpi := x.xftDPI(); dpi > 0 {
		return dpi / 96
	}
	if x.screenWidth > 0 && x.screenWidthMM > 0 {
		dpi := float64(x.screenWidth) * 25.4 / float64(x.screenWidthMM)
		// Under 96 is a television or a made-up number, and over 400 is a
		// screen size the server got wrong by a factor. Neither is a scale
		// anyone chose, and guessing wrong here makes every window the wrong
		// size — so an unbelievable answer is no answer.
		if dpi >= 96 && dpi <= 400 {
			return dpi / 96
		}
	}
	return 0
}

// xftDPI reads Xft.dpi from the resource database on the root window, which
// is where a desktop's display-scaling setting ends up.
func (x *Driver) xftDPI() float64 {
	value, ok := x.property(x.root, xResourceManager, atomString, 16384)
	if !ok {
		return 0
	}
	for line := range strings.Lines(string(value)) {
		name, setting, found := strings.Cut(line, ":")
		if !found || strings.TrimSpace(name) != "Xft.dpi" {
			continue
		}
		dpi, err := strconv.ParseFloat(strings.TrimSpace(setting), 64)
		if err != nil || dpi < 48 || dpi > 480 {
			return 0
		}
		return dpi
	}
	return 0
}

// property reads a property off a window. maxBytes bounds what is asked for
// in one go: a property longer than that is read short rather than in
// several round trips, which is right for everything here — the resource
// database is the only one read at all, and 16 KB of it is far more than the
// one line this wants.
func (x *Driver) property(window, name, typ uint32, maxBytes int) ([]byte, bool) {
	if window == 0 || name == 0 {
		return nil, false
	}
	body := make([]byte, 20)
	binary.LittleEndian.PutUint32(body[0:], window)
	binary.LittleEndian.PutUint32(body[4:], name)
	binary.LittleEndian.PutUint32(body[8:], typ)
	binary.LittleEndian.PutUint32(body[12:], 0) // offset, in 32-bit words
	binary.LittleEndian.PutUint32(body[16:], uint32(maxBytes/4))

	seq := x.send(xGetProperty, 0, body) // detail 0 = do not delete it
	reply, extra, ok := x.awaitReply(seq)
	if !ok || binary.LittleEndian.Uint32(reply[8:]) == 0 {
		return nil, false // no such property, or not of the type asked for
	}
	length := int(binary.LittleEndian.Uint32(reply[16:])) * int(reply[1]) / 8
	if length <= 0 || length > len(extra) {
		return nil, false
	}
	return extra[:length], true
}
//go:build windows

package windows

import (
	"math"
	"unsafe"

	"github.com/gabrielluizsf/antui/backend"
)

// The size lock on Win32, which is the platform that actually enforces one.
//
// Three separate things have to agree or the lock leaks. The style bits
// decide whether the frame can be dragged at all and whether the maximise
// button is there; WM_GETMINMAXINFO answers the bounds while it is dragged
// and when it is maximised; and WM_SIZING is the only place an aspect ratio
// can be kept, because it is the one message that can change the rectangle
// before it is applied.

const (
	wsThickFrame  = 0x00040000 // the frame is a resize handle
	wsMaximizeBox = 0x00010000

	wmSizing = 0x0214

	// Which edge of the rectangle is being dragged, from WM_SIZING's wparam.
	wmszLeft        = 1
	wmszRight       = 2
	wmszTop         = 3
	wmszBottom      = 4
	wmszTopLeft     = 5
	wmszTopRight    = 6
	wmszBottomLeft  = 7
	wmszBottomRight = 8

	swpNoZOrder       = 0x0004
	swpNoActivate     = 0x0010
	logPixelsX        = 88
	defaultWindowsDPI = 96
)

// SetLimits keeps the limits for the two messages that answer with them, and
// rewrites the style bits, which is what the frame and the maximise button
// are made of.
func (d *Driver) SetLimits(limits backend.Limits) {
	d.limits = limits
	if d.hwnd == 0 || d.fullscreen {
		// A full screen is a popup with no frame at all, and putting the
		// windowed style bits back now would draw one over the game.
		return
	}

	style, _, _ := procGetWindowLongPtr.Call(d.hwnd, gwlStyle)
	wanted := style
	if limits.Fixed {
		wanted &^= uintptr(wsThickFrame | wsMaximizeBox)
	} else {
		wanted |= wsThickFrame
		if limits.NoMaximize {
			wanted &^= uintptr(wsMaximizeBox)
		} else {
			wanted |= wsMaximizeBox
		}
	}
	if wanted == style {
		return
	}
	procSetWindowLongPtr.Call(d.hwnd, gwlStyle, wanted)
	// Windows caches the frame it drew; without SWP_FRAMECHANGED the button
	// stays on screen until something else happens to repaint it.
	procSetWindowPos.Call(d.hwnd, hwndTop, 0, 0, 0, 0,
		swpNoMove|swpNoZOrder|swpNoActivate|swpFrameChanged)
}

// frame is how much wider and taller the whole window is than its drawable
// area, which every one of these calls is in terms of and none of them says.
func (d *Driver) frame() (dx, dy int32) {
	style, _, _ := procGetWindowLongPtr.Call(d.hwnd, gwlStyle)
	r := rect{0, 0, 100, 100}
	procAdjustWindowRect.Call(uintptr(unsafe.Pointer(&r)), style, 0)
	return (r.Right - r.Left) - 100, (r.Bottom - r.Top) - 100
}

// SetSize asks the platform for a new drawable size, reporting whether the
// request could be made.
func (d *Driver) SetSize(win backend.Face, width, height int) bool {
	if d.hwnd == 0 {
		return false
	}
	d.windowedWidth, d.windowedHeight = width, height
	dx, dy := d.frame()
	procSetWindowPos.Call(d.hwnd, hwndTop, 0, 0,
		uintptr(int32(width)+dx), uintptr(int32(height)+dy),
		swpNoMove|swpNoZOrder|swpNoActivate)
	return true
}

// minMax answers WM_GETMINMAXINFO. The numbers it wants are the whole
// window's, frame included, so the frame is added to every one of them.
func (d *Driver) minMax(mmi *minMaxInfo) {
	dx, dy := d.frame()
	limits := d.limits
	if d.win != nil {
		limits = d.win.Bounds()
	}

	// A floor under the minimum whatever was asked for: a window that cannot
	// show its own title bar is one nobody can move or close, and the
	// interface stops fitting long before that.
	minW, minH := int32(orElse(limits.MinWidth, 120)), int32(orElse(limits.MinHeight, 80))
	mmi.MinTrackSize = point{minW + dx, minH + dy}

	if limits.MaxWidth > 0 || limits.MaxHeight > 0 {
		maxW := int32(orElse(limits.MaxWidth, 1<<15))
		maxH := int32(orElse(limits.MaxHeight, 1<<15))
		mmi.MaxTrackSize = point{maxW + dx, maxH + dy}
		// MaxSize is what maximising gives, and leaving it alone would let
		// the maximise button walk straight past a limit the drag obeys.
		mmi.MaxSize = point{maxW + dx, maxH + dy}
	}
}

// sizing keeps the aspect ratio while the window is dragged. lparam points at
// the rectangle Windows is about to apply, and writing through it before
// returning is the whole mechanism.
func (d *Driver) sizing(edge uintptr, r *rect) {
	limits := d.limits
	if d.win != nil {
		limits = d.win.Bounds()
	}
	if limits.Aspect <= 0 {
		return
	}

	dx, dy := d.frame()
	width := float64(r.Right-r.Left-dx) / limits.Aspect
	height := float64(r.Bottom - r.Top - dy)

	// The axis that moves is the one the other is derived from: dragging the
	// top or the bottom sets the height and the width follows, and every
	// other edge — including a corner, where a person is thinking in width —
	// sets the width.
	switch edge {
	case wmszTop, wmszBottom:
		wanted := int32(math.Round(height*limits.Aspect)) + dx
		r.Right = r.Left + wanted
	default:
		wanted := int32(math.Round(width)) + dy
		if edge == wmszTopLeft || edge == wmszTopRight {
			r.Top = r.Bottom - wanted
		} else {
			r.Bottom = r.Top + wanted
		}
	}
}

// ContentScale is pixels per point, or 0 when the system does not say.
//
// GetDpiForWindow is the per-monitor answer and is right on a machine with
// two displays at different scales, which is now ordinary. It arrived in
// Windows 10 1607, so a machine older than that falls back to the one number
// the whole desktop shares.
func (d *Driver) ContentScale() float64 {
	if d.hwnd != 0 && procGetDpiForWindow.Find() == nil {
		if dpi, _, _ := procGetDpiForWindow.Call(d.hwnd); dpi > 0 {
			return float64(dpi) / defaultWindowsDPI
		}
	}
	hdc, _, _ := procGetDC.Call(0)
	if hdc == 0 {
		return 0
	}
	defer procReleaseDC.Call(0, hdc)
	dpi, _, _ := procGetDeviceCaps.Call(hdc, logPixelsX)
	if dpi == 0 {
		return 0
	}
	return float64(dpi) / defaultWindowsDPI
}

// orElse is v, or the fallback when v says nothing.
func orElse(v, fallback int) int {
	if v > 0 {
		return v
	}
	return fallback
}
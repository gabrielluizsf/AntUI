package backend

import "math"

// Limits are what a window may be resized to. The zero value is a window
// with no bounds at all, which is what a window is without them — so
// Limits{MinWidth: 640, MinHeight: 480} says only what it looks like it says
// and leaves everything else alone.
type Limits struct {
	// The smallest and largest drawable size, in pixels. Zero means no bound
	// on that side. The common case is a window resizable but not below the
	// point where the interface stops fitting, which is a minimum and nothing
	// else.
	//
	// A maximum below a minimum is taken at its word and the maximum wins: a
	// maximum is what a display or a lock imposes, and a minimum is only what
	// the interface would like.
	MinWidth, MinHeight int
	MaxWidth, MaxHeight int

	// Aspect is the width-to-height ratio the window keeps while it is
	// dragged, or zero to keep none. 16.0/9.0 for a widescreen composition.
	//
	// It cannot always be kept — every full screen on a display of another
	// shape is a case where it cannot — and the answer there is bars at the
	// sides rather than a stretched picture. Stretching is the failure people
	// notice and cannot name.
	Aspect float64

	// Fixed forbids resizing outright: the window is both the smallest and
	// the largest it may be. Pixel art at a fixed multiple, a layout composed
	// for one shape, a game meant to look identical on every machine — all of
	// them want this, and a maximise button that silently breaks the
	// composition is a bug the player finds first.
	Fixed bool

	// NoMaximize takes the maximise button away from a window that can still
	// be dragged to a new size.
	//
	// On X11 there is no way to ask for that on its own. A client cannot set
	// _NET_WM_ALLOWED_ACTIONS — the window manager owns it — and the only
	// lever ICCCM gives is an equal minimum and maximum, which is Fixed. So
	// this is honoured on Win32 and Cocoa and quietly does nothing on X11.
	NoMaximize bool
}

// Resizable reports whether the window may be resized at all.
func (l Limits) Resizable() bool { return !l.Fixed }

// bound holds v between low and high, where zero on either side means no
// bound there. High is applied last, so a maximum below the minimum wins.
func bound(v, low, high int) int {
	if low > 0 && v < low {
		v = low
	}
	if high > 0 && v > high {
		v = high
	}
	return v
}

// orElse is v, or the fallback when v says nothing.
func orElse(v, fallback int) int {
	if v > 0 {
		return v
	}
	return fallback
}

// clamp reduces a size to what the limits allow. It is the whole of the
// arithmetic in this package and is deliberately free of any window: it is
// what the tests can check on every platform, including the ones with no
// window system at all.
func (l Limits) Clamp(width, height int) (int, int) {
	width, height = max(width, 1), max(height, 1)
	width = bound(width, l.MinWidth, l.MaxWidth)
	height = bound(height, l.MinHeight, l.MaxHeight)

	if l.Aspect > 0 {
		// The height follows the width, because a window is nearly always
		// dragged by a side or a corner and the width is the axis a person
		// thinks in. When the height that follows is out of bounds the width
		// follows the height instead, so a ratio can never push a window past
		// a limit it was also given — the limits are the harder promise.
		wanted := max(int(math.Round(float64(width)/l.Aspect)), 1)
		if held := bound(wanted, l.MinHeight, l.MaxHeight); held == wanted {
			height = wanted
		} else {
			height = held
			width = bound(max(int(math.Round(float64(held)*l.Aspect)), 1),
				l.MinWidth, l.MaxWidth)
		}
	}
	return width, height
}

// scaled is the limits with every length in points turned into pixels.
func (l Limits) Scaled(scale float64) Limits {
	l.MinWidth, l.MinHeight = scaleTo(l.MinWidth, scale), scaleTo(l.MinHeight, scale)
	l.MaxWidth, l.MaxHeight = scaleTo(l.MaxWidth, scale), scaleTo(l.MaxHeight, scale)
	return l
}

// scaleTo multiplies a length by the display scale, for a size given in
// points. A zero or absurd scale is left alone rather than believed.
func scaleTo(v int, scale float64) int {
	if v <= 0 || scale <= 0 || scale == 1 {
		return v
	}
	return max(int(math.Round(float64(v)*scale)), 1)
}

// Options is everything a window can be opened with. The zero value plus a
// size is exactly what Open gives: a resizable window with no bounds.
type Options struct {
	Title         string
	Width, Height int
	Limits        Limits

	// Points asks for Width and Height in points rather than pixels. A window
	// asked for at 1280 by 720 opens 2560 by 1440 pixels on a display at 200
	// per cent and appears the same size on both.
	//
	// Without it "fixed size" means a different thing on every machine: a
	// locked 1280 by 720 window is a postage stamp on a 4K laptop. The limits
	// are scaled with it, since a minimum in points is what was meant.
	Points bool

	// Fullscreen opens straight into a full screen, without the frame of a
	// windowed frame appearing first.
	Fullscreen bool
}
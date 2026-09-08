package antui

import (
	"math"

	"github.com/gabrielluizsf/antui/backend"
)

// A window that stays the size it was told to.
//
// Two things have to be said plainly, because both are ways this goes wrong.
//
// **A size lock is a request, not a guarantee.** On X11 the size hints are
// advice: a tiling window manager reads them and resizes the window anyway,
// and i3 does exactly that. Win32 and Cocoa enforce them, so the lock is well
// worth asking for — but everything drawn has to go on following the
// framebuffer rather than the size the window started at, because on the
// machines where the lock does not hold, a renderer that only worked because
// nothing ever changed is wrong. Locking a window must never become the
// reason resizing works, which is why nothing in here touches the resize
// path: a size that arrives anyway is accepted, the canvas follows it, and
// Width and Height report what really happened.
//
// **Resolution is two numbers, not one.** There is a size the game draws at
// and a size the window is, and they are separate settings that happen to
// default to the same value. This file is about the second one; the engine
// above it owns the first.

// Limits are what a window may be resized to. The zero value is a window
// with no bounds at all, which is what a window is without them — so
// Limits{MinWidth: 640, MinHeight: 480} says only what it looks like it says
// and leaves everything else alone.
type Limits = backend.Limits

// Options is everything a window can be opened with. The zero value plus a
// size is exactly what Open gives: a resizable window with no bounds.
type Options = backend.Options

// Open creates and shows a window: resizable, unbounded, its size in pixels.
func Open(title string, width, height int) (*Window, error) {
	return OpenWith(Options{Title: title, Width: width, Height: height})
}

// scaleTo multiplies a length by the display scale, for a size given in
// points. A zero or absurd scale is left alone rather than believed.
func scaleTo(v int, scale float64) int {
	if v <= 0 || scale <= 0 || scale == 1 {
		return v
	}
	return max(int(math.Round(float64(v)*scale)), 1)
}

// fitDisplay reduces a size to what the display can actually show.
//
// Asking for a window larger than the screen is not an error and must not be
// one — a 1920 by 1080 game opened on a 1366 by 768 laptop is an ordinary
// Tuesday. The request is reduced to what fits and the game is told what it
// got, by Width and Height reporting the truth.
func fitDisplay(width, height, displayWidth, displayHeight int) (int, int) {
	if displayWidth > 0 && width > displayWidth {
		width = displayWidth
	}
	if displayHeight > 0 && height > displayHeight {
		height = displayHeight
	}
	return width, height
}

// orElse is v, or the fallback when v says nothing.
func orElse(v, fallback int) int {
	if v > 0 {
		return v
	}
	return fallback
}

// bounds resolves the limits against the size the window is now: a fixed
// window's smallest and largest size are both the size it has. Keeping Fixed
// as a flag rather than writing the numbers in is what lets a window locked
// at no particular size follow a display change instead of being frozen at
// whatever it happened to be when the game asked.
func (win *Window) Bounds() Limits {
	l := win.limits
	if l.Fixed {
		// A lock that names no size locks the window at whatever it is now.
		// One that names a size is that size, and the maximum is the half
		// believed when the two disagree — the same rule clamp follows.
		width := orElse(l.MaxWidth, orElse(l.MinWidth, win.width))
		height := orElse(l.MaxHeight, orElse(l.MinHeight, win.height))
		l.MinWidth, l.MaxWidth = width, width
		l.MinHeight, l.MaxHeight = height, height
	}
	return l
}

// Limits returns the constraints the window was given.
func (win *Window) Limits() Limits {
	if win == nil {
		return Limits{}
	}
	return win.limits
}

// SetLimits changes what the window may be resized to, and takes effect at
// once: a window already outside the new bounds is asked to come back inside
// them.
func (win *Window) SetLimits(limits Limits) {
	if win == nil {
		return
	}
	win.limits = limits
	if win.native != nil {
		win.native.setLimits(win, win.Bounds())
	}
	if width, height := limits.Clamp(win.width, win.height); width != win.width ||
		height != win.height {
		win.SetSize(width, height)
	}
}

// SetFixedSize locks the window to one size: both the smallest and the
// largest it may be, with the maximise button gone.
func (win *Window) SetFixedSize(width, height int) {
	if win == nil || width < 1 || height < 1 {
		return
	}
	limits := win.limits
	limits.Fixed, limits.NoMaximize = true, true
	limits.MinWidth, limits.MaxWidth = width, width
	limits.MinHeight, limits.MaxHeight = height, height
	win.SetLimits(limits)
}

// Resizable reports whether the window may be resized. It says what was
// asked for, not what the window manager will do about it.
func (win *Window) Resizable() bool { return win != nil && !win.limits.Fixed }

// SetSize asks for a new drawable size, reduced to the limits and to what the
// display can show. It reports whether the request could be made at all, not
// whether the window manager honoured it.
//
// The canvas is not resized here. It is resized when the size actually
// arrives, which is what keeps everything drawn following the framebuffer
// rather than following a request that may have been ignored.
func (win *Window) SetSize(width, height int) bool {
	if win == nil || win.native == nil || width < 1 || height < 1 {
		return false
	}
	if dw, dh, ok := win.native.displaySize(); ok {
		width, height = fitDisplay(width, height, dw, dh)
	}
	width, height = win.limits.Clamp(width, height)
	if width == win.width && height == win.height {
		return true
	}
	return win.native.setSize(win, width, height)
}

// ContentScale is how many pixels the window draws for each point the system
// lays it out in: 1 on an ordinary display, 2 on a Mac retina or a Windows
// laptop at 200 per cent, 1.5 at 150.
//
// It is what makes "fixed size" mean the same thing on every machine. A
// system that does not say answers 1 rather than 0, so arithmetic on it is
// always safe.
func (win *Window) ContentScale() float64 {
	if win == nil || win.scale <= 0 {
		return 1
	}
	return win.scale
}

// SizeInPoints is the size the window appears to be, as against the pixels it
// draws — the two are the same number until the display is scaled.
func (win *Window) SizeInPoints() (width, height int) {
	scale := win.ContentScale()
	if win == nil {
		return 0, 0
	}
	return max(int(math.Round(float64(win.Width())/scale)), 1),
		max(int(math.Round(float64(win.Height())/scale)), 1)
}
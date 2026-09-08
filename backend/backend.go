// Package backend carries everything a platform backend and the library
// agree on, without either of them importing the other.
//
// A platform backend (backend/linux, backend/windows, backend/macos) opens
// windows and drives a display. The library's Window
// owns the drawing surface and the input state, and the two halves meet in
// this package:
//
//   - The types a backend produces events with ([Event], [Key], [Mod],
//     [MouseButton], [Tool]) and the window a backend is given to drive
//     ([Face]).
//   - What a window can be opened with ([Options]) and the size it is
//     constrained to ([Limits]).
//
// This package imports nothing but math and the canvas package, so a
// backend never has to import the library — the library imports this
// package, and the backends sit between them. Without it the two halves
// refer to each other by name and Go notices the cycle.
package backend

import "github.com/gabrielluizsf/antui/canvas"

// Face is the window a backend is given to drive: the half of Window that a
// platform implementation is allowed to touch.
//
// The methods are exported, because an interface with an unexported method
// can only be implemented in the package that defines it and a platform
// backend lives in its own folder. They exist for the backends, not for
// programs, whose Window already has every one of them with a better name.
// The library's Window satisfies this interface, and a backend never has to
// name the window's type — which is what keeps a platform package free of an
// import cycle.
type Face interface {
	// Push adds one event to the window's queue, from the backend's own
	// goroutine.
	Push(Event)
	// PushSimple pushes an event with no fields other than its kind.
	PushSimple(EventType)
	// SetMouse relocates the pointer position the window believes in.
	SetMouse(x, y int)
	// ResizeCanvas grows the window's drawing surface to match a resize,
	// and reports whether it managed.
	ResizeCanvas(width, height int) bool
	// SetTouchFirst says there is no other kind of input on this platform.
	SetTouchFirst()
	// SetSafeArea reports the part of the canvas the system is not covering.
	SetSafeArea(a canvas.Area)
	// Canvas is the window's current drawable.
	Canvas() *canvas.Canvas
	// SetShouldClose asks the window to begin closing.
	SetShouldClose()
	// Bounds is the window's size constraints, resolved against its current
	// size.
	Bounds() Limits
}

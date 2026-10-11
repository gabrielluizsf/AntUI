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

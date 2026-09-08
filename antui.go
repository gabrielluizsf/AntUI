// Package antui builds desktop apps the simple way, with no dependency to
// download and nothing to link against beyond what the system already has.
//
//	Windows   →  Win32 (user32 + gdi32), through syscall
//	macOS     →  Cocoa through the Objective-C runtime
//	Linux/BSD →  the X11 protocol spoken straight over the socket
//
// Everything is drawn in software into a Canvas, which is a plain slice of
// pixels — so drawing works with no window at all, and can be tested in a
// terminal.
package antui
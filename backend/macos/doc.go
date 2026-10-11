//go:build darwin

// Package macos drives windows on Cocoa (AppKit) through the Objective-C runtime.
//
// The window core of AntUI lives elsewhere; this package is the platform
// half. It receives the window to drive as a [backend.Face], which is the
// half of Window a platform is allowed to touch, and drives the display with
// the exported calls below.
//
// macOS is the one platform that needs a C toolchain: AppKit is
// Objective-C, and there is no socket or DLL to talk to it through. The
// Cocoa side lives in backend_darwin.m; driver.go is the translation between
// its event struct and the package's own.
package macos

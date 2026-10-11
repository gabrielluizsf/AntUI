//go:build (linux && !android) || freebsd || openbsd || netbsd || dragonfly

package antui

// newBackend is what OpenWith asks for when there is not yet a window.
func newBackend() platform { return &x11Native{} }

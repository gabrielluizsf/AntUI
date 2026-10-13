//go:build android

package ndk

import "unsafe"

// Buffer is a locked window buffer: memory owned by the compositor that the
// app may write to until it posts.
type Buffer struct {
	Width, Height int
	// Stride is the distance between rows **in pixels, not bytes**, and it
	// is not the width — the compositor rounds it up, often to 16 or 32.
	// Ignoring it shears every frame diagonally.
	Stride int
	Format Format
	Bits   unsafe.Pointer
}

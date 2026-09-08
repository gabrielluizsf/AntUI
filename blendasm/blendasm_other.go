//go:build !amd64 || purego

package blendasm

// No assembly for this architecture: the Go loops in canvas are the
// implementation, and these are never called.

const Have = false

func Rows(dst, src *uint32, count int)        { panic("unreachable") }
func Solid(dst *uint32, colour uint32, n int)  { panic("unreachable") }

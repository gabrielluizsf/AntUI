//go:build amd64 && !purego

package blendasm

// SSE2 is part of the amd64 baseline — every 64-bit x86 processor has it —
// so there is nothing to detect and no path that can be missing at runtime.
const Have = true

// Rows composites count pixels of src over dst. count must be a multiple of
// four and both runs must hold at least that many pixels.
//
//go:noescape
func Rows(dst, src *uint32, count int)

// Solid composites one colour over count pixels of dst. count must be a
// multiple of four, and the colour's alpha must be neither 0 nor 255 — those
// two are handled without arithmetic by the caller.
//
//go:noescape
func Solid(dst *uint32, colour uint32, count int)

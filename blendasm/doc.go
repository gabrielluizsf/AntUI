// Package blendasm is the assembly half of row blending, and it is a package
// of its own for one reason: a package that uses cgo may not hold Go
// assembly. The toolchain hands .s files in a cgo package to the C compiler,
// which does not read Go's assembly syntax, and says so.
//
// antui uses cgo on macOS, because AppKit is Objective-C. So the SSE2 kernel
// cannot live beside the window code, and moving it here is what keeps the
// fast path on macOS rather than dropping to the Go loop there.
//
// It works on *uint32 rather than on canvas.Color to keep the dependency one
// way — the colour type is canvas's, and a package underneath it should not
// know about it. The two are the same thirty-two bits.
package blendasm

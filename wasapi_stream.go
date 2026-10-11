//go:build windows

package antui

// wasapiStream is one open playback stream's handles and buffers.
type wasapiStream struct {
	client   uintptr
	render   uintptr
	frames   uint32
	channels int
	asFloats bool
	block    []float32
}

//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import (
	"github.com/gabrielluizsf/antui/backend"
	"github.com/gabrielluizsf/antui/canvas"
)

// Present copies the window's dirty rectangle to the drawable, in as many
// chunks as the server's maximum request size demands. It reports whether
// every chunk got written: a connection that drops partway through leaves
// the drawable holding only the rows sent before it did, and the caller
// needs to know that rather than mistake it for a delivered frame.
func (x *Driver) Present(win backend.Face, dirty canvas.Area) bool {
	if !x.alive || dirty.Width <= 0 || dirty.Height <= 0 {
		return false
	}
	cv := win.Canvas()
	rowPixels := dirty.Width
	rowBytes := rowPixels * 4
	rowsPerChunk := max(int((x.maxRequestBytes-32)/uint32(rowBytes)), 1)

	if len(x.scratch) < rowBytes*rowsPerChunk {
		x.scratch = make([]byte, rowBytes*rowsPerChunk)
	}

	for y := 0; y < dirty.Height; y += rowsPerChunk {
		rows := min(dirty.Height-y, rowsPerChunk)

		for row := range rows {
			at := (dirty.Y+y+row)*cv.Stride + dirty.X
			x.packRow(x.scratch[row*rowBytes:], cv.Pixels[at:at+rowPixels])
		}

		if !x.putChunk(dirty, y, rows, rowBytes) {
			return false
		}
	}
	return true
}

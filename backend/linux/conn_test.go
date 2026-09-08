//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import (
	"encoding/binary"
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
)

// The fast path and the slow one have to agree, because which one runs
// depends on the server: a machine where they differ gets the wrong colours
// and nothing says so.
func TestPackRowAgrees(t *testing.T) {
	pixels := []canvas.Color{
		canvas.RGBA(0x11, 0x22, 0x33, 0x44), canvas.White, canvas.Black,
		canvas.RGB(0xFF, 0x00, 0x00), canvas.RGB(0x00, 0xFF, 0x00), canvas.RGB(0x00, 0x00, 0xFF),
	}

	fast := make([]byte, len(pixels)*4)
	slow := make([]byte, len(pixels)*4)
	(&Driver{samePixels: true}).packRow(fast, pixels)
	(&Driver{samePixels: false}).packRow(slow, pixels)

	for i := range fast {
		if fast[i] != slow[i] {
			t.Fatalf("byte %d is %#02x on the fast path and %#02x on the slow one",
				i, fast[i], slow[i])
		}
	}

	// And the bytes really are the pixels, little end first, which is what
	// LSBFirst means.
	if got := binary.LittleEndian.Uint32(slow); got != uint32(pixels[0]) {
		t.Errorf("the first pixel went out as %#08x, want %#08x", got, uint32(pixels[0]))
	}

	// A server that wants red and blue the other way round gets them that
	// way round.
	swapped := make([]byte, 4)
	(&Driver{swapRB: true}).packRow(swapped, []canvas.Color{canvas.RGB(0xFF, 0x00, 0x00)})
	if got := binary.LittleEndian.Uint32(swapped); got != 0xFF0000FF {
		t.Errorf("red came out as %#08x on a swapped visual", got)
	}
}
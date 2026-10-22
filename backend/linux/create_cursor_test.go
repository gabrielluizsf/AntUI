//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import (
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
)

// TestAverageBrightKeepsThePictureOwnColour is the regression test for a
// cursor that bleached to plain white: every legacy cursor only gets one
// foreground colour, and a fixed white threw away whatever colour the
// picture itself was drawn in — a neon-blue icon came out indistinguishable
// from the system's own default pointer.
func TestAverageBrightKeepsThePictureOwnColour(t *testing.T) {
	cv, err := canvas.NewCanvas(4, 4)
	if err != nil {
		t.Fatalf("NewCanvas: %v", err)
	}
	neon := canvas.Color(0xFF57C1FF) // R 0x57, G 0xC1, B 0xFF
	cv.FillRect(canvas.Area{X: 0, Y: 0, Width: 4, Height: 4}, neon)

	bright := func(px, py int) bool { return true }
	r, g, b := averageBright(cv, 4, 4, bright)
	if r != 0x57 || g != 0xC1 || b != 0xFF {
		t.Fatalf("average = %#02x,%#02x,%#02x, want %#02x,%#02x,%#02x", r, g, b, 0x57, 0xC1, 0xFF)
	}
}

// TestAverageBrightFallsBackToWhite is the other half: an icon with no
// bright pixel at all — plain text-caret black, say — still gets a usable
// cursor instead of one with no foreground colour asked for.
func TestAverageBrightFallsBackToWhite(t *testing.T) {
	cv, err := canvas.NewCanvas(2, 2)
	if err != nil {
		t.Fatalf("NewCanvas: %v", err)
	}
	none := func(px, py int) bool { return false }
	r, g, b := averageBright(cv, 2, 2, none)
	if r != 0xFF || g != 0xFF || b != 0xFF {
		t.Fatalf("average = %#02x,%#02x,%#02x, want white", r, g, b)
	}
}

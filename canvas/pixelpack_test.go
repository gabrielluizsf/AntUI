package canvas

import (
	"encoding/binary"
	"testing"
)

// The point of packRGBA is byte order, so the test is about bytes. Android
// names the format RGBA_8888 and means R, G, B, A **in memory** — which is
// the opposite way round from how the same word is written down.
func TestPackRGBAIsByteOrder(t *testing.T) {
	var got [4]byte
	binary.LittleEndian.PutUint32(got[:], packRGBA(RGBA(0x11, 0x22, 0x33, 0x44)))
	want := [4]byte{0x11, 0x22, 0x33, 0x44} // R, G, B, A
	if got != want {
		t.Fatalf("packRGBA laid the bytes out as %v, want %v (R, G, B, A)", got, want)
	}
}

func TestPackRGBASwapsRedAndBlue(t *testing.T) {
	for _, c := range []Color{Red, Green, Blue, White, Black, Transparent, 0x89ABCDEF} {
		v := packRGBA(c)
		if a := uint8(v >> 24); a != c.A() {
			t.Errorf("packRGBA(%08X): alpha %02X, want %02X", uint32(c), a, c.A())
		}
		if g := uint8(v >> 8); g != c.G() {
			t.Errorf("packRGBA(%08X): green %02X, want %02X", uint32(c), g, c.G())
		}
		// The two that move.
		if r := uint8(v); r != c.R() {
			t.Errorf("packRGBA(%08X): red is %02X in the low byte, want %02X", uint32(c), r, c.R())
		}
		if b := uint8(v >> 16); b != c.B() {
			t.Errorf("packRGBA(%08X): blue is %02X at bits 16-23, want %02X", uint32(c), b, c.B())
		}
	}
}

// A colour that goes out and comes back must be itself: this is what would
// catch a mask written with one F too few, which every channel test above
// would still pass.
func TestPackRGBARoundTrips(t *testing.T) {
	unpack := func(v uint32) Color {
		return RGBA(uint8(v), uint8(v>>8), uint8(v>>16), uint8(v>>24))
	}
	for _, c := range []Color{
		0x00000000, 0xFFFFFFFF, 0x01020304, 0xFF0000FF, 0x00FF00FF, 0x0000FFFF,
		Red, Orange, Yellow, Green, Cyan, Blue, Purple, Magenta,
	} {
		if got := unpack(packRGBA(c)); got != c {
			t.Errorf("packRGBA(%08X) unpacks to %08X", uint32(c), uint32(got))
		}
	}
}

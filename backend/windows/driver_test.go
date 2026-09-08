//go:build windows

package windows

import (
	"testing"

	"github.com/gabrielluizsf/antui/backend"
	"github.com/gabrielluizsf/antui/canvas"
)

// The Win32 backend is syscall with no types of its own to draw on: a window
// handle and a message number, and a wndProc that turns those into events.
// This is the half of it that can be tested without opening a window — the
// key table and the character stream, which are where the wrong answers live.

// testFace is a window that keeps its own record, so the driver can be
// driven without a display of its own and still be asked what happened.
type testFace struct {
	events      []backend.Event
	shouldClose bool
	canvas      *canvas.Canvas
}

func (f *testFace) Push(ev backend.Event) { f.events = append(f.events, ev) }
func (f *testFace) PushSimple(t backend.EventType) {
	f.Push(backend.Event{Type: t})
}
func (f *testFace) SetMouse(x, y int)                   {}
func (f *testFace) ResizeCanvas(width, height int) bool { return true }
func (f *testFace) SetTouchFirst()                      {}
func (f *testFace) SetSafeArea(canvas.Area)             {}
func (f *testFace) Canvas() *canvas.Canvas              { return f.canvas }
func (f *testFace) SetShouldClose()                     { f.shouldClose = true }
func (f *testFace) Bounds() backend.Limits              { return backend.Limits{} }

// wmSize decodes the packings every message works in.
func TestLoHiWord(t *testing.T) {
	cases := []struct {
		packed uintptr
		lo, hi int
	}{
		{0x12345678, 0x5678, 0x1234},
		{0xFFFFFFFF, -1, -1},
		{0x00000000, 0, 0},
		{0x00010064, 100, 1},
	}
	for _, c := range cases {
		if lo := loWord(c.packed); lo != c.lo {
			t.Errorf("loWord(%#x) = %d, want %d", c.packed, lo, c.lo)
		}
		if hi := hiWord(c.packed); hi != c.hi {
			t.Errorf("hiWord(%#x) = %d, want %d", c.packed, hi, c.hi)
		}
	}
}

// winKey is where the physical key codes become program keys. The subtle
// part is the extended bit: lparam bit 24 is the only thing that tells a
// right control from a left one, and a mistake there swaps a pair of
// modifier keys nobody notices until a shortcut stops working.
func TestWinKeyMapsVirtualCodes(t *testing.T) {
	cases := []struct {
		name   string
		vk     uintptr
		lparam uintptr
		want   backend.Key
	}{
		{"A", 'A', 0, backend.KeyA},
		{"9", '9', 0, backend.Key9},
		{"space", vkSpace, 0, backend.KeySpace},
		{"escape", vkEscape, 0, backend.KeyEscape},
		{"enter", vkReturn, 0, backend.KeyEnter},
		{"tab", vkTab, 0, backend.KeyTab},
		{"backspace", vkBack, 0, backend.KeyBackspace},
		{"insert", vkInsert, 0, backend.KeyInsert},
		{"delete", vkDelete, 0, backend.KeyDelete},
		{"right", vkRight, 0, backend.KeyRight},
		{"left", vkLeft, 0, backend.KeyLeft},
		{"down", vkDown, 0, backend.KeyDown},
		{"up", vkUp, 0, backend.KeyUp},
		{"page up", vkPrior, 0, backend.KeyPageUp},
		{"page down", vkNext, 0, backend.KeyPageDown},
		{"home", vkHome, 0, backend.KeyHome},
		{"end", vkEnd, 0, backend.KeyEnd},
		{"caps lock", vkCapital, 0, backend.KeyCapsLock},
		{"shift", vkShift, 0, backend.KeyLeftShift},
		{"left shift", vkLShift, 0, backend.KeyLeftShift},
		{"right shift", vkRShift, 0, backend.KeyRightShift},
		{"left control", vkLCtrl, 0, backend.KeyLeftControl},
		{"right control", vkRCtrl, 0, backend.KeyRightControl},
		{"left alt", vkLMenu, 0, backend.KeyLeftAlt},
		{"right alt", vkRMenu, 0, backend.KeyRightAlt},
		{"left super", vkLWin, 0, backend.KeyLeftSuper},
		{"right super", vkRWin, 0, backend.KeyRightSuper},
		{"comma", vkOEMComma, 0, backend.KeyComma},
		{"period", vkOEMPeriod, 0, backend.KeyPeriod},
		{"minus", vkOEMMinus, 0, backend.KeyMinus},
		{"equal", vkOEMPlus, 0, backend.KeyEqual},
		{"f1", vkF1, 0, backend.KeyF1},
		{"f12", vkF12, 0, backend.KeyF12},
		{"numpad 0", vkNumpad0, 0, backend.Key0},
		{"numpad 9", vkNumpad9, 0, backend.Key9},
		{"a key nobody invented", 0x7E, 0, backend.KeyUnknown},

		{"control, plain", vkControl, 0, backend.KeyLeftControl},
		{"control, extended", vkControl, 1 << 24, backend.KeyRightControl},
		{"alt, plain", vkMenu, 0, backend.KeyLeftAlt},
		{"alt, extended", vkMenu, 1 << 24, backend.KeyRightAlt},
	}
	for _, c := range cases {
		if got := winKey(c.vk, c.lparam); got != c.want {
			t.Errorf("%s: winKey(%#x, %#x) = %v, want %v", c.name, c.vk, c.lparam, got, c.want)
		}
	}
}

// Text arrives as UTF-16 code units, and a character above the BMP arrives
// as two. The window has to hold the high half until its partner shows up,
// and drop the halves that never find one.
func TestHandleCharPairsSurrogates(t *testing.T) {
	face := &testFace{}
	d := &Driver{win: face}

	if got := d.handleChar(0xD83D); got != 0 || len(face.events) != 0 {
		t.Fatalf("a high surrogate alone became %d events", len(face.events))
	}
	d.handleChar(0xDE00)
	if len(face.events) != 1 {
		t.Fatalf("a joined pair became %d events", len(face.events))
	}
	if ev := face.events[0]; ev.Rune != 0x1F600 || ev.Text != "\U0001F600" {
		t.Errorf("the pair came out as %q (%U), want a smiley", ev.Text, ev.Rune)
	}
}

func TestHandleCharFiltersControl(t *testing.T) {
	face := &testFace{}
	d := &Driver{win: face}

	for _, unit := range []uint16{'\t', 0x1B, 0x7F, 0xDC00} {
		before := len(face.events)
		d.handleChar(unit)
		if len(face.events) != before {
			t.Errorf("unit %#04x produced a text event", unit)
		}
	}

	d.handleChar('a')
	if len(face.events) != 1 {
		t.Fatalf("a plain letter became %d events", len(face.events))
	}
	if ev := face.events[0]; ev.Rune != 'a' || ev.Text != "a" {
		t.Errorf("the letter came out as %q (%U)", ev.Text, ev.Rune)
	}
}
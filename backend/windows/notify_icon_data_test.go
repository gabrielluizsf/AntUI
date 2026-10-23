//go:build windows

package windows

import (
	"testing"
	"unsafe"
)

// TestNotifyIconDataSize: the shell picks which NOTIFYICONDATAW it was handed
// by cbSize alone, so a struct one byte off is refused outright.
func TestNotifyIconDataSize(t *testing.T) {
	want := uintptr(956) // NOTIFYICONDATAW_V3 + hBalloonIcon, 32-bit
	if unsafe.Sizeof(uintptr(0)) == 8 {
		want = 976
	}
	if got := unsafe.Sizeof(notifyIconDataW{}); got != want {
		t.Fatalf("sizeof(notifyIconDataW) = %d, want %d", got, want)
	}
}

func TestCopyUTF16(t *testing.T) {
	field := make([]uint16, 4)
	copyUTF16(field, "abcdef")
	if string(rune(field[0]))+string(rune(field[1]))+string(rune(field[2])) != "abc" || field[3] != 0 {
		t.Fatalf("a long string should be cut to fit with its NUL, got %v", field)
	}
	// "ab😀": the emoji is a surrogate pair that would straddle the cut, so
	// it goes whole.
	copyUTF16(field, "ab😀")
	if field[2] != 0 || field[1] != 'b' {
		t.Fatalf("a surrogate pair split by the cut should be dropped, got %v", field)
	}
	copyUTF16(field, "")
	if field[0] != 0 {
		t.Fatalf("an empty string should leave an empty field, got %v", field)
	}
}

//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import (
	"encoding/binary"
	"net"
	"testing"
)

// TestCreateGCForSetsForegroundAndBackground is the regression test for a
// cursor that came out as its own negative: PutImage's format Bitmap draws a
// 1 bit as the GC's foreground pixel and a 0 bit as its background, and a
// fresh GC's own defaults are foreground 0, background 1 — backwards from
// the plain pass-through [Driver.packBitmap]'s bits are meant to be. Without
// an explicit foreground of 1 and background of 0, every cursor's source
// and mask pixmap comes out inverted: a crosshair on a transparent square
// turns into a transparent crosshair cut out of a solid square.
func TestCreateGCForSetsForegroundAndBackground(t *testing.T) {
	client, server := net.Pipe()
	defer client.Close()
	defer server.Close()

	received := make(chan []byte, 1)
	go func() {
		buf := make([]byte, 64)
		n, _ := server.Read(buf)
		received <- buf[:n]
	}()

	x := &Driver{conn: client, alive: true, idMask: 0x1FFFFF}
	x.createGCFor(7)

	body := <-received
	if len(body) < 24 {
		t.Fatalf("CreateGC request was %d bytes, want at least 24", len(body))
	}
	if body[0] != xCreateGC {
		t.Fatalf("opcode = %d, want %d", body[0], xCreateGC)
	}
	if mask := binary.LittleEndian.Uint32(body[12:]); mask != gcForeground|gcBackground {
		t.Fatalf("value-mask = %#x, want foreground and background set", mask)
	}
	if fg := binary.LittleEndian.Uint32(body[16:]); fg != 1 {
		t.Errorf("foreground = %d, want 1", fg)
	}
	if bg := binary.LittleEndian.Uint32(body[20:]); bg != 0 {
		t.Errorf("background = %d, want 0", bg)
	}
}

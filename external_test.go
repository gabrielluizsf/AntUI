package antui

import (
	"testing"
	"time"

	"github.com/gabrielluizsf/antui/canvas"
)

// While another renderer owns the surface, End must not send the canvas: a
// copy would land on top of the renderer's frame, and the scan for what
// changed is time a 3D frame does not have.
func TestExternalWindowNeverPresentsTheCanvas(t *testing.T) {
	win, stub := newTestWindow(t, 64, 64)
	win.SetFPS(0)
	win.sleep = func(time.Duration) {}

	if !win.SetExternal(true) || !stub.external {
		t.Fatal("SetExternal did not reach the platform")
	}
	for range 3 {
		if !win.Begin() {
			t.Fatal("window closed")
		}
		win.FillRect(canvas.Area{X: 0, Y: 0, Width: 64, Height: 64}, canvas.White)
		if win.End() {
			t.Error("End reported a present while another renderer owns the surface")
		}
	}
	if len(stub.presented) != 0 {
		t.Errorf("the canvas was presented %d times over the other renderer", len(stub.presented))
	}
}

// Falling back to software drawing has to repaint everything: the screen
// holds the other renderer's last frame, not the canvas the diff remembers.
func TestLeavingExternalSendsTheWholeCanvas(t *testing.T) {
	win, stub := newTestWindow(t, 64, 64)
	win.SetFPS(0)
	win.sleep = func(time.Duration) {}

	win.Begin()
	drawTestFrame(win)
	win.End()

	win.SetExternal(true)
	win.Begin()
	drawTestFrame(win)
	win.End()

	win.SetExternal(false)
	if stub.external {
		t.Fatal("SetExternal(false) did not reach the platform")
	}
	stub.presented = nil
	win.Begin()
	drawTestFrame(win) // the same picture the diff already holds
	if !win.End() {
		t.Fatal("the first software frame after an external stretch sent nothing")
	}
	whole := canvas.Area{Width: 64, Height: 64}
	if len(stub.presented) != 1 || stub.presented[0] != whole {
		t.Errorf("presented %v, want the whole canvas %v", stub.presented, whole)
	}
}

func TestNativeHandleComesFromThePlatform(t *testing.T) {
	win, stub := newTestWindow(t, 8, 8)
	if _, ok := win.NativeHandle(); ok {
		t.Error("a platform with no handle reported one")
	}
	stub.handle = NativeHandle{System: NativeX11, Window: 42, Display: ":0"}
	h, ok := win.NativeHandle()
	if !ok || h != stub.handle {
		t.Errorf("NativeHandle = %+v, %v; want %+v", h, ok, stub.handle)
	}

	off := &Window{}
	if _, ok := off.NativeHandle(); ok {
		t.Error("a window with no platform reported a handle")
	}
	if off.SetExternal(true) || !off.External() {
		t.Error("SetExternal with no platform should remember the mode and report false")
	}
}

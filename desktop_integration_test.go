//go:build linux || darwin || windows

package antui

import (
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/gabrielluizsf/antui/assert"
	"github.com/gabrielluizsf/antui/canvas"
)

// uiTestQueue is a channel used to dispatch functions to the main OS thread.
var uiTestQueue = make(chan func())

func init() {
	// Lock the main OS thread for macOS AppKit UI requirements.
	runtime.LockOSThread()
}

// TestMain acts as the entry point for the tests, keeping the main thread
// alive and processing UI closures while the test runner works in the background.
func TestMain(m *testing.M) {
	code := make(chan int)
	go func() {
		code <- m.Run()
	}()

	for {
		select {
		case f := <-uiTestQueue:
			f()
		case c := <-code:
			os.Exit(c)
		}
	}
}

// RunOnMain synchronously executes the provided function on the main OS thread.
func RunOnMain(f func()) {
	done := make(chan struct{})
	uiTestQueue <- func() {
		f()
		close(done)
	}
	<-done
}

// TestDesktopWindowRunsAnApp is the integration test: it opens a real window
// on the machine the suite is running on and drives it the way an app would,
// so the backends that cannot be exercised from a cross-compile — the X11
// socket, the Win32 message loop, the Cocoa run loop — are heard from by the
// only thing that can make them say anything, which is running them.
func TestDesktopWindowRunsAnApp(t *testing.T) {
	RunOnMain(func() {
		win, err := OpenWith(Options{
			Title:  "AntUI integration test",
			Width:  640,
			Height: 480,
			Limits: Limits{MaxWidth: 1600, MaxHeight: 1200},
		})
		if err != nil {
			t.Skipf("no desktop to open a window on: %v", err)
			return
		}
		defer win.Close()

		assert.Equal(t, 640, win.Width())
		assert.Equal(t, 480, win.Height())

		win.SetTitle("AntUI integration test")

		big, err := canvas.NewCanvas(32, 32)
		assert.NoError(t, err)
		big.Clear(canvas.RGB(0x3E, 0x63, 0xDD))

		small, err := canvas.NewCanvas(16, 16)
		assert.NoError(t, err)
		small.Clear(canvas.RGB(0x3E, 0x63, 0xDD))

		assert.True(t, win.SetIcon(big, small))

		swatch := canvas.RGB(0x3E, 0x63, 0xDD)
		frames, events, drew := 0, 0, false
		resizeSeen := false
		deadline := time.Now().Add(4 * time.Second)

		for win.Begin() && frames < 240 {
			win.Clear(canvas.RGB(0x10, 0x12, 0x18))
			win.FillRect(20, 20, 200, 120, swatch)
			win.Circle(win.Width()/2, win.Height()/2, 60, canvas.White)
			win.Line(0, win.Height()-1, win.Width()-1, 0, canvas.RGBA(0xFF, 0, 0, 0xFF))
			win.Rect(4, 4, win.Width()-8, win.Height()-8, canvas.RGBA(0xFF, 0xFF, 0xFF, 0x80))
			win.Text(20, 160, "AntUI builds desktop apps the simple way.", canvas.White)
			win.End()
			frames++

			if win.Canvas().At(25, 25) == swatch {
				drew = true
			}
			win.Events(func(ev Event) bool {
				events++
				if ev.Type == EventResize {
					resizeSeen = true
				}
				return true
			})

			switch frames {
			case 15:
				assert.True(t, win.SetFullscreen(true))
				assert.True(t, win.Fullscreen())
			case 20:
				win.SetFullscreen(false)
				win.SetSize(800, 600)
			}
			if time.Now().After(deadline) {
				break
			}
		}

		assert.True(t, frames > 0)
		assert.True(t, drew)
		t.Logf("ran %d frames and heard %d events", frames, events)

		if !resizeSeen {
			assert.Equal(t, 800, win.Width())
		}
		assert.False(t, win.Fullscreen())

		if w, h, ok := win.DisplaySize(); ok {
			t.Logf("the display is %dx%d", w, h)
		} else {
			t.Log("the display did not say its size")
		}
		if hz := win.DisplayRefresh(); hz > 0 {
			t.Logf("the display refreshes at %d Hz", hz)
		}
		if scale := win.ContentScale(); scale > 1 {
			t.Logf("the display is scaled at %g", scale)
		}

		const clip = "AntUI clipboard round-trip"
		assert.True(t, win.SetClipboardText(clip))
		assert.Equal(t, clip, win.ClipboardText())

		win.SetFixedSize(640, 480)
		assert.False(t, win.Resizable())

		win.SetLimits(Limits{})
		assert.True(t, win.Resizable())

		if _, err := SystemFace(14); err != nil {
			t.Logf("no system font to draw with: %v", err)
		} else {
			t.Log("read the platform's own font")
		}
	})
}
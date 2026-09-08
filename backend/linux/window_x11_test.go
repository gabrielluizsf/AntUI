//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import (
	"encoding/binary"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gabrielluizsf/antui/backend"
	"github.com/gabrielluizsf/antui/canvas"
)

// These need a real X server, so they skip without one. They exist because
// the whole input path once broke without a single unit test noticing:
// drawing goes out over the socket and input comes back over it, and a
// backend that only ever writes looks perfectly healthy in a screenshot.

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
func (f *testFace) SetMouse(x, y int)                  {}
func (f *testFace) ResizeCanvas(width, height int) bool { return true }
func (f *testFace) SetTouchFirst()                     {}
func (f *testFace) SetSafeArea(canvas.Area)             {}
func (f *testFace) Canvas() *canvas.Canvas              { return f.canvas }
func (f *testFace) SetShouldClose()                    { f.shouldClose = true }
func (f *testFace) Bounds() backend.Limits             { return backend.Limits{} }

func openTestWindow(t *testing.T, title string) (*Driver, *testFace) {
	t.Helper()
	if os.Getenv("DISPLAY") == "" {
		t.Skip("no DISPLAY: this needs a real X server")
	}
	face := &testFace{}
	cv, err := canvas.NewCanvas(320, 200)
	if err != nil {
		t.Fatal(err)
	}
	face.canvas = cv
	d, err := Open(face, backend.Options{Title: title, Width: 320, Height: 200})
	if err != nil {
		t.Skipf("could not open a window: %v", err)
	}
	t.Cleanup(d.Close)
	return d, face
}

// drain runs a few frames so that whatever was just asked has been sent and
// whatever came back has been read.
func drain(t *testing.T, d *Driver, face *testFace, frames int) {
	t.Helper()
	full := canvas.Area{X: 0, Y: 0, Width: face.canvas.Width, Height: face.canvas.Height}
	for range frames {
		d.Pump(face)
		d.Present(face, full)
	}
	time.Sleep(50 * time.Millisecond)
}

// A newly mapped window always hears something back — an expose, a focus
// change, a configure. Hearing nothing means the read path is dead, which is
// what a window that cannot be closed feels like from the outside.
//
// What arrives unasked depends on the server and on whether there is a window
// manager at all: a bare X server with nothing else running on it says very
// little. So the window asks for something after a second — a size it is not
// — which every server answers with a configure. Either way a backend that
// only ever writes hears nothing and fails here, which is the whole point.
func TestEventsArriveFromTheServer(t *testing.T) {
	d, face := openTestWindow(t, "antui event test")

	asked := false
	full := canvas.Area{X: 0, Y: 0, Width: face.canvas.Width, Height: face.canvas.Height}
	deadline := time.Now().Add(3 * time.Second)
	for d.alive {
		d.Pump(face)
		if len(face.events) > 0 {
			return // anything at all is enough
		}
		d.Present(face, full)

		if !asked && time.Now().After(deadline.Add(-2*time.Second)) {
			d.SetSize(face, 400, 260)
			asked = true
		}
		if time.Now().After(deadline) {
			break
		}
		time.Sleep(2 * time.Millisecond)
	}

	t.Fatalf("no event arrived in three seconds; the connection is write-only "+
		"(alive %v %q, asked for a resize: %v)", d.alive, d.dead, asked)
}

// The connection has to survive being idle. A read deadline left over from
// the handshake once killed it exactly five seconds in, which looked like a
// window that froze for no reason.
func TestConnectionSurvivesIdling(t *testing.T) {
	if testing.Short() {
		t.Skip("takes seven seconds")
	}
	d, face := openTestWindow(t, "antui idle test")

	full := canvas.Area{X: 0, Y: 0, Width: face.canvas.Width, Height: face.canvas.Height}
	deadline := time.Now().Add(7 * time.Second)
	for d.alive {
		d.Pump(face)
		d.Present(face, full)
		if time.Now().After(deadline) {
			return
		}
		time.Sleep(5 * time.Millisecond)
	}

	t.Fatalf("the window closed itself while idling (%q)", d.dead)
}

// The size hints go on the window as a property, so the way to know they
// really landed is to read them back off the server. This is the same check
// `xprop -name <title> WM_NORMAL_HINTS` prints, done without the tool.
func TestSizeHintsReachTheServer(t *testing.T) {
	d, face := openTestWindow(t, "antui size hints")

	d.SetLimits(backend.Limits{MinWidth: 640, MinHeight: 480, MaxWidth: 1280,
		MaxHeight: 960, Aspect: 16.0 / 9.0})
	drain(t, d, face, 3)

	value, got := d.property(d.window, xWMNormalHints, xWMSizeHints, 4096)
	if !got || len(value) < 18*4 {
		t.Fatalf("WM_NORMAL_HINTS came back as %d bytes", len(value))
	}
	field := func(i int) uint32 { return binary.LittleEndian.Uint32(value[i*4:]) }

	if flags := field(0); flags&xSizeHintMinSize == 0 || flags&xSizeHintMaxSize == 0 ||
		flags&xSizeHintAspect == 0 {
		t.Errorf("the flags say %#x", flags)
	}
	if field(5) != 640 || field(6) != 480 {
		t.Errorf("the minimum reached the server as %dx%d", field(5), field(6))
	}
	if field(7) != 1280 || field(8) != 960 {
		t.Errorf("the maximum reached the server as %dx%d", field(7), field(8))
	}
	// The ratio goes as a fraction, with both ends equal so the window
	// manager keeps exactly it rather than a range around it.
	if field(11) != 1778 || field(12) != 1000 || field(13) != 1778 || field(14) != 1000 {
		t.Errorf("the ratio reached the server as %d/%d..%d/%d",
			field(11), field(12), field(13), field(14))
	}

	// A window unlocked afterwards must not leave the old hints behind, or it
	// stays locked to numbers nothing in the program remembers setting.
	d.SetLimits(backend.Limits{})
	drain(t, d, face, 3)
	value, got = d.property(d.window, xWMNormalHints, xWMSizeHints, 4096)
	if !got || len(value) < 4 {
		t.Fatal("the hints vanished rather than being cleared")
	}
	if flags := binary.LittleEndian.Uint32(value); flags != 0 {
		t.Errorf("an unlocked window still says %#x", flags)
	}
}

// Xft.dpi is where a desktop's display-scaling setting ends up, and reading
// it is the only answer X11 has to how much the display is scaled by. This
// checks the reader against what xrdb reports, when there is one to ask.
func TestReadsTheDisplayScale(t *testing.T) {
	d, _ := openTestWindow(t, "antui display scale")

	out, err := exec.Command("xrdb", "-query").Output()
	if err != nil {
		t.Skip("no xrdb to check against")
	}
	var want float64
	for line := range strings.Lines(string(out)) {
		name, setting, found := strings.Cut(line, ":")
		if found && strings.TrimSpace(name) == "Xft.dpi" {
			want, _ = strconv.ParseFloat(strings.TrimSpace(setting), 64)
		}
	}
	if want == 0 {
		t.Skip("this display has no Xft.dpi set")
	}

	if got := d.xftDPI(); got != want {
		t.Errorf("read Xft.dpi as %v, xrdb says %v", got, want)
	}
	if got := d.ContentScale(); got != want/96 {
		t.Errorf("the scale is %v, want %v", got, want/96)
	}
}
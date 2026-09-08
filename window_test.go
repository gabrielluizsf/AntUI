package antui

import (
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
)

// stubBackend stands in for a window system, so the frame loop, the input
// state and the widgets can be tested without a display.
type stubBackend struct {
	clip      string
	clipFiles []string
	icons     []*canvas.Canvas

	presented []canvas.Area
	title     string
	full      bool

	// What the stub was asked for, and what it answers with. A zero display
	// size stands for the 1920 by 1080 the older tests were written against.
	limits             Limits
	limitsSet          int
	sized              []canvas.Area
	refuseSize         bool
	displayW, displayH int
	scale              float64
}

func (s *stubBackend) open(*Window, string, int, int) error { return nil }
func (s *stubBackend) close()                               {}
func (s *stubBackend) pump(*Window)                         {}
func (s *stubBackend) present(_ *Window, dirty canvas.Area) { s.presented = append(s.presented, dirty) }
func (s *stubBackend) setTitle(title string)                { s.title = title }
func (s *stubBackend) setFullscreen(on bool) bool           { s.full = on; return true }
func (s *stubBackend) displayRefresh() int                  { return 60 }

func (s *stubBackend) clipboard() (string, []string) { return s.clip, s.clipFiles }

func (s *stubBackend) setIcon(images []*canvas.Canvas) bool {
	s.icons = images
	return true
}

func (s *stubBackend) setClipboard(text string) bool {
	s.clip = text
	return true
}

func (s *stubBackend) displaySize() (int, int, bool) {
	if s.displayW > 0 && s.displayH > 0 {
		return s.displayW, s.displayH, true
	}
	return 1920, 1080, true
}

func (s *stubBackend) setLimits(_ *Window, limits Limits) {
	s.limits = limits
	s.limitsSet++
}

func (s *stubBackend) setSize(_ *Window, width, height int) bool {
	s.sized = append(s.sized, canvas.Area{Width: width, Height: height})
	return !s.refuseSize
}

func (s *stubBackend) contentScale() float64 { return s.scale }

// newTestWindow builds a window backed by the stub, at a known size.
func newTestWindow(t *testing.T, w, h int) (*Window, *stubBackend) {
	t.Helper()
	cv, err := canvas.NewCanvas(w, h)
	if err != nil {
		t.Fatalf("NewCanvas: %v", err)
	}
	stub := &stubBackend{}
	win := &Window{
		cv:             cv,
		native:         stub,
		width:          w,
		height:         h,
		windowedWidth:  w,
		windowedHeight: h,
		fullRedraw:     true,
		queue:          make([]Event, 0, 8),
		theme:          LightTheme(),
	}
	return win, stub
}

func TestBeginClearsPerFrameState(t *testing.T) {
	win, _ := newTestWindow(t, 64, 64)

	win.Push(Event{Type: EventKeyDown, Key: KeyA})
	win.Push(Event{Type: EventMouseDown, Button: MouseLeft, X: 5, Y: 5})
	win.Push(Event{Type: EventMouseWheel, Wheel: 3})
	win.Push(Event{Type: EventText, Text: "a"})

	if !win.KeyPressed(KeyA) || !win.MousePressed(MouseLeft) {
		t.Fatal("the press should be visible in the frame it arrived")
	}
	if win.Wheel() != 3 || win.TextInput() != "a" {
		t.Fatalf("wheel = %d, text = %q, want 3 and \"a\"", win.Wheel(), win.TextInput())
	}

	win.Begin()

	// Held state survives the frame boundary; edges and text do not.
	if !win.KeyDown(KeyA) {
		t.Error("a key held down should stay down across frames")
	}
	if win.KeyPressed(KeyA) || win.MousePressed(MouseLeft) {
		t.Error("presses belong to one frame only")
	}
	if win.Wheel() != 0 || win.TextInput() != "" {
		t.Error("the wheel and the typed text belong to one frame only")
	}
	if !win.MouseDown(MouseLeft) {
		t.Error("a button held down should stay down across frames")
	}
}

func TestLosingFocusReleasesEverything(t *testing.T) {
	win, _ := newTestWindow(t, 64, 64)
	win.Push(Event{Type: EventKeyDown, Key: KeyA, Mods: ModShift})
	win.Push(Event{Type: EventMouseDown, Button: MouseLeft})

	win.Push(Event{Type: EventFocus, Focused: false})

	if win.KeyDown(KeyA) {
		t.Error("a key must not stay stuck down once the window loses focus")
	}
	if win.MouseDown(MouseLeft) {
		t.Error("a button must not stay stuck down once the window loses focus")
	}
	if win.Mods() != 0 {
		t.Errorf("mods = %v, want them cleared with the focus", win.Mods())
	}
}

func TestCloseEventStopsTheLoop(t *testing.T) {
	win, _ := newTestWindow(t, 64, 64)
	if !win.Begin() {
		t.Fatal("a fresh window should be running")
	}
	win.Push(Event{Type: EventClose})
	if win.Begin() {
		t.Error("Begin should report false once the window was asked to close")
	}
	if win.Running() {
		t.Error("Running should report false once the window was asked to close")
	}
}

func TestQuit(t *testing.T) {
	win, _ := newTestWindow(t, 64, 64)
	win.Quit()
	if win.Begin() {
		t.Error("Begin should report false after Quit")
	}
}

func TestEventQueueDrains(t *testing.T) {
	win, _ := newTestWindow(t, 64, 64)
	win.Push(Event{Type: EventExpose})
	win.Push(Event{Type: EventResize, Width: 10, Height: 20})

	var got []EventType
	for ev := range win.Events {
		got = append(got, ev.Type)
	}
	if len(got) != 2 || got[0] != EventExpose || got[1] != EventResize {
		t.Errorf("drained %v, want [Expose Resize]", got)
	}
	if _, ok := win.NextEvent(); ok {
		t.Error("the queue should be empty after draining it")
	}
}

func TestEventQueueDropsRatherThanGrowing(t *testing.T) {
	win, _ := newTestWindow(t, 64, 64)
	for range maxEvents + 50 {
		win.Push(Event{Type: EventExpose})
	}
	if len(win.queue) != maxEvents {
		t.Errorf("queue length = %d, want it capped at %d", len(win.queue), maxEvents)
	}
	if win.droppedEvents != 50 {
		t.Errorf("dropped = %d, want 50", win.droppedEvents)
	}
	// The state still has to be right: the events were seen, only unqueued.
	if !win.fullRedraw {
		t.Error("an expose past the cap should still have been folded into the state")
	}
}

func TestMouseMoveAccumulatesDelta(t *testing.T) {
	win, _ := newTestWindow(t, 64, 64)
	win.Begin()
	win.Push(Event{Type: EventMouseMove, X: 10, Y: 10})
	win.Push(Event{Type: EventMouseMove, X: 14, Y: 7})

	if win.MouseX() != 14 || win.MouseY() != 7 {
		t.Errorf("position = %d,%d, want 14,7", win.MouseX(), win.MouseY())
	}
	if win.mouseDX != 14 || win.mouseDY != 7 {
		t.Errorf("delta = %d,%d, want the whole frame's movement", win.mouseDX, win.mouseDY)
	}
}

func TestResizeReplacesTheCanvas(t *testing.T) {
	win, _ := newTestWindow(t, 64, 64)
	win.Clear(canvas.Red)
	if !win.ResizeCanvas(100, 50) {
		t.Fatal("ResizeCanvas failed")
	}
	if win.Width() != 100 || win.Height() != 50 {
		t.Errorf("size = %dx%d, want 100x50", win.Width(), win.Height())
	}
	if !win.fullRedraw {
		t.Error("a resize must force the next frame to be sent whole")
	}
	if win.shadow != nil {
		t.Error("the previous frame's copy is the wrong size now and must be dropped")
	}
}

func TestDirtyRegion(t *testing.T) {
	win, _ := newTestWindow(t, 32, 32)

	// The first frame has nothing to compare against, so all of it is dirty.
	win.Clear(canvas.Black)
	if got := win.dirtyRegion(); got != (canvas.Area{X: 0, Y: 0, Width: 32, Height: 32}) {
		t.Errorf("first frame = %+v, want the whole window", got)
	}

	// An unchanged frame sends nothing at all.
	if got := win.dirtyRegion(); got != (canvas.Area{}) {
		t.Errorf("unchanged frame = %+v, want nothing", got)
	}

	// One pixel changing sends exactly that pixel.
	win.Pixel(10, 20, canvas.White)
	if got := win.dirtyRegion(); got != (canvas.Area{X: 10, Y: 20, Width: 1, Height: 1}) {
		t.Errorf("one pixel = %+v, want a 1x1 area at 10,20", got)
	}

	// A rectangle changing sends its bounding box.
	win.FillRect(4, 6, 5, 3, canvas.Blue)
	if got := win.dirtyRegion(); got != (canvas.Area{X: 4, Y: 6, Width: 5, Height: 3}) {
		t.Errorf("rectangle = %+v, want its bounding box", got)
	}
}

func TestSetFullscreenRemembersTheWindowedSize(t *testing.T) {
	win, stub := newTestWindow(t, 800, 600)
	if !win.SetFullscreen(true) {
		t.Fatal("SetFullscreen(true) failed")
	}
	if !win.Fullscreen() || !stub.full {
		t.Error("the window should report that full screen was asked for")
	}
	// The window manager grants it and resizes us.
	win.ResizeCanvas(1920, 1080)
	// Asking again must not overwrite the size to come back out to.
	win.SetFullscreen(true)
	if win.windowedWidth != 800 || win.windowedHeight != 600 {
		t.Errorf("windowed size = %dx%d, want the 800x600 it went in at",
			win.windowedWidth, win.windowedHeight)
	}
}

func TestSetFPSNeverGoesNegative(t *testing.T) {
	win, _ := newTestWindow(t, 8, 8)
	win.SetFPS(-30)
	if win.targetFPS != 0 {
		t.Errorf("targetFPS = %d, want a negative cap to mean no cap", win.targetFPS)
	}
}
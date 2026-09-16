package template

import (
	"testing"

	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/event"
)

func TestCyberpunkButtonClicked(t *testing.T) {
	win, _, err := antui.Offscreen(400, 200)
	if err != nil {
		t.Fatal(err)
	}
	tpl := Cyberpunk(win)

	// No click before interaction.
	win.Begin()
	if e := tpl.Button(win, 40, 80, 120, 36, "Go"); e.Ok() {
		t.Errorf("no click should produce no event, got %v", e)
	}
	win.End()

	// Full press-release in one frame.
	win.Begin()
	testClick(t, win, 80, 100)
	e := tpl.Button(win, 40, 80, 120, 36, "Go")
	win.End()
	if !e.Is(event.Button, event.Click) {
		t.Errorf("expected Click event, got %v", e)
	}
}

func TestBuiltinButtonClicked(t *testing.T) {
	win, _, err := antui.Offscreen(400, 200)
	if err != nil {
		t.Fatal(err)
	}
	tpl := Builtin(win)

	win.Begin()
	testClick(t, win, 100, 55)
	e := tpl.Button(win, 60, 40, 100, 30, "OK")
	win.End()
	if !e.Is(event.Button, event.Click) {
		t.Errorf("expected Click event, got %v", e)
	}
}

func TestButtonDragOffCancels(t *testing.T) {
	win, _, _ := antui.Offscreen(400, 200)
	tpl := Builtin(win)

	// Press inside.
	win.Begin()
	testClick(t, win, 100, 55)
	tpl.Button(win, 60, 40, 100, 30, "OK")
	win.End()

	// Release outside the button.
	win.Begin()
	win.Push(antui.Event{Type: antui.EventMouseMove, X: 300, Y: 150})
	win.Push(antui.Event{Type: antui.EventMouseUp, Button: antui.MouseLeft, X: 300, Y: 150})
	e := tpl.Button(win, 60, 40, 100, 30, "OK")
	win.End()
	if e.Ok() {
		t.Error("releasing away from the button should cancel the click")
	}
}

func TestCheckboxToggles(t *testing.T) {
	win, _, _ := antui.Offscreen(400, 200)
	tpl := Builtin(win)
	value := false

	win.Begin()
	testClick(t, win, 20, 20)
	e := tpl.Checkbox(win, 10, 10, "Enable", &value)
	win.End()
	if !value || !e.Is(event.Checkbox, event.Toggle) {
		t.Errorf("checkbox should toggle on, got value=%v event=%v", value, e)
	}
}

func TestInputFiresFocusAndType(t *testing.T) {
	win, _, _ := antui.Offscreen(400, 200)
	tpl := Cyberpunk(win)
	text := ""

	// Click the field to focus it.
	win.Begin()
	testClick(t, win, 100, 105)
	focus := tpl.Input(win, 50, 90, 200, 30, &text)
	win.End()
	if !focus.Is(event.TextInput, event.Focus) {
		t.Errorf("clicking a field should produce a Focus event, got %v", focus)
	}

	// Type a character.
	win.Begin()
	win.Push(antui.Event{Type: antui.EventText, Text: "A"})
	typed := tpl.Input(win, 50, 90, 200, 30, &text)
	win.End()
	if !typed.Is(event.TextInput, event.Type) {
		t.Errorf("typing should produce a Type event, got %v", typed)
	}
	if text != "A" {
		t.Errorf("typed text = %q, want %q", text, "A")
	}
}

func TestOffscreenCanvasIsDrawable(t *testing.T) {
	win, cv, err := antui.Offscreen(64, 64)
	if err != nil {
		t.Fatal(err)
	}
	tpl := Cyberpunk(win)
	win.Begin()
	tpl.Button(win, 10, 10, 44, 44, "X")
	win.End()
	hot := false
	for _, c := range cv.Pixels {
		if c != 0 {
			hot = true
			break
		}
	}
	if !hot {
		t.Error("Cyberpunk.Button painted nothing onto the canvas")
	}
}

// busyboxStyle satisfies Style for testing New() without any drawing.
type busyboxStyle struct{}

func (busyboxStyle) Background(*antui.Window) canvas.Color                             { return 0 }
func (busyboxStyle) Label(*antui.Window, int, int, string)                             {}
func (busyboxStyle) Button(*antui.Window, State, int, int, int, int, string)           {}
func (busyboxStyle) Checkbox(*antui.Window, State, int, int, string, bool)             {}
func (busyboxStyle) Radio(*antui.Window, State, int, int, string, bool)                {}
func (busyboxStyle) Slider(*antui.Window, State, int, int, int, int, float32)          {}
func (busyboxStyle) Input(*antui.Window, State, int, int, int, int, string, int, bool) {}

func TestNewCreatesATemplate(t *testing.T) {
	win, _, _ := antui.Offscreen(200, 200)
	tpl := New(win, busyboxStyle{}, nil)
	if tpl == nil {
		t.Fatal("New returned nil")
	}
	if tpl.Background(win) != 0 {
		t.Error("Background should return the style's colour")
	}
}

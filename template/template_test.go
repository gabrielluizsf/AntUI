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
	tpl := Cyberpunk()

	// No click before interaction.
	win.Begin()
	if e := tpl.Button(win, canvas.Area{X: 40, Y: 80, Width: 120, Height: 36}, "Go"); e.Ok() {
		t.Errorf("no click should produce no event, got %v", e)
	}
	win.End()

	// Full press-release in one frame.
	win.Begin()
	testClick(t, win, 80, 100)
	e := tpl.Button(win, canvas.Area{X: 40, Y: 80, Width: 120, Height: 36}, "Go")
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
	tpl := Builtin()

	win.Begin()
	testClick(t, win, 100, 55)
	e := tpl.Button(win, canvas.Area{X: 60, Y: 40, Width: 100, Height: 30}, "OK")
	win.End()
	if !e.Is(event.Button, event.Click) {
		t.Errorf("expected Click event, got %v", e)
	}
}

func TestButtonDragOffCancels(t *testing.T) {
	win, _, _ := antui.Offscreen(400, 200)
	tpl := Builtin()

	// Press inside.
	win.Begin()
	testClick(t, win, 100, 55)
	tpl.Button(win, canvas.Area{X: 60, Y: 40, Width: 100, Height: 30}, "OK")
	win.End()

	// Release outside the button.
	win.Begin()
	win.Push(antui.Event{Type: antui.EventMouseMove, X: 300, Y: 150})
	win.Push(antui.Event{Type: antui.EventMouseUp, Button: antui.MouseLeft, X: 300, Y: 150})
	e := tpl.Button(win, canvas.Area{X: 60, Y: 40, Width: 100, Height: 30}, "OK")
	win.End()
	if e.Ok() {
		t.Error("releasing away from the button should cancel the click")
	}
}

func TestCheckboxToggles(t *testing.T) {
	win, _, _ := antui.Offscreen(400, 200)
	tpl := Builtin()
	value := false

	win.Begin()
	testClick(t, win, 20, 20)
	e := tpl.Checkbox(win, canvas.Area{X: 10, Y: 10}, antui.CheckboxParams{Label: "Enable", Value: &value})
	win.End()
	if !value || !e.Is(event.Checkbox, event.Toggle) {
		t.Errorf("checkbox should toggle on, got value=%v event=%v", value, e)
	}
}

func TestInputFiresFocusAndType(t *testing.T) {
	win, _, _ := antui.Offscreen(400, 200)
	tpl := Cyberpunk()
	text := ""

	// Click the field to focus it.
	win.Begin()
	testClick(t, win, 100, 105)
	focus := tpl.Input(win, canvas.Area{X: 50, Y: 90, Width: 200, Height: 30}, &text)
	win.End()
	if !focus.Is(event.TextInput, event.Focus) {
		t.Errorf("clicking a field should produce a Focus event, got %v", focus)
	}

	// Type a character.
	win.Begin()
	win.Push(antui.Event{Type: antui.EventText, Text: "A"})
	typed := tpl.Input(win, canvas.Area{X: 50, Y: 90, Width: 200, Height: 30}, &text)
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
	tpl := Cyberpunk()
	win.Begin()
	tpl.Button(win, canvas.Area{X: 10, Y: 10, Width: 44, Height: 44}, "X")
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

func (busyboxStyle) Background(*antui.Window) canvas.Color     { return 0 }
func (busyboxStyle) Label(LabelParams)                         {}
func (busyboxStyle) Button(ButtonParams)                       {}
func (busyboxStyle) Checkbox(CheckboxParams)                   {}
func (busyboxStyle) Radio(RadioParams)                         {}
func (busyboxStyle) Slider(SliderParams)                       {}
func (busyboxStyle) Input(InputParams)                         {}
func (busyboxStyle) InputTextPos(InputTextPosParams) int       { return 0 }
func (busyboxStyle) TextAreaTextPos(TextAreaTextPosParams) int { return 0 }
func (busyboxStyle) Select(SelectParams)                       {}
func (busyboxStyle) SelectOption(SelectOptionParams)           {}
func (busyboxStyle) TextArea(TextAreaParams)                   {}
func (busyboxStyle) Switch(SwitchParams)                       {}
func (busyboxStyle) Progress(ProgressParams)                   {}
func (busyboxStyle) DatePicker(DatePickerParams)               {}
func (busyboxStyle) DatePickerBox(DatePickerBoxParams)         {}

func TestNewCreatesATemplate(t *testing.T) {
	win, _, _ := antui.Offscreen(200, 200)
	tpl := New(busyboxStyle{}, nil)
	if tpl == nil {
		t.Fatal("New returned nil")
	}
	if tpl.Background(win) != 0 {
		t.Error("Background should return the style's colour")
	}
}

func TestSelectOpensPicksAndCloses(t *testing.T) {
	win, _, _ := antui.Offscreen(400, 300)
	tpl := Builtin()
	options := []string{"Apple", "Banana", "Cherry"}
	index := 0

	// Click the box: the menu opens.
	win.Begin()
	testClick(t, win, 100, 40)
	e := tpl.Select(win, canvas.Area{X: 50, Y: 25, Width: 140, Height: 30}, SelectParams{Index: &index, Options: options})
	win.End()
	if !e.Is(event.Select, event.Open) {
		t.Fatalf("expected Select/Open, got %v", e)
	}

	// Click the second option: it is picked and the menu closes. The option
	// boxes hang below the box at oy = y + h*(i+1); option 1 spans
	// y+h*2 .. y+h*3.
	win.Begin()
	testClick(t, win, 100, 25+30*2+15)
	e = tpl.Select(win, canvas.Area{X: 50, Y: 25, Width: 140, Height: 30}, SelectParams{Index: &index, Options: options})
	win.End()
	if !e.Is(event.Select, event.Pick) {
		t.Fatalf("expected Select/Pick, got %v", e)
	}
	if index != 1 {
		t.Errorf("index = %d, want 1", index)
	}
}

// TestSelectStaysClosedAfterAnOutsideClick guards a bug where closing the
// menu with a press outside it only changed the local variable this frame's
// Select drew with, never the field that survives to the next one — so a
// menu an outside click had just closed opened again on its own the very
// next frame, as if the click had never landed.
func TestSelectStaysClosedAfterAnOutsideClick(t *testing.T) {
	win, _, _ := antui.Offscreen(400, 300)
	tpl := Builtin()
	options := []string{"Apple", "Banana", "Cherry"}
	index := 0
	area := canvas.Area{X: 50, Y: 25, Width: 140, Height: 30}

	// Click the box: the menu opens.
	win.Begin()
	testClick(t, win, 100, 40)
	e := tpl.Select(win, area, SelectParams{Index: &index, Options: options})
	win.End()
	if !e.Is(event.Select, event.Open) {
		t.Fatalf("expected Select/Open, got %v", e)
	}

	// Click well outside the box and the menu both: it closes.
	win.Begin()
	testClick(t, win, 350, 280)
	e = tpl.Select(win, area, SelectParams{Index: &index, Options: options})
	win.End()
	if !e.Is(event.Select, event.Close) {
		t.Fatalf("expected Select/Close, got %v", e)
	}

	// A frame later, with no further input, it is still closed.
	win.Begin()
	e = tpl.Select(win, area, SelectParams{Index: &index, Options: options})
	win.End()
	if e.Ok() {
		t.Errorf("an untouched frame after the outside click produced %v, want no event", e)
	}
}

func TestSelectKeyboardPicks(t *testing.T) {
	win, _, _ := antui.Offscreen(400, 300)
	tpl := Builtin()
	options := []string{"Alpha", "Beta", "Gamma"}
	index := 0

	// Click the box opens the menu and grabs the keyboard.
	win.Begin()
	testClick(t, win, 60, 40)
	e := tpl.Select(win, canvas.Area{X: 50, Y: 25, Width: 120, Height: 30}, SelectParams{Index: &index, Options: options})
	win.End()
	if !e.Is(event.Select, event.Open) {
		t.Fatalf("expected Select/Open, got %v", e)
	}

	// Arrows walk the menu; nothing is committed yet.
	win.Begin()
	win.Push(antui.Event{Type: antui.EventKeyDown, Key: antui.KeyDown})
	e = tpl.Select(win, canvas.Area{X: 50, Y: 25, Width: 120, Height: 30}, SelectParams{Index: &index, Options: options})
	win.End()
	if e.Ok() {
		t.Fatalf("walking the menu should not pick yet, got %v", e)
	}
	if index != 0 {
		t.Errorf("index = %d, want 0 before commit", index)
	}

	// Up wraps nowhere (top), then Enter commits index 1.
	win.Begin()
	win.Push(antui.Event{Type: antui.EventKeyDown, Key: antui.KeyUp})
	e = tpl.Select(win, canvas.Area{X: 50, Y: 25, Width: 120, Height: 30}, SelectParams{Index: &index, Options: options})
	win.End()
	if e.Ok() {
		t.Fatalf("arrow between picks should not commit, got %v", e)
	}

	win.Begin()
	win.Push(antui.Event{Type: antui.EventKeyDown, Key: antui.KeyDown})
	e = tpl.Select(win, canvas.Area{X: 50, Y: 25, Width: 120, Height: 30}, SelectParams{Index: &index, Options: options})
	win.End()

	win.Begin()
	win.Push(antui.Event{Type: antui.EventKeyDown, Key: antui.KeyEnter})
	e = tpl.Select(win, canvas.Area{X: 50, Y: 25, Width: 120, Height: 30}, SelectParams{Index: &index, Options: options})
	win.End()
	if !e.Is(event.Select, event.Pick) {
		t.Fatalf("expected Select/Pick on Enter, got %v", e)
	}
	if index != 1 {
		t.Errorf("index = %d, want 1", index)
	}
}

func TestSelectAutoclosesOnOutsideClick(t *testing.T) {
	win, _, _ := antui.Offscreen(400, 300)
	tpl := Builtin()
	options := []string{"Only"}
	index := 0

	win.Begin()
	testClick(t, win, 60, 40)
	tpl.Select(win, canvas.Area{X: 50, Y: 25, Width: 100, Height: 30}, SelectParams{Index: &index, Options: options})
	win.End()

	// Click well below the open menu: it closes.
	win.Begin()
	testClick(t, win, 300, 250)
	e := tpl.Select(win, canvas.Area{X: 50, Y: 25, Width: 100, Height: 30}, SelectParams{Index: &index, Options: options})
	win.End()
	if !e.Is(event.Select, event.Close) {
		t.Fatalf("expected Select/Close, got %v", e)
	}
}

func TestTextAreaTypesAndWrapsCursor(t *testing.T) {
	win, _, _ := antui.Offscreen(400, 300)
	tpl := Cyberpunk()
	text := ""

	// Click to focus.
	win.Begin()
	testClick(t, win, 60, 60)
	focus := tpl.TextArea(win, canvas.Area{X: 20, Y: 30, Width: 300, Height: 100}, &text)
	win.End()
	if !focus.Is(event.TextArea, event.Focus) {
		t.Fatalf("expected TextArea/Focus, got %v", focus)
	}

	win.Begin()
	win.Push(antui.Event{Type: antui.EventText, Text: "hi"})
	typed := tpl.TextArea(win, canvas.Area{X: 20, Y: 30, Width: 300, Height: 100}, &text)
	win.End()
	if !typed.Is(event.TextArea, event.Type) {
		t.Fatalf("expected TextArea/Type, got %v", typed)
	}
	if text != "hi" {
		t.Errorf("text = %q, want %q", text, "hi")
	}

	// Enter inserts a line break at the cursor.
	win.Begin()
	win.Push(antui.Event{Type: antui.EventKeyDown, Key: antui.KeyEnter})
	tpl.TextArea(win, canvas.Area{X: 20, Y: 30, Width: 300, Height: 100}, &text)
	win.End()
	if text != "hi\n" {
		t.Errorf("after Enter text = %q, want %q", text, "hi\n")
	}
}

func TestSwitchToggles(t *testing.T) {
	win, _, _ := antui.Offscreen(400, 200)
	tpl := Builtin()
	on := false

	win.Begin()
	testClick(t, win, 50, 60)
	e := tpl.Switch(win, canvas.Area{X: 20, Y: 50}, antui.CheckboxParams{Label: "Boost", Value: &on})
	win.End()
	if !on || !e.Is(event.Switch, event.Toggle) {
		t.Fatalf("switch should toggle on, got on=%v event=%v", on, e)
	}
}

func TestProgressIsInert(t *testing.T) {
	win, _, _ := antui.Offscreen(400, 200)
	tpl := Builtin()
	win.Begin()
	// A click on the bar must not blow up.
	testClick(t, win, 200, 100)
	tpl.Progress(win, canvas.Area{X: 100, Y: 90, Width: 200, Height: 16}, 0.5)
	win.End()
}

func TestDatePickerPicksADay(t *testing.T) {
	win, _, _ := antui.Offscreen(400, 400)
	tpl := Builtin()

	// The box sits at (10,10), 150x30. A click on it opens the calendar
	// popup below it at (10,40).
	value := Date{Year: 2026, Month: 9, Day: 1}
	win.Begin()
	testClick(t, win, 85, 25)
	e := tpl.DatePicker(win, canvas.Area{X: 10, Y: 10, Width: 150, Height: 30}, &value)
	win.End()
	if !e.Is(event.DatePicker, event.Open) {
		t.Fatalf("clicking the box should open the calendar, got %v", e)
	}

	// With the popup centred in the window and grown a step for it, the grid
	// starts where the geometry helpers say it does. Sept 1 2026 is a Tuesday
	// (firstWD=2), so day 10 is at grid index 11: row 1, col 4. Click its
	// centre — the value starts on day 1, so this must change it and close the
	// popup.
	g := dateGeomFor(win, 10, 10, 150, 30)
	cx, cy, cc, _ := dayCellRect(g.px, g.py, g.u, (2+10-1)%7, (2+10-1)/7)
	win.Begin()
	testClick(t, win, cx+cc/2, cy+cc/2)
	e = tpl.DatePicker(win, canvas.Area{X: 10, Y: 10, Width: 150, Height: 30}, &value)
	win.End()
	if !e.Is(event.DatePicker, event.Change) {
		t.Fatalf("expected DatePicker/Change, got %v", e)
	}
	if value.Day != 10 {
		t.Errorf("day = %d, want 10", value.Day)
	}

	// The popup closed, so a non-drawing frame reports nothing.
	win.Begin()
	e = tpl.DatePicker(win, canvas.Area{X: 10, Y: 10, Width: 150, Height: 30}, &value)
	win.End()
	if e.Ok() {
		t.Errorf("closed datepicker should report nothing, got %v", e)
	}
}

func TestDatePickerMarkerFollowsMouse(t *testing.T) {
	win, _, _ := antui.Offscreen(400, 400)
	tpl := Builtin()
	value := Date{Year: 2026, Month: 9, Day: 1}

	// Open the calendar: clicking the box focuses it.
	win.Begin()
	testClick(t, win, 85, 25)
	tpl.DatePicker(win, canvas.Area{X: 10, Y: 10, Width: 150, Height: 30}, &value)
	win.End()

	// Wandering over the grid moves the marker along, without clicking.
	// Sept 1 2026 is a Tuesday (firstWD=2), so day 10 sits at row 1 col 4.
	g := dateGeomFor(win, 10, 10, 150, 30)
	cx, cy, cc, _ := dayCellRect(g.px, g.py, g.u, (2+10-1)%7, (2+10-1)/7)
	win.Begin()
	win.Push(antui.Event{Type: antui.EventMouseMove, X: cx + cc/2, Y: cy + cc/2})
	e := tpl.DatePicker(win, canvas.Area{X: 10, Y: 10, Width: 150, Height: 30}, &value)
	win.End()
	if e.Ok() {
		t.Fatalf("hovering a day should not pick it, got %v", e)
	}
	if value.Day != 1 {
		t.Errorf("day = %d, want 1 until the marker commits", value.Day)
	}

	// Enter commits whatever the marker is sitting on: day 10 now.
	win.Begin()
	win.Push(antui.Event{Type: antui.EventKeyDown, Key: antui.KeyEnter})
	e = tpl.DatePicker(win, canvas.Area{X: 10, Y: 10, Width: 150, Height: 30}, &value)
	win.End()
	if !e.Is(event.DatePicker, event.Change) {
		t.Fatalf("expected DatePicker/Change on Enter, got %v", e)
	}
	if value.Day != 10 {
		t.Errorf("day = %d, want 10 picked where the marker was", value.Day)
	}
}

// TestDatePickerArrowsWorkWhileMouseResting makes sure a parked pointer over
// the calendar does not steal the marker back from the arrow keys.
func TestDatePickerArrowsWorkWhileMouseResting(t *testing.T) {
	win, _, _ := antui.Offscreen(400, 400)
	tpl := Builtin()
	value := Date{Year: 2026, Month: 9, Day: 1}

	// Open the calendar: clicking the box focuses it.
	win.Begin()
	testClick(t, win, 85, 25)
	tpl.DatePicker(win, canvas.Area{X: 10, Y: 10, Width: 150, Height: 30}, &value)
	win.End()

	// Rest the pointer over day 10, then walk one day right with the arrow.
	g := dateGeomFor(win, 10, 10, 150, 30)
	cx, cy, cc, _ := dayCellRect(g.px, g.py, g.u, (2+10-1)%7, (2+10-1)/7)
	win.Begin()
	win.Push(antui.Event{Type: antui.EventMouseMove, X: cx + cc/2, Y: cy + cc/2})
	tpl.DatePicker(win, canvas.Area{X: 10, Y: 10, Width: 150, Height: 30}, &value)
	win.End()

	// The mouse sits still now; Right must keep the keyboard's day.
	win.Begin()
	win.Push(antui.Event{Type: antui.EventKeyDown, Key: antui.KeyRight})
	tpl.DatePicker(win, canvas.Area{X: 10, Y: 10, Width: 150, Height: 30}, &value)
	win.End()

	win.Begin()
	win.Push(antui.Event{Type: antui.EventKeyDown, Key: antui.KeyEnter})
	e := tpl.DatePicker(win, canvas.Area{X: 10, Y: 10, Width: 150, Height: 30}, &value)
	win.End()
	if !e.Is(event.DatePicker, event.Change) {
		t.Fatalf("expected DatePicker/Change on Enter, got %v", e)
	}
	// The hover put the marker on 10; the arrow walked it to 11. The parked
	// pointer must not have pulled it back to 10.
	if value.Day != 11 {
		t.Errorf("day = %d, want 11 (the arrow walked the hovered 10 forward)", value.Day)
	}
}

func TestDatePickerKeyboardPicksACloser(t *testing.T) {
	win, _, _ := antui.Offscreen(400, 400)
	tpl := Builtin()
	value := Date{Year: 2026, Month: 9, Day: 1}

	// Open the calendar: clicking the box focuses it.
	win.Begin()
	testClick(t, win, 85, 25)
	tpl.DatePicker(win, canvas.Area{X: 10, Y: 10, Width: 150, Height: 30}, &value)
	win.End()

	// Arrow right highlights day 2; Enter commits it and closes.
	win.Begin()
	win.Push(antui.Event{Type: antui.EventKeyDown, Key: antui.KeyRight})
	e := tpl.DatePicker(win, canvas.Area{X: 10, Y: 10, Width: 150, Height: 30}, &value)
	win.End()
	if e.Ok() {
		t.Errorf("walking days should not pick yet, got %v", e)
	}
	if value.Day != 1 {
		t.Errorf("day = %d, want 1 before commit", value.Day)
	}

	win.Begin()
	win.Push(antui.Event{Type: antui.EventKeyDown, Key: antui.KeyEnter})
	e = tpl.DatePicker(win, canvas.Area{X: 10, Y: 10, Width: 150, Height: 30}, &value)
	win.End()
	if !e.Is(event.DatePicker, event.Change) {
		t.Fatalf("expected DatePicker/Change on Enter, got %v", e)
	}
	if value.Day != 2 {
		t.Errorf("day = %d, want 2", value.Day)
	}
}

func TestDatePickerEscapeCancels(t *testing.T) {
	win, _, _ := antui.Offscreen(400, 400)
	tpl := Builtin()
	value := Date{Year: 2026, Month: 9, Day: 1}

	// Open, walk to another day, then Escape without committing.
	win.Begin()
	testClick(t, win, 85, 25)
	tpl.DatePicker(win, canvas.Area{X: 10, Y: 10, Width: 150, Height: 30}, &value)
	win.End()

	win.Begin()
	win.Push(antui.Event{Type: antui.EventKeyDown, Key: antui.KeyRight})
	tpl.DatePicker(win, canvas.Area{X: 10, Y: 10, Width: 150, Height: 30}, &value)
	win.End()

	win.Begin()
	win.Push(antui.Event{Type: antui.EventKeyDown, Key: antui.KeyEscape})
	e := tpl.DatePicker(win, canvas.Area{X: 10, Y: 10, Width: 150, Height: 30}, &value)
	win.End()
	if !e.Is(event.DatePicker, event.Close) {
		t.Fatalf("expected DatePicker/Close on Escape, got %v", e)
	}
	if value.Day != 1 {
		t.Errorf("Escape must not change the date, got day %d", value.Day)
	}
}

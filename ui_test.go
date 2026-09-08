package antui

import "testing"

// click drives one full press-and-release over a point, one frame each, which
// is what an immediate-mode widget needs to see to count a click.
func click(win *Window, x, y int) {
	win.Begin()
	win.Push(Event{Type: EventMouseMove, X: x, Y: y})
	win.Push(Event{Type: EventMouseDown, Button: MouseLeft, X: x, Y: y})
}

func release(win *Window, x, y int) {
	win.Push(Event{Type: EventMouseUp, Button: MouseLeft, X: x, Y: y})
}

func TestButtonClick(t *testing.T) {
	win, _ := newTestWindow(t, 200, 100)

	// Pressing alone is not a click.
	click(win, 50, 30)
	if win.Button(20, 20, 100, 30, "OK") {
		t.Error("a press with no release should not count as a click")
	}
	// Releasing on the button completes it.
	release(win, 50, 30)
	if !win.Button(20, 20, 100, 30, "OK") {
		t.Error("press then release on the button should count as a click")
	}
	// And it does not repeat on the next frame.
	win.Begin()
	if win.Button(20, 20, 100, 30, "OK") {
		t.Error("a click should be reported once, not every frame after")
	}
}

func TestButtonDragOffCancels(t *testing.T) {
	win, _ := newTestWindow(t, 200, 100)
	click(win, 50, 30)
	win.Button(20, 20, 100, 30, "OK")

	// The pointer leaves the button before the release.
	release(win, 180, 90)
	if win.Button(20, 20, 100, 30, "OK") {
		t.Error("releasing away from the button should cancel the click")
	}
}

func TestCheckboxToggles(t *testing.T) {
	win, _ := newTestWindow(t, 200, 100)
	value := false

	click(win, 15, 15)
	win.Checkbox(10, 10, "Enable", &value)
	release(win, 15, 15)
	if !win.Checkbox(10, 10, "Enable", &value) {
		t.Fatal("the checkbox should report the change")
	}
	if !value {
		t.Error("clicking should have turned it on")
	}

	click(win, 15, 15)
	win.Checkbox(10, 10, "Enable", &value)
	release(win, 15, 15)
	win.Checkbox(10, 10, "Enable", &value)
	if value {
		t.Error("clicking again should have turned it back off")
	}
}

func TestRadioSelectsOnceOnly(t *testing.T) {
	win, _ := newTestWindow(t, 200, 200)
	choice := 0

	click(win, 15, 15)
	win.Radio(10, 10, "Second", &choice, 1)
	release(win, 15, 15)
	if !win.Radio(10, 10, "Second", &choice, 1) {
		t.Fatal("selecting a new option should report the change")
	}
	if choice != 1 {
		t.Errorf("choice = %d, want 1", choice)
	}

	// Clicking the option that is already selected changes nothing.
	click(win, 15, 15)
	win.Radio(10, 10, "Second", &choice, 1)
	release(win, 15, 15)
	if win.Radio(10, 10, "Second", &choice, 1) {
		t.Error("re-clicking the selected option should not report a change")
	}
}

func TestSliderDragsAndClamps(t *testing.T) {
	win, _ := newTestWindow(t, 300, 100)
	value := float32(0)

	// Grab the knob, then drag past the right-hand end.
	click(win, 150, 40)
	win.Slider(20, 30, 200, 20, &value, 0, 10)
	win.Push(Event{Type: EventMouseMove, X: 400, Y: 40})
	if !win.Slider(20, 30, 200, 20, &value, 0, 10) {
		t.Fatal("dragging should report a change")
	}
	if value != 10 {
		t.Errorf("value = %v, want it clamped to the maximum", value)
	}

	// Dragging past the left-hand end clamps the other way.
	win.Push(Event{Type: EventMouseMove, X: -100, Y: 40})
	win.Slider(20, 30, 200, 20, &value, 0, 10)
	if value != 0 {
		t.Errorf("value = %v, want it clamped to the minimum", value)
	}
}

func TestSliderCorrectsAnOutOfRangeValue(t *testing.T) {
	win, _ := newTestWindow(t, 300, 100)
	win.Begin()

	value := float32(99)
	win.Slider(20, 30, 200, 20, &value, 0, 10)
	if value != 10 {
		t.Errorf("value = %v, want it brought into range", value)
	}
}

func TestSliderRejectsAnEmptyRange(t *testing.T) {
	win, _ := newTestWindow(t, 300, 100)
	win.Begin()
	value := float32(5)
	if win.Slider(20, 30, 200, 20, &value, 10, 10) {
		t.Error("a range with no span cannot be dragged")
	}
	if value != 5 {
		t.Errorf("value = %v, want it left alone", value)
	}
}

func TestHovered(t *testing.T) {
	win, _ := newTestWindow(t, 200, 200)
	win.Begin()
	win.Push(Event{Type: EventMouseMove, X: 50, Y: 50})

	if !win.Hovered(40, 40, 20, 20) {
		t.Error("the pointer is inside the rectangle")
	}
	// The right and bottom edges are outside, so neighbouring widgets that
	// share an edge cannot both be hovered.
	if win.Hovered(30, 30, 20, 20) {
		t.Error("the far edge of a rectangle is not inside it")
	}
	if win.Hovered(51, 51, 10, 10) {
		t.Error("the pointer is outside the rectangle")
	}
}

func TestInputTyping(t *testing.T) {
	win, _ := newTestWindow(t, 300, 100)
	text := ""

	// Focus the field.
	click(win, 50, 30)
	win.Input(20, 20, 200, 30, &text)
	release(win, 50, 30)

	win.Begin()
	win.Push(Event{Type: EventText, Text: "ol"})
	if !win.Input(20, 20, 200, 30, &text) {
		t.Fatal("typing should report a change")
	}
	if text != "ol" {
		t.Errorf("text = %q, want %q", text, "ol")
	}

	// An accented character arrives as one composed rune.
	win.Begin()
	win.Push(Event{Type: EventText, Text: "á"})
	win.Input(20, 20, 200, 30, &text)
	if text != "olá" {
		t.Errorf("text = %q, want %q", text, "olá")
	}
}

func TestInputBackspaceDeletesWholeCharacters(t *testing.T) {
	win, _ := newTestWindow(t, 300, 100)
	text := "ação"

	click(win, 50, 30)
	win.Input(20, 20, 200, 30, &text)
	release(win, 50, 30)

	// Backspace takes the whole "o", which is one byte here.
	win.Begin()
	win.Push(Event{Type: EventKeyDown, Key: KeyBackspace})
	win.Input(20, 20, 200, 30, &text)
	if text != "açã" {
		t.Errorf("text = %q, want %q", text, "açã")
	}

	// And again over the two-byte "ã", which must go whole rather than
	// leaving half a character behind.
	win.Begin()
	win.Push(Event{Type: EventKeyDown, Key: KeyBackspace})
	win.Input(20, 20, 200, 30, &text)
	if text != "aç" {
		t.Errorf("text = %q, want %q", text, "aç")
	}
}

func TestInputCursorMovesByCharacters(t *testing.T) {
	win, _ := newTestWindow(t, 300, 100)
	text := "ação"

	click(win, 50, 30)
	win.Input(20, 20, 200, 30, &text)
	release(win, 50, 30)
	if win.uiCursor != len(text) {
		t.Fatalf("cursor = %d, want it at the end on focus", win.uiCursor)
	}

	// Left steps over the one-byte "o", landing after "açã".
	win.Begin()
	win.Push(Event{Type: EventKeyDown, Key: KeyLeft})
	win.Input(20, 20, 200, 30, &text)
	if win.uiCursor != len("açã") {
		t.Errorf("cursor = %d, want %d", win.uiCursor, len("açã"))
	}

	// And again over the two-byte "ã".
	win.Begin()
	win.Push(Event{Type: EventKeyDown, Key: KeyLeft})
	win.Input(20, 20, 200, 30, &text)
	if win.uiCursor != len("aç") {
		t.Errorf("cursor = %d, want %d", win.uiCursor, len("aç"))
	}

	win.Begin()
	win.Push(Event{Type: EventKeyDown, Key: KeyHome})
	win.Input(20, 20, 200, 30, &text)
	if win.uiCursor != 0 {
		t.Errorf("cursor = %d, want Home to take it to the start", win.uiCursor)
	}

	// Delete at the start takes the whole first character.
	win.Begin()
	win.Push(Event{Type: EventKeyDown, Key: KeyDelete})
	win.Input(20, 20, 200, 30, &text)
	if text != "ção" {
		t.Errorf("text = %q, want %q", text, "ção")
	}
}

func TestInputLosesFocusOnAClickOutside(t *testing.T) {
	win, _ := newTestWindow(t, 300, 200)
	text := ""

	click(win, 50, 30)
	win.Input(20, 20, 200, 30, &text)
	release(win, 50, 30)

	click(win, 250, 180)
	win.Input(20, 20, 200, 30, &text)
	if win.uiFocus != 0 {
		t.Error("clicking outside should take the focus away")
	}

	// With no focus, typing goes nowhere.
	win.Begin()
	win.Push(Event{Type: EventText, Text: "x"})
	if win.Input(20, 20, 200, 30, &text) || text != "" {
		t.Errorf("text = %q, want an unfocused field to ignore typing", text)
	}
}

func TestUTF8CursorSteps(t *testing.T) {
	const s = "aç€o" // one, two, three and one byte

	tests := []struct {
		name  string
		start int
		next  int
		prev  int
	}{
		{"ascii", 0, 1, 0},
		{"two-byte", 1, 3, 0},
		{"three-byte", 3, 6, 1},
		{"last", 6, 7, 3},
		{"end", 7, 7, 6},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := utf8Next(s, tt.start); got != tt.next {
				t.Errorf("utf8Next(%d) = %d, want %d", tt.start, got, tt.next)
			}
			if got := utf8Prev(s, tt.start); got != tt.prev {
				t.Errorf("utf8Prev(%d) = %d, want %d", tt.start, got, tt.prev)
			}
		})
	}
}

func TestWidgetIDIsStableAndDistinct(t *testing.T) {
	a := widgetID("button", 10, 20, 30, 40, "OK")
	if a != widgetID("button", 10, 20, 30, 40, "OK") {
		t.Error("the same widget must hash the same every frame")
	}
	for _, other := range []uint32{
		widgetID("button", 11, 20, 30, 40, "OK"),
		widgetID("button", 10, 20, 30, 40, "Cancel"),
		widgetID("checkbox", 10, 20, 30, 40, "OK"),
	} {
		if a == other {
			t.Error("widgets that differ must not share an id")
		}
	}
	if widgetID("", 0, 0, 0, 0, "") == 0 {
		t.Error("zero is the no-widget marker and must never be a real id")
	}
}
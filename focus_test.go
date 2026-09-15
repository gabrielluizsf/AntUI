package antui

import "testing"

// tabFrame presses Tab between Begin and End, so the widget draws, the tab
// order is known, and moveFocus in End walks it.
func tabFrame(win *Window) {
	win.Begin()
	win.Push(Event{Type: EventKeyDown, Key: KeyTab})
}

func TestTabMovesBetweenBuiltinInputs(t *testing.T) {
	win, _ := newTestWindow(t, 400, 120)
	title := ""
	body := ""

	// Focus the first field by clicking it.
	click(win, 50, 30)
	win.Input(20, 20, 200, 30, &title)
	release(win, 50, 30)
	if win.uiFocus != widgetID("input", 20, 20, 200, 30, "") {
		t.Fatal("the title field should have focus after clicking it")
	}

	// Tab goes to the second field, not back onto itself.
	tabFrame(win)
	win.Input(20, 20, 200, 30, &title)
	win.Input(20, 60, 200, 30, &body)
	win.End()
	if win.uiFocus != widgetID("input", 20, 60, 200, 30, "") {
		t.Fatalf("Tab should move focus to the body field, still on %d", win.uiFocus)
	}

	// Typing now lands in the body field, and the title is left alone.
	win.Begin()
	win.Push(Event{Type: EventText, Text: "x"})
	win.Input(20, 20, 200, 30, &title)
	win.Input(20, 60, 200, 30, &body)
	win.End()
	if title != "" || body != "x" {
		t.Errorf("typing after Tab went to wrong field: title=%q body=%q", title, body)
	}
}

func TestTabFocusesFirstWidgetWhenNothingFocused(t *testing.T) {
	win, _ := newTestWindow(t, 300, 100)
	clicked := false

	tabFrame(win)
	clicked = win.Button(20, 20, 100, 30, "OK")
	win.End()
	if clicked {
		t.Error("Tab alone must not click a button")
	}
	if win.uiFocus != widgetID("button", 20, 20, 100, 30, "OK") {
		t.Error("Tab should focus the first widget when nothing has focus")
	}
}

func TestEnterActivatesFocusedButton(t *testing.T) {
	win, _ := newTestWindow(t, 300, 100)

	// Focus the button with Tab, then press Enter in a later frame.
	tabFrame(win)
	win.Button(20, 20, 100, 30, "OK")
	win.End()
	if win.uiFocus == 0 {
		t.Fatal("Tab should have focused the button")
	}

	win.Begin()
	win.Push(Event{Type: EventKeyDown, Key: KeyEnter})
	clicked := win.Button(20, 20, 100, 30, "OK")
	win.End()
	if !clicked {
		t.Error("Enter should click the focused button")
	}
}

func TestSpaceTogglesFocusedCheckbox(t *testing.T) {
	win, _ := newTestWindow(t, 300, 100)
	value := false

	tabFrame(win)
	win.Checkbox(20, 20, "Enable", &value)
	win.End()
	if win.uiFocus == 0 {
		t.Fatal("Tab should have focused the checkbox")
	}

	win.Begin()
	win.Push(Event{Type: EventKeyDown, Key: KeySpace})
	win.Checkbox(20, 20, "Enable", &value)
	win.End()
	if !value {
		t.Error("Space should toggle the focused checkbox")
	}
}

func TestLeftRightArrowsStepFocusedSlider(t *testing.T) {
	win, _ := newTestWindow(t, 300, 100)
	value := float32(5)

	tabFrame(win)
	win.Slider(20, 30, 200, 20, &value, 0, 10)
	win.End()
	if win.uiFocus == 0 {
		t.Fatal("Tab should have focused the slider")
	}

	win.Begin()
	win.Push(Event{Type: EventKeyDown, Key: KeyRight})
	changed := win.Slider(20, 30, 200, 20, &value, 0, 10)
	win.End()
	if !changed || value != 6 {
		t.Errorf("right arrow should step up by 1/10, got value=%v changed=%v", value, changed)
	}
}

func TestShiftTabMovesBackwards(t *testing.T) {
	win, _ := newTestWindow(t, 400, 120)
	one := ""
	two := ""

	// Focus the second field directly.
	click(win, 50, 70)
	win.Input(20, 20, 200, 30, &one)
	win.Input(20, 60, 200, 30, &two)
	release(win, 50, 70)
	win.Input(20, 20, 200, 30, &one)
	win.Input(20, 60, 200, 30, &two)
	if win.uiFocus != widgetID("input", 20, 60, 200, 30, "") {
		t.Fatal("the second field should have focus after clicking it")
	}

	// Shift+Tab goes back to the first, and wraps when there is no previous.
	tabFr := func() {
		win.Begin()
		win.Push(Event{Type: EventKeyDown, Key: KeyTab, Mods: ModShift})
	}
	tabFr()
	win.Input(20, 20, 200, 30, &one)
	win.Input(20, 60, 200, 30, &two)
	win.End()
	if win.uiFocus != widgetID("input", 20, 20, 200, 30, "") {
		t.Errorf("Shift+Tab should move to the first field, still on %d", win.uiFocus)
	}

	tabFr()
	win.Input(20, 20, 200, 30, &one)
	win.Input(20, 60, 200, 30, &two)
	win.End()
	if win.uiFocus != widgetID("input", 20, 60, 200, 30, "") {
		t.Errorf("Shift+Tab should wrap to the last field, still on %d", win.uiFocus)
	}
}
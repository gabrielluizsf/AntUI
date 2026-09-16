package template

import (
	"testing"

	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/template/event"
)

func TestScaleFollowsTheSmallerEdge(t *testing.T) {
	tests := []struct {
		w, h int
		want int
	}{
		{200, 200, 1},   // below the reference
		{360, 360, 1},   // the reference itself
		{520, 420, 1},   // portrait, just over three hundred
		{720, 360, 1},   // the smallest edge is the reference
		{720, 720, 2},   // both edges twice the reference
		{960, 720, 2},   // 720/360
		{1080, 1920, 3}, // a tall phone
		{1440, 2880, 4}, // a phone at 4x
		{1970, 3000, 5}, // a gruesomely big display
	}
	for _, tt := range tests {
		win, _, err := antui.Offscreen(tt.w, tt.h)
		if err != nil {
			t.Fatalf("Offscreen(%d,%d): %v", tt.w, tt.h, err)
		}
		if got := Scale(win); got != tt.want {
			t.Errorf("Scale(%dx%d) = %d, want %d", tt.w, tt.h, got, tt.want)
		}
	}
}

// TestScaledScreensStillPaint draws a whole screen at scale 1 and at scale 2
// and makes sure the ink grows: the controls are drawn in scaled units, so a
// 2x window must show clearly more pixels than a 1x window of the same shape.
func TestScaledScreensStillPaint(t *testing.T) {
	count := func(w, h int) int {
		win, cv, err := antui.Offscreen(w, h)
		if err != nil {
			t.Fatal(err)
		}
		tpl := Cyberpunk(win)
		u := Scale(win)
		win.Begin()
		tpl.Button(win, 40*u, 40*u, 200*u, 80*u, "Launch")
		tpl.Input(win, 40*u, 160*u, 300*u, 80*u, new(string))
		win.End()
		n := 0
		for _, c := range cv.Pixels {
			if c != 0 {
				n++
			}
		}
		return n
	}
	one := count(360, 360) // scale 1
	two := count(720, 720) // scale 2
	if one == 0 || two == 0 {
		t.Fatal("a scaled screen painted nothing")
	}
	if two <= one {
		t.Errorf("scale 2 window painted no more ink than scale 1: %d vs %d", two, one)
	}
}

// TestTemplateTabMovesBetweenInputs is the regression for the bug where Tab,
// met by a tab order that only held the focused field, wrapped onto itself and
// put the caret at the start of the very field being edited.
func TestTemplateTabMovesBetweenInputs(t *testing.T) {
	win, _, err := antui.Offscreen(520, 420)
	if err != nil {
		t.Fatal(err)
	}
	tpl := Builtin(win)
	title := ""
	body := ""

	// Focus the title field by clicking it.
	win.Begin()
	testClick(t, win, 100, 100) // the "Title" input at 24,80 (scale 1)
	tpl.Input(win, 24, 80, 472, 32, &title)
	tpl.Input(win, 24, 154, 472, 32, &body)
	win.End()
	if win.WidgetFocus() == 0 {
		t.Fatal("the title field should have focus after clicking it")
	}
	if title != "" {
		t.Fatal("title should be empty before typing")
	}

	// Type into the title, then Tab to the body.
	win.Begin()
	win.Push(antui.Event{Type: antui.EventText, Text: "A"})
	win.Push(antui.Event{Type: antui.EventKeyDown, Key: antui.KeyTab})
	tpl.Input(win, 24, 80, 472, 32, &title)
	bodyEv := tpl.Input(win, 24, 154, 472, 32, &body)
	win.End()
	if title != "A" {
		t.Errorf("title = %q, want it to hold the typed text", title)
	}
	titleFocused := win.WidgetFocus() == win.WidgetID("template:input", 24, 80, 472, 32, "")
	if titleFocused {
		t.Error("Tab should have left the title field")
	}
	if win.WidgetFocus() != win.WidgetID("template:input", 24, 154, 472, 32, "") {
		t.Errorf("Tab should have focused the body field, got %d", win.WidgetFocus())
	}
	if bodyErr := checkNothing(bodyEv); bodyErr != "" {
		t.Error(bodyErr)
	}

	// And now typing lands in the body.
	win.Begin()
	win.Push(antui.Event{Type: antui.EventText, Text: "B"})
	tpl.Input(win, 24, 80, 472, 32, &title)
	bodyEv = tpl.Input(win, 24, 154, 472, 32, &body)
	win.End()
	if title != "A" || body != "B" {
		t.Errorf("typing after Tab went to the wrong field: title=%q body=%q", title, body)
	}
	if bodyEv.Ok() {
		t.Logf("body event: %v", bodyEv)
	}
}

func checkNothing(e event.Event) string {
	if e.Ok() {
		return "typing a key must not itself produce a Change event"
	}
	return ""
}

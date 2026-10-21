package template

import (
	"testing"

	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
)

// TestTemplateRadioClearsItsGroup checks the template's own radio: clicking
// the picked option empties the group, so a form can go back to "nothing
// chosen".
func TestTemplateRadioClearsItsGroup(t *testing.T) {
	win, _, _ := antui.Offscreen(200, 200)
	ut := cssUI(t, win, "")
	area := canvas.Area{X: 10, Y: 10, Width: 18, Height: 18}
	choice := 0

	win.Begin()
	testClick(t, win, 15, 15)
	ut.Radio(win, area, antui.RadioParams{Label: "Second", Value: &choice, Option: 1})
	win.End()
	if choice != 1 {
		t.Fatalf("choice = %d, want 1 after picking", choice)
	}

	win.Begin()
	testClick(t, win, 15, 15)
	ut.Radio(win, area, antui.RadioParams{Label: "Second", Value: &choice, Option: 1})
	win.End()
	if choice != antui.NoOption {
		t.Errorf("choice = %d, want %d after clearing", choice, antui.NoOption)
	}
}

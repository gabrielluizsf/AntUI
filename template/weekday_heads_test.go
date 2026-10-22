package template

import (
	"testing"

	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
)

// TestWeekdayHeadsStayFixedAcrossMonths guards a bug where the S M T W T F S
// column heads rotated to start at whichever weekday the month happened to
// begin on, instead of staying Sunday-through-Saturday every month the way
// every calendar a person has actually used does — with the day grid below
// it, not the heads above it, leaving empty columns before day 1.
func TestWeekdayHeadsStayFixedAcrossMonths(t *testing.T) {
	render := func(firstWD int) *canvas.Canvas {
		win, cv, err := antui.Offscreen(160, 50)
		if err != nil {
			t.Fatalf("offscreen: %v", err)
		}
		if !win.Begin() {
			t.Fatal("window closed")
		}
		win.Clear(canvas.White)
		p := datePaint{x: 0, y: 0, u: 1, firstWD: firstWD}
		paintWeekdayHeads(win.Canvas(), p, canvas.Color(0xFF000000))
		win.End()
		return cv
	}

	// A month starting on Sunday (firstWD=0) is the case the heads already
	// read correctly for either way; one starting on Thursday (firstWD=4) is
	// where a rotating header would visibly disagree with it.
	sunday, thursday := render(0), render(4)
	identical := true
	for y := 0; y < sunday.Height && identical; y++ {
		for x := 0; x < sunday.Width; x++ {
			if sunday.At(x, y) != thursday.At(x, y) {
				identical = false
				break
			}
		}
	}
	if !identical {
		t.Error("the weekday header drew differently for a month starting on Sunday than one starting on Thursday: the heads are rotating with firstWD instead of staying fixed")
	}
}

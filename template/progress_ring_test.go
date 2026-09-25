package template

import (
	"testing"

	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
)

// progressRingColors finds which RGBA channel dominates a pixel, so the tests
// can tell the red fill arc apart from the green track.
func progressDominant(c canvas.Color) byte {
	if c.A() == 0 {
		return 0
	}
	if c.R() >= c.G() && c.R() >= c.B() {
		return 'R'
	}
	if c.G() >= c.R() && c.G() >= c.B() {
		return 'G'
	}
	return 'B'
}

// TestCSSProgressBarStaysABar a rounded bar keeps drawing as a growing bar:
// a ring only appears when the box is a square whose radius rounds it fully.
func TestCSSProgressBarStaysABar(t *testing.T) {
	win, _, err := antui.Offscreen(200, 120)
	if err != nil {
		t.Fatal(err)
	}
	tpl := TemplateWithCSS(win)
	classes, path := cssTable(t, `
		progress {
			background-color: #ff0000;
			width: 160px;
			height: 16px;
			border-radius: 8px;
		}
	`)
	if err := tpl.SetStyle(path, classes); err != nil {
		t.Fatal(err)
	}

	tpl.Progress(0.5)
	// Bar centred horizontally; CSS flow starts at the top of the window.
	barX, barY := (200-160)/2, 0
	cv := win.Canvas()
	if got := progressDominant(cv.At(barX+40, barY+8)); got != 'R' {
		t.Errorf("filled half of the bar = %v, want the red fill", cv.At(barX+40, barY+8))
	}
	track := cv.At(barX+148, barY+8)
	if progressDominant(track) != 'R' || track.R() >= cv.At(barX+40, barY+8).R() {
		t.Errorf("track half = %v, want a darker red than the fill", track)
	}
}

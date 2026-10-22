package template

import (
	"testing"

	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
)

// TestSelectAndDateBoxDoNotPaintPastTheirBox guards a bug where a value too
// wide for its box — a long option, a long formatted date, or either at a
// large enough window scale — ran straight past the caret and out of the
// control instead of stopping at its own edge, the way Input already did.
func TestSelectAndDateBoxDoNotPaintPastTheirBox(t *testing.T) {
	long := "a value with far more text in it than this narrow box could ever hold"
	cases := []struct {
		name string
		fn   func(win *antui.Window, area canvas.Area)
	}{
		{"builtin select", func(win *antui.Window, area canvas.Area) {
			builtinStyle{}.Select(SelectParams{WidgetGeometry: WidgetGeometry{Win: win, X: area.X, Y: area.Y, W: area.Width, H: area.Height}, Value: long})
		}},
		{"builtin date box", func(win *antui.Window, area canvas.Area) {
			builtinStyle{}.DatePickerBox(DatePickerBoxParams{WidgetGeometry: WidgetGeometry{Win: win, X: area.X, Y: area.Y, W: area.Width, H: area.Height}, Value: long})
		}},
		{"cyberpunk select", func(win *antui.Window, area canvas.Area) {
			cyberStyle{}.Select(SelectParams{WidgetGeometry: WidgetGeometry{Win: win, X: area.X, Y: area.Y, W: area.Width, H: area.Height}, Value: long})
		}},
		{"cyberpunk date box", func(win *antui.Window, area canvas.Area) {
			cyberStyle{}.DatePickerBox(DatePickerBoxParams{WidgetGeometry: WidgetGeometry{Win: win, X: area.X, Y: area.Y, W: area.Width, H: area.Height}, Value: long})
		}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			win, cv, err := antui.Offscreen(400, 100)
			if err != nil {
				t.Fatalf("offscreen: %v", err)
			}
			if !win.Begin() {
				t.Fatal("window closed")
			}
			background := cv.At(0, 0)
			area := canvas.Area{X: 10, Y: 10, Width: 80, Height: 30}
			c.fn(win, area)
			win.End()

			for x := area.X + area.Width + 2; x < cv.Width; x++ {
				if got := cv.At(x, area.Y+area.Height/2); got != background {
					t.Fatalf("pixel at x=%d, past the box (ends at x=%d), is %v: the value painted over its neighbour instead of stopping at its own edge",
						x, area.X+area.Width, got)
				}
			}
		})
	}
}

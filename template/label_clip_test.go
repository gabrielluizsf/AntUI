package template

import (
	"testing"

	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
)

// TestLabelDoesNotPaintPastItsBox guards a bug where a label ignored its own
// width and height entirely: text long enough, or scaled up enough, to
// measure wider than the box it was given painted straight over whatever sat
// next to it instead of stopping at its own edge.
func TestLabelDoesNotPaintPastItsBox(t *testing.T) {
	for _, paint := range []struct {
		name string
		fn   func(LabelParams)
	}{
		{"builtin", builtinStyle{}.Label},
		{"cyberpunk", cyberStyle{}.Label},
	} {
		t.Run(paint.name, func(t *testing.T) {
			win, cv, err := antui.Offscreen(400, 100)
			if err != nil {
				t.Fatalf("offscreen: %v", err)
			}
			if !win.Begin() {
				t.Fatal("window closed")
			}
			background := cv.At(0, 0)
			area := canvas.Area{X: 10, Y: 10, Width: 60, Height: 20}
			paint.fn(LabelParams{
				WidgetGeometry: WidgetGeometry{Win: win, X: area.X, Y: area.Y, W: area.Width, H: area.Height},
				Text:           "a label with far more text in it than this narrow box could ever hold",
			})
			win.End()

			for x := area.X + area.Width + 2; x < cv.Width; x++ {
				if got := cv.At(x, area.Y+area.Height/2); got != background {
					t.Fatalf("pixel at x=%d, past the label's box (ends at x=%d), is %v: the label painted over its neighbour instead of stopping at its own edge",
						x, area.X+area.Width, got)
				}
			}
		})
	}
}

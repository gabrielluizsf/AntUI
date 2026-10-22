package template

import (
	"testing"

	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
)

// TestTextAreaPlaceholderDoesNotRepeatPerWrappedLine guards a bug where a
// text area with no text wrote the whole placeholder string again on every
// wrapped line, instead of that line's own slice of it — invisible while the
// placeholder fit on one line, and only showing up once a narrower or more
// heavily scaled box forced it onto two or more.
func TestTextAreaPlaceholderDoesNotRepeatPerWrappedLine(t *testing.T) {
	for _, paint := range []struct {
		name string
		fn   func(TextAreaParams)
	}{
		{"builtin", builtinStyle{}.TextArea},
		{"cyberpunk", cyberStyle{}.TextArea},
	} {
		t.Run(paint.name, func(t *testing.T) {
			win, cv, err := antui.Offscreen(400, 300)
			if err != nil {
				t.Fatalf("offscreen: %v", err)
			}
			if !win.Begin() {
				t.Fatal("window closed")
			}
			area := canvas.Area{X: 10, Y: 10, Width: 90, Height: 200}
			paint.fn(TextAreaParams{WidgetGeometry: WidgetGeometry{
				Win: win, X: area.X, Y: area.Y, W: area.Width, H: area.Height,
			}})
			win.End()

			u := Scale(win)
			padding := 8 * u
			lineH := textHeight(u)
			textW := max((area.Width-2*padding)/u, 12)
			lines := textLines("Escreva uma anotação...", textW)
			if len(lines) < 2 {
				t.Fatalf("test setup: the placeholder did not wrap (%d line(s)) — narrow the box further", len(lines))
			}

			band := func(li int) []canvas.Color {
				top := area.Y + padding + li*lineH
				row := make([]canvas.Color, (area.Width-2*padding)*lineH)
				i := 0
				for y := top; y < top+lineH; y++ {
					for x := area.X + padding; x < area.X+area.Width-padding; x++ {
						row[i] = cv.At(x, y)
						i++
					}
				}
				return row
			}
			first, second := band(0), band(1)
			identical := true
			for i := range first {
				if first[i] != second[i] {
					identical = false
					break
				}
			}
			if identical {
				t.Error("the second wrapped line of the placeholder is pixel-identical to the first: the whole placeholder is being redrawn on every line instead of its own wrapped slice")
			}
		})
	}
}

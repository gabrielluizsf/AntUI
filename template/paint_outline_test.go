package template

import (
	"testing"

	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

func TestSolidOutlineKeepsSquareCornersSquare(t *testing.T) {
	win, cv, err := antui.Offscreen(80, 60)
	if err != nil {
		t.Fatal(err)
	}
	st := css.Style{
		OutlineWidth:  2,
		OutlineStyle:  css.BorderSolid,
		OutlineColor:  canvas.RGBA(canvas.RGBAParams{R: 23, G: 23, B: 28, A: 36}),
		OutlineOffset: 1,
	}
	cs := &cssStyle{win: win}

	win.Begin()
	win.Clear(canvas.White)
	cs.paintOutline(win, st, 10, 10, 40, 20)
	win.End()

	if got := cv.At(7, 7); got.R() >= 250 {
		t.Fatalf("square outline rounded its outer corner: %v", got)
	}
}

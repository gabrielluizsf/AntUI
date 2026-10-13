package template

import (
	"testing"

	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

// The border painters take a boxSide instead of eight loose numbers, and the
// dashed and dotted cases now share one walk. These paint through the real
// canvas and check each style covers the side it is supposed to, side by side
// — a border drawn one pixel off is invisible to every other test here.

// painted runs one side through a painter and hands back the canvas it landed
// on, so a test can ask what was drawn where.
func painted(t *testing.T, s boxSide, style uint8) *canvas.Canvas {
	t.Helper()
	win, cv, err := antui.Offscreen(s.x+s.w+8, s.y+s.h+8)
	if err != nil {
		t.Fatal(err)
	}
	cs := &cssStyle{win: win}
	win.Begin()
	win.Clear(canvas.White)
	cs.paintSide(cv, s, style)
	win.End()
	return cv
}

// inked counts the pixels of a colour that differ from the white the canvas was
// cleared to.
func inked(cv *canvas.Canvas) int {
	n := 0
	for y := 0; y < cv.Height; y++ {
		for x := 0; x < cv.Width; x++ {
			if c := cv.At(x, y); c.R() < 250 {
				n++
			}
		}
	}
	return n
}

func TestSolidBorderCoversTheWholeSideOnEverySide(t *testing.T) {
	const bw = 3
	// The box is deliberately not square, so a side that measured the wrong
	// extent shows up instead of matching by accident.
	sides := []int{css.SideTop, css.SideRight, css.SideBottom, css.SideLeft}
	for _, side := range sides {
		s := boxSide{side: side, x: 5, y: 5, w: 60, h: 30, width: bw,
			colour: canvas.Black}
		cv := painted(t, s, css.BorderSolid)

		want := bw * 60 // the top and bottom sides run the box's width
		if side == css.SideRight || side == css.SideLeft {
			want = bw * 30 // the left and right sides run its height
		}
		if got := inked(cv); got != want {
			t.Errorf("side %d: painted %d pixels, want %d", side, got, want)
		}
	}
}

func TestDoubleBorderDrawsTwoLinesWithAGap(t *testing.T) {
	const bw = 6
	line := bw / 3
	s := boxSide{side: css.SideTop, x: 5, y: 5, w: 60, h: 30, width: bw,
		colour: canvas.Black}
	cv := painted(t, s, css.BorderDouble)

	// Two lines of the full width, with the gap between them empty.
	want := 2 * line * 60
	if got := inked(cv); got != want {
		t.Fatalf("double top: painted %d pixels, want %d", got, want)
	}
	// The first line is flush with the top of the box; the second sits a
	// stroke-width in from the bottom of the band.
	if c := cv.At(10, 5); c.R() >= 250 {
		t.Errorf("double top: outer line missing at the box's edge: %v", c)
	}
	if c := cv.At(10, 5+line); c.R() < 250 {
		t.Errorf("double top: the gap between the lines is painted: %v", c)
	}
	if c := cv.At(10, 5+bw-line); c.R() >= 250 {
		t.Errorf("double top: inner line missing: %v", c)
	}
}

func TestDashedBorderLeavesGapsAndDottedLeavesDots(t *testing.T) {
	const bw = 3
	mk := func() boxSide {
		return boxSide{side: css.SideTop, x: 5, y: 5, w: 60, h: 30, width: bw,
			colour: canvas.Black}
	}
	dash := painted(t, mk(), css.BorderDashed)
	dot := painted(t, mk(), css.BorderDotted)

	dashedPixels := inked(dash)
	dotPixels := inked(dot)
	if dashedPixels == 0 || dotPixels == 0 {
		t.Fatalf("nothing painted: dash %d, dot %d", dashedPixels, dotPixels)
	}
	// Both are broken up, so neither fills the side the way a solid one does.
	if solid := bw * 60; dashedPixels >= solid || dotPixels >= solid {
		t.Errorf("a broken border filled its side: dash %d, dot %d, solid %d",
			dashedPixels, dotPixels, solid)
	}
	// Dashes are rectangles and dots are discs, so a dash covers more of each
	// step it takes than a dot of the same thickness does.
	if dotPixels >= dashedPixels {
		t.Errorf("dotted border painted %d pixels, dashed only %d", dotPixels, dashedPixels)
	}
}

func TestDashAndDotNeverHangOffTheEndOfTheSide(t *testing.T) {
	// A side short enough that the first step nearly fits: the painter must stop
	// rather than draw past the box.
	const bw = 3
	for _, style := range []uint8{css.BorderDashed, css.BorderDotted} {
		for _, side := range []int{css.SideTop, css.SideRight, css.SideBottom, css.SideLeft} {
			s := boxSide{side: side, x: 5, y: 5, w: 9, h: 7, width: bw,
				colour: canvas.Black}
			cv := painted(t, s, style)
			for y := 0; y < cv.Height; y++ {
				for x := 0; x < cv.Width; x++ {
					if cv.At(x, y).R() >= 250 {
						continue
					}
					if x < s.x || y < s.y || x >= s.x+s.w || y >= s.y+s.h {
						t.Errorf("style %d side %d: painted outside the box at %d,%d",
							style, side, x, y)
					}
				}
			}
		}
	}
}

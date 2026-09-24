package template

import (
	"testing"

	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

// bgRule builds a computed style from one label rule, the direct way a
// painter test wants it without going through the flow.
func bgRule(t *testing.T, rule string) css.Style {
	t.Helper()
	sh, err := css.Parse("label {" + rule + "}")
	if err != nil {
		t.Fatal(err)
	}
	return sh.Style("label", nil, css.StateNone, 800)
}

// blankWin is a window cleared to a colour that nothing else will paint, so a
// test can tell "untouched" apart from "painted".
func blankWin(t *testing.T, w, h int) *antui.Window {
	t.Helper()
	win, _, err := antui.Offscreen(w, h)
	if err != nil {
		t.Fatal(err)
	}
	win.Canvas().FillRect(0, 0, w, h, canvas.RGBA(30, 40, 50, 255))
	return win
}

func TestCSSBackgroundGradientPaint(t *testing.T) {
	win := blankWin(t, 40, 40)
	st := bgRule(t, "background-image: linear-gradient(red, blue);")
	cs := newCSSStyle(win)
	cs.paintBackground(win, st, 10, 10, 20, 20, 0, 0)

	top := win.Canvas().At(20, 11)
	bottom := win.Canvas().At(20, 29)
	if top.R() <= top.B() || top.R() <= 200 {
		t.Errorf("top of the box should be red, got %v", top)
	}
	if bottom.B() <= bottom.R() || bottom.B() <= 200 {
		t.Errorf("bottom of the box should be blue, got %v", bottom)
	}
	if got := win.Canvas().At(2, 2); got != canvas.RGBA(30, 40, 50, 255) {
		t.Errorf("outside the box should be untouched, got %v", got)
	}
}

func TestCSSBackgroundColorUnderTranslucentGradient(t *testing.T) {
	win := blankWin(t, 40, 40)
	st := bgRule(t, "background-color: #FF0000; background-image: linear-gradient(rgba(0,0,255,0.55), rgba(0,0,255,0.55));")
	cs := newCSSStyle(win)
	cs.paintBackground(win, st, 10, 10, 20, 20, 0, 0)

	px := win.Canvas().At(20, 20)
	if px.R() == 0 || px.R() < 100 {
		t.Errorf("a translucent blue over the red colour must keep red showing, got %v", px)
	}
	if px.B() < 100 {
		t.Errorf("the blue layer must tint the surface, got %v", px)
	}
}

func TestCSSBackgroundImagePaint(t *testing.T) {
	win := blankWin(t, 40, 40)
	st := bgRule(t, "background-image: url(logo); background-repeat: no-repeat;")
	cs := newCSSStyle(win)
	logo, err := canvas.NewCanvas(5, 5)
	if err != nil {
		t.Fatal(err)
	}
	logo.FillRect(0, 0, 5, 5, canvas.Red)
	cs.images["logo"] = logo

	cs.paintBackground(win, st, 10, 10, 20, 20, 0, 0)
	if got := win.Canvas().At(12, 12); got != canvas.Red {
		t.Errorf("image pixel = %v, want red", got)
	}
	if got := win.Canvas().At(22, 22); got != canvas.RGBA(30, 40, 50, 255) {
		t.Errorf("a no-repeat image must not tile, got %v", got)
	}
}

func TestCSSBackgroundImageRepeats(t *testing.T) {
	win := blankWin(t, 40, 40)
	st := bgRule(t, "background-image: url(tile); background-repeat: repeat;")
	cs := newCSSStyle(win)
	tile, _ := canvas.NewCanvas(5, 5)
	tile.FillRect(0, 0, 5, 5, canvas.Red)
	cs.images["tile"] = tile

	cs.paintBackground(win, st, 10, 10, 21, 21, 0, 0)
	for _, where := range [][2]int{{12, 12}, {18, 12}, {12, 18}, {18, 18}, {29, 29}, {29, 12}} {
		if got := win.Canvas().At(where[0], where[1]); got != canvas.Red {
			t.Errorf("repeat tile at %v = %v, want red", where, got)
		}
	}
}

func TestCSSBackgroundPositionAndSize(t *testing.T) {
	win := blankWin(t, 40, 40)
	st := bgRule(t, "background-image: url(pic); background-repeat: no-repeat; background-size: 10px 10px; background-position: right 2px top 3px;")
	cs := newCSSStyle(win)
	pic, _ := canvas.NewCanvas(5, 5)
	pic.FillRect(0, 0, 5, 5, canvas.Red)
	cs.images["pic"] = pic
	u := cs.u()
	box := 20 * u

	cs.paintBackground(win, st, 10, 10, box, box, 0, 0)
	// right 2px top 3px puts a 10px tile 8 across and 3 down, all scaled.
	if got := win.Canvas().At(10+9*u, 10+4*u); got != canvas.Red {
		t.Errorf("inside the placed tile = %v, want red", got)
	}
	if got := win.Canvas().At(10, 10+3*u); got != canvas.RGBA(30, 40, 50, 255) {
		t.Errorf("left of the offset tile must stay blank, got %v", got)
	}
}

func TestCSSBackgroundRoundedCornersClipped(t *testing.T) {
	win := blankWin(t, 40, 40)
	st := bgRule(t, "background-image: url(pill); background-repeat: no-repeat; border-radius: 5px;")
	cs := newCSSStyle(win)
	img, _ := canvas.NewCanvas(20, 20)
	img.FillRect(0, 0, 20, 20, canvas.Red)
	cs.images["pill"] = img

	cs.paintBackground(win, st, 10, 10, 20, 20, 5, 5)
	// The corner of the box is outside the rounded clip: the gradient there
	// must have been cut away.
	if got := win.Canvas().At(11, 11); got == canvas.Red {
		t.Errorf("the rounded corner should be clipped to transparency, got %v", got)
	}
	if got := win.Canvas().At(15, 17); got != canvas.Red {
		t.Errorf("inside the rounded clip the layer should paint, got %v", got)
	}
}

func TestCSSBackgroundClipInsideBorder(t *testing.T) {
	win := blankWin(t, 40, 40)
	rule := "background-image: url(dot); background-repeat: no-repeat; background-clip: padding-box; border: 4px solid #000000;"
	cs := newCSSStyle(win)
	img, _ := canvas.NewCanvas(12, 12)
	img.FillRect(0, 0, 12, 12, canvas.Red)
	cs.images["dot"] = img
	st := bgRule(t, rule)

	// The image anchors to the padding box, so the border ring stays blank.
	cs.paintBackground(win, st, 10, 10, 20, 20, 0, 0)
	if got := win.Canvas().At(10, 12); got != canvas.RGBA(30, 40, 50, 255) {
		t.Errorf("under the border the clip must clear, got %v", got)
	}
	if got := win.Canvas().At(14, 14); got != canvas.Red {
		t.Errorf("inside the padding box the image paints, got %v", got)
	}
}

func TestCSSBackgroundLayersStackAlpha(t *testing.T) {
	win := blankWin(t, 40, 40)
	// The first layer is topmost: an opaque red under a translucent blue, so
	// the top colour adds blue to the red beneath only when the stacking is
	// right.
	st := bgRule(t, `background-image:
		linear-gradient(rgba(0, 0, 255, 0.55), rgba(0, 0, 255, 0.55)),
		linear-gradient(rgba(255, 0, 0, 1), rgba(255, 0, 0, 1));`)
	cs := newCSSStyle(win)
	cs.paintBackground(win, st, 10, 10, 20, 20, 0, 0)

	px := win.Canvas().At(18, 18)
	// The translucent blue top layer blends 55% over the opaque red beneath:
	// red stays at ~115 and blue lands at ~140. Reversed stacking would read
	// either pure blue or pure red, never the mix.
	if px.R() < 110 || px.R() > 125 || px.B() < 128 || px.B() > 160 {
		t.Errorf("stacking layers should blend in declaration order, got %v", px)
	}
}

func TestCSSBackgroundUnknownImagePaintsNothing(t *testing.T) {
	win := blankWin(t, 40, 40)
	st := bgRule(t, "background-image: url(missing);")
	cs := newCSSStyle(win)
	cs.paintBackground(win, st, 10, 10, 20, 20, 0, 0)
	if got := win.Canvas().At(15, 15); got != canvas.RGBA(30, 40, 50, 255) {
		t.Errorf("an unregistered image must paint nothing, got %v", got)
	}
}

func TestCSSBackgroundConicAndRadialPaint(t *testing.T) {
	for _, rule := range []string{
		"background-image: radial-gradient(red, blue);",
		"background-image: conic-gradient(red, blue);",
	} {
		win := blankWin(t, 40, 40)
		st := bgRule(t, rule)
		cs := newCSSStyle(win)
		cs.paintBackground(win, st, 10, 10, 20, 20, 0, 0)

		px := win.Canvas().At(20, 20)
		if px.A() == 0 && px == canvas.RGBA(30, 40, 50, 255) {
			t.Errorf("%s painted nothing where a gradient should sit", rule)
		}
	}
}

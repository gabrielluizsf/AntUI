package template

import (
	"testing"

	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
	"github.com/gabrielluizsf/antui/template/event"
)

func TestCSSTransformTranslates(t *testing.T) {
	win := blankWin(t, 40, 40)
	st := bgRule(t, "transform: translate(10px 0); background-color: red;")
	cs := newCSSStyle(win)
	cs.transformed(st, 10, 10, 20, 20, func() {
		cs.paintBackground(win, st, 10, 10, 20, 20, 0, 0)
	})

	if got := win.Canvas().At(25, 15); got != canvas.RGBA(255, 0, 0, 255) {
		t.Errorf("pixel inside the moved box = %v, want red", got)
	}
	if got := win.Canvas().At(10, 15); got != canvas.RGBA(30, 40, 50, 255) {
		t.Errorf("the box must leave its old spot blank, got %v", got)
	}
}

func TestCSSTransformScales(t *testing.T) {
	win := blankWin(t, 40, 40)
	// Scaling from the top-left corner doubles the box outward: 10..30
	// becomes 10..50, hungrier than the window accepts but still clipped.
	st := bgRule(t, "transform: scale(2); transform-origin: left top; background-color: red;")
	cs := newCSSStyle(win)
	cs.transformed(st, 10, 10, 20, 20, func() {
		cs.paintBackground(win, st, 10, 10, 20, 20, 0, 0)
	})

	for _, p := range [][2]int{{10, 10}, {20, 20}, {38, 38}} {
		if got := win.Canvas().At(p[0], p[1]); got != canvas.RGBA(255, 0, 0, 255) {
			t.Errorf("scaled pixel at %v = %v, want red", p, got)
		}
	}
	if got := win.Canvas().At(8, 8); got != canvas.RGBA(30, 40, 50, 255) {
		t.Errorf("before the scaled box = %v, want blank", got)
	}
}

func TestCSSTransformDefaultOriginIsCentre(t *testing.T) {
	win := blankWin(t, 40, 40)
	// A 30-wide box rotated a quarter turn keeps its centre still and spins
	// the corners: the original top-left corner leaves (10,10). A pivot at
	// the corner would keep that pixel red, so blank proves the default.
	st := bgRule(t, "transform: rotate(90deg); background-color: red;")
	cs := newCSSStyle(win)
	cs.transformed(st, 10, 10, 30, 10, func() {
		cs.paintBackground(win, st, 10, 10, 30, 10, 0, 0)
	})

	if got := win.Canvas().At(25, 15); got != canvas.RGBA(255, 0, 0, 255) {
		t.Errorf("centre of the rotated box = %v, want red", got)
	}
	if got := win.Canvas().At(10, 10); got == canvas.RGBA(255, 0, 0, 255) {
		t.Error("the corner must rotate away from the pivot, still red at (10,10)")
	}
	if got := win.Canvas().At(25, 20); got != canvas.RGBA(255, 0, 0, 255) {
		t.Errorf("under the rotated box = %v, want red", got)
	}
}

func TestCSSTransformWithoutTransformPaintsDirect(t *testing.T) {
	win := blankWin(t, 40, 40)
	st := bgRule(t, "background-color: red;")
	cs := newCSSStyle(win)
	cs.transformed(st, 10, 10, 20, 20, func() {
		cs.paintBackground(win, st, 10, 10, 20, 20, 0, 0)
	})

	if got := win.Canvas().At(15, 15); got != canvas.RGBA(255, 0, 0, 255) {
		t.Errorf("a box without a transform must paint in place, got %v", got)
	}
}

func TestCSSTransformPadCoversEffects(t *testing.T) {
	win := blankWin(t, 40, 40)
	st := bgRule(t, "box-shadow: 0 4px 0 3px black; outline: 2px solid black; outline-offset: 1px;")
	cs := newCSSStyle(win)
	// The outmost reach is the shadow: 3 spread + 4 offset + blur 0 + the 2
	// pixel sampling margin. The outline's 3 is small fry next to it.
	if got := cs.transformPad(st); got != 9 {
		t.Errorf("transformPad = %d, want 9", got)
	}

	shadow := bgRule(t, "transform: translate(0 0); box-shadow: 0 4px 0 3px black; background-color: red;")
	cs.transformed(shadow, 10, 10, 20, 20, func() {
		cs.paintBackground(win, shadow, 10, 10, 20, 20, 0, 0)
	})
	// The shadow pool sits just below the box, offset by the pad beyond it.
	if got := win.Canvas().At(15, 10+4+3-1); got == canvas.RGBA(30, 40, 50, 255) {
		t.Errorf("the box-shadow must follow the box through the composit, got blank")
	}
}

func TestCSSTransformButtonThroughPaint(t *testing.T) {
	win, _, err := antui.Offscreen(400, 200)
	if err != nil {
		t.Fatal(err)
	}
	tpl := TemplateWithCSS(win)
	classes, path := cssTable(t, "button { transform: translate(30px 0); background-color: red; }")
	if err := tpl.SetStyle(path, classes); err != nil {
		t.Fatal(err)
	}
	u := Scale(win)

	win.Begin()
	tpl.Button("Go")
	win.End()

	// The button painted as the first flow slot. Deciding its box by asking
	// for a second slot would measure the slot below, so recover the slot
	// it used from the difference.
	x, y, _, h, _, ok := tpl.layout(css.RoleButton, "Go")
	if !ok {
		t.Fatal("button should be laid out")
	}
	y = y - h

	if got := win.Canvas().At(x+10, y+h/2); got == canvas.RGBA(30, 40, 50, 255) {
		t.Errorf("the original button spot must be blank, got %v", got)
	}
	if got := win.Canvas().At(x+30*u+2, y+h/2); got != canvas.RGBA(255, 0, 0, 255) {
		t.Errorf("the button must paint 30px right, got %v", got)
	}
}

func TestCSSTransformCompositesDrawnOverspill(t *testing.T) {
	win := blankWin(t, 40, 40)
	// A select's menu paints below the box, well past where any shadow pad
	// reaches. It must ride the composite instead of being cut at the pad.
	st := bgRule(t, "transform: translate(10px 0);")
	cs := newCSSStyle(win)
	cs.transformed(st, 10, 10, 20, 20, func() {
		win.Canvas().FillRect(10, 20, 15, 8, canvas.RGBA(255, 0, 0, 255))
	})

	// The menu block lands 10px right, unmangled, through the transform.
	if got := win.Canvas().At(25, 26); got != canvas.RGBA(255, 0, 0, 255) {
		t.Errorf("the overspill must come over with the box, got %v", got)
	}
	if got := win.Canvas().At(25, 22); got != canvas.RGBA(255, 0, 0, 255) {
		t.Errorf("the overspill interior must come over too, got %v", got)
	}
}

func TestCSSTransformShadowStaysSoftInLayer(t *testing.T) {
	win := blankWin(t, 40, 40)
	// The transform records the whole widget on a transparent layer; the
	// box-shadow's blur is stamped there first. Alpha-forcing blending used
	// to make the shadow's whole padded layer opaque, so it came back as a
	// solid black slab around the box — bigger than the box.
	st := bgRule(t, "transform: translate(0px 0); box-shadow: 0 0 8px black; background-color: red;")
	cs := newCSSStyle(win)
	cs.transformed(st, 10, 10, 20, 20, func() {
		cs.paintBox(win, st, 10, 10, 20, 20)
	})
	// A point well left of the box, inside the composite pad but past any
	// real glow, must keep the window's colour — never opaque black.
	for _, p := range [][2]int{{2, 20}, {6, 20}} {
		if got := win.Canvas().At(p[0], p[1]); got == canvas.RGBA(0, 0, 0, 255) {
			t.Errorf("the shadow smeared as an opaque slab at %v, got solid black", p)
		}
	}
}

func TestCSSTransformBackdropSeesPageBehind(t *testing.T) {
	win := blankWin(t, 40, 40)
	for y := 10; y < 20; y++ {
		win.Canvas().FillRect(10, y, 20, 1, canvas.RGBA(255, 0, 0, 255))
	}
	st := bgRule(t, "transform: translate(10px 10px); backdrop-filter: invert(1);")
	cs := newCSSStyle(win)
	cs.transformed(st, 0, 0, 20, 10, func() {
		cs.paintBox(win, st, 0, 0, 20, 10)
	})

	// The box lands over the red page. Its backdrop inverts it: cyan here
	// proves the filter saw the window behind, not an empty layer.
	if got := win.Canvas().At(25, 15); got != canvas.RGBA(0, 255, 255, 255) {
		t.Errorf("the backdrop must invert the window behind the box, got %v", got)
	}
}

func TestCSSTransformHitTestFollowsBox(t *testing.T) {
	win, _, err := antui.Offscreen(400, 200)
	if err != nil {
		t.Fatal(err)
	}
	tpl := TemplateWithCSS(win)
	classes, path := cssTable(t, "button { transform: translate(30px 0); background-color: red; }")
	if err := tpl.SetStyle(path, classes); err != nil {
		t.Fatal(err)
	}
	u := Scale(win)

	// The button the frame draws is the flow slot under the probe: probing
	// drains a slot, the frame's button lands directly below it.
	box := func() (int, int) {
		x, y, w, h, _, ok := tpl.layout(css.RoleButton, "Go")
		if !ok {
			t.Fatal("button should be laid out")
		}
		return x + w/2, y + h + h/2
	}

	wx, wy := box()
	win.Begin()
	testClick(t, win, wx+30*u, wy)
	e1 := tpl.Button("Go")
	win.End()
	if !e1.Is(event.Button, event.Click) {
		t.Errorf("a click where the box moved to must hit, got %v", e1)
	}

	wx, wy = box()
	win.Begin()
	testClick(t, win, wx, wy)
	e2 := tpl.Button("Go")
	win.End()
	if e2.Ok() {
		t.Errorf("a click on the untransformed spot must miss, got %v", e2)
	}
}

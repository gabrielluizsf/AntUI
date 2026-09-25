package template

import (
	"testing"

	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
	"github.com/gabrielluizsf/antui/template/event"
)

// box is one box the draw pass hands back for a flex item.
type box struct{ x, y, w, h int }

// flexRun draws the widgets the callback asks for inside a Flex block on an
// offscreen window and returns the boxes the callback stored in dst, skipping
// the boxes the silent measure pass returns.
func flexRun(t *testing.T, sheet string, draw func(*CSS, *[]box)) []box {
	t.Helper()
	win, _, err := antui.Offscreen(400, 200)
	if err != nil {
		t.Fatal(err)
	}
	tpl := TemplateWithCSS(win)
	classes, path := cssTable(t, sheet)
	if err := tpl.SetStyle(path, classes); err != nil {
		t.Fatal(err)
	}
	var got []box
	win.Begin()
	tpl.Flex(func(f *CSS) {
		draw(f, &got)
	})
	win.End()
	return got
}

// flexLayout tracks a button's box into dst when the draw pass hands it back.
func flexLayout(dst *[]box) func(*CSS, string) {
	return func(f *CSS, label string) {
		x, y, w, h, _, ok := f.layout(css.RoleButton, label)
		if ok {
			*dst = append(*dst, box{x, y, w, h})
		}
	}
}

func TestFlexRowPlacesLeftToRight(t *testing.T) {
	got := flexRun(t, "flex { display: flex; }", func(f *CSS, dst *[]box) {
		lay := flexLayout(dst)
		lay(f, "A")
		lay(f, "B")
	})
	if len(got) != 2 {
		t.Fatalf("got %d boxes, want 2", len(got))
	}
	if got[0].x != 0 || got[0].y != 0 {
		t.Errorf("first item at %d,%d want 0,0", got[0].x, got[0].y)
	}
	if got[1].x != got[0].w {
		t.Errorf("second item x=%d want %d (packed right after the first)", got[1].x, got[0].w)
	}
	if got[1].y != 0 {
		t.Errorf("second item y=%d want 0", got[1].y)
	}
	if got[0].w != got[0].w-got[0].x+got[0].x {
		t.Errorf("first box should be a sane width, got %d", got[0].w)
	}
}

// flexLayout2 draws two buttons, A and B, capturing both boxes into dst.
func flexLayout2(dst *[]box) func(*CSS) {
	lay := flexLayout(dst)
	return func(f *CSS) {
		lay(f, "A")
		lay(f, "B")
	}
}

func TestFlexGapSpacesItems(t *testing.T) {
	got := flexRun(t, "flex { display: flex; gap: 10px; }", func(f *CSS, dst *[]box) {
		flexLayout2(dst)(f)
	})
	if got[1].x != got[0].w+10 {
		t.Errorf("second item x=%d want %d (10px gap)", got[1].x, got[0].w+10)
	}
}

func TestFlexJustifyCenterPacksMiddle(t *testing.T) {
	got := flexRun(t, "flex { display: flex; justify-content: center; }", func(f *CSS, dst *[]box) {
		flexLayout2(dst)(f)
	})
	left := (400 - got[0].w - got[1].w) / 2
	if got[0].x != left {
		t.Errorf("first item x=%d want %d (centred)", got[0].x, left)
	}
	if got[1].x != got[0].x+got[0].w {
		t.Errorf("second item x=%d want %d", got[1].x, got[0].x+got[0].w)
	}
}

func TestFlexJustifySpaceBetween(t *testing.T) {
	got := flexRun(t, "flex { display: flex; justify-content: space-between; }", func(f *CSS, dst *[]box) {
		flexLayout2(dst)(f)
	})
	if got[0].x != 0 {
		t.Errorf("first item x=%d want 0", got[0].x)
	}
	if got[1].x+got[1].w != 400 {
		t.Errorf("second item right edge=%d want 400 (space-between pushes to the far edge)", got[1].x+got[1].w)
	}
}

func TestFlexGrowSharesFreeSpace(t *testing.T) {
	got := flexRun(t, "flex { display: flex; } button { flex-grow: 1; }", func(f *CSS, dst *[]box) {
		flexLayout2(dst)(f)
	})
	if got[0].w != got[1].w {
		t.Errorf("equal grow should give equal widths, got %d vs %d", got[0].w, got[1].w)
	}
	if got[0].w+got[1].w != 400 {
		t.Errorf("grown widths should fill the container, got %d+%d", got[0].w, got[1].w)
	}
	if got[0].x != 0 || got[1].x != got[0].w {
		t.Errorf("grown items misplaced: %d,%d want 0,%d", got[0].x, got[1].x, got[0].w)
	}
}

func TestFlexColumnStacksTopDown(t *testing.T) {
	got := flexRun(t, "flex { display: flex; flex-direction: column; }", func(f *CSS, dst *[]box) {
		flexLayout2(dst)(f)
	})
	if got[1].y != got[0].y+got[0].h {
		t.Errorf("column second item y=%d want %d", got[1].y, got[0].y+got[0].h)
	}
	if got[1].x != got[0].x {
		t.Errorf("column items share the x, got %d vs %d", got[1].x, got[0].x)
	}
}

func TestFlexRowReverseMirrors(t *testing.T) {
	got := flexRun(t, "flex { display: flex; flex-direction: row-reverse; }", func(f *CSS, dst *[]box) {
		flexLayout2(dst)(f)
	})
	if got[0].x+got[0].w != 400 {
		t.Errorf("first-drawn item should sit against the right edge, right=%d want 400", got[0].x+got[0].w)
	}
	if got[1].x != got[0].x-got[1].w {
		t.Errorf("second item x=%d want %d (to the left of the first)", got[1].x, got[0].x-got[1].w)
	}
}

func TestFlexAlignItemsCenter(t *testing.T) {
	got := flexRun(t, `
		flex { display: flex; align-items: center; }
		button { height: 30px; }
		progress { height: 50px; }
	`, func(f *CSS, dst *[]box) {
		x, y, w, h, _, ok := f.layout(css.RoleButton, "A")
		if ok {
			*dst = append(*dst, box{x, y, w, h})
		}
		x, y, w, h, _, ok = f.layout(css.RoleProgress, "")
		if ok {
			*dst = append(*dst, box{x, y, w, h})
		}
	})
	if len(got) != 2 {
		t.Fatalf("got %d boxes, want 2", len(got))
	}
	cb, pb := got[0], got[1]
	cbC := cb.y + cb.h/2
	pbC := pb.y + pb.h/2
	if pb.y != cbC-pb.h/2 {
		t.Errorf("progress should centre on the button: centre %d want %d", pbC, cbC)
	}
}

func TestFlexAlignSelfOverrides(t *testing.T) {
	got := flexRun(t, `
		flex { display: flex; align-items: center; }
		button { height: 30px; }
		progress { height: 50px; align-self: flex-end; }
	`, func(f *CSS, dst *[]box) {
		x, y, w, h, _, ok := f.layout(css.RoleButton, "A")
		if ok {
			*dst = append(*dst, box{x, y, w, h})
		}
		x, y, w, h, _, ok = f.layout(css.RoleProgress, "")
		if ok {
			*dst = append(*dst, box{x, y, w, h})
		}
	})
	if len(got) != 2 {
		t.Fatalf("got %d boxes, want 2", len(got))
	}
	cb, pb := got[0], got[1]
	cross := pb.h
	if pb.y != cross-pb.h {
		t.Errorf("align-self flex-end y=%d want %d (flush with the line bottom)", pb.y, cross-pb.h)
	}
	if cb.y != (cross-cb.h)/2 {
		t.Errorf("button should stay centre, y=%d want %d (cross=%d)", cb.y, (cross-cb.h)/2, cross)
	}
}

func TestFlexWrapMovesToNextLine(t *testing.T) {
	got := flexRun(t, `
		flex { display: flex; flex-wrap: wrap; gap: 4px; }
		button { width: 150px; }
	`, func(f *CSS, dst *[]box) {
		lay := flexLayout(dst)
		lay(f, "A")
		lay(f, "B")
		lay(f, "C")
	})
	if len(got) != 3 {
		t.Fatalf("got %d boxes, want 3", len(got))
	}
	if got[1].y != 0 {
		t.Errorf("second item should share the first line, y=%d", got[1].y)
	}
	if got[2].y == got[0].y {
		t.Error("third item should wrap to a new line")
	}
	if got[2].x != 0 {
		t.Errorf("wrapped item should start a fresh line at x=0, got %d", got[2].x)
	}
}

func TestFlexOrderReorders(t *testing.T) {
	got := map[string]box{}
	flexRun(t, `
		flex { display: flex; }
		button { order: 2; }
		progress { order: 1; }
	`, func(f *CSS, _ *[]box) {
		if x, y, w, h, _, ok := f.layout(css.RoleButton, "A"); ok {
			got["button"] = box{x, y, w, h}
		}
		if x, y, w, h, _, ok := f.layout(css.RoleProgress, ""); ok {
			got["progress"] = box{x, y, w, h}
		}
	})
	if len(got) != 2 {
		t.Fatalf("got %d boxes, want 2", len(got))
	}
	if got["progress"].x != 0 {
		t.Errorf("order 1 should draw first at x=0, got %d", got["progress"].x)
	}
	if got["button"].x != got["progress"].w {
		t.Errorf("order 2 should draw after order 1: x=%d want %d", got["button"].x, got["progress"].w)
	}
}

func TestFlexContainerFlowsBelow(t *testing.T) {
	win, _, _ := antui.Offscreen(400, 200)
	tpl := TemplateWithCSS(win)
	classes, path := cssTable(t, "flex { display: flex; }")
	if err := tpl.SetStyle(path, classes); err != nil {
		t.Fatal(err)
	}
	win.Begin()
	blockH := 0
	tpl.Flex(func(f *CSS) {
		_, _, _, h, _, ok := f.layout(css.RoleButton, "A")
		if ok {
			blockH = h
		}
	})
	_, y, _, _, _, ok := tpl.layout(css.RoleLabel, "after")
	win.End()
	if !ok {
		t.Fatal("label should be visible")
	}
	if blockH == 0 {
		t.Fatal("flex block should have a height")
	}
	if y != blockH {
		t.Errorf("label should flow right under the flex block, y=%d want %d", y, blockH)
	}
}

func TestFlexContainerBoxPaints(t *testing.T) {
	win, _, _ := antui.Offscreen(400, 200)
	tpl := TemplateWithCSS(win)
	classes, path := cssTable(t, `flex { display: flex; background-color: #ff0000; }`)
	if err := tpl.SetStyle(path, classes); err != nil {
		t.Fatal(err)
	}
	win.Begin()
	tpl.Flex(func(f *CSS) { f.layout(css.RoleButton, "A") })
	win.End()
	if got := win.Canvas().At(10, 10); got != canvas.RGBA(255, 0, 0, 255) {
		t.Errorf("container background pixel = %v, want red", got)
	}
}

func TestFlexButtonClickInside(t *testing.T) {
	win, _, _ := antui.Offscreen(400, 200)
	tpl := TemplateWithCSS(win)
	classes, path := cssTable(t, "flex { display: flex; }")
	if err := tpl.SetStyle(path, classes); err != nil {
		t.Fatal(err)
	}
	u := Scale(win)
	w := textWidth(u, "Go")
	win.Begin()
	testClick(t, win, w/2, textHeight(u)/2)
	var e event.Event
	tpl.Flex(func(f *CSS) {
		e = f.Button("Go")
	})
	win.End()
	if !e.Is(event.Button, event.Click) {
		t.Errorf("click inside a flexed button must produce a Click, got %v", e)
	}
}

func TestFlexShrinksToFitLine(t *testing.T) {
	got := flexRun(t, `
		flex { display: flex; }
		button { width: 250px; }
	`, func(f *CSS, dst *[]box) {
		flexLayout2(dst)(f)
	})
	if len(got) != 2 {
		t.Fatalf("got %d boxes, want 2", len(got))
	}
	total := got[0].w + got[1].w
	if total != 400 {
		t.Errorf("two 250px items shrink to fit 400px, got %d", total)
	}
}

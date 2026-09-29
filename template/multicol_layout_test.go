package template

import (
	"testing"

	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

// widget is one thing a layout test asks to be drawn: the role that styles it
// and the text it is identified by.
type widget struct {
	role, text string
}

func labels(n int) []widget {
	out := make([]widget, n)
	for i := range out {
		out[i] = widget{role: css.RoleLabel, text: string(rune('A' + i))}
	}
	return out
}

// multiRun draws the widgets it is given inside a MultiCol block on an
// offscreen window and returns the boxes the draw pass handed back, skipping
// the ones the silent measure pass returns.
func multiRun(t *testing.T, w, h int, sheet string, want ...widget) []box {
	t.Helper()
	win, _, err := antui.Offscreen(w, h)
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
	tpl.MultiCol(func(m *CSS) {
		for _, it := range want {
			x, y, w, h, _, ok := m.layout(it.role, it.text)
			if ok {
				got = append(got, box{x, y, w, h})
			}
		}
	})
	win.End()
	return got
}

// TestMultiColCountSplitsTheContentInColumns is the plain case: a count of two
// and four children of one height puts two children in each column, and the
// columns are half the container wide.
func TestMultiColCountSplitsTheContentInColumns(t *testing.T) {
	got := multiRun(t, 400, 200, "multicolumn { column-count: 2; } label { height: 20px; }", labels(4)...)
	if len(got) != 4 {
		t.Fatalf("got %d boxes, want 4", len(got))
	}
	for i, want := range []box{
		{x: 0, y: 0, w: 200, h: 20},
		{x: 0, y: 20, w: 200, h: 20},
		{x: 200, y: 0, w: 200, h: 20},
		{x: 200, y: 20, w: 200, h: 20},
	} {
		if got[i] != want {
			t.Errorf("box %d = %+v want %+v", i, got[i], want)
		}
	}
}

// TestMultiColGapIsTheGutter: column-gap is the same property flex and grid
// read, so the gutter is taken off the columns before they are shared out.
func TestMultiColGapIsTheGutter(t *testing.T) {
	got := multiRun(t, 400, 200, "multicolumn { column-count: 2; column-gap: 20px; } label { height: 10px; }", labels(2)...)
	if len(got) != 2 {
		t.Fatalf("got %d boxes, want 2", len(got))
	}
	if got[0].w != 190 || got[1].w != 190 {
		t.Errorf("column width = %d and %d, want 190 and 190", got[0].w, got[1].w)
	}
	if got[1].x != 210 {
		t.Errorf("second column at x=%d want 210, a 20px gutter after 190", got[1].x)
	}
}

// TestMultiColWidthDecidesTheCount: with no count, the width a column asks for
// says how many of them fit the container.
func TestMultiColWidthDecidesTheCount(t *testing.T) {
	got := multiRun(t, 400, 200, "multicolumn { column-width: 100px; } label { height: 10px; }", labels(4)...)
	if len(got) != 4 {
		t.Fatalf("got %d boxes, want 4", len(got))
	}
	for i, at := range []int{0, 100, 200, 300} {
		if got[i].x != at || got[i].w != 100 {
			t.Errorf("box %d at x=%d width=%d, want x=%d width=100", i, got[i].x, got[i].w, at)
		}
	}
}

// TestMultiColCountAndWidthTakeTheSmaller: a count of three with a width only
// two of which fit makes two columns, the way a browser does.
func TestMultiColCountAndWidthTakeTheSmaller(t *testing.T) {
	got := multiRun(t, 400, 200, "multicolumn { column-count: 3; column-width: 200px; } label { height: 10px; }", labels(2)...)
	if len(got) != 2 {
		t.Fatalf("got %d boxes, want 2", len(got))
	}
	if got[0].x != 0 || got[1].x != 200 {
		t.Errorf("columns at x=%d and %d, want 0 and 200: two of them fit", got[0].x, got[1].x)
	}
	if got[0].w != 200 {
		t.Errorf("column width = %d want 200", got[0].w)
	}
}

// TestMultiColWithoutACountStaysOneColumn: a column block with nothing to say
// about its columns is the block it would have been.
func TestMultiColWithoutACountStaysOneColumn(t *testing.T) {
	got := multiRun(t, 400, 200, "multicolumn { } label { height: 10px; }", labels(2)...)
	if len(got) != 2 {
		t.Fatalf("got %d boxes, want 2", len(got))
	}
	if got[0] != (box{x: 0, y: 0, w: 400, h: 10}) || got[1] != (box{x: 0, y: 10, w: 400, h: 10}) {
		t.Errorf("boxes %+v %+v want one full-width column", got[0], got[1])
	}
}

// TestMultiColBalanceGrowsToTheTallestColumn: a container with no height of its
// own is as tall as its columns, and the content is shared out evenly — four
// children of one height in two columns make two columns two children tall.
func TestMultiColBalanceGrowsToTheTallestColumn(t *testing.T) {
	win, _, err := antui.Offscreen(400, 200)
	if err != nil {
		t.Fatal(err)
	}
	tpl := TemplateWithCSS(win)
	classes, path := cssTable(t, "multicolumn { column-count: 2; } label { height: 20px; }")
	if err := tpl.SetStyle(path, classes); err != nil {
		t.Fatal(err)
	}
	win.Begin()
	tpl.MultiCol(func(m *CSS) {
		m.Label("one")
		m.Label("two")
		m.Label("three")
	})
	// A label after the block rests below it, so where the flow cursor is
	// afterwards is the height of the columns.
	_, after, _, _, _, _ := tpl.layout(css.RoleLabel, "after")
	win.End()
	if after != 40 {
		t.Errorf("the flow resumed at y=%d, want 40: two columns of two labels", after)
	}
}

// TestMultiColDeclaredHeightCutsABox: a container with a height has one to cut
// a box against, so a child too tall for a column is cut and continued in the
// next one — one draw pass per fragment, so the pass hands back both.
func TestMultiColDeclaredHeightCutsABox(t *testing.T) {
	got := multiRun(t, 400, 200, "multicolumn { column-count: 2; height: 100px; column-fill: auto; } label { height: 150px; }", labels(1)...)
	if len(got) != 2 {
		t.Fatalf("got %d fragments, want 2", len(got))
	}
	if got[0] != (box{x: 0, y: 0, w: 200, h: 100}) {
		t.Errorf("first fragment %+v want the whole column", got[0])
	}
	if got[1] != (box{x: 200, y: 0, w: 200, h: 50}) {
		t.Errorf("second fragment %+v want the rest of the box in the next column", got[1])
	}
}

// TestMultiColBreaksBeforeAnAvoidedBox: a box that asked not to be cut starts
// the next column whole instead of being split, and the box before it is the
// only thing in the column it had room for.
func TestMultiColBreaksBeforeAnAvoidedBox(t *testing.T) {
	sheet := "multicolumn { column-count: 2; height: 100px; column-fill: auto; } " +
		"label { height: 60px; } button { height: 60px; break-inside: avoid; }"
	got := multiRun(t, 400, 200, sheet,
		widget{css.RoleLabel, "first"},
		widget{css.RoleButton, "second"},
	)
	if len(got) != 2 {
		t.Fatalf("got %d boxes, want 2", len(got))
	}
	if got[0] != (box{x: 0, y: 0, w: 200, h: 60}) {
		t.Errorf("first box %+v want the top of the first column", got[0])
	}
	if got[1] != (box{x: 200, y: 0, w: 200, h: 60}) {
		t.Errorf("the uncut box %+v want the whole next column", got[1])
	}
}

// TestMultiColContentPastTheLastColumnHangsBelow: the column count is what the
// container asked for, not a hint, so content that does not fit them all ends
// in the last one and hangs below the container. The boxes come back one pass
// at a time — the first fragment of each box, then the second — because a
// fragment is what a pass of the callback hands back.
func TestMultiColContentPastTheLastColumnHangsBelow(t *testing.T) {
	got := multiRun(t, 400, 200, "multicolumn { column-count: 2; height: 60px; column-fill: auto; } label { height: 100px; }", labels(2)...)
	for i, want := range []box{
		{x: 0, y: 0, w: 200, h: 60},
		{x: 200, y: 40, w: 200, h: 20},
		{x: 200, y: 0, w: 200, h: 40},
		{x: 200, y: 60, w: 200, h: 80},
	} {
		if i >= len(got) {
			break
		}
		if got[i] != want {
			t.Errorf("box %d = %+v want %+v", i, got[i], want)
		}
	}
	if len(got) != 4 {
		t.Errorf("got %d fragments, want 4", len(got))
	}
}

// TestMultiColRuleIsPaintedInTheGutter: the rule is a border's three
// declarations, drawn in the middle of the gutter between two columns and no
// wider than the gutter.
func TestMultiColRuleIsPaintedInTheGutter(t *testing.T) {
	win := blankWin(t, 200, 60)
	tpl := TemplateWithCSS(win)
	classes, path := cssTable(t, "multicolumn { column-count: 2; column-gap: 20px; column-rule: 2px solid rgb(255,0,0); } label { height: 40px; }")
	if err := tpl.SetStyle(path, classes); err != nil {
		t.Fatal(err)
	}
	win.Begin()
	tpl.MultiCol(func(m *CSS) {
		m.Label("a")
		m.Label("b")
	})
	win.End()

	rule := win.Canvas().At(99, 10)
	if rule.R() < 200 || rule.G() != 0 {
		t.Errorf("the gutter at x=99 is %v, want the red rule", rule)
	}
	if edge := win.Canvas().At(95, 10); edge != canvas.RGBA(30, 40, 50, 255) {
		t.Errorf("the rule must not spill out of the gutter, got %v at x=95", edge)
	}
	if below := win.Canvas().At(99, 50); below != canvas.RGBA(30, 40, 50, 255) {
		t.Errorf("the rule must stop with the column, got %v below it", below)
	}
}

// TestMultiColRulesStandInEveryGutter: a third column adds a second gutter, and
// its rule stands in that gutter — a column and a half along, not a gap further
// on, which would put it inside the column beside.
func TestMultiColRulesStandInEveryGutter(t *testing.T) {
	win := blankWin(t, 220, 60)
	tpl := TemplateWithCSS(win)
	classes, path := cssTable(t, "multicolumn { column-count: 3; column-gap: 20px; column-rule: 2px solid rgb(255,0,0); } label { height: 40px; }")
	if err := tpl.SetStyle(path, classes); err != nil {
		t.Fatal(err)
	}
	win.Begin()
	tpl.MultiCol(func(m *CSS) {
		for i := 0; i < 3; i++ {
			m.Label("a")
		}
	})
	win.End()

	cv := win.Canvas()
	// Three columns of 60px in a 220px box: gutters of 20px at x=60 and x=140,
	// so the rules sit at x=69 and x=149.
	for _, x := range []int{69, 149} {
		if rule := cv.At(x, 10); rule.R() < 200 || rule.G() != 0 {
			t.Errorf("the gutter at x=%d is %v, want the red rule", x, rule)
		}
	}
	// Inside the middle column, where a rule stepped by the gap alone would be.
	if inside := cv.At(99, 10); inside == canvas.RGBA(255, 0, 0, 255) {
		t.Errorf("a rule at x=99 stands inside a column, not in a gutter")
	}
}

// TestMultiColInAMultiColIsAnItem: a column block drawn inside another one is an
// item of a column, with the columns of its own inside the width that column has.
func TestMultiColInAMultiColIsAnItem(t *testing.T) {
	win, _, err := antui.Offscreen(400, 200)
	if err != nil {
		t.Fatal(err)
	}
	tpl := TemplateWithCSS(win)
	classes, path := cssTable(t, "multicolumn { column-count: 2; } button { height: 10px; }")
	if err := tpl.SetStyle(path, classes); err != nil {
		t.Fatal(err)
	}
	var got []box
	win.Begin()
	tpl.MultiCol(func(outer *CSS) {
		outer.MultiCol(func(inner *CSS) {
			for i := range 2 {
				if x, y, w, h, _, ok := inner.layout(css.RoleButton, string(rune('A'+i))); ok {
					got = append(got, box{x, y, w, h})
				}
			}
		})
	})
	win.End()
	if len(got) != 2 {
		t.Fatalf("got %d boxes, want 2", len(got))
	}
	if got[0].w != 100 || got[1].x != 100 {
		t.Errorf("inner columns %+v %+v want two of 100 in a 200px column", got[0], got[1])
	}
}

// TestMultiColInAFlexItemIsAnItem: a column block drawn in a flex container is an
// item of it, and its own columns answer to the width the flex solver gave it.
func TestMultiColInAFlexItemIsAnItem(t *testing.T) {
	win, _, err := antui.Offscreen(200, 200)
	if err != nil {
		t.Fatal(err)
	}
	tpl := TemplateWithCSS(win)
	classes, path := cssTable(t, "flex { display: flex; } multicolumn { column-count: 2; } button { height: 10px; }")
	if err := tpl.SetStyle(path, classes); err != nil {
		t.Fatal(err)
	}
	var got []box
	win.Begin()
	tpl.Flex(func(f *CSS) {
		f.MultiCol(func(m *CSS) {
			for i := range 2 {
				if x, y, w, h, _, ok := m.layout(css.RoleButton, string(rune('A'+i))); ok {
					got = append(got, box{x, y, w, h})
				}
			}
		})
	})
	win.End()
	if len(got) != 2 {
		t.Fatalf("got %d boxes, want 2", len(got))
	}
	if got[0].w != 100 || got[1].x != 100 {
		t.Errorf("columns %+v %+v want two of 100 in a 200px item", got[0], got[1])
	}
}


// TestMultiColInAFlexItemStartsInsideIt: a column block placed by a parent is
// drawn where the parent put it, so the boxes its children get back carry the
// parent's y and not only the offset inside the block.
func TestMultiColInAFlexItemStartsInsideIt(t *testing.T) {
	win, _, err := antui.Offscreen(200, 200)
	if err != nil {
		t.Fatal(err)
	}
	tpl := TemplateWithCSS(win)
	classes, path := cssTable(t, `
		flex { display: flex; flex-direction: column; }
		label { height: 30px; }
		multicolumn { column-count: 1; padding: 5px; }
		button { height: 10px; }
	`)
	if err := tpl.SetStyle(path, classes); err != nil {
		t.Fatal(err)
	}
	var got []box
	win.Begin()
	tpl.Flex(func(f *CSS) {
		f.Label("above")
		f.MultiCol(func(m *CSS) {
			for i := range 2 {
				if x, y, w, h, _, ok := m.layout(css.RoleButton, string(rune('A'+i))); ok {
					got = append(got, box{x, y, w, h})
				}
			}
		})
	})
	win.End()
	if len(got) != 2 {
		t.Fatalf("got %d boxes, want 2", len(got))
	}
	// 30 for the label above, then the block's own 5px padding.
	for i, at := range got {
		if at.y != 35+i*10 {
			t.Errorf("box %d at y=%d want %d: inside the block, not at the top of the window", i, at.y, 35+i*10)
		}
	}
}

// TestMultiColWidthAnswersToTheScale: a window past the scale base is drawn
// larger, and a column-width is a length like any other, so the same sheet asks
// for a column twice as many pixels wide and the columns of the two windows are
// the same width to the eye.
func TestMultiColWidthAnswersToTheScale(t *testing.T) {
	sheet := "multicolumn { column-width: 100px; } label { height: 10px; }"
	small := multiRun(t, 400, 200, sheet, labels(4)...)
	big := multiRun(t, 900, 760, sheet, labels(4)...)
	if len(small) != 4 || len(big) != 4 {
		t.Fatalf("got %d and %d boxes, want 4 each", len(small), len(big))
	}
	if small[0].w != 100 {
		t.Errorf("a 400px window: column %d wide, want 100", small[0].w)
	}
	// The width a column asks for is a minimum and the columns share what is
	// left over, so four of them divide the 900px window: 225 each where 200
	// is the minimum. Reading the length without the scale would ask for 100px
	// columns and cut nine of them out of the same window.
	if big[0].w != 225 {
		t.Errorf("a 900x760 window: column %d wide, want 225: 200 is the minimum, the rest is shared", big[0].w)
	}
	for i, at := range []int{0, 225, 450, 675} {
		if big[i].x != at {
			t.Errorf("column %d at x=%d want %d", i, big[i].x, at)
		}
	}
}

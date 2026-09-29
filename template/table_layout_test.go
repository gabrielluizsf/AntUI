package template

import (
	"testing"

	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

// tableRun draws a table of the given rows on an offscreen window of the given
// size and returns the box of every cell, in the order the draw pass handed them
// back: the first cell of every row, then the second, and so on. A row is the
// list of roles its cells carry.
func tableRun(t *testing.T, winW, winH int, sheet string, rows ...[]string) []box {
	t.Helper()
	win, _, err := antui.Offscreen(winW, winH)
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
	tpl.Table(func(tb *CSS) {
		for _, cells := range rows {
			tb.TableRow(func(r *CSS) {
				for i, role := range cells {
					x, y, w, h, _, ok := r.layout(role, string(rune('A'+i)))
					if ok {
						got = append(got, box{x, y, w, h})
					}
				}
			})
		}
	})
	win.End()
	return got
}

func TestTableRowsLineUpInColumns(t *testing.T) {
	got := tableRun(t, 100, 200,
		"table { } label { width: 50px; height: 20px; }",
		[]string{css.RoleLabel, css.RoleLabel},
		[]string{css.RoleLabel, css.RoleLabel},
	)
	if len(got) != 4 {
		t.Fatalf("got %d cells, want 4", len(got))
	}
	for i, want := range []box{
		{x: 0, y: 0, w: 50, h: 20},
		{x: 50, y: 0, w: 50, h: 20},
		{x: 0, y: 20, w: 50, h: 20},
		{x: 50, y: 20, w: 50, h: 20},
	} {
		if got[i] != want {
			t.Errorf("cell %d = %+v want %+v", i, got[i], want)
		}
	}
}

// TestTableWidestCellWinsItsColumn: a column is as wide as the widest cell that
// lands in it, over every row, and the rows below it start in the same column.
func TestTableWidestCellWinsItsColumn(t *testing.T) {
	got := tableRun(t, 200, 200,
		"table { } label { width: 50px; height: 20px; } button { width: 90px; height: 20px; }",
		[]string{css.RoleLabel, css.RoleButton},
		[]string{css.RoleLabel},
	)
	if len(got) != 3 {
		t.Fatalf("got %d cells, want 3", len(got))
	}
	for i, want := range []box{
		{x: 30, y: 0, w: 50, h: 20},  // the first column
		{x: 80, y: 0, w: 90, h: 20},  // the second is as wide as its widest cell
		{x: 30, y: 20, w: 50, h: 20}, // the next row reuses the columns
	} {
		if got[i] != want {
			t.Errorf("cell %d = %+v want %+v", i, got[i], want)
		}
	}
}

// TestTableRowIsAsTallAsItsTallestCell: a cell with a height of its own gives
// the row its height, and a cell that asked for the top takes the whole row.
func TestTableRowIsAsTallAsItsTallestCell(t *testing.T) {
	got := tableRun(t, 200, 200,
		"table { } label { width: 50px; height: 20px; } button { width: 50px; height: 60px; }",
		[]string{css.RoleLabel, css.RoleButton},
	)
	if len(got) != 2 {
		t.Fatalf("got %d cells, want 2", len(got))
	}
	for i, want := range []box{
		{x: 50, y: 0, w: 50, h: 60},
		{x: 100, y: 0, w: 50, h: 60},
	} {
		if got[i] != want {
			t.Errorf("cell %d = %+v want %+v: the row is as tall as its tallest cell", i, got[i], want)
		}
	}
}

// TestTableDeclaredWidthIsSharedOverTheColumns: a table that declares a width
// shares it out over the columns in proportion to what they asked for, so a wide
// column keeps more of the extra room than a narrow one.
func TestTableDeclaredWidthIsSharedOverTheColumns(t *testing.T) {
	got := tableRun(t, 400, 200,
		"table { width: 200px; } label { width: 50px; height: 10px; } button { width: 150px; height: 10px; }",
		[]string{css.RoleLabel, css.RoleButton},
	)
	if len(got) != 2 {
		t.Fatalf("got %d cells, want 2", len(got))
	}
	for i, want := range []box{
		{x: 100, y: 0, w: 50, h: 10},
		{x: 150, y: 0, w: 150, h: 10},
	} {
		if got[i] != want {
			t.Errorf("cell %d = %+v want %+v", i, got[i], want)
		}
	}
}

// TestTableStaysCentred: a table narrower than the flow is centred in it, the
// way any block box is.
func TestTableStaysCentred(t *testing.T) {
	got := tableRun(t, 400, 200,
		"table { } label { width: 100px; height: 10px; }",
		[]string{css.RoleLabel},
	)
	if len(got) != 1 {
		t.Fatalf("got %d cells, want 1", len(got))
	}
	if got[0] != (box{x: 150, y: 0, w: 100, h: 10}) {
		t.Errorf("cell %+v want the 100px table centred in 400", got[0])
	}
}

// TestTableVerticalAlignPlacesTheContentInItsRow: a cell that asks to sit in
// the middle or at the bottom of its row takes only the height it needs and is
// placed there, while a cell that asked for the top takes the whole row.
func TestTableVerticalAlignPlacesTheContentInItsRow(t *testing.T) {
	for _, tt := range []struct {
		align string
		y, h  int
	}{
		{align: "", y: 0, h: 60},
		{align: "vertical-align: middle;", y: 20, h: 20},
		{align: "vertical-align: bottom;", y: 40, h: 20},
		{align: "vertical-align: sub;", y: 20, h: 20},
	} {
		got := tableRun(t, 200, 200, "table { } label { width: 50px; height: 60px; } button { width: 50px; height: 20px; "+tt.align+" }",
			[]string{css.RoleLabel, css.RoleButton})
		if len(got) != 2 {
			t.Fatalf("%q: got %d cells, want 2", tt.align, len(got))
		}
		if got[1] != (box{x: 100, y: tt.y, w: 50, h: tt.h}) {
			t.Errorf("%q: cell %+v want y=%d h=%d", tt.align, got[1], tt.y, tt.h)
		}
	}
}

// TestTableWidgetIsARowOfItsOwn: a widget drawn straight in a table, without a
// row around it, becomes a row of one cell — what a browser's anonymous boxes do
// with a child that did not ask to be a row.
func TestTableWidgetIsARowOfItsOwn(t *testing.T) {
	win, _, err := antui.Offscreen(100, 200)
	if err != nil {
		t.Fatal(err)
	}
	tpl := TemplateWithCSS(win)
	classes, path := cssTable(t, "table { } label { width: 50px; height: 20px; }")
	if err := tpl.SetStyle(path, classes); err != nil {
		t.Fatal(err)
	}
	var got []box
	win.Begin()
	tpl.Table(func(tb *CSS) {
		for i := range 2 {
			x, y, w, h, _, ok := tb.layout(css.RoleLabel, string(rune('A'+i)))
			if ok {
				got = append(got, box{x, y, w, h})
			}
		}
	})
	win.End()
	if len(got) != 2 {
		t.Fatalf("got %d cells, want 2", len(got))
	}
	for i, want := range []box{{x: 25, y: 0, w: 50, h: 20}, {x: 25, y: 20, w: 50, h: 20}} {
		if got[i] != want {
			t.Errorf("cell %d = %+v want %+v: a row each", i, got[i], want)
		}
	}
}

// TestTableCaptionSitsAboveTheRows: a widget the stylesheet marks as a caption
// is not a cell: it is given the table's whole width and the run above the rows,
// and the rows start below it.
func TestTableCaptionSitsAboveTheRows(t *testing.T) {
	win, _, err := antui.Offscreen(100, 200)
	if err != nil {
		t.Fatal(err)
	}
	tpl := TemplateWithCSS(win)
	classes, path := cssTable(t, "table { } label { display: table-caption; height: 20px; } button { width: 50px; height: 30px; }")
	if err := tpl.SetStyle(path, classes); err != nil {
		t.Fatal(err)
	}
	var got []box
	win.Begin()
	tpl.Table(func(tb *CSS) {
		if x, y, w, h, _, ok := tb.layout(css.RoleLabel, "caption"); ok {
			got = append(got, box{x, y, w, h})
		}
		tb.TableRow(func(r *CSS) {
			if x, y, w, h, _, ok := r.layout(css.RoleButton, "A"); ok {
				got = append(got, box{x, y, w, h})
			}
		})
	})
	win.End()
	if len(got) != 2 {
		t.Fatalf("got %d boxes, want 2", len(got))
	}
	if got[0] != (box{x: 25, y: 0, w: 50, h: 20}) {
		t.Errorf("caption %+v want the 20px above the rows", got[0])
	}
	if got[1] != (box{x: 25, y: 20, w: 50, h: 30}) {
		t.Errorf("cell %+v want to start below the caption", got[1])
	}
}

// TestTablePaintsItsOwnBox: a table with a background paints the face a block
// would, around the cells the rows drew.
func TestTablePaintsItsOwnBox(t *testing.T) {
	win := blankWin(t, 200, 60)
	tpl := TemplateWithCSS(win)
	classes, path := cssTable(t, "table { background-color: rgb(0,255,0); width: 100px; } button { width: 100px; height: 40px; }")
	if err := tpl.SetStyle(path, classes); err != nil {
		t.Fatal(err)
	}
	win.Begin()
	tpl.Table(func(tb *CSS) {
		tb.TableRow(func(r *CSS) { r.Button("a") })
	})
	win.End()
	if px := win.Canvas().At(100, 5); px != canvas.RGBA(0, 255, 0, 255) {
		t.Errorf("the table face at the top of the box is %v, want the green background", px)
	}
	if px := win.Canvas().At(5, 5); px != canvas.RGBA(30, 40, 50, 255) {
		t.Errorf("outside the table the window is %v, want it untouched", px)
	}
}

// TestTableInACellIsACell: a table drawn inside a cell is that cell's content.
// A table is as wide as its own columns ask for, so it keeps the width it
// measured and centres in the cell the way a block box centres in the flow.
func TestTableInACellIsACell(t *testing.T) {
	win, _, err := antui.Offscreen(200, 200)
	if err != nil {
		t.Fatal(err)
	}
	tpl := TemplateWithCSS(win)
	classes, path := cssTable(t,
		"table { } label { width: 200px; height: 100px; } button { width: 40px; height: 10px; }")
	if err := tpl.SetStyle(path, classes); err != nil {
		t.Fatal(err)
	}
	var got []box
	win.Begin()
	tpl.Table(func(outer *CSS) {
		outer.TableRow(func(r *CSS) {
			r.Table(func(inner *CSS) {
				inner.TableRow(func(c *CSS) {
					if x, y, w, h, _, ok := c.layout(css.RoleButton, "A"); ok {
						got = append(got, box{x, y, w, h})
					}
				})
			})
		})
	})
	win.End()
	if len(got) != 1 {
		t.Fatalf("got %d cells, want 1", len(got))
	}
	if got[0].w != 40 {
		t.Errorf("the inner cell is %d wide, want the 40px its column asks for", got[0].w)
	}
	if got[0].x != 80 {
		t.Errorf("the inner cell is at x=%d, want 80: a 40px table centred in a 200px cell", got[0].x)
	}
}

// TestTableInAFlexItemIsAnItem: a table drawn in a flex container is an item of
// it, placed by the flex solver, and it lines its own rows up inside the box it
// was given.
func TestTableInAFlexItemIsAnItem(t *testing.T) {
	win, _, err := antui.Offscreen(200, 200)
	if err != nil {
		t.Fatal(err)
	}
	tpl := TemplateWithCSS(win)
	classes, path := cssTable(t,
		"flex { display: flex; } label { width: 60px; height: 20px; } button { width: 30px; height: 10px; }")
	if err := tpl.SetStyle(path, classes); err != nil {
		t.Fatal(err)
	}
	var got []box
	win.Begin()
	tpl.Flex(func(f *CSS) {
		f.Table(func(tb *CSS) {
			tb.TableRow(func(r *CSS) {
				if x, y, w, h, _, ok := r.layout(css.RoleButton, "A"); ok {
					got = append(got, box{x, y, w, h})
				}
			})
		})
		f.Label("after")
	})
	win.End()
	if len(got) != 1 {
		t.Fatalf("got %d cells, want 1", len(got))
	}
	if got[0] != (box{x: 0, y: 0, w: 30, h: 10}) {
		t.Errorf("the cell of a table in a flex item %+v want the 30px column at the top", got[0])
	}
}

// TestTableInAColumnIsAnItem: a table drawn in a multi-column container is an
// item of one column — as wide as its own cells, not as wide as the column.
func TestTableInAColumnIsAnItem(t *testing.T) {
	win, _, err := antui.Offscreen(400, 200)
	if err != nil {
		t.Fatal(err)
	}
	tpl := TemplateWithCSS(win)
	classes, path := cssTable(t,
		"multicolumn { column-count: 2; } button { width: 30px; height: 10px; }")
	if err := tpl.SetStyle(path, classes); err != nil {
		t.Fatal(err)
	}
	var got []box
	win.Begin()
	tpl.MultiCol(func(m *CSS) {
		for i := range 2 {
			m.Table(func(tb *CSS) {
				tb.TableRow(func(r *CSS) {
					if x, y, w, h, _, ok := r.layout(css.RoleButton, string(rune('A'+i))); ok {
						got = append(got, box{x, y, w, h})
					}
				})
			})
		}
	})
	win.End()
	if len(got) != 2 {
		t.Fatalf("got %d cells, want 2", len(got))
	}
	// Each table is as wide as its own cell asks for and centres in the 200px
	// column it landed in, so the two sit 85px in from either column edge.
	for i, want := range []box{{x: 85, y: 0, w: 30, h: 10}, {x: 285, y: 0, w: 30, h: 10}} {
		if got[i] != want {
			t.Errorf("cell %d = %+v want %+v", i, got[i], want)
		}
	}
}

// TestTableWithAPercentageWidthAsksForItsColumns: a table whose width is a
// percentage has nothing definite to resolve it against while its parent is still
// measuring, so it says what its own columns add up to. That is what sizes the
// cell holding it: without it a nested table claims the whole window, the column
// of the cell around it grows to match and the cells beside it are squeezed into
// what is left.
func TestTableWithAPercentageWidthAsksForItsColumns(t *testing.T) {
	win, _, err := antui.Offscreen(400, 200)
	if err != nil {
		t.Fatal(err)
	}
	tpl := TemplateWithCSS(win)
	classes, path := cssTable(t, `
		table { width: 100%; }
		label { width: 100px; height: 20px; }
		button { width: 30px; height: 10px; }
	`)
	if err := tpl.SetStyle(path, classes); err != nil {
		t.Fatal(err)
	}
	var got []box
	win.Begin()
	tpl.Table(func(outer *CSS) {
		outer.TableRow(func(r *CSS) {
			if x, y, w, h, _, ok := r.layout(css.RoleLabel, "outer"); ok {
				got = append(got, box{x, y, w, h})
			}
			r.Table(func(inner *CSS) {
				inner.TableRow(func(c *CSS) {
					if x, y, w, h, _, ok := c.layout(css.RoleButton, "A"); ok {
						got = append(got, box{x, y, w, h})
					}
				})
			})
		})
	})
	win.End()
	if len(got) != 2 {
		t.Fatalf("got %d boxes, want 2: the label cell and the cell of the nested table", len(got))
	}
	label, inner := got[0], got[1]
	// The table is as wide as the window, so the cell holding the nested table
	// is what the row has left of the row after the label's column.
	cell := box{x: label.x + label.w, w: 400 - label.x - label.w}
	if label.x != 0 {
		t.Errorf("the label column starts at x=%d want 0: the table is as wide as the window", label.x)
	}
	// The label asked for 100px and the nested table for the 30px of its own
	// column, and the row shares its 400px in proportion, so the label's column
	// keeps more than three times what the cell around the table got. A nested
	// table that claimed the whole window while being measured would take the row
	// instead, and the label would be left with a fifth of it.
	if label.w < 3*cell.w {
		t.Errorf("the label column is %d wide and the cell around the nested table %d, "+
			"want the label's more than three times as wide", label.w, cell.w)
	}
	// The table answers its percentage to the cell it was given this time, which
	// is the whole of it.
	if inner.x != cell.x || inner.w != cell.w {
		t.Errorf("the nested table is at x=%d width=%d, want it filling the cell at x=%d width=%d",
			inner.x, inner.w, cell.x, cell.w)
	}
}

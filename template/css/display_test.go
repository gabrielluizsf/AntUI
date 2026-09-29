package css

import "testing"

// TestDisplayKeepsTheTableInteriorsApart pins the one thing the table layout
// leans on: a stylesheet that says display: table-cell is heard as a cell, not
// folded into display: table the way a browser does, where anonymous boxes make
// the difference invisible.
func TestDisplayKeepsTheTableInteriorsApart(t *testing.T) {
	for raw, want := range map[string]uint8{
		"block":              DisplayBlock,
		"flow-root":          DisplayBlock,
		"list-item":          DisplayBlock,
		"none":               DisplayNone,
		"inline":             DisplayInline,
		"inline-block":       DisplayInlineBlock,
		"inline-flex":        DisplayInlineBlock,
		"inline-grid":        DisplayInlineBlock,
		"flex":               DisplayFlex,
		"grid":               DisplayGrid,
		"table":              DisplayTable,
		"inline-table":       DisplayInlineTable,
		"table-row":          DisplayTableRow,
		"table-row-group":    DisplayTableRowGroup,
		"table-header-group": DisplayTableHeaderGroup,
		"table-footer-group": DisplayTableFooterGroup,
		"table-cell":         DisplayTableCell,
		"table-caption":      DisplayTableCaption,
		"table-column":       DisplayTableColumn,
		"table-column-group": DisplayTableColumn,
	} {
		v, ok := parseDisplay(raw)
		if !ok || v != want {
			t.Errorf("%q = %d,%v want %d,true", raw, v, ok, want)
		}
	}
	if _, ok := parseDisplay("table-rowgroup"); ok {
		t.Error(`"table-rowgroup" is not a display keyword`)
	}
}

// TestDisplayTableExtrasOnStyle reaches the same keywords the way a template
// reads them: off a computed style, as the helpers that name the boxes.
func TestDisplayTableExtrasOnStyle(t *testing.T) {
	for _, raw := range []string{"table-row", "table-row-group", "table-header-group", "table-footer-group"} {
		st := flexRule(t, "label", "display: "+raw+";")
		if !st.TableBox() {
			t.Errorf("%q must be a table box", raw)
		}
		if st.Cell() || st.Caption() {
			t.Errorf("%q must not read as a cell or a caption", raw)
		}
	}
	st := flexRule(t, "label", "display: table-cell;")
	if !st.Cell() || st.TableBox() || st.Caption() {
		t.Error("table-cell must read as a cell and nothing else")
	}
	st = flexRule(t, "label", "display: table-caption;")
	if !st.Caption() || st.Cell() {
		t.Error("table-caption must read as a caption and nothing else")
	}
	st = flexRule(t, "label", "display: table;")
	if st.TableBox() || st.Cell() || st.Caption() {
		t.Error("display: table is the container, not one of the boxes inside it")
	}
}

// TestInlineTableFlowsAlongTheLine pins the display value an inline table is
// given: it stays on the line, the way inline-block does, and it is not
// mistaken for a plain block.
func TestInlineTableFlowsAlongTheLine(t *testing.T) {
	st := flexRule(t, "label", "display: inline-table;")
	if !st.Inline() {
		t.Error("inline-table must flow along a line")
	}
	if st.Display == DisplayInlineBlock {
		t.Error("inline-table must keep its own display value")
	}
	st = flexRule(t, "label", "display: inline-block;")
	if !st.Inline() {
		t.Error("inline-block must flow along a line")
	}
	st = flexRule(t, "label", "display: table;")
	if st.Inline() {
		t.Error("display: table starts a line of its own")
	}
}

package template

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
	"github.com/gabrielluizsf/antui/template/event"
)

// cssTable writes a stylesheet to a temp file and loads it into a fresh,
// empty class table. It returns the table and the file to pass to
// TemplateWithCSS.SetStyle.
func cssTable(t *testing.T, sheet string) (*css.CSSClasses, string) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "app.css")
	if err := os.WriteFile(path, []byte(sheet), 0o644); err != nil {
		t.Fatal(err)
	}
	classes := css.NewTable()
	if err := classes.SetStyle(path); err != nil {
		t.Fatal(err)
	}
	return classes, path
}

func TestCSSButtonClicked(t *testing.T) {
	win, _, err := antui.Offscreen(400, 200)
	if err != nil {
		t.Fatal(err)
	}
	tpl := TemplateWithCSS(win)
	classes, path := cssTable(t, "")
	if err := tpl.SetStyle(path, classes); err != nil {
		t.Fatal(err)
	}

	win.Begin()
	testClick(t, win, 200, 10)
	e := tpl.Button("Go")
	win.End()
	if !e.Is(event.Button, event.Click) {
		t.Errorf("expected Click event, got %v", e)
	}
}

func TestCSSLayoutCentresAndFlows(t *testing.T) {
	win, _, _ := antui.Offscreen(400, 200)
	tpl := TemplateWithCSS(win)
	classes, path := cssTable(t, "")
	if err := tpl.SetStyle(path, classes); err != nil {
		t.Fatal(err)
	}

	u := Scale(win)

	x, y, w, h, _, ok := tpl.layout(css.RoleLabel, "Hi")
	if !ok {
		t.Fatal("label should be visible")
	}
	if x != (400-w)/2 {
		t.Errorf("label should be centred, got x=%d w=%d", x, w)
	}
	if y != 0 {
		t.Errorf("first widget starts at the top, got y=%d", y)
	}

	_, y2, _, h2, _, ok := tpl.layout(css.RoleButton, "Go")
	if !ok {
		t.Fatal("button should be visible")
	}
	if y2 != y+h {
		t.Errorf("widget should flow below the previous one, got %d want %d", y2, y+h)
	}
	if h2 != textHeight(u) {
		t.Errorf("button natural height = text height, got %d want %d", h2, textHeight(u))
	}
}

func TestCSSDisplayNoneHides(t *testing.T) {
	win, _, _ := antui.Offscreen(400, 200)
	tpl := TemplateWithCSS(win)
	classes, path := cssTable(t, "button { display: none; }")
	if err := tpl.SetStyle(path, classes); err != nil {
		t.Fatal(err)
	}

	if _, _, _, _, _, ok := tpl.layout(css.RoleButton, "Hidden"); ok {
		t.Error("a display:none widget must not be laid out")
	}
	win.Begin()
	e := tpl.Button("Hidden")
	win.End()
	if e.Ok() {
		t.Errorf("a hidden button must report nothing, got %v", e)
	}
}

func TestCSSPseudoSelectors(t *testing.T) {
	win, _, _ := antui.Offscreen(400, 200)
	tpl := TemplateWithCSS(win)
	classes, path := cssTable(t, `
		button { background-color: #111111; }
		button:hover { background-color: #222222; }
		button:focus { border: 2px solid #333333; }
	`)
	if err := tpl.SetStyle(path, classes); err != nil {
		t.Fatal(err)
	}

	winner := classes.GetStyle(css.RoleButton, nil, css.StateWith(true, false, false, false), 400)
	if winner.Background.R() != 0x22 {
		t.Errorf("hover should win over resting background, got %v", winner.Background)
	}
}

func TestCSSResponsiveWidth(t *testing.T) {
	win, _, _ := antui.Offscreen(800, 400)
	tpl := TemplateWithCSS(win)
	classes, path := cssTable(t, `button { width: 50%; }`)
	if err := tpl.SetStyle(path, classes); err != nil {
		t.Fatal(err)
	}

	_, _, w, _, _, ok := tpl.layout(css.RoleButton, "Wide")
	if !ok {
		t.Fatal("button should be visible")
	}
	if w != 400 {
		t.Errorf("50%% of 800 wide should be 400, got %d", w)
	}
}

func TestCSSDatePickerFlows(t *testing.T) {
	win, _, _ := antui.Offscreen(400, 200)
	tpl := TemplateWithCSS(win)
	classes, path := cssTable(t, "")
	if err := tpl.SetStyle(path, classes); err != nil {
		t.Fatal(err)
	}

	// The picker is a select-like box with an empty stylesheet: no padding,
	// no border, no label, so width is the caret allowance or the font floor
	// and the height is one text line.
	u := Scale(win)
	wantW := max(16*u, 24*canvas.FontWidth*u)
	wantH := textHeight(u)
	_, _, w, h, _, ok := tpl.layout(css.RoleDatePicker, "")
	if !ok {
		t.Fatal("datepicker should be visible")
	}
	if w != wantW || h != wantH {
		t.Errorf("datepicker natural size, got %dx%d want %dx%d", w, h, wantW, wantH)
	}
}

func TestCSSMarginsCollapse(t *testing.T) {
	win, _, _ := antui.Offscreen(400, 200)
	tpl := TemplateWithCSS(win)
	classes, path := cssTable(t, `
		label  { margin: 12px; }
		button { margin: 12px; }
		button { margin-top: 30px; }
	`)
	if err := tpl.SetStyle(path, classes); err != nil {
		t.Fatal(err)
	}
	u := Scale(win)
	th := textHeight(u)

	// Sibling margins of 12px share one gap: the first label's own
	// margin-top (no previous bottom margin to collapse with)…
	_, y1, _, _, _, ok := tpl.layout(css.RoleLabel, "One")
	if !ok {
		t.Fatal("label should be laid out")
	}
	if y1 != 12*u {
		t.Errorf("first box keeps its own margin-top, y=%d want %d", y1, 12*u)
	}

	// …and only the larger of 12px and the button's 30px margin-top is
	// taken between the two boxes.
	_, y2, _, _, _, ok := tpl.layout(css.RoleButton, "Two")
	if !ok {
		t.Fatal("button should be laid out")
	}
	if y2 != 12*u+th+30*u {
		t.Errorf("collapsed gap must be the bigger margin, y=%d want %d", y2, 12*u+th+30*u)
	}
}

func TestCSSMarginsAndPosition(t *testing.T) {
	win, _, _ := antui.Offscreen(400, 200)
	tpl := TemplateWithCSS(win)
	classes, path := cssTable(t, `
		button { margin: 10px 20px; }
		label  { position: absolute; top: 5px; left: 15px; }
	`)
	if err := tpl.SetStyle(path, classes); err != nil {
		t.Fatal(err)
	}
	u := Scale(win)

	// The label is taken out of the flow by position:absolute.
	x, y, _, _, _, ok := tpl.layout(css.RoleLabel, "Pin")
	if !ok {
		t.Fatal("label should be laid out")
	}
	if x != 15*u || y != 5*u {
		t.Errorf("absolute position, got (%d,%d) want (15,5)", x, y)
	}

	// The button flows below the label's slot: it is absolute, so it does
	// not push the flow, but the button's own margin-top does.
	_, y2, _, _, _, ok := tpl.layout(css.RoleButton, "Go")
	if !ok {
		t.Fatal("button should be laid out")
	}
	if y2 != 10*u {
		t.Errorf("margin-top should offset the button, got y=%d want %d", y2, 10*u)
	}
}

// cssUI builds a bare uiTemplate whose style is the CSS style bound to the
// given sheet, for hit-testing the interaction layer without the flow layout.
func cssUI(t *testing.T, sheet string) *uiTemplate {
	t.Helper()
	win, _, _ := antui.Offscreen(900, 800) // u=2, 16 px glyphs at 15px
	classes, _ := cssTable(t, sheet)
	st := &cssStyle{win: win, classes: classes}
	ut := &uiTemplate{win: win, style: st}
	return ut
}

func TestCSSInputClickPositionsCaret(t *testing.T) {
	ut := cssUI(t, "input { font-size: 15px; }")
	win := ut.win
	glyph := 16 // one 15px glyph at scale 2

	text := "abcdef"
	x, y, w, h := 20, 30, 400, 40

	// A click in the middle of the fifth glyph ('e') puts the caret after it.
	win.Begin()
	testClick(t, win, x+4*glyph+glyph/2, y+h/2)
	ut.Input(win, x, y, w, h, &text)
	win.End()
	if c := win.WidgetCursor(); c != 5 {
		t.Fatalf("clicking 'e' should drop the caret at 5, got %d", c)
	}

	// The caret stays where the click put it: typing lands between e and f.
	win.Begin()
	win.Push(antui.Event{Type: antui.EventText, Text: "X"})
	ut.Input(win, x, y, w, h, &text)
	win.End()
	if text != "abcdeXf" {
		t.Fatalf("typing after the click = %q, want %q", text, "abcdeXf")
	}
}

func TestCSSTextAreaClickPositionsCaret(t *testing.T) {
	ut := cssUI(t, "textarea { font-size: 15px; }")
	win := ut.win

	text := "abc\ndef"
	x, y, w, h := 20, 30, 400, 100
	lineH := 32 // one 15px line at scale 2

	// Second visual line sits at y+lineH; clicking the left half of its
	// second glyph ('e', index 5 in the whole string) puts the caret right
	// before the 'e'.
	win.Begin()
	testClick(t, win, x+16+4, y+lineH+16)
	ut.TextArea(win, x, y, w, h, &text)
	win.End()
	if c := win.WidgetCursor(); c != 5 {
		t.Fatalf("clicking 'e' on the second line should give cursor 5, got %d", c)
	}

	// And typing goes in at the clicked spot.
	win.Begin()
	win.Push(antui.Event{Type: antui.EventText, Text: "Y"})
	ut.TextArea(win, x, y, w, h, &text)
	win.End()
	if text != "abc\ndYef" {
		t.Fatalf("typing after the click = %q, want %q", text, "abc\ndYef")
	}
}

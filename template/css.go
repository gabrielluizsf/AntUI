package template

import "github.com/gabrielluizsf/antui"

// CSS is the coord-free template [TemplateWithCSS] makes: the developer says
// what to draw, and the template says where. Each widget is laid out by the
// stylesheet in [CSSClasses], flowing top-to-bottom and centred horizontally,
// exactly as a web document flows: margins space widgets apart, padding and
// borders grow each box, percentages are of the window, and display:none
// drops a widget out of the document entirely. Inline and inline-block boxes
// pack along a line and wrap at the window edge; relative and sticky boxes
// keep their flow slot and only shift where they paint; absolutely and fixed
// positioned boxes leave the flow and take their inset edges.
//
//	app := template.TemplateWithCSS(win)
//	if err := app.SetStyle("app.css", css.NewTable()); err != nil {
//		// ...
//	}
//	app.Button("button").SetProperty("button", nil, "background-color", "#3E63DD")
//
//	app.Reset()          // restart the flow at the top of the window
//	app.Label("Welcome") // the CSS decides position and style
//	if app.Button("Start").Is(event.Button, event.Click) {
//		app.Go() // or whatever the click means
//	}
//
// It shares the interaction machinery of the coordinate templates — pressing,
// focusing, typing, the events — so a CSS screen plays the same music as
// anything [New] makes.
type CSS struct {
	win     *antui.Window
	style   *cssStyle
	ui      *uiTemplate
	cursorY int // the line the flow is filling
	lineX   int // where the next inline box lands in the line
	lineH   int // how tall the line has grown
	lastMB  int // the previous box's margin-bottom, banked for collapsing

	// The flex block currently drawing its children. flexing is true while a
	// [CSS.Flex] callback runs; measure marks the silent first pass, which
	// collects the children's boxes so the draw pass can place them. A
	// [CSS.Flex] inside one flattens into the same batch, while a [CSS.Grid]
	// inside one is a real item of it.
	flex    *flexBatch
	flexing bool
	grid    *gridBatch
	griding bool
	multi   *multiBatch
	table   *tableBatch
	// multicoling and tabling are true while a [CSS.MultiCol] or a [CSS.Table]
	// callback runs, and row is the row collecting cells inside a [CSS.TableRow],
	// which is how a widget drawn in a table tells it is a cell and not a block.
	multicoling bool
	tabling     bool
	row         *tableRow
	measure     bool
}

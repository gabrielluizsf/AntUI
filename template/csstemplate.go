package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/audio"
	"github.com/gabrielluizsf/antui/template/css"
	"github.com/gabrielluizsf/antui/template/event"
)

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
	// [CSS.Flex] inside one flattens into the same batch.
	flex    *flexBatch
	flexing bool
	measure bool
}

// TemplateWithCSS makes the coord-free template: the developer passes no
// coordinates, and the template decides where each widget lives from an
// external stylesheet. Call [CSS.SetStyle] with the CSS file and a class
// table before drawing anything.
func TemplateWithCSS(win *antui.Window) *CSS {
	cs := newCSSStyle(win)
	return &CSS{
		win:   win,
		style: cs,
		ui:    &uiTemplate{win: win, style: cs, sound: audio.Simple},
	}
}

// SetStyle loads the stylesheet into the table and binds it to every widget
// the template draws afterwards. It replaces whatever the table held, and any
// compile error stops the load.
//
//	classes := css.NewTable()
//	classes.Button = "primary wide"
//	if err := tpl.SetStyle("app.css", classes); err != nil { ... }
func (c *CSS) SetStyle(cssFile string, classes *css.CSSClasses) error {
	if err := classes.SetStyle(cssFile); err != nil {
		return err
	}
	c.style.classes = classes
	c.cursorY = 0
	c.lineX = 0
	c.lineH = 0
	return nil
}

// Style exposes the underlying class table, so the developer can reach the
// whole [css.CSSClasses] API — Apply, SetProperty, Property, Warnings — after
// a stylesheet is loaded. The returned pointer is shared with the painter.
func (c *CSS) Style() *css.CSSClasses { return c.style.classes }

// SetSound swaps the sound bank the components ring their events through. The
// default is the [github.com/gabrielluizsf/antui/template/audio.Simple]
// bank; a nil bank makes the template silent.
func (c *CSS) SetSound(s event.SoundBank) { c.ui.sound = s }

// SetImage registers a canvas a stylesheet's background-image url() can point
// at. The stylesheet names images by the literal text of the url(), so it is
// lowercase, exactly as the declaration wrote it. A registered image paints at
// its own pixel size until background-size scales it.
func (c *CSS) SetImage(name string, cv *canvas.Canvas) { c.style.images[name] = cv }

// Sounds is the bank the components' events ring through right now.
func (c *CSS) Sounds() event.SoundBank { return c.ui.sound }

// Background is the colour the window is cleared to. It comes from the body
// of the stylesheet when one is loaded, and from the theme otherwise.
func (c *CSS) Background(win *antui.Window) canvas.Color {
	return c.style.Background(win)
}

// Reset starts a new "page": the vertical cursor returns to the top of the
// window, the widget timeline begins a fresh run, and the widgets drawn
// afterwards flow down from it again. Call it once per frame before drawing.
func (c *CSS) Reset() {
	c.style.resetTimeline()
	c.cursorY = 0
	c.lineX = 0
	c.lineH = 0
	c.lastMB = 0
	c.flex = nil
	c.flexing = false
	c.measure = false
}

// layout measures and places one widget in the flow, opening the widget's
// timeline entry and hosting its style on it. ok is false for a widget the
// stylesheet hides (display:none), which draws nothing and hears nothing.
// An absolutely- or fixed-positioned widget is placed by its inset edges and
// does not advance the flow; a relative or sticky one keeps its slot and only
// shifts; every other widget flows down, or along a line when inline. Inside
// a [CSS.Flex] block a widget is a flex item instead: the silent pass
// records it and reports ok=false, and the draw pass hands back the box the
// solver placed it in.
func (c *CSS) layout(role, label string) (x, y, w, h int, st css.Style, ok bool) {
	if c.flexing {
		return c.layoutFlex(role, label)
	}
	e := c.style.beginWidget(role, label)
	return c.layoutBox(e, role, label, e.state)
}

// layoutBox is layout's engine: it resolves the entry's style for the given
// interaction state and places the widget it describes.
func (c *CSS) layoutBox(e *widgetEntry, role, label string, s State) (x, y, w, h int, st css.Style, ok bool) {
	st = c.style.entryStyle(e, role, s)
	if st.Display == css.DisplayNone {
		return 0, 0, 0, 0, st, false
	}
	u := c.style.u()
	winW, winH := c.win.Width(), c.win.Height()

	w, h = c.style.natural(role, label, st, u)
	if st.Has("width") {
		w = c.style.length(st, st.Width, winW)
	}
	if st.Has("min-width") {
		w = max(w, c.style.length(st, st.MinWidth, winW))
	}
	if st.Has("max-width") {
		w = min(w, c.style.length(st, st.MaxWidth, winW))
	}
	if st.Has("height") {
		h = c.style.length(st, st.Height, winH)
	}
	if st.Has("min-height") {
		h = max(h, c.style.length(st, st.MinHeight, winH))
	}
	if st.Has("max-height") {
		h = min(h, c.style.length(st, st.MaxHeight, winH))
	}

	ml := c.style.length(st, st.Margin[3], winW)
	mr := c.style.length(st, st.Margin[1], winW)
	mt := c.style.length(st, st.Margin[0], winW)
	mb := c.style.length(st, st.Margin[2], winW)

	if st.OutOfFlow() {
		x, y = c.placeOutOfFlow(st, w, h, winW, winH)
		return x, y, w, h, st, true
	}
	if st.Inline() {
		x, y = c.placeInline(w, h, winW, ml, mr)
	} else {
		x, y = c.placeBlock(w, h, winW, ml, mr, mt, mb)
	}
	x, y = c.shiftInset(st, x, y, winW, winH)
	return x, y, w, h, st, true
}

// placeBlock starts a new line for a block box: it closes whatever inline line
// was open, collapses the vertical margins and centres the box.
func (c *CSS) placeBlock(w, h, winW, ml, mr, mt, mb int) (int, int) {
	c.cursorY += c.lineH
	c.lineH = 0
	c.lineX = 0
	c.cursorY += max(c.lastMB, mt)
	x := ml + (winW-ml-mr-w)/2
	y := c.cursorY
	c.cursorY += h
	c.lastMB = mb
	return x, y
}

// placeInline packs an inline or inline-block box along the current line,
// wrapping to a new line when it would cross the window edge. Vertical
// margins of inline boxes are ignored, as they are in CSS.
func (c *CSS) placeInline(w, h, winW, ml, mr int) (int, int) {
	if c.lineH > 0 && c.lineX+ml+w+mr > winW {
		c.cursorY += c.lineH
		c.lineH = 0
		c.lineX = 0
	}
	x := c.lineX + ml
	y := c.cursorY
	c.lineX += ml + w + mr
	c.lineH = max(c.lineH, h)
	c.lastMB = 0
	return x, y
}

// placeOutOfFlow positions an absolute or fixed box from its inset edges:
// left/top when given, otherwise right/bottom measured from the far edge of
// the window (the viewport). With neither, the box rests at the window origin.
func (c *CSS) placeOutOfFlow(st css.Style, w, h, winW, winH int) (int, int) {
	x, y := 0, 0
	switch {
	case st.Has("left"):
		x = c.style.length(st, st.Left, winW)
	case st.Has("right"):
		x = winW - w - c.style.length(st, st.Right, winW)
	}
	switch {
	case st.Has("top"):
		y = c.style.length(st, st.Top, winH)
	case st.Has("bottom"):
		y = winH - h - c.style.length(st, st.Bottom, winH)
	}
	return x, y
}

// shiftInset moves a relative or sticky box away from its flow position by its
// inset edges, leaving the slot it occupies untouched.
func (c *CSS) shiftInset(st css.Style, x, y, winW, winH int) (int, int) {
	if st.Position != css.PositionRelative && st.Position != css.PositionSticky {
		return x, y
	}
	switch {
	case st.Has("left"):
		x += c.style.length(st, st.Left, winW)
	case st.Has("right"):
		x -= c.style.length(st, st.Right, winW)
	}
	switch {
	case st.Has("top"):
		y += c.style.length(st, st.Top, winH)
	case st.Has("bottom"):
		y -= c.style.length(st, st.Bottom, winH)
	}
	return x, y
}

// visible reports whether a widget painted by the stylesheet should draw and
// hear input. visibility:hidden keeps its flow slot but is skipped here.
func visible(st css.Style) bool {
	return st.Visibility == css.VisibilityVisible
}

// paint runs a widget's painter clipped to its box when overflow asks for it,
// then adds the scrollbar indicators for scroll/auto. A transform on the
// style records the whole pass on a layer and composites it back, so the box
// and its shadows rotate and scale together.
func (c *CSS) paint(st css.Style, role, label string, x, y, w, h int, draw func()) {
	if !visible(st) {
		return
	}
	c.style.transformed(st, x, y, w, h, func() {
		if interactive(st) && c.style.boxOn(st) {
			c.style.paintBox(c.win, st, x, y, w, h)
		}
		if st.Overflow == [2]uint8{css.OverflowVisible, css.OverflowVisible} {
			draw()
			return
		}
		u := c.style.u()
		nw, nh := c.style.natural(role, label, st, u)
		cv := c.win.Canvas()
		saved := cv.Clip
		cv.SetClip(x, y, w, h)
		draw()
		cv.Clip = saved
		c.style.paintScrollbars(c.win, st, x, y, w, h, nw, nh)
	})
}

// interactionState reads the state the shared interaction machinery assigns a
// widget: hovered when the pointer is in its hover box, focused when the
// keyboard arrived, pressed when a press is landing on it, and active while
// it is the one being held. The identity box matches the widget's own
// WidgetID — for focus to agree — and idW/idH are those calling arguments;
// hovW/hovH are the box the widget itself counts a hover in. extra carries
// the states the widget decided itself, the On of a checked box included.
func (c *CSS) interactionState(kind string, x, y, idW, idH int, label string, hovW, hovH int, extra State) State {
	s := extra
	id := c.win.WidgetID("template:"+kind, x, y, idW, idH, label)
	hovered := c.win.Hovered(x, y, hovW, hovH)
	s.Hovered = hovered
	s.Focused = c.win.WidgetFocus() == id
	s.Active = c.win.WidgetActive() == id
	if s.Active && hovered {
		s.Pressed = true
	}
	return s
}

// interactive reports whether a widget takes pointer input. pointer-events:
// none still paints, so the caller falls back to a paint-only path.
func interactive(st css.Style) bool {
	return st.PointerEvents != css.PointerEventsNone
}

// Label paints one line of text in the window, placed and styled by CSS.
// The widget never reports an event.
func (c *CSS) Label(text string) {
	x, y, w, h, st, ok := c.layout(css.RoleLabel, text)
	if !ok || !visible(st) {
		return
	}
	c.style.transformed(st, x, y, w, h, func() {
		c.style.Label(c.win, x, y, text)
	})
}

// Button draws a button and reports whether it was clicked.
func (c *CSS) Button(label string) event.Event {
	x, y, w, h, st, ok := c.layout(css.RoleButton, label)
	if !ok {
		return event.Nothing
	}
	s := c.interactionState("button", x, y, w, h, label, w, h, State{})
	st = c.style.entryStyle(c.style.cur, css.RoleButton, s)
	if !interactive(st) {
		c.paint(st, css.RoleButton, label, x, y, w, h, func() {
			c.style.Button(c.win, s, x, y, w, h, label)
		})
		return event.Nothing
	}
	var e event.Event
	c.paint(st, css.RoleButton, label, x, y, w, h, func() {
		e = c.ui.Button(c.win, x, y, w, h, label)
	})
	return e
}

// Checkbox draws a labelled checkbox, flips value when clicked, and reports
// whether it flipped this frame.
func (c *CSS) Checkbox(label string, value *bool) event.Event {
	x, y, w, h, st, ok := c.layout(css.RoleCheckbox, label)
	if !ok {
		return event.Nothing
	}
	u := c.style.u()
	box := 18 * u
	on := value != nil && *value
	s := c.interactionState("checkbox", x, y, box, box, label, box+8*u+textWidth(u, label), box, State{On: on})
	st = c.style.entryStyle(c.style.cur, css.RoleCheckbox, s)
	if !interactive(st) {
		c.paint(st, css.RoleCheckbox, label, x, y, w, h, func() {
			c.style.Checkbox(c.win, s, x, y, label, on)
		})
		return event.Nothing
	}
	var e event.Event
	c.paint(st, css.RoleCheckbox, label, x, y, w, h, func() {
		e = c.ui.Checkbox(c.win, x, y, label, value)
	})
	return e
}

// Radio draws one option of a group. Clicking it sets value to option and
// reports true; clicking the already-selected option changes nothing.
func (c *CSS) Radio(label string, value *int, option int) event.Event {
	x, y, w, h, st, ok := c.layout(css.RoleRadio, label)
	if !ok {
		return event.Nothing
	}
	u := c.style.u()
	size := 18 * u
	on := value != nil && *value == option
	s := c.interactionState("radio", x, y, size, option, label, size+8*u+textWidth(u, label), size, State{On: on})
	st = c.style.entryStyle(c.style.cur, css.RoleRadio, s)
	if !interactive(st) {
		c.paint(st, css.RoleRadio, label, x, y, w, h, func() {
			c.style.Radio(c.win, s, x, y, label, on)
		})
		return event.Nothing
	}
	var e event.Event
	c.paint(st, css.RoleRadio, label, x, y, w, h, func() {
		e = c.ui.Radio(c.win, x, y, label, value, option)
	})
	return e
}

// Slider draws a horizontal slider and reports whether its value changed.
func (c *CSS) Slider(value *float32, minValue, maxValue float32) event.Event {
	x, y, w, h, st, ok := c.layout(css.RoleSlider, "")
	if !ok {
		return event.Nothing
	}
	s := c.interactionState("slider", x, y, w, h, "", w, h, State{})
	st = c.style.entryStyle(c.style.cur, css.RoleSlider, s)
	var e event.Event
	c.paint(st, css.RoleSlider, "", x, y, w, h, func() {
		e = c.ui.Slider(c.win, x, y, w, h, value, minValue, maxValue)
	})
	return e
}

// Input draws a single-line text field and reports whether its text changed.
func (c *CSS) Input(text *string) event.Event {
	x, y, w, h, st, ok := c.layout(css.RoleInput, "")
	if !ok {
		return event.Nothing
	}
	s := c.interactionState("input", x, y, w, h, "", w, h, State{})
	st = c.style.entryStyle(c.style.cur, css.RoleInput, s)
	var e event.Event
	c.paint(st, css.RoleInput, "", x, y, w, h, func() {
		e = c.ui.Input(c.win, x, y, w, h, text)
	})
	return e
}

// Select draws a dropdown picker: a box showing the currently chosen option,
// and a menu of ALL options when it is open. Clicking the box opens or closes
// it; picking an option sets index to its position and closes the menu.
func (c *CSS) Select(index *int, options []string) event.Event {
	x, y, w, h, st, ok := c.layout(css.RoleSelect, "")
	if !ok {
		return event.Nothing
	}
	s := c.interactionState("select", x, y, w, h, "", w, h, State{})
	st = c.style.entryStyle(c.style.cur, css.RoleSelect, s)
	var e event.Event
	c.paint(st, css.RoleSelect, "", x, y, w, h, func() {
		e = c.ui.Select(c.win, x, y, w, h, index, options)
	})
	return e
}

// TextArea draws a multi-line text field and reports whether its text
// changed. Enter inserts a line break; word wrap and scrolling follow.
func (c *CSS) TextArea(text *string) event.Event {
	x, y, w, h, st, ok := c.layout(css.RoleTextArea, "")
	if !ok {
		return event.Nothing
	}
	s := c.interactionState("textarea", x, y, w, h, "", w, h, State{})
	st = c.style.entryStyle(c.style.cur, css.RoleTextArea, s)
	var e event.Event
	c.paint(st, css.RoleTextArea, "", x, y, w, h, func() {
		e = c.ui.TextArea(c.win, x, y, w, h, text)
	})
	return e
}

// Switch draws a flip toggle like a checkbox but a sliding thumb. It flips
// value when clicked and reports whether it flipped this frame.
func (c *CSS) Switch(label string, value *bool) event.Event {
	x, y, w, h, st, ok := c.layout(css.RoleSwitch, label)
	if !ok {
		return event.Nothing
	}
	u := c.style.u()
	on := value != nil && *value
	s := c.interactionState("switch", x, y, 36*u, 20*u, label, 36*u+8*u+textWidth(u, label), 20*u, State{On: on})
	st = c.style.entryStyle(c.style.cur, css.RoleSwitch, s)
	if !interactive(st) {
		c.paint(st, css.RoleSwitch, label, x, y, w, h, func() {
			c.style.Switch(c.win, s, x, y, label, on)
		})
		return event.Nothing
	}
	var e event.Event
	c.paint(st, css.RoleSwitch, label, x, y, w, h, func() {
		e = c.ui.Switch(c.win, x, y, label, value)
	})
	return e
}

// Progress paints a read-only progress bar with progress from 0 to 1. It
// never reports an event: a bar has nothing to say.
func (c *CSS) Progress(progress float32) {
	x, y, w, h, st, ok := c.layout(css.RoleProgress, "")
	if !ok {
		return
	}
	c.paint(st, css.RoleProgress, "", x, y, w, h, func() {
		c.ui.Progress(c.win, x, y, w, h, progress)
	})
}

// DatePicker draws a picker that works like a select: a box showing the chosen
// date, and a calendar that pops up below the box when it is clicked.
func (c *CSS) DatePicker(value *Date) event.Event {
	x, y, w, h, st, ok := c.layout(css.RoleDatePicker, "")
	if !ok {
		return event.Nothing
	}
	s := c.interactionState("date", x, y, w, h, "", w, h, State{})
	st = c.style.entryStyle(c.style.cur, css.RoleDatePicker, s)
	var e event.Event
	c.paint(st, css.RoleDatePicker, "", x, y, w, h, func() {
		e = c.ui.DatePicker(c.win, x, y, w, h, value)
	})
	return e
}

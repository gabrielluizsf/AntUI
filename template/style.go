package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
)

// State is everything the interaction layer knows about a widget this frame,
// handed to a style so it can paint the right state with no guessing.
type State struct {
	Hovered bool // the pointer is over the widget
	Pressed bool // the widget is being pressed right now
	Active  bool // the widget is the one being dragged (a slider)
	Focused bool // the widget has the keyboard; Tab walks between these
	On      bool // for checkboxes and radios: the value it shows
}

// A Style is the purely visual half of a template: how each component is
// painted. The interaction — pressing, releasing, focusing, typing, and the
// events that come out of it — is shared by every template, so a style only
// decides what it all looks like. Implement Style, pick an
// [github.com/gabrielluizsf/antui/template/event.SoundBank], and [New] makes
// a whole template out of them.
//
// Every component is painted on the window's canvas with the window's own
// helper methods, so a style may use the theme, the built-in face, or any
// drawing primitive it pleases.
type Style interface {
	// Background is the colour the window is cleared to each frame. It is
	// what "empty" looks like in this style.
	Background(win *antui.Window) canvas.Color
	// Label paints one line of text at its top-left corner.
	Label(win *antui.Window, x, y int, text string)
	// Button paints a button. state.Hovered and state.Pressed say how the
	// user is holding it; label is centered.
	Button(win *antui.Window, state State, x, y, w, h int, label string)
	// Checkbox paints a labelled checkbox. on is the value it currently
	// shows, on the box's own row.
	Checkbox(win *antui.Window, state State, x, y int, label string, on bool)
	// Radio paints one radio option. on says whether this option is picked.
	Radio(win *antui.Window, state State, x, y int, label string, on bool)
	// Slider paints a horizontal slider showing the given value between 0
	// and 1, normalised by the interaction layer.
	Slider(win *antui.Window, state State, x, y, w, h int, value float32)
	// Input paints a single-line text field. text is the whole field, cursor
	// where the caret sits within it, and caret whether the caret should be
	// visible this frame. Scrolling and clipping are the style's to do.
	Input(win *antui.Window, state State, x, y, w, h int, text string, cursor int, caret bool)
	// InputTextPos maps a click at (mx,my), in window coordinates, inside a
	// single-line input onto the byte index of the character it lands on,
	// honouring the same inside padding and scroll the box draws with. cursor
	// is the caret's current position, used to compute that scroll.
	InputTextPos(win *antui.Window, state State, x, y, w, h int, text string, cursor, mx, my int) int
	// TextAreaTextPos is the wrapped, multi-line cousin of InputTextPos: it
	// picks the line under my and then the character under mx, with the same
	// wrap, viewport and row heights the box paints.
	TextAreaTextPos(win *antui.Window, state State, x, y, w, h int, text string, cursor, mx, my int) int
	// Select paints a dropdown picker. value is the option currently shown
	// in the box. The open list is painted by SelectOption, one call per
	// entry.
	Select(win *antui.Window, state State, x, y, w, h int, value string, open bool)
	// SelectOption paints one entry of an open dropdown menu. selected says
	// whether it is the currently-chosen option.
	SelectOption(win *antui.Window, state State, x, y, w, h int, label string, selected bool)
	// TextArea paints a multi-line text field, drawing all of text wrapped
	// at w. cursor and caret behave exactly as for Input.
	TextArea(win *antui.Window, state State, x, y, w, h int, text string, cursor int, caret bool)
	// Switch paints a flip toggle. on is the value it shows.
	Switch(win *antui.Window, state State, x, y int, label string, on bool)
	// Progress paints a read-only progress bar, value 0..1.
	Progress(win *antui.Window, x, y, w, h int, value float32)
	// DatePicker paints the calendar popup of an open date picker, below the
	// box. year and month name the month on show, firstWD is the weekday of
	// its first day, days its day count, and selected/today/hover are day
	// numbers in that month (0 when not applicable). The style places the day
	// cells with this package's shared calendar geometry, so the painting and
	// the hit-testing agree.
	DatePicker(win *antui.Window, state State, x, y, w, h int, year, month, firstWD, days, selected, today, hover int)
	// DatePickerBox paints the closed date picker: a box showing the chosen
	// date, with the open flag drawing the caret flipped the way a dropdown's
	// is. It looks like Select by design, since the picker is a select that
	// opens a calendar.
	DatePickerBox(win *antui.Window, state State, x, y, w, h int, value string, open bool)
}

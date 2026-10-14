package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
)

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
	Label(p LabelParams)
	// Button paints a button. state.Hovered and state.Pressed say how the
	// user is holding it; label is centered.
	Button(p ButtonParams)
	// Checkbox paints a labelled checkbox. on is the value it currently
	// shows, on the box's own row.
	Checkbox(p CheckboxParams)
	// Radio paints one radio option. on says whether this option is picked.
	Radio(p RadioParams)
	// Slider paints a horizontal slider showing the given value between 0
	// and 1, normalised by the interaction layer.
	Slider(p SliderParams)
	// Input paints a single-line text field. text is the whole field, cursor
	// where the caret sits within it, and caret whether the caret should be
	// visible this frame. Scrolling and clipping are the style's to do.
	Input(p InputParams)
	// InputTextPos maps a click at (mx,my), in window coordinates, inside a
	// single-line input onto the byte index of the character it lands on,
	// honouring the same inside padding and scroll the box draws with. cursor
	// is the caret's current position, used to compute that scroll.
	InputTextPos(p InputTextPosParams) int
	// TextAreaTextPos is the wrapped, multi-line cousin of InputTextPos: it
	// picks the line under my and then the character under mx, with the same
	// wrap, viewport and row heights the box paints.
	TextAreaTextPos(p TextAreaTextPosParams) int
	// Select paints a dropdown picker. value is the option currently shown
	// in the box. The open list is painted by SelectOption, one call per
	// entry.
	Select(p SelectParams)
	// SelectOption paints one entry of an open dropdown menu. selected says
	// whether it is the currently-chosen option.
	SelectOption(p SelectOptionParams)
	// TextArea paints a multi-line text field, drawing all of text wrapped
	// at w. cursor and caret behave exactly as for Input.
	TextArea(p TextAreaParams)
	// Switch paints a flip toggle. on is the value it shows.
	Switch(p SwitchParams)
	// Progress paints a read-only progress bar, value 0..1.
	Progress(p ProgressParams)
	// DatePicker paints the calendar popup of an open date picker, below the
	// box. year and month name the month on show, firstWD is the weekday of
	// its first day, days its day count, and selected/today/hover are day
	// numbers in that month (0 when not applicable). The style places the day
	// cells with this package's shared calendar geometry, so the painting and
	// the hit-testing agree.
	DatePicker(p DatePickerParams)
	// DatePickerBox paints the closed date picker: a box showing the chosen
	// date, with the open flag drawing the caret flipped the way a dropdown's
	// is. It looks like Select by design, since the picker is a select that
	// opens a calendar.
	DatePickerBox(p DatePickerBoxParams)
}

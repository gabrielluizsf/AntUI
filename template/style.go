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
}
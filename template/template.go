package template

import (
	"time"

	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/event"
)

// A Template is one whole look-and-sound, bound to a window. Each interactive
// component draws itself, reports what the user did to it this frame as an
// [event.Event], and — with a sound bank attached — plays its own sound. The
// ready-made ones are [Cyberpunk] and [Builtin]; other templates are made
// with [New].
type Template interface {
	// Background is the colour the window is cleared to before a screen
	// draws. It is the template saying what empty should look like.
	Background(win *antui.Window) canvas.Color
	// Label paints one line of text in the template's own voice, at its
	// top-left corner.
	Label(win *antui.Window, x, y int, text string)
	// Button draws a button and reports whether it was clicked.
	Button(win *antui.Window, x, y, w, h int, label string) event.Event
	// Checkbox draws a labelled checkbox, flips value when clicked, and
	// reports whether it flipped this frame.
	Checkbox(win *antui.Window, x, y int, label string, value *bool) event.Event
	// Radio draws one option of a group. Clicking it sets value to option and
	// reports true; clicking the already-selected option changes nothing.
	Radio(win *antui.Window, x, y int, label string, value *int, option int) event.Event
	// Slider draws a horizontal slider and reports whether its value changed.
	Slider(win *antui.Window, x, y, w, h int, value *float32, minValue, maxValue float32) event.Event
	// Input draws a single-line text field and reports whether its text
	// changed.
	Input(win *antui.Window, x, y, w, h int, text *string) event.Event
	// Select draws a dropdown picker: a box showing the currently chosen
	// option, and a menu of ALL options when it is open. Clicking the box
	// opens or closes it; picking an option sets index to its position and
	// closes the menu.
	Select(win *antui.Window, x, y, w, h int, index *int, options []string) event.Event
	// TextArea draws a multi-line text field. Enter inserts a line break;
	// the responsive parts — word wrap and scrolling — follow, but the
	// paging keys (Home, End and the arrows) already move the cursor.
	TextArea(win *antui.Window, x, y, w, h int, text *string) event.Event
	// Switch draws a flip toggle like a checkbox but a sliding thumb. It
	// flips value when clicked and reports whether it flipped this frame.
	Switch(win *antui.Window, x, y int, label string, value *bool) event.Event
	// Progress paints a read-only progress bar with progress from 0 to 1.
	// It never reports an event: a bar has nothing to say.
	Progress(win *antui.Window, x, y, w, h int, progress float32)
	// DatePicker draws a picker that works like a select: a box showing the
	// chosen date, and a calendar that pops up below it when the box is
	// clicked, for flipping months and picking a day.
	DatePicker(win *antui.Window, x, y, w, h int, value *Date) event.Event
	// Sounds is the bank the components' events ring through. A nil bank
	// makes the template silent, which is a sound too.
	Sounds() event.SoundBank
	// Ready says everything the program loads at start-up has finished
	// loading, so the loading screen can go away and the program itself can
	// show. Until it is called — or the loading timeout passes and closes
	// the app — the loading screen is what the window wears. It may be
	// called from any goroutine, which is the point: whichever goroutine
	// finished loading says so. Calling it twice, or after the screen is
	// already gone, changes nothing: a start-up happens once.
	Ready()
	// SetDisableLoading turns the start-up loading screen off when disabled
	// is true and back on when it is false. It is for a program with
	// nothing worth waiting for, or one that does all of its loading before
	// the window opens; by default the screen is on, and a program that
	// never says Ready is closed by the timeout rather than left staring at
	// a bar that never ends. See [CSS.SetDisableLoading] and
	// [Context.SetDisableLoading] for the same flag on the two ways a
	// program is usually wired.
	SetDisableLoading(disabled bool)
	// SetLoadingTimeout changes how long the loading screen waits before
	// the app closes itself, which is also the span the bar fills over, so
	// that 100% and the closing land on the same instant. The default is
	// five seconds; a duration of zero or less changes nothing, and a
	// program that wants the screen to stay without a deadline wants
	// [SetDisableLoading] instead.
	SetLoadingTimeout(d time.Duration)
	// loadingState is the screen's own state, for this package to arm,
	// wipe and read. It is unexported so that a Template is something only
	// this package can answer for: the loading screen is not a decoration
	// hung on a template but half of what a template is, and a type from
	// outside could not keep that promise. A program makes templates with
	// [New].
	loadingState() *loading
}

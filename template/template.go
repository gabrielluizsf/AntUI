// Package template is how a whole screen of AntUI is made without drawing a
// single control from scratch. A template is one look and one sound bound
// together: its components draw themselves the way the template wants and,
// when the user does something to one, play the sound the template carries
// for that component. The machinery in between — pressing, focusing, typing,
// the events — is shared by every template, so a template supplies only the
// painting and the hearing, and the package supplies the rest.
//
// Two templates come ready-made and match the two sound banks in
// [github.com/gabrielluizsf/antui/template/audio]:
//
//	cyber := template.Cyberpunk(win)   // futuristic look, Tech sound
//	built := template.Builtin(win)     // built-in look, Simple (cartoon) sound
//
//	clicked := cyber.Button(win, 20, 20, 120, 36, "Launch") // click clicks
//	typed := built.Input(win, 20, 70, 180, 30, &name)       // typing tick-tocks
//
// Both return an [event.Event] that says exactly what happened this frame, so
// a screen can tell a click from a drag and a typed key from a backspace.
//
// The whole point is that a game does not build its own widgets: it names a
// template, draws screens with it, and the screens point at one another and
// come back. See [Context] for that.
//
// Making a third template is filling in the two halves. The [Style]
// interface is the painting; any [event.SoundBank] — an
// [github.com/gabrielluizsf/antui/template/audio.Template], or a type of your
// own — is the hearing; and [New] joins them.
package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/audio"
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
	// Sounds is the bank the components' events ring through. A nil bank
	// makes the template silent, which is a sound too.
	Sounds() event.SoundBank
}

// New joins a painting and a hearing into a working template: everything the
// [Style] does not say — how a press feels, how focus moves, what event an
// interaction counts as — is the same for every template, and the package
// does it. The sound bank hears every event the components report.
//
//	myStyle := MyStyle{}
//	mySound := MySoundBank{}
//	tpl := template.New(win, myStyle, mySound)
func New(win *antui.Window, style Style, sound event.SoundBank) Template {
	return &uiTemplate{win: win, style: style, sound: sound}
}

// Cyberpunk is the futuristic template: neon edges and glowing controls, the
// level of finish a modern website skins itself in, voiced by
// [github.com/gabrielluizsf/antui/template/audio.Tech].
func Cyberpunk(win *antui.Window) Template {
	return New(win, cyberStyle{}, audio.Tech)
}

// Builtin is the template that paints exactly like the built-in widgets — the
// light and dark themes, the built-in face — voiced by
// [github.com/gabrielluizsf/antui/template/audio.Simple]. It is the zero-cost
// starting point, and the Cyberpunk template's plainspoken neighbour.
func Builtin(win *antui.Window) Template {
	return New(win, builtinStyle{}, audio.Simple)
}

// uiTemplate is [New]'s ordinary answer to [Template]: it runs the shared
// interaction machinery and leaves the painting to the style and the hearing
// to the bank.
type uiTemplate struct {
	win              *antui.Window
	style            Style
	sound            event.SoundBank
	sliderDragOffset int
	sliderDragID     uint64
}

func (t *uiTemplate) Background(win *antui.Window) canvas.Color {
	return t.style.Background(win)
}

func (t *uiTemplate) Label(win *antui.Window, x, y int, text string) {
	t.style.Label(win, x, y, text)
}

func (t *uiTemplate) Sounds() event.SoundBank { return t.sound }

// fire reports an event and, when one is attached, rings it through the
// sound bank before answering.
func (t *uiTemplate) fire(e event.Event) event.Event {
	if e.Ok() && t.sound != nil {
		t.sound.On(e)
	}
	return e
}

func (t *uiTemplate) Button(win *antui.Window, x, y, w, h int, label string) event.Event {
	id := win.WidgetID("template:button", x, y, w, h, label)
	hovered := win.Hovered(x, y, w, h)
	clicked := win.WidgetClick(id, hovered)
	focused := win.WidgetFocus() == id
	win.WidgetTabStop(id)
	keyed := win.WidgetKeyActivate(id)
	if keyed {
		clicked = true
	}
	t.style.Button(win, State{
		Hovered: hovered,
		Pressed: (win.WidgetActive() == id && hovered) || keyed,
		Focused: focused,
	}, x, y, w, h, label)
	if clicked {
		return t.fire(event.Event{Component: event.Button, Kind: event.Click})
	}
	return event.Nothing
}

func (t *uiTemplate) Checkbox(win *antui.Window, x, y int, label string, value *bool) event.Event {
	if value == nil {
		return event.Nothing
	}
	u := Scale(win)
	box := 18 * u
	w := box + 8*u + textWidth(u, label)
	id := win.WidgetID("template:checkbox", x, y, box, box, label)
	hovered := win.Hovered(x, y, w, box)
	clicked := win.WidgetClick(id, hovered)
	focused := win.WidgetFocus() == id
	win.WidgetTabStop(id)
	if win.WidgetKeyActivate(id) {
		clicked = true
	}
	changed := false
	if clicked {
		*value = !*value
		changed = true
	}
	t.style.Checkbox(win, State{Hovered: hovered, Focused: focused, On: *value}, x, y, label, *value)
	if changed {
		return t.fire(event.Event{Component: event.Checkbox, Kind: event.Toggle})
	}
	return event.Nothing
}

func (t *uiTemplate) Radio(win *antui.Window, x, y int, label string, value *int, option int) event.Event {
	if value == nil {
		return event.Nothing
	}
	u := Scale(win)
	size := 18 * u
	w := size + 8*u + textWidth(u, label)
	id := win.WidgetID("template:radio", x, y, size, option, label)
	hovered := win.Hovered(x, y, w, size)
	clicked := win.WidgetClick(id, hovered)
	focused := win.WidgetFocus() == id
	win.WidgetTabStop(id)
	if win.WidgetKeyActivate(id) {
		clicked = true
	}
	selected := *value == option

	if clicked && !selected {
		*value = option
		selected = true
	} else {
		clicked = false
	}
	t.style.Radio(win, State{Hovered: hovered, Focused: focused, On: selected}, x, y, label, selected)
	if clicked {
		return t.fire(event.Event{Component: event.Radio, Kind: event.Select})
	}
	return event.Nothing
}

func (t *uiTemplate) Slider(
	win *antui.Window,
	x, y, w, h int,
	value *float32,
	minValue, maxValue float32,
) event.Event {
	if value == nil || w <= 0 || h <= 0 || maxValue <= minValue {
		return event.Nothing
	}

	id := win.WidgetID("template:slider", x, y, w, h, "")

	u := Scale(win)

	// The knob is the only clickable part of the slider.
	knobRadius := max(h/2, 6*u)

	span := maxValue - minValue

	// The usable movement is from the center of the left knob
	// to the center of the right knob.
	usable := max(w-2*knobRadius, 1)

	// Convert the current value to a normalized position.
	position := (*value - minValue) / span
	position = min(max(position, 0), 1)

	// Current knob center.
	knobX := x + knobRadius + int(position*float32(usable)+0.5)
	knobY := y + h/2

	mouseX := win.MouseX()
	mouseY := win.MouseY()

	// ---------------------------------------------------------------------
	// Hit testing
	// ---------------------------------------------------------------------

	dx := mouseX - knobX
	dy := mouseY - knobY

	knobHit := dx*dx+dy*dy <= knobRadius*knobRadius

	if knobHit {
		win.WidgetHot(id)
	}

	// ---------------------------------------------------------------------
	// Start dragging
	// ---------------------------------------------------------------------

	if knobHit && win.MousePressed(antui.MouseLeft) {
		win.WidgetStartPress(id)
		win.WidgetFocusSet(id)

		// Preserve the exact point where the mouse grabbed the knob.
		//
		// Example:
		//     knob center = 200
		//     mouse       = 197
		//
		// We remember -3 so the knob does not jump to the mouse center.
		t.sliderDragOffset = mouseX - knobX
	}

	active := win.WidgetActive() == id

	// ---------------------------------------------------------------------
	// End dragging
	// ---------------------------------------------------------------------

	if active && !win.MouseDown(antui.MouseLeft) {
		win.WidgetEndPress(id)
		active = false
		t.sliderDragOffset = 0
	}

	focused := win.WidgetFocus() == id
	win.WidgetTabStop(id)

	changed := false

	// ---------------------------------------------------------------------
	// Keyboard
	// ---------------------------------------------------------------------

	if focused {
		step := span / 10

		if win.KeyPressed(antui.KeyLeft) {
			wanted := *value - step
			wanted = min(max(wanted, minValue), maxValue)

			if wanted != *value {
				*value = wanted
				changed = true
			}
		}

		if win.KeyPressed(antui.KeyRight) {
			wanted := *value + step
			wanted = min(max(wanted, minValue), maxValue)

			if wanted != *value {
				*value = wanted
				changed = true
			}
		}
	}

	// ---------------------------------------------------------------------
	// Mouse dragging
	// ---------------------------------------------------------------------

	if active {
		// Compensate for the point where the mouse grabbed the knob.
		//
		// Without this:
		//
		//     mouse ────────●
		//                    ↑
		//               knob jumps
		//
		// With the offset:
		//
		//     mouse ────●
		//               ↑
		//          knob follows naturally
		//
		adjustedMouseX := mouseX - t.sliderDragOffset

		position := float32(
			adjustedMouseX-(x+knobRadius),
		) / float32(usable)

		position = min(max(position, 0), 1)

		wanted := minValue + position*span
		wanted = min(max(wanted, minValue), maxValue)

		if wanted != *value {
			*value = wanted
			changed = true
		}
	}

	// ---------------------------------------------------------------------
	// Final clamp
	// ---------------------------------------------------------------------

	*value = min(max(*value, minValue), maxValue)

	// ---------------------------------------------------------------------
	// Draw
	// ---------------------------------------------------------------------

	visualValue := (*value - minValue) / span
	visualValue = min(max(visualValue, 0), 1)

	t.style.Slider(
		win,
		State{
			Hovered: knobHit,
			Active:  active,
			Focused: focused,
		},
		x, y, w, h,
		visualValue,
	)

	if changed {
		return t.fire(event.Event{
			Component: event.Slider,
			Kind:      event.Change,
		})
	}

	return event.Nothing
}

func (t *uiTemplate) Input(win *antui.Window, x, y, w, h int, text *string) event.Event {
	if text == nil {
		return event.Nothing
	}
	id := win.WidgetID("template:input", x, y, w, h, "")
	hovered := win.Hovered(x, y, w, h)
	if hovered {
		win.WidgetHot(id)
	}

	justFocused := false
	if win.MousePressed(antui.MouseLeft) {
		switch {
		case hovered:
			if win.WidgetFocus() != id {
				win.WidgetFocusSet(id)
				win.WidgetCursorSet(len(*text))
				justFocused = true
			}
		case win.WidgetFocus() == id:
			win.WidgetFocusClear(id)
		}
	}
	focused := win.WidgetFocus() == id
	win.WidgetTabStop(id)

	changed := false
	if focused {
		changed = win.WidgetEdit(text)
	}

	t.style.Input(win, State{
		Hovered: hovered,
		Focused: focused,
	}, x, y, w, h, *text, win.WidgetCursor(), focused && win.WidgetCaret())

	switch {
	case justFocused:
		return t.fire(event.Event{Component: event.TextInput, Kind: event.Focus})
	case changed:
		return t.fire(event.Event{
			Component: event.TextInput,
			Kind:      event.Type,
			Step:      win.WidgetCursor(),
		})
	}
	return event.Nothing
}

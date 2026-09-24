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
	"unicode/utf8"

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
	selectOpen       uint32 // the widget id of an open dropdown menu, 0 when none
	selectCursor     int    // which option the keyboard is walking in an open menu
	dateOpen         uint32 // the widget id of an open calendar popup, 0 when none
	dateYear         int    // the month an open calendar popup shows
	dateMonth        int
	dateCursor       int // the day an open calendar popup highlights
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
		return t.fire(event.Event{Component: event.Radio, Kind: event.Pick})
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
				justFocused = true
			}
			// The caret goes where the pointer pointed, even on a re-click
			// into the field.
			win.WidgetCursorSet(t.style.InputTextPos(win, State{
				Hovered: hovered,
				Focused: win.WidgetFocus() == id,
			}, x, y, w, h, *text, win.WidgetCursor(), win.MouseX(), win.MouseY()))
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

// TextArea draws a multi-line text field: Enter inserts a line break, the
// arrows and Home/End move the cursor between and along wrapped lines, and
// text flows onto as many lines as will fit in the box. It reports when the
// text changed or the field gained focus.
func (t *uiTemplate) TextArea(win *antui.Window, x, y, w, h int, text *string) event.Event {
	if text == nil {
		return event.Nothing
	}
	id := win.WidgetID("template:textarea", x, y, w, h, "")
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
				justFocused = true
			}
			// The caret goes where the pointer pointed, even on a re-click
			// into the field.
			win.WidgetCursorSet(t.style.TextAreaTextPos(win, State{
				Hovered: hovered,
				Focused: win.WidgetFocus() == id,
			}, x, y, w, h, *text, win.WidgetCursor(), win.MouseX(), win.MouseY()))
		case win.WidgetFocus() == id:
			win.WidgetFocusClear(id)
		}
	}
	focused := win.WidgetFocus() == id
	win.WidgetTabStop(id)

	changed := false
	if focused {
		u := Scale(win)
		maxRef := max(w/(8*u), 12)
		changed = t.editTextArea(win, text, maxRef)
	}

	t.style.TextArea(win, State{
		Hovered: hovered,
		Focused: focused,
	}, x, y, w, h, *text, win.WidgetCursor(), focused && win.WidgetCaret())

	switch {
	case justFocused:
		return t.fire(event.Event{Component: event.TextArea, Kind: event.Focus})
	case changed:
		return t.fire(event.Event{
			Component: event.TextArea,
			Kind:      event.Type,
			Step:      win.WidgetCursor(),
		})
	}
	return event.Nothing
}

// editTextArea is the multi-line cousin of WidgetEdit: it applies this frame's
// keys and typing to the focused field. maxRef is the wrap width in reference
// pixels (the width the text is measured at scale 1).
func (t *uiTemplate) editTextArea(win *antui.Window, text *string, maxRef int) bool {
	win.WidgetCursorSet(clampInt(win.WidgetCursor(), 0, len(*text)))
	changed := false

	if win.KeyPressed(antui.KeyBackspace) && win.WidgetCursor() > 0 {
		start := utf8PrevCursor(*text, win.WidgetCursor())
		*text = (*text)[:start] + (*text)[win.WidgetCursor():]
		win.WidgetCursorSet(start)
		changed = true
	}
	if win.KeyPressed(antui.KeyDelete) && win.WidgetCursor() < len(*text) {
		end := utf8NextCursor(*text, win.WidgetCursor())
		*text = (*text)[:win.WidgetCursor()] + (*text)[end:]
		changed = true
	}
	if win.KeyPressed(antui.KeyLeft) && win.WidgetCursor() > 0 {
		win.WidgetCursorSet(utf8PrevCursor(*text, win.WidgetCursor()))
	}
	if win.KeyPressed(antui.KeyRight) && win.WidgetCursor() < len(*text) {
		win.WidgetCursorSet(utf8NextCursor(*text, win.WidgetCursor()))
	}
	if win.KeyPressed(antui.KeyEnter) {
		c := win.WidgetCursor()
		*text = (*text)[:c] + "\n" + (*text)[c:]
		win.WidgetCursorSet(utf8NextCursor(*text, c))
		changed = true
	}

	// Navigation that depends on the wrapped layout: Up/Down move between
	// visual lines, Home/End snap to the edges of the current visual line.
	lines := textLines(*text, maxRef)
	li, in := cursorLine(lines, win.WidgetCursor())
	if win.KeyPressed(antui.KeyUp) {
		win.WidgetCursorSet(moveLineUp(*text, lines, li, in))
	}
	if win.KeyPressed(antui.KeyDown) {
		win.WidgetCursorSet(moveLineDown(*text, lines, li, in))
	}
	if win.KeyPressed(antui.KeyHome) {
		win.WidgetCursorSet(lines[li][0])
	}
	if win.KeyPressed(antui.KeyEnd) {
		win.WidgetCursorSet(lines[li][1])
	}

	if typed := win.TextInput(); typed != "" {
		c := win.WidgetCursor()
		*text = (*text)[:c] + typed + (*text)[c:]
		win.WidgetCursorSet(c + len(typed))
		changed = true
	}
	return changed
}

// caretAt maps a click at mx onto a caret position in text: each glyph's cell
// is split in half and snapped to its nearer edge, so a click lands where a
// text editor would put the caret. left is where the line's drawing starts
// and advance measures a text prefix exactly the way the box draws it.
// Positions before the first glyph return 0; past the last, len(text).
func caretAt(advance func(string) int, text string, left, mx int) int {
	if mx <= left {
		return 0
	}
	var (
		bounds []int // advance of each glyph boundary, in text order
		at     []int // byte index of each glyph boundary
	)
	for pos := 0; pos < len(text); {
		pos = utf8NextCursor(text, pos)
		bounds = append(bounds, advance(text[:pos]))
		at = append(at, pos)
	}
	for i, b := range bounds {
		prev := 0
		if i > 0 {
			prev = bounds[i-1]
		}
		edge := left + b
		if mx < edge {
			if i > 0 && mx < left+(prev+b)/2 {
				return at[i-1]
			}
			if i == 0 && mx < left+b/2 {
				return 0
			}
			return at[i]
		}
	}
	return len(text)
}

// textViewTop is the first line the box should show for cursor, so the caret
// stays visible as it moves. lines holds the wrapped layout, visible how many
// rows fit in the box.
func textViewTop(lines [][2]int, cursor, visible int) int {
	viewTop := 0
	for li, l := range lines {
		if cursor >= l[0] && cursor < l[1] || (li == len(lines)-1 && cursor == l[1]) {
			viewTop = max(li-visible+1, 0)
			break
		}
	}
	return max(min(viewTop, max(len(lines)-visible, 0)), 0)
}

// textAreaTextPos maps a click inside a wrapped multi-line text box onto the
// byte index of the character it lands on. It replicates where the box's
// drawing puts each line, so the caret lands exactly under the pointer.
// padding and lineH describe the box's inside margin and row height in
// window pixels, advance measures a prefix the way the box draws, and mx,my
// are the click in window coordinates.
func textAreaTextPos(text string, lines [][2]int, padding, x, y, h, cursor int, lineH int, advance func(string) int, mx, my int) int {
	visible := max((h-2*padding)/lineH, 1)
	viewTop := textViewTop(lines, cursor, visible)
	for lif, l := range lines {
		if lif < viewTop || lif >= viewTop+visible {
			continue
		}
		cy := y + padding + (lif-viewTop)*lineH
		if my >= cy && my < cy+lineH {
			return l[0] + caretAt(advance, text[l[0]:l[1]], x+padding, mx)
		}
	}
	return len(text)
}

// utf8PrevCursor and utf8NextCursor step a cursor by whole characters.
func utf8PrevCursor(text string, index int) int {
	if index <= 0 {
		return 0
	}
	index--
	for index > 0 && text[index]&0xC0 == 0x80 {
		index--
	}
	return index
}

func utf8NextCursor(text string, index int) int {
	if index >= len(text) {
		return len(text)
	}
	index++
	for index < len(text) && text[index]&0xC0 == 0x80 {
		index++
	}
	return index
}

// textLines breaks multiline text into wrapped visual lines. Each entry is
// the [start, end) byte range of one visual line; a '\n' ends a line and is
// not part of it. maxRef is the wrap width at reference scale.
func textLines(text string, maxRef int) [][2]int {
	var lines [][2]int
	i := 0
	for i <= len(text) {
		start := i
		w := 0
		brk := -1
		for i < len(text) {
			if text[i] == '\n' {
				brk = i
				break
			}
			r, size := utf8.DecodeRuneInString(text[i:])
			gw := canvas.TextWidth(string(r))
			if w+gw > maxRef && w > 0 {
				brk = i
				break
			}
			w += gw
			i += size
		}
		if brk == -1 {
			lines = append(lines, [2]int{start, len(text)})
			return lines
		}
		lines = append(lines, [2]int{start, brk})
		if text[brk] == '\n' {
			i = brk + 1
		} else {
			i = brk
		}
	}
	return lines
}

// cursorLine finds which visual line a byte cursor sits on and its offset
// within that line.
func cursorLine(lines [][2]int, cursor int) (li, in int) {
	for i, l := range lines {
		if cursor >= l[0] && cursor <= l[1] {
			return i, cursor - l[0]
		}
	}
	if len(lines) == 0 {
		return 0, 0
	}
	last := lines[len(lines)-1]
	return len(lines) - 1, last[1] - last[0]
}

// moveLineUp moves the cursor one visual line up, and moveLineDown one down,
// keeping the horizontal position as close as possible.
func moveLineUp(text string, lines [][2]int, li, in int) int {
	if li == 0 {
		return lines[0][0]
	}
	want := canvas.TextWidth(text[lines[li][0] : lines[li][0]+in])
	ls, le := lines[li-1][0], lines[li-1][1]
	best := ls
	for i := ls; i < le; i = utf8NextCursor(text, i) {
		n := utf8NextCursor(text, i)
		w := canvas.TextWidth(text[ls:n])
		if w <= want {
			best = n
		} else {
			break
		}
	}
	return best
}

func moveLineDown(text string, lines [][2]int, li, in int) int {
	if li >= len(lines)-1 {
		return lines[len(lines)-1][1]
	}
	want := canvas.TextWidth(text[lines[li][0] : lines[li][0]+in])
	ls, le := lines[li+1][0], lines[li+1][1]
	best := ls
	for i := ls; i < le; i = utf8NextCursor(text, i) {
		n := utf8NextCursor(text, i)
		w := canvas.TextWidth(text[ls:n])
		if w <= want {
			best = n
		} else {
			break
		}
	}
	return best
}

// Switch draws a flip toggle: a track with a thumb that slides between off
// and on. It flips value when clicked and reports whether it flipped.
func (t *uiTemplate) Switch(win *antui.Window, x, y int, label string, value *bool) event.Event {
	if value == nil {
		return event.Nothing
	}
	u := Scale(win)
	trackW := 36 * u
	trackH := 20 * u
	w := trackW + 8*u + textWidth(u, label)
	id := win.WidgetID("template:switch", x, y, trackW, trackH, label)
	hovered := win.Hovered(x, y, w, trackH)
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
	t.style.Switch(win, State{Hovered: hovered, Focused: focused, On: *value}, x, y, label, *value)
	if changed {
		return t.fire(event.Event{Component: event.Switch, Kind: event.Toggle})
	}
	return event.Nothing
}

// Progress paints a read-only progress bar with progress from 0 to 1. It
// reports nothing: a bar has no events.
func (t *uiTemplate) Progress(win *antui.Window, x, y, w, h int, progress float32) {
	t.style.Progress(win, x, y, w, h, progress)
}

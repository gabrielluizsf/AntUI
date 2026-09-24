package antui

// The exported helpers here are the same machinery the ready-made widgets
// work with, opened up for widgets that live outside this package — in
// particular the templates. A widget built elsewhere gets the same press,
// focus and typing behaviour the built-in ones have, down to sharing the hot,
// active and focused identities with them.

// WidgetID is the identity a widget's state is kept under. Two calls with the
// same kind, place and label are the same widget; changing any of them makes
// a new one. The built-in widgets use it, and a template's widgets should
// too, so that a template widget and a built-in one never mistake each other
// for the same thing.
func (win *Window) WidgetID(kind string, x, y, w, h int, label string) uint32 {
	return widgetID(kind, x, y, w, h, label)
}

// WidgetFrontmost defers draw until the very end of the frame, after every
// widget painted in flow order, so an open dropdown or calendar paints above
// whatever sits under it. The window runs the draws just before presenting,
// and forgets them at the next Begin, so a caller registers one every frame.
func (win *Window) WidgetFrontmost(draw func()) {
	win.frontmost = append(win.frontmost, draw)
}

// WidgetClick drives the press-and-release logic every clickable widget
// shares and reports whether it was clicked this frame. See [Window.Clicked].
func (win *Window) WidgetClick(id uint32, hovered bool) bool {
	return win.widgetClick(id, hovered)
}

// WidgetHot marks the widget as the one under the pointer, which is how a
// widget that is not drawn under the cursor still keeps others from acting on
// presses meant for it.
func (win *Window) WidgetHot(id uint32) { win.uiHot = id }

// WidgetActive is the widget that is in the middle of being pressed, or 0.
func (win *Window) WidgetActive() uint32 { return win.uiActive }

// WidgetStartPress begins a press on the widget. It takes the keyboard away
// from any text field, the same way the built-in widgets do when pressed.
func (win *Window) WidgetStartPress(id uint32) {
	win.uiActive = id
	win.uiFocus = 0
}

// WidgetEndPress ends the press on the widget, whatever it was.
func (win *Window) WidgetEndPress(id uint32) {
	if win.uiActive == id {
		win.uiActive = 0
	}
}

// WidgetFocus is the widget the keyboard is going to, or 0.
func (win *Window) WidgetFocus() uint32 { return win.uiFocus }

// WidgetFocusSet hands the keyboard to the widget.
func (win *Window) WidgetFocusSet(id uint32) {
	win.uiFocus = id
	win.uiBlink = 0
}

// WidgetFocusClear takes the keyboard away again.
func (win *Window) WidgetFocusClear(id uint32) {
	if win.uiFocus == id {
		win.uiFocus = 0
	}
}

// WidgetCursor is where the cursor sits in the focused text field.
func (win *Window) WidgetCursor() int { return win.uiCursor }

// WidgetCursorSet moves the cursor in the focused text field.
func (win *Window) WidgetCursorSet(c int) { win.uiCursor = c }

// WidgetCaret reports whether the text cursor is lit right now. It blinks on
// the same clock the built-in text field blinks on.
func (win *Window) WidgetCaret() bool { return int64(win.uiBlink*2)%2 == 0 }

// WidgetEdit applies this frame's keys and typing to the focused field,
// exactly as the built-in Input does.
func (win *Window) WidgetEdit(text *string) bool { return win.editText(text) }

// WidgetTabStop says the widget takes the keyboard, at its place in the
// frame's focus order. Widgets register themselves, in the order they draw,
// and Tab walks that order.
func (win *Window) WidgetTabStop(id uint32) { win.uiTab = append(win.uiTab, id) }

// WidgetKeyActivate reports whether Enter or Space was pressed this frame
// while the widget held the keyboard. A focused button, checkbox or radio is
// "clicked" by Enter or Space, the way the web does it.
func (win *Window) WidgetKeyActivate(id uint32) bool {
	if win.uiFocus != id || (win.uiActive != 0 && win.uiActive != id) {
		return false
	}
	return win.KeyPressed(KeyEnter) || win.KeyPressed(KeySpace)
}

// moveFocus walks the frame's focus order when Tab or Shift+Tab is pressed.
// It runs at the top of End, no matter how the frame was drawn, so the same
// frame that moved the focus draws the ring.
func (win *Window) moveFocus() {
	if !win.KeyPressed(KeyTab) {
		return
	}
	if len(win.uiTab) == 0 {
		return
	}
	shift := win.mods&ModShift != 0
	for i, id := range win.uiTab {
		if win.uiFocus != id {
			continue
		}
		if shift {
			i = (i - 1 + len(win.uiTab)) % len(win.uiTab)
		} else {
			i = (i + 1) % len(win.uiTab)
		}
		win.uiFocus = win.uiTab[i]
		win.uiCursor = 0
		win.uiBlink = 0
		return
	}
	if shift {
		win.uiFocus = win.uiTab[len(win.uiTab)-1]
	} else {
		win.uiFocus = win.uiTab[0]
	}
	win.uiCursor = 0
	win.uiBlink = 0
}

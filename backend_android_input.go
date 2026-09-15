//go:build android

package antui

import (
	"github.com/gabrielluizsf/antui/backend/android/ndk"
)

// drainInput takes everything the system has and turns it into events the
// window core understands. It runs once a frame, from pump.
func (b *androidWindow) drainInput(win *Window) {
	if b.activity == nil {
		return
	}
	b.activity.Input(func(e *ndk.InputEvent) bool {
		switch e.Kind() {
		case ndk.EventKey:
			return b.key(win, e)
		case ndk.EventMotion:
			return b.motion(win, e)
		}
		return false
	})
}

// motion sorts a pointer event by what produced it. The low bits of every
// pointer source are the same, so this cannot be a switch on the source: a
// mouse has to be recognised by its own bits before anything else, and
// everything left that points is treated as a finger.
func (b *androidWindow) motion(win *Window, e *ndk.InputEvent) bool {
	if e.Source().Is(ndk.SourceMouse) {
		return b.mouse(win, e)
	}
	if e.Source().Is(ndk.SourceJoystick) {
		// A controller's sticks arrive as motion on axes rather than as a
		// position. The buttons already work, through key events; the axes
		// need an API antui does not have yet.
		return false
	}
	return b.touch(win, e)
}

// touch turns one motion report into touch events, and drives the mouse from
// the first finger down.
//
// The action says what happened and, for the two "pointer" variants, which
// finger it happened to — as an **index into this event**, not as an id. A
// move reports every finger at once; everything else reports one.
func (b *androidWindow) touch(win *Window, e *ndk.InputEvent) bool {
	switch e.MotionAction() {
	case ndk.MotionDown, ndk.MotionPointerDown:
		b.touchAt(win, e, e.MotionIndex(), EventTouchDown)
	case ndk.MotionUp, ndk.MotionPointerUp:
		b.touchAt(win, e, e.MotionIndex(), EventTouchUp)
	case ndk.MotionMove:
		for i := range e.PointerCount() {
			b.touchAt(win, e, i, EventTouchMove)
		}
	case ndk.MotionCancel:
		for i := range e.PointerCount() {
			b.touchAt(win, e, i, EventTouchCancel)
		}
	default:
		return false
	}
	return true
}

// touchAt pushes one finger's event, and the mouse event that shadows it.
func (b *androidWindow) touchAt(win *Window, e *ndk.InputEvent, index int, kind EventType) {
	if index < 0 || index >= e.PointerCount() {
		return
	}
	id := e.PointerID(index)
	x, y := int(e.X(index)), int(e.Y(index))

	win.Push(Event{
		Type:      kind,
		TouchID:   id,
		X:         x,
		Y:         y,
		Pressure:  e.Pressure(index),
		TouchSize: e.Size(index),
		Tool:      toolOf(e.ToolType(index)),
	})

	// One finger — the first one down — is also reported as the left mouse
	// button, so that a program written for a desktop works on a phone with
	// nothing added. The rest are touches only: a mouse has one position,
	// and pretending otherwise would make two fingers look like one jumping
	// between them.
	switch kind {
	case EventTouchDown:
		if b.primary >= 0 {
			return
		}
		b.primary = id
		win.Push(Event{Type: EventMouseDown, Button: MouseLeft, X: x, Y: y})
	case EventTouchMove:
		if id != b.primary {
			return
		}
		win.Push(Event{Type: EventMouseMove, X: x, Y: y})
	case EventTouchUp, EventTouchCancel:
		if id != b.primary {
			return
		}
		b.primary = -1
		win.Push(Event{Type: EventMouseUp, Button: MouseLeft, X: x, Y: y})
	}
}

func toolOf(t int) Tool {
	switch t {
	case ndk.ToolFinger:
		return ToolFinger
	case ndk.ToolStylus:
		return ToolStylus
	case ndk.ToolMouse:
		return ToolMouse
	case ndk.ToolEraser:
		return ToolEraser
	}
	return ToolUnknown
}

// mouse handles a real mouse or a trackpad — a Chromebook, a phone in
// desktop mode, a tablet with a keyboard case.
//
// Android reports which buttons are *held*, not which one changed, so the
// change has to be worked out by comparing with the last report. Getting
// this wrong produces a button that sticks down, which looks like a bug in
// the app rather than in the translation.
func (b *androidWindow) mouse(win *Window, e *ndk.InputEvent) bool {
	x, y := int(e.X(0)), int(e.Y(0))
	switch e.MotionAction() {
	case ndk.MotionScroll:
		// A wheel step is 1.0 per notch, and positive is away from the user,
		// which is the same direction antui calls up.
		if v := e.Axis(ndk.AxisVScroll, 0); v != 0 {
			win.Push(Event{Type: EventMouseWheel, Wheel: int(v), X: x, Y: y})
		}
		return true

	case ndk.MotionMove, ndk.MotionHoverMove:
		win.Push(Event{Type: EventMouseMove, X: x, Y: y, Mods: modsOf(e.Meta())})
		return true

	case ndk.MotionDown, ndk.MotionUp, ndk.MotionButtonPress, ndk.MotionButtonUp:
		now := e.Buttons()
		changed := now ^ b.buttons
		b.buttons = now
		for _, m := range []struct {
			bit    int32
			button MouseButton
		}{
			{ndk.ButtonPrimary, MouseLeft},
			{ndk.ButtonSecondary, MouseRight},
			{ndk.ButtonTertiary, MouseMiddle},
		} {
			if changed&m.bit == 0 {
				continue
			}
			kind := EventMouseUp
			if now&m.bit != 0 {
				kind = EventMouseDown
			}
			win.Push(Event{Type: kind, Button: m.button, X: x, Y: y, Mods: modsOf(e.Meta())})
		}
		return true
	}
	return false
}

// key turns a key event into a key event.
func (b *androidWindow) key(win *Window, e *ndk.InputEvent) bool {
	code := e.KeyCode()
	k := keyOf(code)
	mods := modsOf(e.Meta())
	down := e.KeyAction() == ndk.KeyDown

	// The back button is the system's "go back". Treating it as a close is
	// what the platform does when nothing handles it, and matches what the
	// window's close button does on a desktop — so a program that already
	// copes with being closed copes with this.
	if code == ndk.KeyCodeBack && down {
		win.Push(Event{Type: EventKeyDown, Key: KeyBack, Mods: mods})
		win.Push(Event{Type: EventClose})
		return true
	}
	// Volume is the system's, always. Taking it would leave the user unable
	// to turn the sound down, which no app is entitled to.
	if code == ndk.KeyCodeVolumeUp || code == ndk.KeyCodeVolumeDown {
		return false
	}
	if k == KeyUnknown {
		return false
	}

	if down {
		win.Push(Event{
			Type:   EventKeyDown,
			Key:    k,
			Mods:   mods,
			Repeat: e.KeyRepeat() > 0,
		})
		if text := textOf(k, e.Meta()); text != "" {
			win.Push(Event{Type: EventText, Text: text, Rune: []rune(text)[0]})
		}
	} else {
		win.Push(Event{Type: EventKeyUp, Key: k, Mods: mods})
	}
	return true
}

func modsOf(m ndk.Meta) Mod {
	var mods Mod
	if m.Has(ndk.MetaShift) {
		mods |= ModShift
	}
	if m.Has(ndk.MetaCtrl) {
		mods |= ModControl
	}
	if m.Has(ndk.MetaAlt) {
		mods |= ModAlt
	}
	if m.Has(ndk.MetaSuper) {
		mods |= ModSuper
	}
	return mods
}

// textOf is what a key types, for a hardware keyboard.
//
// It is deliberately small. Android knows what a key produces on the user's
// layout, but only through Java's KeyCharacterMap — there is no NDK call for
// it — so this covers the letters, the digits and space, which is a US
// layout and nothing more. Anything typed on a screen keyboard, in any other
// layout, or in a language that composes goes through the IME instead, and
// that needs the JNI layer.
func textOf(k Key, meta ndk.Meta) string {
	if meta.Has(ndk.MetaCtrl) || meta.Has(ndk.MetaAlt) || meta.Has(ndk.MetaSuper) {
		return ""
	}
	upper := meta.Has(ndk.MetaShift) != meta.Has(ndk.MetaCaps)
	switch {
	case k >= KeyA && k <= KeyZ:
		if upper {
			return string(rune(k))
		}
		return string(rune(k) + 32)
	case k >= Key0 && k <= Key9 && !meta.Has(ndk.MetaShift):
		return string(rune(k))
	case k == KeySpace:
		return " "
	}
	return ""
}

// keyOf maps Android's key code onto antui's.
//
// Android numbers keys by what is printed on them rather than by where they
// sit, so this is a straight table and not a layout question. The three
// ranges are contiguous on both sides, which is the only reason they are
// arithmetic rather than another sixty lines.
func keyOf(code int) Key {
	switch {
	case code >= ndk.KeyCodeA && code <= ndk.KeyCodeZ:
		return KeyA + Key(code-ndk.KeyCodeA)
	case code >= ndk.KeyCode0 && code <= ndk.KeyCode9:
		return Key0 + Key(code-ndk.KeyCode0)
	case code >= ndk.KeyCodeF1 && code <= ndk.KeyCodeF12:
		return KeyF1 + Key(code-ndk.KeyCodeF1)
	}
	switch code {
	case ndk.KeyCodeSpace:
		return KeySpace
	case ndk.KeyCodeEnter:
		return KeyEnter
	case ndk.KeyCodeTab:
		return KeyTab
	case ndk.KeyCodeDel:
		return KeyBackspace // "DEL" is the platform's name for backspace
	case ndk.KeyCodeForwardDel:
		return KeyDelete
	case ndk.KeyCodeEscape:
		return KeyEscape
	case ndk.KeyCodeInsert:
		return KeyInsert
	case ndk.KeyCodePageUp:
		return KeyPageUp
	case ndk.KeyCodePageDown:
		return KeyPageDown
	case ndk.KeyCodeMoveHome:
		return KeyHome
	case ndk.KeyCodeMoveEnd:
		return KeyEnd
	case ndk.KeyCodeCapsLock:
		return KeyCapsLock
	case ndk.KeyCodeDpadUp:
		return KeyUp
	case ndk.KeyCodeDpadDown:
		return KeyDown
	case ndk.KeyCodeDpadLeft:
		return KeyLeft
	case ndk.KeyCodeDpadRight:
		return KeyRight
	case ndk.KeyCodeDpadCenter:
		return KeyEnter
	case ndk.KeyCodeShiftLeft:
		return KeyLeftShift
	case ndk.KeyCodeShiftRight:
		return KeyRightShift
	case ndk.KeyCodeCtrlLeft:
		return KeyLeftControl
	case ndk.KeyCodeCtrlRight:
		return KeyRightControl
	case ndk.KeyCodeAltLeft:
		return KeyLeftAlt
	case ndk.KeyCodeAltRight:
		return KeyRightAlt
	case ndk.KeyCodeMetaLeft:
		return KeyLeftSuper
	case ndk.KeyCodeMetaRight:
		return KeyRightSuper
	case ndk.KeyCodeComma:
		return KeyComma
	case ndk.KeyCodePeriod:
		return KeyPeriod
	case ndk.KeyCodeMinus:
		return KeyMinus
	case ndk.KeyCodeEquals:
		return KeyEqual
	case ndk.KeyCodeLeftBracket:
		return KeyLeftBracket
	case ndk.KeyCodeRightBracket:
		return KeyRightBracket
	case ndk.KeyCodeBackslash:
		return KeyBackslash
	case ndk.KeyCodeSemicolon:
		return KeySemicolon
	case ndk.KeyCodeApostrophe:
		return KeyApostrophe
	case ndk.KeyCodeSlash:
		return KeySlash
	case ndk.KeyCodeGrave:
		return KeyGrave
	case ndk.KeyCodeMenu:
		return KeyMenu
	case ndk.KeyCodeSearch:
		return KeySearch
	case ndk.KeyCodeButtonA:
		return KeyPadA
	case ndk.KeyCodeButtonB:
		return KeyPadB
	case ndk.KeyCodeButtonX:
		return KeyPadX
	case ndk.KeyCodeButtonY:
		return KeyPadY
	case ndk.KeyCodeButtonL1:
		return KeyPadL1
	case ndk.KeyCodeButtonR1:
		return KeyPadR1
	case ndk.KeyCodeButtonL2:
		return KeyPadL2
	case ndk.KeyCodeButtonR2:
		return KeyPadR2
	case ndk.KeyCodeButtonStart:
		return KeyPadStart
	case ndk.KeyCodeButtonSelect:
		return KeyPadSelect
	case ndk.KeyCodeButtonThumbL:
		return KeyPadThumbL
	case ndk.KeyCodeButtonThumbR:
		return KeyPadThumbR
	case ndk.KeyCodeButtonMode:
		return KeyPadMode
	}
	return KeyUnknown
}

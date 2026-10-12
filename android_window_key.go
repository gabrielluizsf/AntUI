//go:build android

package antui

import (
	"github.com/gabrielluizsf/antui/backend/android/ndk"
)

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

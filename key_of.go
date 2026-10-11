//go:build android

package antui

import (
	"github.com/gabrielluizsf/antui/backend/android/ndk"
)

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
	if k := keyOfNav(code); k != KeyUnknown {
		return k
	}
	if k := keyOfModifier(code); k != KeyUnknown {
		return k
	}
	if k := keyOfPunct(code); k != KeyUnknown {
		return k
	}
	return keyOfPad(code)
}

//go:build android

package antui

import "github.com/gabrielluizsf/antui/backend/android/ndk"

// keyOfNav maps the navigation and editing keys.
func keyOfNav(code int) Key {
	switch code {
	case ndk.KeyCodeSpace:
		return KeySpace
	case ndk.KeyCodeEnter, ndk.KeyCodeDpadCenter:
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
	}
	return KeyUnknown
}

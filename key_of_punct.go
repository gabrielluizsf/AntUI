//go:build android

package antui

import "github.com/gabrielluizsf/antui/backend/android/ndk"

// keyOfPunct maps the punctuation keys.
func keyOfPunct(code int) Key {
	switch code {
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
	}
	return KeyUnknown
}

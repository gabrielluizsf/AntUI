//go:build android

package antui

import "github.com/gabrielluizsf/antui/backend/android/ndk"

// keyOfModifier maps the shift, control, alt and meta keys.
func keyOfModifier(code int) Key {
	switch code {
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
	}
	return KeyUnknown
}

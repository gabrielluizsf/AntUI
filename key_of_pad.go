//go:build android

package antui

import "github.com/gabrielluizsf/antui/backend/android/ndk"

// keyOfPad maps the menu, search and gamepad keys.
func keyOfPad(code int) Key {
	switch code {
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

//go:build android

package ndk

/*
#include <android/keycodes.h>
*/
import "C"

// The key codes this library maps. They are taken from the platform's own
// header rather than written out as numbers, so a value cannot drift.
//
// Android numbers keys by what is printed on them, not by where they sit, so
// there is no scan-code business here: AKEYCODE_A is the A key on any
// layout. The list is the subset antui has a Key for, plus the ones only a
// phone has.
const (
	KeyCodeUnknown = int(C.AKEYCODE_UNKNOWN)

	// The buttons a phone has and a desktop does not.
	KeyCodeBack       = int(C.AKEYCODE_BACK)
	KeyCodeHome       = int(C.AKEYCODE_HOME)
	KeyCodeMenu       = int(C.AKEYCODE_MENU)
	KeyCodeSearch     = int(C.AKEYCODE_SEARCH)
	KeyCodeVolumeUp   = int(C.AKEYCODE_VOLUME_UP)
	KeyCodeVolumeDown = int(C.AKEYCODE_VOLUME_DOWN)
	KeyCodePower      = int(C.AKEYCODE_POWER)

	KeyCodeDpadUp     = int(C.AKEYCODE_DPAD_UP)
	KeyCodeDpadDown   = int(C.AKEYCODE_DPAD_DOWN)
	KeyCodeDpadLeft   = int(C.AKEYCODE_DPAD_LEFT)
	KeyCodeDpadRight  = int(C.AKEYCODE_DPAD_RIGHT)
	KeyCodeDpadCenter = int(C.AKEYCODE_DPAD_CENTER)

	KeyCode0 = int(C.AKEYCODE_0)
	KeyCode9 = int(C.AKEYCODE_9)
	KeyCodeA = int(C.AKEYCODE_A)
	KeyCodeZ = int(C.AKEYCODE_Z)

	KeyCodeSpace        = int(C.AKEYCODE_SPACE)
	KeyCodeEnter        = int(C.AKEYCODE_ENTER)
	KeyCodeTab          = int(C.AKEYCODE_TAB)
	KeyCodeDel          = int(C.AKEYCODE_DEL)         // backspace, despite the name
	KeyCodeForwardDel   = int(C.AKEYCODE_FORWARD_DEL) // delete
	KeyCodeEscape       = int(C.AKEYCODE_ESCAPE)
	KeyCodeInsert       = int(C.AKEYCODE_INSERT)
	KeyCodePageUp       = int(C.AKEYCODE_PAGE_UP)
	KeyCodePageDown     = int(C.AKEYCODE_PAGE_DOWN)
	KeyCodeMoveHome     = int(C.AKEYCODE_MOVE_HOME)
	KeyCodeMoveEnd      = int(C.AKEYCODE_MOVE_END)
	KeyCodeCapsLock     = int(C.AKEYCODE_CAPS_LOCK)
	KeyCodeShiftLeft    = int(C.AKEYCODE_SHIFT_LEFT)
	KeyCodeShiftRight   = int(C.AKEYCODE_SHIFT_RIGHT)
	KeyCodeAltLeft      = int(C.AKEYCODE_ALT_LEFT)
	KeyCodeAltRight     = int(C.AKEYCODE_ALT_RIGHT)
	KeyCodeCtrlLeft     = int(C.AKEYCODE_CTRL_LEFT)
	KeyCodeCtrlRight    = int(C.AKEYCODE_CTRL_RIGHT)
	KeyCodeMetaLeft     = int(C.AKEYCODE_META_LEFT)
	KeyCodeMetaRight    = int(C.AKEYCODE_META_RIGHT)
	KeyCodeComma        = int(C.AKEYCODE_COMMA)
	KeyCodePeriod       = int(C.AKEYCODE_PERIOD)
	KeyCodeMinus        = int(C.AKEYCODE_MINUS)
	KeyCodeEquals       = int(C.AKEYCODE_EQUALS)
	KeyCodeLeftBracket  = int(C.AKEYCODE_LEFT_BRACKET)
	KeyCodeRightBracket = int(C.AKEYCODE_RIGHT_BRACKET)
	KeyCodeBackslash    = int(C.AKEYCODE_BACKSLASH)
	KeyCodeSemicolon    = int(C.AKEYCODE_SEMICOLON)
	KeyCodeApostrophe   = int(C.AKEYCODE_APOSTROPHE)
	KeyCodeSlash        = int(C.AKEYCODE_SLASH)
	KeyCodeGrave        = int(C.AKEYCODE_GRAVE)

	KeyCodeF1  = int(C.AKEYCODE_F1)
	KeyCodeF12 = int(C.AKEYCODE_F12)

	// A game controller's buttons, in the order the platform numbers them.
	KeyCodeButtonA      = int(C.AKEYCODE_BUTTON_A)
	KeyCodeButtonB      = int(C.AKEYCODE_BUTTON_B)
	KeyCodeButtonX      = int(C.AKEYCODE_BUTTON_X)
	KeyCodeButtonY      = int(C.AKEYCODE_BUTTON_Y)
	KeyCodeButtonL1     = int(C.AKEYCODE_BUTTON_L1)
	KeyCodeButtonR1     = int(C.AKEYCODE_BUTTON_R1)
	KeyCodeButtonL2     = int(C.AKEYCODE_BUTTON_L2)
	KeyCodeButtonR2     = int(C.AKEYCODE_BUTTON_R2)
	KeyCodeButtonThumbL = int(C.AKEYCODE_BUTTON_THUMBL)
	KeyCodeButtonThumbR = int(C.AKEYCODE_BUTTON_THUMBR)
	KeyCodeButtonStart  = int(C.AKEYCODE_BUTTON_START)
	KeyCodeButtonSelect = int(C.AKEYCODE_BUTTON_SELECT)
	KeyCodeButtonMode   = int(C.AKEYCODE_BUTTON_MODE)
)

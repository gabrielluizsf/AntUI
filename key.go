package antui

import "github.com/gabrielluizsf/antui/backend"

// Key identifies a physical key, independent of the text it produces.
// Printable keys carry their own uppercase ASCII code, so KeyA is 'A' and
// Key1 is '1'; everything else is numbered from 256 up.
type Key = backend.Key

// Printable keys.
const (
	KeyUnknown      = backend.KeyUnknown
	KeySpace        = backend.KeySpace
	KeyApostrophe   = backend.KeyApostrophe
	KeyComma        = backend.KeyComma
	KeyMinus        = backend.KeyMinus
	KeyPeriod       = backend.KeyPeriod
	KeySlash        = backend.KeySlash
	Key0            = backend.Key0
	Key1            = backend.Key1
	Key2            = backend.Key2
	Key3            = backend.Key3
	Key4            = backend.Key4
	Key5            = backend.Key5
	Key6            = backend.Key6
	Key7            = backend.Key7
	Key8            = backend.Key8
	Key9            = backend.Key9
	KeySemicolon    = backend.KeySemicolon
	KeyEqual        = backend.KeyEqual
	KeyA            = backend.KeyA
	KeyB            = backend.KeyB
	KeyC            = backend.KeyC
	KeyD            = backend.KeyD
	KeyE            = backend.KeyE
	KeyF            = backend.KeyF
	KeyG            = backend.KeyG
	KeyH            = backend.KeyH
	KeyI            = backend.KeyI
	KeyJ            = backend.KeyJ
	KeyK            = backend.KeyK
	KeyL            = backend.KeyL
	KeyM            = backend.KeyM
	KeyN            = backend.KeyN
	KeyO            = backend.KeyO
	KeyP            = backend.KeyP
	KeyQ            = backend.KeyQ
	KeyR            = backend.KeyR
	KeyS            = backend.KeyS
	KeyT            = backend.KeyT
	KeyU            = backend.KeyU
	KeyV            = backend.KeyV
	KeyW            = backend.KeyW
	KeyX            = backend.KeyX
	KeyY            = backend.KeyY
	KeyZ            = backend.KeyZ
	KeyLeftBracket  = backend.KeyLeftBracket
	KeyBackslash    = backend.KeyBackslash
	KeyRightBracket = backend.KeyRightBracket
	KeyGrave        = backend.KeyGrave
)

// Special keys.
const (
	KeyEscape         = backend.KeyEscape
	KeyEnter          = backend.KeyEnter
	KeyTab            = backend.KeyTab
	KeyBackspace      = backend.KeyBackspace
	KeyInsert         = backend.KeyInsert
	KeyDelete         = backend.KeyDelete
	KeyRight          = backend.KeyRight
	KeyLeft           = backend.KeyLeft
	KeyDown           = backend.KeyDown
	KeyUp             = backend.KeyUp
	KeyPageUp         = backend.KeyPageUp
	KeyPageDown       = backend.KeyPageDown
	KeyHome           = backend.KeyHome
	KeyEnd            = backend.KeyEnd
	KeyCapsLock       = backend.KeyCapsLock
	KeyF1             = backend.KeyF1
	KeyF2             = backend.KeyF2
	KeyF3             = backend.KeyF3
	KeyF4             = backend.KeyF4
	KeyF5             = backend.KeyF5
	KeyF6             = backend.KeyF6
	KeyF7             = backend.KeyF7
	KeyF8             = backend.KeyF8
	KeyF9             = backend.KeyF9
	KeyF10            = backend.KeyF10
	KeyF11            = backend.KeyF11
	KeyF12            = backend.KeyF12
	KeyLeftShift      = backend.KeyLeftShift
	KeyLeftControl    = backend.KeyLeftControl
	KeyLeftAlt        = backend.KeyLeftAlt
	KeyLeftSuper      = backend.KeyLeftSuper
	KeyRightShift     = backend.KeyRightShift
	KeyRightControl   = backend.KeyRightControl
	KeyRightAlt       = backend.KeyRightAlt
	KeyRightSuper     = backend.KeyRightSuper

	KeyBack        = backend.KeyBack
	KeyMenu        = backend.KeyMenu
	KeySearch      = backend.KeySearch
	KeyVolumeUp    = backend.KeyVolumeUp
	KeyVolumeDown  = backend.KeyVolumeDown

	KeyPadA       = backend.KeyPadA
	KeyPadB       = backend.KeyPadB
	KeyPadX       = backend.KeyPadX
	KeyPadY       = backend.KeyPadY
	KeyPadL1      = backend.KeyPadL1
	KeyPadR1      = backend.KeyPadR1
	KeyPadL2      = backend.KeyPadL2
	KeyPadR2      = backend.KeyPadR2
	KeyPadStart   = backend.KeyPadStart
	KeyPadSelect  = backend.KeyPadSelect
	KeyPadThumbL  = backend.KeyPadThumbL
	KeyPadThumbR  = backend.KeyPadThumbR
	KeyPadMode    = backend.KeyPadMode
)

// keyCount bounds the key-state tables. Keys outside it are ignored rather
// than indexed, so an unmapped keysym cannot walk off the array.
const keyCount = 320

// Mod is a set of held modifier keys.
type Mod = backend.Mod

// The modifiers, combined with |.
const (
	ModShift   = backend.ModShift
	ModControl = backend.ModControl
	ModAlt     = backend.ModAlt
	ModSuper   = backend.ModSuper
)

// MouseButton identifies a mouse button.
type MouseButton = backend.MouseButton

// The mouse buttons.
const (
	MouseLeft   = backend.MouseLeft
	MouseRight  = backend.MouseRight
	MouseMiddle = backend.MouseMiddle
	// mouseCount bounds the mouse-state tables. It comes after the buttons.
	mouseCount = 3
)
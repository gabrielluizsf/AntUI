package backend

// Key identifies a physical key, independent of the text it produces.
// Printable keys carry their own uppercase ASCII code, so KeyA is 'A' and
// Key1 is '1'; everything else is numbered from 256 up.
type Key int

// Printable keys.
const (
	KeyUnknown      Key = 0
	KeySpace        Key = ' '
	KeyApostrophe   Key = '\''
	KeyComma        Key = ','
	KeyMinus        Key = '-'
	KeyPeriod       Key = '.'
	KeySlash        Key = '/'
	Key0            Key = '0'
	Key1            Key = '1'
	Key2            Key = '2'
	Key3            Key = '3'
	Key4            Key = '4'
	Key5            Key = '5'
	Key6            Key = '6'
	Key7            Key = '7'
	Key8            Key = '8'
	Key9            Key = '9'
	KeySemicolon    Key = ';'
	KeyEqual        Key = '='
	KeyA            Key = 'A'
	KeyB            Key = 'B'
	KeyC            Key = 'C'
	KeyD            Key = 'D'
	KeyE            Key = 'E'
	KeyF            Key = 'F'
	KeyG            Key = 'G'
	KeyH            Key = 'H'
	KeyI            Key = 'I'
	KeyJ            Key = 'J'
	KeyK            Key = 'K'
	KeyL            Key = 'L'
	KeyM            Key = 'M'
	KeyN            Key = 'N'
	KeyO            Key = 'O'
	KeyP            Key = 'P'
	KeyQ            Key = 'Q'
	KeyR            Key = 'R'
	KeyS            Key = 'S'
	KeyT            Key = 'T'
	KeyU            Key = 'U'
	KeyV            Key = 'V'
	KeyW            Key = 'W'
	KeyX            Key = 'X'
	KeyY            Key = 'Y'
	KeyZ            Key = 'Z'
	KeyLeftBracket  Key = '['
	KeyBackslash    Key = '\\'
	KeyRightBracket Key = ']'
	KeyGrave        Key = '`'
)

// Special keys.
const (
	KeyEscape Key = 256 + iota
	KeyEnter
	KeyTab
	KeyBackspace
	KeyInsert
	KeyDelete
	KeyRight
	KeyLeft
	KeyDown
	KeyUp
	KeyPageUp
	KeyPageDown
	KeyHome
	KeyEnd
	KeyCapsLock
	KeyF1
	KeyF2
	KeyF3
	KeyF4
	KeyF5
	KeyF6
	KeyF7
	KeyF8
	KeyF9
	KeyF10
	KeyF11
	KeyF12
	KeyLeftShift
	KeyLeftControl
	KeyLeftAlt
	KeyLeftSuper
	KeyRightShift
	KeyRightControl
	KeyRightAlt
	KeyRightSuper

	// Keys a phone has and a desktop does not. KeyBack is the one that
	// matters: it is the system's "go back", and a program that does not
	// handle it gets the platform's answer, which is to close.
	KeyBack
	KeyMenu
	KeySearch
	KeyVolumeUp
	KeyVolumeDown

	// A game controller's buttons, named the way the platform names them
	// rather than the way any one pad is printed — what is drawn on the
	// button varies by manufacturer and the code does not.
	KeyPadA
	KeyPadB
	KeyPadX
	KeyPadY
	KeyPadL1
	KeyPadR1
	KeyPadL2
	KeyPadR2
	KeyPadStart
	KeyPadSelect
	KeyPadThumbL
	KeyPadThumbR
	KeyPadMode
)

// Mod is a set of held modifier keys.
type Mod int

// The modifiers, combined with |.
const (
	ModShift Mod = 1 << iota
	ModControl
	ModAlt
	ModSuper
)

// Has reports whether every modifier in want is held.
func (m Mod) Has(want Mod) bool { return m&want == want }

// MouseButton identifies a mouse button.
type MouseButton int

// The mouse buttons.
const (
	MouseLeft MouseButton = iota
	MouseRight
	MouseMiddle
)
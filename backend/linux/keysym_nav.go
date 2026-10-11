//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import "github.com/gabrielluizsf/antui/backend"

// keysymNav maps the navigation and editing keys.
func keysymNav(sym uint32) backend.Key {
	switch sym {
	case ' ':
		return backend.KeySpace
	case symBackspace:
		return backend.KeyBackspace
	case symTab:
		return backend.KeyTab
	case symReturn, symKPEnter:
		return backend.KeyEnter
	case symEscape:
		return backend.KeyEscape
	case symHome:
		return backend.KeyHome
	case symLeft:
		return backend.KeyLeft
	case symUp:
		return backend.KeyUp
	case symRight:
		return backend.KeyRight
	case symDown:
		return backend.KeyDown
	case symPageUp:
		return backend.KeyPageUp
	case symPageDown:
		return backend.KeyPageDown
	case symEnd:
		return backend.KeyEnd
	case symInsert:
		return backend.KeyInsert
	case symDelete:
		return backend.KeyDelete
	}
	return backend.KeyUnknown
}

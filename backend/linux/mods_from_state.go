//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import (
	"github.com/gabrielluizsf/antui/backend"
)

func modsFromState(state uint16) backend.Mod {
	var mods backend.Mod
	if state&stateShift != 0 {
		mods |= backend.ModShift
	}
	if state&stateCtrl != 0 {
		mods |= backend.ModControl
	}
	if state&stateAlt != 0 {
		mods |= backend.ModAlt
	}
	if state&stateSuper != 0 {
		mods |= backend.ModSuper
	}
	return mods
}

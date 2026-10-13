//go:build windows

package windows

import "github.com/gabrielluizsf/antui/backend"

func winMods() backend.Mod {
	var mods backend.Mod
	down := func(vk uintptr) bool {
		state, _, _ := procGetKeyState.Call(vk)
		return state&0x8000 != 0
	}
	if down(vkShift) {
		mods |= backend.ModShift
	}
	if down(vkControl) {
		mods |= backend.ModControl
	}
	if down(vkMenu) {
		mods |= backend.ModAlt
	}
	if down(vkLWin) || down(vkRWin) {
		mods |= backend.ModSuper
	}
	return mods
}

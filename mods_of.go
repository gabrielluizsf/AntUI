//go:build android

package antui

import "github.com/gabrielluizsf/antui/backend/android/ndk"

func modsOf(m ndk.Meta) Mod {
	var mods Mod
	if m.Has(ndk.MetaShift) {
		mods |= ModShift
	}
	if m.Has(ndk.MetaCtrl) {
		mods |= ModControl
	}
	if m.Has(ndk.MetaAlt) {
		mods |= ModAlt
	}
	if m.Has(ndk.MetaSuper) {
		mods |= ModSuper
	}
	return mods
}

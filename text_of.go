//go:build android

package antui

import "github.com/gabrielluizsf/antui/backend/android/ndk"

// textOf is what a key types, for a hardware keyboard.
//
// It is deliberately small. Android knows what a key produces on the user's
// layout, but only through Java's KeyCharacterMap — there is no NDK call for
// it — so this covers the letters, the digits and space, which is a US
// layout and nothing more. Anything typed on a screen keyboard, in any other
// layout, or in a language that composes goes through the IME instead, and
// that needs the JNI layer.
func textOf(k Key, meta ndk.Meta) string {
	if meta.Has(ndk.MetaCtrl) || meta.Has(ndk.MetaAlt) || meta.Has(ndk.MetaSuper) {
		return ""
	}
	upper := meta.Has(ndk.MetaShift) != meta.Has(ndk.MetaCaps)
	switch {
	case k >= KeyA && k <= KeyZ:
		if upper {
			return string(rune(k))
		}
		return string(rune(k) + 32)
	case k >= Key0 && k <= Key9 && !meta.Has(ndk.MetaShift):
		return string(rune(k))
	case k == KeySpace:
		return " "
	}
	return ""
}

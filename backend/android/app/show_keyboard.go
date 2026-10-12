//go:build android

package app

/*
#cgo LDFLAGS: -landroid

#include <android/configuration.h>
#include <android/native_activity.h>
*/
import "C"

// ShowKeyboard asks for the on-screen keyboard.
//
// This is the whole of what the NDK offers: it can be asked for and asked to
// go away. **How tall it is, and what is typed on it, are not here** — the
// height is a window inset and the text goes through an input connection,
// and both of those are Java. Until the JNI layer lands, a native app can
// raise the keyboard and will not hear a word of what it produces.
func ShowKeyboard(flags uint32) {
	stateMu.Lock()
	act := activity
	stateMu.Unlock()
	if act != nil {
		C.ANativeActivity_showSoftInput(act, C.uint32_t(flags))
	}
}

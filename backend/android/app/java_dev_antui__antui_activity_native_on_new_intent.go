//go:build android

package app

/*
#include <jni.h>
*/
import "C"

import (
	"github.com/gabrielluizsf/antui/backend/android/jni"
	"github.com/gabrielluizsf/antui/backend/android/ndk"
)

// Java_dev_antui_AntuiActivity_nativeOnNewIntent is an intent arriving at an
// activity that is already running: a deep link, a notification being
// tapped, a file being shared to the app.
//
//export Java_dev_antui_AntuiActivity_nativeOnNewIntent
func Java_dev_antui_AntuiActivity_nativeOnNewIntent(env *C.JNIEnv, cls C.jclass, intent C.jobject) {
	intentMu.Lock()
	f := intentFn
	intentMu.Unlock()
	if f == nil {
		return
	}
	var held jni.Object
	if err := jni.Do(func(e *jni.Env) error {
		if i := jni.Wrap(uintptr(intent)); !i.IsNil() {
			held = e.Global(i)
		}
		return nil
	}); err != nil {
		ndk.Errorf("cannot hold on to the intent: %v", err)
		return
	}
	// On a goroutine of its own: this is the UI thread, and the handler is
	// the app's code, which may take as long as it likes and may want the UI
	// thread itself.
	go f(held)
}

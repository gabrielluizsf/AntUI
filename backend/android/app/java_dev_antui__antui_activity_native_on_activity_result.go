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

// Java_dev_antui_AntuiActivity_nativeOnActivityResult is the answer to
// [StartForResult].
//
//export Java_dev_antui_AntuiActivity_nativeOnActivityResult
func Java_dev_antui_AntuiActivity_nativeOnActivityResult(env *C.JNIEnv, cls C.jclass, code C.jint, result C.jint, data C.jobject) {
	out := ActivityResult{Code: int(result)}
	// The Intent is a local reference belonging to this call. Whoever is
	// waiting will read it later and on another thread, so it has to be made
	// global first — and they have to give it back.
	if err := jni.Do(func(e *jni.Env) error {
		if d := jni.Wrap(uintptr(data)); !d.IsNil() {
			out.Data = e.Global(d)
		}
		return nil
	}); err != nil {
		ndk.Errorf("cannot hold on to the activity result: %v", err)
	}

	waitMu.Lock()
	ch := actWait[int32(code)]
	delete(actWait, int32(code))
	waitMu.Unlock()
	if ch != nil {
		ch <- out
		return
	}
	// Nobody was waiting, so nobody will free it.
	if !out.Data.IsNil() {
		jni.Do(func(e *jni.Env) error { e.DeleteGlobal(out.Data); return nil })
	}
}

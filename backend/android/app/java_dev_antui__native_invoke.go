//go:build android

package app

/*
#include <jni.h>
#include "glue.h"
*/
import "C"

import (
	"github.com/gabrielluizsf/antui/backend/android/jni"
	"github.com/gabrielluizsf/antui/backend/android/ndk"
)

// Java_dev_antui_Native_invoke is where every Java listener in the shim ends
// up. It always answers null; see [Handler].
//
//export Java_dev_antui_Native_invoke
func Java_dev_antui_Native_invoke(env *C.JNIEnv, cls C.jclass, token C.jlong, method C.jstring, args C.jobjectArray) C.jobject {
	handlerMu.Lock()
	h := handlers[int64(token)]
	handlerMu.Unlock()
	if h == nil {
		// A listener that outlived its handler. It is not an error — the
		// platform may deliver one more after it was unregistered — but it
		// is worth seeing when a whole feature has gone quiet.
		ndk.Warnf("a Java listener called back with no handler (token %d)", int64(token))
		return 0
	}
	err := jni.Do(func(e *jni.Env) error {
		name, err := e.GoString(jni.Wrap(uintptr(method)))
		if err != nil {
			return err
		}
		var values []jni.Object
		if a := jni.Wrap(uintptr(args)); !a.IsNil() {
			n, err := e.Len(a)
			if err != nil {
				return err
			}
			values = make([]jni.Object, n)
			for i := range n {
				if values[i], err = e.Index(a, i); err != nil {
					return err
				}
			}
		}
		h(e, name, values)
		return nil
	})
	if err != nil {
		ndk.Errorf("a Java listener failed: %v", err)
	}
	return 0
}

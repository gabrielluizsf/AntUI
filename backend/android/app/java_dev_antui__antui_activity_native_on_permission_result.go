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

// Java_dev_antui_AntuiActivity_nativeOnPermissionResult is the answer to
// [RequestPermissions], delivered on the UI thread.
//
//export Java_dev_antui_AntuiActivity_nativeOnPermissionResult
func Java_dev_antui_AntuiActivity_nativeOnPermissionResult(env *C.JNIEnv, cls C.jclass, code C.jint, perms C.jobjectArray, results C.jintArray) {
	var out PermissionResult
	err := jni.Do(func(e *jni.Env) error {
		names, err := e.GoStrings(jni.Wrap(uintptr(perms)))
		if err != nil {
			return err
		}
		answers, err := e.GoInts(jni.Wrap(uintptr(results)))
		if err != nil {
			return err
		}
		out.Permissions = names
		out.Granted = make([]bool, len(answers))
		for i, a := range answers {
			// PackageManager.PERMISSION_GRANTED is 0 and DENIED is -1, which
			// is the other way round from how it reads.
			out.Granted[i] = a == 0
		}
		return nil
	})
	if err != nil {
		ndk.Errorf("cannot read the permission result: %v", err)
	}

	waitMu.Lock()
	ch := permWait[int32(code)]
	delete(permWait, int32(code))
	waitMu.Unlock()
	if ch != nil {
		ch <- out
	}
}

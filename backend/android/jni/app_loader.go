//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

// appLoader finds the class loader that loaded the activity, which is the
// one that can see the app's classes.
func (e *Env) appLoader() (C.jobject, C.jmethodID, error) {
	classMu.RLock()
	l, m := loader, loadWith
	classMu.RUnlock()
	if l != 0 {
		return l, m, nil
	}

	global, load, err := e.findAppLoader()
	if err != nil {
		return 0, nil, err
	}

	classMu.Lock()
	if loader == 0 {
		loader, loadWith = global, load
	} else {
		C.mw_DeleteGlobalRef(e.ptr, global)
		global, load = loader, loadWith
	}
	l, m = loader, loadWith
	classMu.Unlock()
	return l, m, nil
}

//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

import "sync"

var (
	classMu sync.RWMutex
	classes = map[string]C.jclass{}
	// The app's own class loader, and its loadClass method. See Class.
	loader   C.jobject
	loadWith C.jmethodID
)

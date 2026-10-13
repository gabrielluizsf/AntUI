//go:build android

package app

/*
#cgo LDFLAGS: -landroid -llog

#include <android/native_activity.h>
#include "glue.h"
*/
import "C"

import "github.com/gabrielluizsf/antui/backend/android/ndk"

//export antuiOnContentRectChanged
func antuiOnContentRectChanged(_ *C.ANativeActivity, r *C.ARect) {
	if a := current(); a != nil && r != nil {
		a.mu.Lock()
		a.rect = ndk.Rect{
			Left:   int32(r.left),
			Top:    int32(r.top),
			Right:  int32(r.right),
			Bottom: int32(r.bottom),
		}
		a.mu.Unlock()
	}
	post(ContentRect, nil, false)
}

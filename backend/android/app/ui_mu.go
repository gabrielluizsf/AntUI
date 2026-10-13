//go:build android

package app

/*
#include <android/looper.h>
#include "glue.h"
*/
import "C"

import "sync"

// Most of Android's Java API may only be touched from the thread the
// activity was made on. Reading the window insets from anywhere else throws;
// changing the system bars from anywhere else is undefined.
//
// The usual answer is Activity.runOnUiThread, which takes a Runnable — a
// Java object, which needs a Java class, which needs a dex file in the APK.
// None of that is necessary. The UI thread already has a looper, that looper
// already knows how to wait on a file descriptor, and a pipe is a file
// descriptor. So the runner is a pipe and a callback, with no Java at all.
var (
	uiMu     sync.Mutex
	uiQueue  []func()
	uiWrite  int = -1
	uiLooper *C.ALooper
)

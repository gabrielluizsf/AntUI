//go:build android

package app

import (
	"github.com/gabrielluizsf/antui/backend/android/ndk"
)

// attach hooks the input queue up to this thread's looper, and unhooks it.
//
// It has to happen here and not in a callback: a queue is delivered to the
// looper of the thread that attached it, and the UI thread's looper is not
// the one this app reads. Attaching from the wrong thread does not fail — it
// simply means no input ever arrives.
func (a *App) attach(m message) {
	switch m.ev {
	case InputUp:
		if a.loop == nil {
			a.loop = ndk.PrepareLooper()
		}
		a.queue = ndk.WrapQueue(m.ptr)
		if a.loop != nil && a.queue != nil {
			a.queue.AttachLooper(a.loop)
		}
	case InputDown:
		if a.queue != nil {
			a.queue.DetachLooper()
			a.queue = nil
		}
	}
}

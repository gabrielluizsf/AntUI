//go:build android

package app

/*
#include <jni.h>
#include "glue.h"
*/
import "C"

import (
	"sync"
	"sync/atomic"

	"github.com/gabrielluizsf/antui/backend/android/jni"
	"github.com/gabrielluizsf/antui/backend/android/ndk"
)

// Handler is a Java listener, written in Go.
//
// method is the name of the method that was called — "onClick", "onLost",
// "onReceive" — and args are what it was given. **The arguments are local
// references belonging to the call**: anything kept past the return has to be
// made global first, with [jni.Env.Global].
//
// It runs on whatever thread Android chose, which for a dialog is the UI
// thread and for a broadcast is not. Two rules follow. Do not take long:
// something is waiting. And do not call [RunOnUISync] without asking
// [OnUIThread] first, or the call may be waiting for the thread it is
// already on.
//
// A handler returns nothing. Every listener this library needs returns void,
// and a Java method that has to be answered would need its answer to outlive
// the reference frame it was made in — a rule worth not having to remember
// for a case that has not come up.
type Handler func(e *jni.Env, method string, args []jni.Object)

// The registry. Java holds a number and nothing else: a Java field holding a
// Go pointer is something neither side's collector could keep honest.
var (
	handlerMu sync.Mutex
	handlers  = map[int64]Handler{}
	lastToken atomic.Int64
)

// Listen registers a handler and gives back the token that identifies it, and
// the function that takes it off again.
//
// **The release function has to be called.** A handler that is never released
// is a Go closure that lives for the life of the process, and — worse — a
// Java object still registered with the platform, which will go on calling it.
func Listen(h Handler) (token int64, release func()) {
	token = lastToken.Add(1)
	handlerMu.Lock()
	handlers[token] = h
	handlerMu.Unlock()
	var once sync.Once
	return token, func() {
		once.Do(func() {
			handlerMu.Lock()
			delete(handlers, token)
			handlerMu.Unlock()
		})
	}
}

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

// Proxy makes a Java object that implements the named interfaces and sends
// every call to the handler registered under token.
//
// The names are dotted, as Class.forName takes them:
//
//	app.Proxy(e, token, "android.content.DialogInterface$OnClickListener")
//
// It only works for **interfaces**. An abstract class cannot be proxied, and
// the two Android insists on — BroadcastReceiver and NetworkCallback — are
// written out in the shim instead. The object comes back as a local
// reference; anything that outlives the call needs [jni.Env.Global].
func Proxy(e *jni.Env, token int64, interfaces ...string) (jni.Object, error) {
	names, err := e.NewStrings(interfaces)
	if err != nil {
		return jni.Object{}, err
	}
	return e.Static("dev/antui/AntuiProxy", "create",
		jni.Sig(jni.TObject, jni.TLong, jni.TArray(jni.TString)),
		jni.Long(token), jni.Ref(names))
}

// Listener makes one of the shim's own listener classes — the ones that had
// to be written out because Android wanted an abstract class rather than an
// interface.
func Listener(e *jni.Env, class string, token int64) (jni.Object, error) {
	return e.Static(class, "create", jni.Sig(jni.TObject, jni.TLong), jni.Long(token))
}

// The listener classes the shim carries.
const (
	NetworkCallbackClass = "dev/antui/AntuiNetworkCallback"
	AuthCallbackClass    = "dev/antui/AntuiAuth"
	ReceiverClass        = "dev.antui.AntuiReceiver"
	// ReceiverToken is the intent extra a broadcast's token travels in. The
	// system builds the receiver itself, so the token cannot live on the
	// object.
	ReceiverToken = "dev.antui.token"
)

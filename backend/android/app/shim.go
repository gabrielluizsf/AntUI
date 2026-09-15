//go:build android

package app

/*
#include <jni.h>
*/
import "C"

import (
	"errors"
	"sync"
	"sync/atomic"

	"github.com/gabrielluizsf/antui/backend/android/jni"
	"github.com/gabrielluizsf/antui/backend/android/ndk"
)

// The three things the Java shim exists for. Each arrives on the UI thread,
// carrying arguments that are only valid for the length of the call, and
// each is matched back to whoever asked by a request code.
//
// The Go functions below are found by the machine by **name**: a native
// method on dev.antui.AntuiActivity called nativeOnPermissionResult is
// looked up as Java_dev_antui_AntuiActivity_nativeOnPermissionResult, which
// is exactly what cgo exports. Nothing registers them.

// PermissionResult is the answer to a permission request.
type PermissionResult struct {
	// Permissions is what was asked for, in the order it was asked.
	Permissions []string
	// Granted is the answer for each, in the same order.
	Granted []bool
}

// All reports whether every permission asked for was granted.
func (r PermissionResult) All() bool {
	if len(r.Granted) == 0 {
		return false
	}
	for _, ok := range r.Granted {
		if !ok {
			return false
		}
	}
	return true
}

// Has reports whether one particular permission was granted.
func (r PermissionResult) Has(name string) bool {
	for i, p := range r.Permissions {
		if p == name && i < len(r.Granted) {
			return r.Granted[i]
		}
	}
	return false
}

// ActivityResult is what another activity sent back.
type ActivityResult struct {
	// Code is Activity.RESULT_OK (-1) when the user finished, and
	// RESULT_CANCELED (0) when they backed out.
	Code int
	// Data is the Intent that came back, as a **global reference the
	// receiver owns**: it has to be given back with jni.Env.DeleteGlobal, or
	// the reference table fills up. It is nil when nothing came back, which
	// is the usual case for a cancel.
	Data jni.Object
}

// OK reports whether the other activity finished rather than being backed
// out of.
func (r ActivityResult) OK() bool { return r.Code == ResultOK }

// The two result codes android.app.Activity defines. RESULT_OK is -1 and not
// 0, which is the opposite of every other convention and is worth stating.
const (
	ResultOK       = -1
	ResultCanceled = 0
)

// The registry. A request goes out with a code, and the callback that comes
// back carries the same code — which is the only way to tell one answer from
// another when two are outstanding.
var (
	waitMu   sync.Mutex
	permWait = map[int32]chan PermissionResult{}
	actWait  = map[int32]chan ActivityResult{}
	lastCode atomic.Int32

	intentMu sync.Mutex
	intentFn func(jni.Object)
)

// nextRequestCode is a code nothing else is using. It stays under 32768
// because startActivityForResult only carries the low bits of it on some
// platform versions, and zero is skipped because parts of the platform use
// that to mean "no request".
func nextRequestCode() int32 {
	for {
		if c := lastCode.Add(1) & 0x7FFF; c != 0 {
			return c
		}
	}
}

func awaitPermission(code int32) chan PermissionResult {
	ch := make(chan PermissionResult, 1)
	waitMu.Lock()
	permWait[code] = ch
	waitMu.Unlock()
	return ch
}

func awaitActivity(code int32) chan ActivityResult {
	ch := make(chan ActivityResult, 1)
	waitMu.Lock()
	actWait[code] = ch
	waitMu.Unlock()
	return ch
}

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

// OnNewIntent registers what to do when an intent reaches an activity that
// is already running. The Intent handed to f is a **global reference f
// owns** and must give back.
//
// It is called on a goroutine of its own, not on the UI thread and not on
// the app's.
func OnNewIntent(f func(jni.Object)) {
	intentMu.Lock()
	intentFn = f
	intentMu.Unlock()
}

// ErrNoShim is what the calls here report when the app was built without the
// Java shim, so there is nothing to deliver the answer.
var ErrNoShim = errors.New("app: this app was built without the Java shim, so " +
	"permission and activity results cannot be delivered; build without -no-shim")

// HasPermission reports whether a permission has already been granted.
//
// Below API 23 every permission an app declares is granted when it is
// installed, and there is nothing to ask, so this is always true there.
func HasPermission(name string) (bool, error) {
	a := Current()
	if a == nil {
		return false, ErrNoUIThread
	}
	if a.Config().SDK < 23 {
		return true, nil
	}
	var granted bool
	err := jni.Do(func(e *jni.Env) error {
		act := jni.Activity()
		cls, err := e.ClassOf(act)
		if err != nil {
			return err
		}
		defer e.Delete(cls.Object())
		m, err := e.Method(cls, "checkSelfPermission", jni.Sig(jni.TInt, jni.TString))
		if err != nil {
			return err
		}
		js, err := e.String(name)
		if err != nil {
			return err
		}
		r, err := e.CallInt(act, m, jni.Ref(js))
		granted = r == 0
		return err
	})
	return granted, err
}

// RequestPermissions asks the user and waits for the answer.
//
// **Do not call it from the frame loop.** Asking puts a dialog in front of
// the app, which pauses it — and Android waits for the app to acknowledge
// that pause. A frame loop that is waiting here is a frame loop that is not
// acknowledging, and both sides stop. Run it on a goroutine of its own and
// let the frame loop carry on drawing.
//
// A permission that is already granted comes back straight away, with no
// dialog. One the user has refused twice comes back refused, again with no
// dialog, and that is the platform's decision rather than this library's.
func RequestPermissions(names ...string) (PermissionResult, error) {
	if len(names) == 0 {
		return PermissionResult{}, nil
	}
	a := Current()
	if a == nil {
		return PermissionResult{}, ErrNoUIThread
	}
	if a.Config().SDK < 23 {
		granted := make([]bool, len(names))
		for i := range granted {
			granted[i] = true
		}
		return PermissionResult{Permissions: names, Granted: granted}, nil
	}

	code := nextRequestCode()
	ch := awaitPermission(code)
	var callErr error
	if err := RunOnUISync(func() {
		callErr = jni.Do(func(e *jni.Env) error {
			act := jni.Activity()
			cls, err := e.ClassOf(act)
			if err != nil {
				return err
			}
			defer e.Delete(cls.Object())
			m, err := e.Method(cls, "requestPermissions",
				jni.Sig(jni.TVoid, jni.TArray(jni.TString), jni.TInt))
			if err != nil {
				return err
			}
			arr, err := e.NewStrings(names)
			if err != nil {
				return err
			}
			return e.CallVoid(act, m, jni.Ref(arr), jni.Int(code))
		})
	}); err != nil {
		callErr = err
	}
	if callErr != nil {
		waitMu.Lock()
		delete(permWait, code)
		waitMu.Unlock()
		return PermissionResult{}, callErr
	}
	return <-ch, nil
}

// StartForResult starts another activity and waits for what it sends back —
// a photo from the camera, a file from the picker, a document from a cloud
// drive.
//
// The same warning as [RequestPermissions] applies, and more strongly: the
// other activity is in front for as long as the user wants it to be.
//
// The Intent in the result is a global reference the caller owns.
func StartForResult(intent jni.Object) (ActivityResult, error) {
	if intent.IsNil() {
		return ActivityResult{}, errors.New("app: there is no intent to start")
	}
	code := nextRequestCode()
	ch := awaitActivity(code)
	var callErr error
	if err := RunOnUISync(func() {
		callErr = jni.Do(func(e *jni.Env) error {
			act := jni.Activity()
			cls, err := e.ClassOf(act)
			if err != nil {
				return err
			}
			defer e.Delete(cls.Object())
			m, err := e.Method(cls, "startActivityForResult",
				jni.Sig(jni.TVoid, jni.TClass("android/content/Intent"), jni.TInt))
			if err != nil {
				return err
			}
			return e.CallVoid(act, m, jni.Ref(intent), jni.Int(code))
		})
	}); err != nil {
		callErr = err
	}
	if callErr != nil {
		waitMu.Lock()
		delete(actWait, code)
		waitMu.Unlock()
		return ActivityResult{}, callErr
	}
	return <-ch, nil
}

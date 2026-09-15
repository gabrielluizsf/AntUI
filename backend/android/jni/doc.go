// Package jni calls Java from Go.
//
// Most of Android is not in the NDK. Permissions, notifications, the
// clipboard, the share sheet, storage, location, vibration — all of it is a
// Java API and nothing else, and the only way to reach it from a native app
// is JNI. This package is that way, made safe enough to use: the attaching,
// the reference frames and the exception checking that JNI leaves to the
// caller are done here, once, rather than at every call site.
//
// # Everything happens inside Do
//
//	err := jni.Do(func(e *jni.Env) error {
//		cls, err := e.Class("android/os/Build")
//		if err != nil {
//			return err
//		}
//		f, err := e.StaticField(cls, "MODEL", jni.TString)
//		if err != nil {
//			return err
//		}
//		model, err := e.StaticString(cls, f)
//		...
//	})
//
// [Do] pins the goroutine to its thread, attaches that thread to the virtual
// machine if it is not already, opens a local reference frame, and closes
// all of it afterwards. **The Env must not escape**: it belongs to one
// thread and one frame, and so does every [Object] made with it. Anything
// that has to outlive the call is made global with [Env.Global].
//
// # What is dangerous here
//
// JNI has no type checking and no memory safety. A wrong signature string is
// not an error, it is a crash inside the runtime with a stack that names
// nothing. A local reference used after its frame has closed is the same. An
// exception left pending makes the *next* call fail in a way that has
// nothing to do with what is wrong. This package checks for the exception
// after every call and turns it into an error, which removes the third of
// those; the other two are what the typed helpers exist to reduce.
package jni

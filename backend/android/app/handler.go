//go:build android

package app

import "github.com/gabrielluizsf/antui/backend/android/jni"

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

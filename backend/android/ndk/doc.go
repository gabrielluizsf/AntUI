// Package ndk is the thin part: the C APIs Android exposes to native code,
// bound one to one and nothing more.
//
// It has no opinion about how an app is shaped. Everything here is a handle
// the system handed us and the calls that can be made on it — the native
// window, the looper, the log. The shape of an app is [antui/backend/android/app]'s
// business, and the features Android only offers through Java are the
// feature packages'.
//
// Nothing in here builds off Android. On any other system the package is
// empty rather than absent, so that "go build ./..." on a desktop still
// covers the tree.
package ndk

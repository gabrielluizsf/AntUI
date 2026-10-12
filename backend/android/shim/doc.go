package shim

// Package shim holds the only Java in this library, and the compiled form of
// it that ships inside an APK.
//
// Three things Android delivers by overriding a method on the activity —
// the result of a permission request, the result of another activity, and an
// intent arriving at one already running — cannot be reached from native
// code, because there is no way to override a Java method from C. So there
// is one small subclass of NativeActivity that overrides those three and
// hands each one's arguments straight to Go.
//
// The dex is **committed**, not built. An app made with this library needs
// no JDK and no Java toolchain; only changing the shim itself does, and
// [Build] is what does it. d8 is deterministic, so the committed file can be
// checked against the source rather than trusted.

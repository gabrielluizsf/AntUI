package shim

import _ "embed"

// ProGuard is what R8 must be told to leave alone.
//
// AntUI's own pipeline never runs R8: antuiapk compiles the shim with javac
// and dexes it with d8, and neither shrinks anything. This is for an app
// that merges the shim into a Gradle build, where R8 is on for every release
// build and would otherwise remove most of the shim as dead code — nothing
// in it is called from Java, so to a shrinker none of it is reachable.
//
// The failure it prevents is a quiet one: the build succeeds, the app
// installs, and a JNI call throws NoSuchMethodError the first time that code
// path runs.
//
//	os.WriteFile("proguard-rules.pro", shim.ProGuard, 0o644)
//
//go:embed proguard-rules.pro
var ProGuard []byte

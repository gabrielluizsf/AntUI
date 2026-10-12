package sdk

// Require says what a caller cannot do without. The zero value asks for a
// usable SDK with build tools and a platform, which is what building an APK
// takes.
type Require struct {
	NDK bool // the app has native code, which every Go app does
	JDK bool // only for rebuilding the Java shim

	// MinBuildTools and MinPlatform are floors, not preferences: the newest
	// found is always the one chosen, and these only decide whether it is
	// good enough.
	MinBuildTools Version
	MinPlatform   int
}

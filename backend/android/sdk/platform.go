package sdk

// Platform is one entry under the SDK's platforms directory. It matters for
// two things only: the android.jar the Java shim compiles against, and the
// API level an app says it targets.
type Platform struct {
	Dir  string
	Name string // "android-37.2", as the directory is named
	API  int    // 37
	// Minor is the second number of a platform like android-36.1, which
	// Android started shipping when it began releasing an SDK mid-year. It
	// is 0 on a platform named by a single number.
	Minor int
	// Ext is the extension level of an android-34-ext11, and 0 otherwise.
	// These are not chosen automatically: they are a strictly older platform
	// with newer APIs bolted on, which is not what a build wants by default.
	Ext int
}

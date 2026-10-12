package sdk

// Problem is one thing wrong with the install.
type Problem struct {
	What string // "build-tools", "ndk", "a JDK"
	// Have and Want are set when the thing is there but too old, and empty
	// when it is not there at all.
	Have, Want string
	// Install is the sdkmanager package that fixes it, empty when nothing
	// sdkmanager has would.
	Install string
}

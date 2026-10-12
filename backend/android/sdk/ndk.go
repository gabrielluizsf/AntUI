package sdk

// NDK is the native toolchain: the clang that compiles Go's cgo for a phone.
type NDK struct {
	Dir     string
	Version Version
	// Host is the prebuilt directory the toolchain for this machine is in —
	// "linux-x86_64" and the like. It is read off the disk rather than
	// assumed, because which host builds the NDK ships has changed.
	Host string
}

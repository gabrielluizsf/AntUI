//go:build android

package app

// Info is what the activity knows about itself and never changes.
type Info struct {
	// InternalData is the app's private directory: what os.UserCacheDir
	// would be if there were a home directory, and the only place writing is
	// certain to work.
	InternalData string
	// ExternalData is the app's directory on shared storage, which may not
	// exist and may not be writable.
	ExternalData string
	// OBB is where an expansion file would be mounted.
	OBB string
	// Cache is where things that can be thrown away go. The system deletes
	// it when the device runs short, without asking and without warning, so
	// nothing that matters belongs here — and everything that does not
	// belongs nowhere else.
	Cache string
	// SDK is the API level of the device — 21 and up.
	SDK int
}

package shim

// Activity is the class an app's manifest points at when the shim is in, and
// Receiver is the one a notification's action buttons broadcast to.
const (
	Activity = "dev.antui.AntuiActivity"
	Receiver = "dev.antui.AntuiReceiver"
	// Provider is the one that hands a file to another app, and
	// ProviderSuffix is what its authority is the app's package plus.
	Provider       = "dev.antui.AntuiFileProvider"
	ProviderSuffix = ".antuifiles"
)

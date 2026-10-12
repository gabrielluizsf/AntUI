//go:build android

package app

// The listener classes the shim carries.
const (
	NetworkCallbackClass = "dev/antui/AntuiNetworkCallback"
	AuthCallbackClass    = "dev/antui/AntuiAuth"
	ReceiverClass        = "dev.antui.AntuiReceiver"
	// ReceiverToken is the intent extra a broadcast's token travels in. The
	// system builds the receiver itself, so the token cannot live on the
	// object.
	ReceiverToken = "dev.antui.token"
)

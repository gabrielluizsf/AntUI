//go:build android

package network

// State is what the device's connection is like now.
type State struct {
	// Online is whether there is a network that reaches the internet.
	Online bool
	// Validated is stronger: the system has actually checked that the
	// internet is reachable, rather than merely being attached to something.
	// A captive portal in a hotel is Online and not Validated.
	Validated bool
	// Metered is a connection the user pays for by the byte. An app that
	// downloads anything large should ask, and should not do it here.
	Metered bool
	// Transport is what it is.
	Transport Transport
}

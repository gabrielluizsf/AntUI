//go:build android

package network

// The transports, numbered as NetworkCapabilities numbers them.
const (
	Cellular  Transport = 0
	WiFi      Transport = 1
	Bluetooth Transport = 2
	Ethernet  Transport = 3
	VPN       Transport = 4
	None      Transport = -1
)

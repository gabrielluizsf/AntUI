//go:build android

package network

// String names how the device is connected: wifi, mobile, ethernet.
func (t Transport) String() string {
	switch t {
	case Cellular:
		return "cellular"
	case WiFi:
		return "wifi"
	case Bluetooth:
		return "bluetooth"
	case Ethernet:
		return "ethernet"
	case VPN:
		return "vpn"
	}
	return "none"
}

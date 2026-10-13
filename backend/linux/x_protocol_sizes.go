//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

// X11 puts these sizes on the wire repeatedly; naming them keeps a packet's
// shape visible without making each request explain the same number again.
const (
	xProtocolWordBytes             = 4
	xProtocolHeaderBytes           = 8
	xHandshakeRequestBytes         = 12
	xEventPacketBytes              = 32
	xClientMessageBytes            = 40
	xClientMessageHeaderBytes      = 8
	xPropertyRequestBytes          = 20
	xSetupHeaderBytes              = 32
	xFormatRecordBytes             = 8
	xScreenMinimumBytes            = 40
	xMinRequestBytes               = 4096
	xFormat8Bits              byte = 8
	xFormat32Bits             byte = 32
	xPixel32Bits                   = 32
)

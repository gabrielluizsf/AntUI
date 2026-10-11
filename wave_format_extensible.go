//go:build windows

package antui

type waveFormatExtensible struct {
	waveFormatEx
	samples     uint16 // a union; the valid bits per sample
	channelMask uint32
	subFormat   guid
}

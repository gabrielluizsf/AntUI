//go:build windows

package antui

// A waveFormatEx is the shape of a stream. The extensible form adds three
// fields after it, and a mix format is always the extensible one.
type waveFormatEx struct {
	tag            uint16
	channels       uint16
	samplesPerSec  uint32
	avgBytesPerSec uint32
	blockAlign     uint16
	bitsPerSample  uint16
	size           uint16
}

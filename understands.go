//go:build windows

package antui

// understands reports whether the samples are ones this can write.
func (w *wasapi) understands(f *waveFormatExtensible) bool {
	if f.channels == 0 || f.samplesPerSec == 0 {
		return false
	}
	switch f.tag {
	case formatFloat:
		return f.bitsPerSample == 32
	case formatPCM:
		return f.bitsPerSample == 16
	case formatExtensible:
		if f.subFormat == subtypeFloat {
			return f.bitsPerSample == 32
		}
		return f.subFormat == subtypePCM && f.bitsPerSample == 16
	}
	return false
}

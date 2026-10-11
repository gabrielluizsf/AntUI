//go:build windows

package antui

// floats reports whether the device takes samples as they already are.
func (w *wasapi) floats(f *waveFormatExtensible) bool {
	return f.tag == formatFloat || (f.tag == formatExtensible && f.subFormat == subtypeFloat)
}

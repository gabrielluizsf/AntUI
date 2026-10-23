//go:build windows

package windows

import "unicode/utf16"

// copyUTF16 writes s into a fixed WCHAR field, cut short to leave room for
// the NUL. A cut never splits a surrogate pair: half a character is a broken
// one, where dropping it whole is just a shorter string.
func copyUTF16(dst []uint16, s string) {
	units := utf16.Encode([]rune(s))
	n := min(len(units), len(dst)-1)
	if n > 0 && n < len(units) && utf16.IsSurrogate(rune(units[n-1])) && units[n-1] < 0xDC00 {
		n--
	}
	copy(dst, units[:n])
	clear(dst[n:])
}

//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

// compose joins a pending dead accent to the character that followed it,
// returning 0 when the pair does not make a letter this font can draw.
func compose(dead uint32, base rune) rune {
	if r, ok := composeDead[deadPair{dead, base}]; ok {
		return r
	}
	if base == ' ' {
		return deadAlone[dead]
	}
	return 0
}

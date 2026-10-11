//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

// deadAlone is the accent on its own, which is what a space after it means.
var deadAlone = map[uint32]rune{
	0xFE50: '`', 0xFE51: '\'', 0xFE52: '^', 0xFE53: '~', 0xFE57: '"',
}

package canvas

// The pixel formats.
const (
	ARGB32 Format = iota // four bytes, the default
	RGB565               // two bytes: 5 red, 6 green, 5 blue, no alpha
	Pal8                 // one byte, an index into a 256-entry palette
)

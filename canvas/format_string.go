package canvas

import "fmt"

// String names the format.
func (f Format) String() string {
	switch f {
	case ARGB32:
		return "ARGB32"
	case RGB565:
		return "RGB565"
	case Pal8:
		return "PAL8"
	default:
		return fmt.Sprintf("Format(%d)", int(f))
	}
}

package canvas

// Bytes is how much one pixel takes in this format.
func (f Format) Bytes() int {
	switch f {
	case RGB565:
		return 2
	case Pal8:
		return 1
	default:
		return 4
	}
}

package canvas

// R returns the red channel.
func (c Color) R() uint8 { return uint8(c >> 16) }

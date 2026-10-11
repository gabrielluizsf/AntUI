package canvas

// A returns the alpha channel.
func (c Color) A() uint8 { return uint8(c >> 24) }

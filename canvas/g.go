package canvas

// G returns the green channel.
func (c Color) G() uint8 { return uint8(c >> 8) }

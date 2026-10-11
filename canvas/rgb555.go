package canvas

func rgb555(c Color) uint32 {
	return (uint32(c)>>9)&0x7C00 | (uint32(c)>>6)&0x03E0 | (uint32(c)>>3)&0x001F
}

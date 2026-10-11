package canvas

func to565Dithered(c Color, x, y int) uint16 {
	t := bayer[(y&3)*4+(x&3)] * 16
	r5 := min(div255(uint32(c.R())*31+t), 31)
	g6 := min(div255(uint32(c.G())*63+t), 63)
	b5 := min(div255(uint32(c.B())*31+t), 31)
	return uint16(r5<<11 | g6<<5 | b5)
}

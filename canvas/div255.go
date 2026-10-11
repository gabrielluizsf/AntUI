package canvas

// div255 divides by 255 without dividing: exact for everything below 65535.
func div255(t uint32) uint32 { return (t + (t >> 8) + 1) >> 8 }

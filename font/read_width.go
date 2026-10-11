package font

// readWidth reads one offset from an index subtable, 16 or 32 bits wide.
func readWidth(b []byte, i, width int) uint32 {
	if width == 2 {
		return uint32(be16(b, i))
	}
	return be32(b, i)
}

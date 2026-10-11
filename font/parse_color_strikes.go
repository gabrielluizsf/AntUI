package font

// parseColorStrikes reads a CBLC table and the CBDT table it points into,
// into one entry per strike per glyph. Anything that does not read — a
// truncated table, an index format with nothing this engine draws, an image
// format that carries no PNG — is left out rather than reported.
func parseColorStrikes(cblc, cbdt []byte) []colorStrike {
	if len(cblc) < 8 || len(cbdt) < 4 {
		return nil
	}
	numSizes := int(be32(cblc, 4))
	if numSizes <= 0 || numSizes > (len(cblc)-8)/48 {
		return nil
	}
	var strikes []colorStrike
	for i := range numSizes {
		if s, ok := parseColorStrike(cblc, cbdt, 8+i*48); ok {
			strikes = append(strikes, s)
		}
	}
	return strikes
}

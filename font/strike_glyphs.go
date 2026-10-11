package font

// strikeGlyphs reads one IndexSubTable: the glyphs its range covers and where
// each one's data begins. Only image formats that hold a PNG are read; the
// rest are the monochrome tables this engine draws no colour from.
func strikeGlyphs(cblc, cbdt []byte, sub, first, last int) map[int]colorBitmap {
	if sub < 0 || sub+8 > len(cblc) || first > last {
		return nil
	}
	indexFormat, imageFormat := int(be16(cblc, sub)), int(be16(cblc, sub+2))
	if !isPNGImage(imageFormat) {
		return nil
	}
	s := &strikeReader{
		cblc: cblc, cbdt: cbdt, sub: sub, first: first, last: last,
		imageFormat: imageFormat, base: int(be32(cblc, sub+4)),
		out: map[int]colorBitmap{},
	}
	var ok bool
	switch indexFormat {
	case 1, 3:
		ok = s.readOffsets(indexFormat)
	case 2:
		ok = s.readUniform()
	case 4:
		ok = s.readSparse()
	case 5:
		ok = s.readSparseUniform()
	default:
		return nil
	}
	if !ok {
		return nil
	}
	if s.hasConst {
		for g, b := range s.out {
			b.metrics, b.inIndex = s.constant, true
			s.out[g] = b
		}
	}
	return s.out
}

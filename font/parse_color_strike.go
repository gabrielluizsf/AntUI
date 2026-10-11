package font

// parseColorStrike reads one BitmapSize record and the index subtables it
// points at. A strike that holds no picture is not a strike this font has.
func parseColorStrike(cblc, cbdt []byte, rec int) (colorStrike, bool) {
	s := colorStrike{ppem: int(cblc[rec+44]), glyphs: map[int]colorBitmap{}}
	listOff, numSubs := int(be32(cblc, rec)), int(be32(cblc, rec+8))
	if s.ppem == 0 || listOff < 0 || numSubs <= 0 ||
		numSubs > (len(cblc)-listOff)/8 {
		return s, false
	}
	for i := range numSubs {
		// The record names its subtable from the start of the list, and
		// the range it covers comes off the record itself.
		at := listOff + i*8
		first, last := int(be16(cblc, at)), int(be16(cblc, at+2))
		sub := listOff + int(be32(cblc, at+4))
		for g, b := range strikeGlyphs(cblc, cbdt, sub, first, last) {
			s.glyphs[g] = b
		}
	}
	return s, len(s.glyphs) > 0
}

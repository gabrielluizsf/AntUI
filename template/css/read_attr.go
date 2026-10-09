package css

// readAttr consumes an attribute selector "…]", which the flat widget model
// cannot evaluate, and marks the selector unsupported.
func (s *selectorScan) readAttr() {
	s.warns = append(s.warns, fmtErrf("ignoring attribute selector in %q", s.text).Error())
	s.sel.Unsupported = markUnsupported(s.sel.Unsupported, "attribute")
	s.i = max(attrEnd(s.t, s.i), s.i+1)
	s.saw = true
}

// readID consumes an id selector "#name", which the flat widget model cannot
// evaluate, and marks the selector unsupported.
func (s *selectorScan) readID() {
	_, n := readName(s.t[s.i+1:])
	s.warns = append(s.warns, fmtErrf("ignoring id selector in %q", s.text).Error())
	s.sel.Unsupported = markUnsupported(s.sel.Unsupported, "id")
	s.i += 1 + n
	s.saw = true
}

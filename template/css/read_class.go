package css

import "strings"

// readClass consumes ".name"; a bad one warns and the selector is marked
// unsupported.
func (s *selectorScan) readClass() {
	s.i++
	for s.i < len(s.t) && isSpace(s.t[s.i]) {
		s.i++
	}
	name, n := readName(s.t[s.i:])
	if name == "" {
		s.warns = append(s.warns, fmtErrf("bad class in selector %q", s.text).Error())
		s.sel.Unsupported = markUnsupported(s.sel.Unsupported, "class")
		s.i += max(n, 1)
		s.saw = true
		return
	}
	s.sel.Classes = append(s.sel.Classes, strings.ToLower(name))
	s.i += n
	s.saw = true
}

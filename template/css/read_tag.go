package css

import "strings"

// readTag consumes a type name; finding a second one makes a combinator, so
// the selector never matches the flat widget model. An unrecognised byte (an
// @ inside a selector, a stray glyph) is skipped — the index always advances,
// so there is no hang.
func (s *selectorScan) readTag() {
	name, n := readName(s.t[s.i:])
	if n == 0 {
		s.i++
		return
	}
	if s.sel.Tag != "" || s.sel.All {
		s.comb = true
	} else {
		s.sel.Tag = strings.ToLower(name)
	}
	s.i += n
	s.saw = true
}

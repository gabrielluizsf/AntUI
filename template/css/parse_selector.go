package css

import "strings"

// selectorScan is one selector being read: the source text, its comment-free
// body, the cursor, the Selector under construction, the warnings gathered
// so far, and whether a token or a combinator has been seen.
type selectorScan struct {
	text  string
	t     string
	i     int
	sel   Selector
	warns []string
	saw   bool
	comb  bool
}

// parseSelector reads one selector: a tag and/or classes and pseudo-classes,
// in any order. Whatever the flat widget model cannot evaluate — combinators,
// attribute selectors, ids, pseudo-elements, structural pseudo-classes — is
// recorded as [Selector.Unsupported] with a warning, so the rule parses and
// then never matches rather than failing the sheet.
func parseSelector(text string) (Selector, []string) {
	s := &selectorScan{text: text}
	s.t = strings.TrimSpace(stripComments(text))
	if s.t == "" {
		return Selector{Unsupported: "empty"}, nil
	}
	for s.i < len(s.t) {
		switch c := s.t[s.i]; {
		case isSpace(c):
			s.skipSpace()
		case c == '.':
			s.readClass()
		case c == ':':
			s.readPseudo()
		case c == '[':
			s.readAttr()
		case c == '#':
			s.readID()
		case c == '*':
			s.sel.All = true
			s.i++
			s.saw = true
		case c == '>', c == '+', c == '~', c == '|':
			s.comb = true
			s.i++
		default:
			s.readTag()
		}
	}
	if s.comb {
		s.warns = append(s.warns, fmtErrf("ignoring combinators in %q", text).Error())
		s.sel.Unsupported = markUnsupported(s.sel.Unsupported, "combinators")
	}
	return s.sel, s.warns
}

package css

import "strings"

// readPseudo consumes ":" and the pseudo-class or pseudo-element after it.
// The state pseudo-classes (:hover, :focus, :active, :checked) and :root are
// understood; anything else — a pseudo-element, a structural or functional
// pseudo-class — is recorded as unsupported with a warning, its parenthesised
// argument list balanced and consumed.
func (s *selectorScan) readPseudo() {
	if s.i+1 < len(s.t) && s.t[s.i+1] == ':' {
		s.readPseudoElement()
		return
	}
	s.i++
	name, n := readName(s.t[s.i:])
	s.i += n
	switch strings.ToLower(name) {
	case "hover":
		s.sel.State |= StateHover
	case "focus":
		s.sel.State |= StateFocus
	case "active", "pressed":
		s.sel.State |= StateActive
	case "checked":
		s.sel.State |= StateChecked
	case "root":
		// :root is the document root; on a window the template paints that
		// as the body element, so the rule targets body.
		s.sel.Tag = "body"
	default:
		s.warns = append(s.warns, fmtErrf("ignoring pseudo-class :%s in %q", name, s.text).Error())
		s.sel.Unsupported = markUnsupported(s.sel.Unsupported, ":"+name)
		if s.i < len(s.t) && s.t[s.i] == '(' {
			s.i = scanParen(s.t, s.i)
		}
	}
	s.saw = true
}

// readPseudoElement consumes a "::" pseudo-element, which the flat widget
// model cannot paint.
func (s *selectorScan) readPseudoElement() {
	s.i += 2
	name, n := readName(s.t[s.i:])
	s.i += n
	s.warns = append(s.warns, fmtErrf("ignoring pseudo-element :%s in %q", name, s.text).Error())
	s.sel.Unsupported = markUnsupported(s.sel.Unsupported, "pseudo-element")
	if s.i < len(s.t) && s.t[s.i] == '(' {
		s.i = scanParen(s.t, s.i)
	}
	s.saw = true
}

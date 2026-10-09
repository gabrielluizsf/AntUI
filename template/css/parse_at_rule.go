package css

import "strings"

// parseAtRule handles every at-rule: @media is evaluated against the window
// the frame is drawn at, @supports against the engine itself, the
// rule-group at-rules that have nothing to test (@layer, @container) apply
// their body as written, @import reads another stylesheet in at the place the
// statement stands, and @font-face names a font file to draw the family it
// declares with. Every other at-rule is consumed wholesale and dropped with a
// warning. Whatever a stylesheet throws, the parser walks on.
func (p *parser) parseAtRule(sh *Sheet, outer Media) error {
	// Read the word after '@'.
	start := p.i + 1
	for !p.eof() {
		c := p.src[p.i]
		if isSpace(c) || c == '{' || c == ';' || c == '(' {
			break
		}
		p.i++
	}
	word := strings.ToLower(p.src[start:p.i])
	switch word {
	case "media":
		return p.parseMediaAtRule(sh, outer)
	case "supports":
		return p.parseSupportsAtRule(sh, outer)
	case "layer", "container":
		return p.parseLayerAtRule(sh, outer)
	case "keyframes", "-webkit-keyframes", "-moz-keyframes", "-o-keyframes":
		prelude, _, err := p.readHeader(false)
		if err != nil {
			return err
		}
		if p.peek() == '{' {
			p.next()
		}
		return p.parseKeyframes(sh, strings.TrimSpace(prelude))
	case "import":
		return p.parseImport(sh, outer)
	case "font-face":
		return p.parseFontFace(sh)
	default:
		// Statement and declaration at-rules — @charset, @namespace, @page,
		// @property, vendor and future ones — name things no box can paint:
		// fonts, faces, fallbacks, namespaces. Their prelude and body are
		// consumed, and the rule is reported.
		sh.Warn = append(sh.Warn, fmtErrf("ignoring @%s", word).Error())
		return p.skipAtRuleBody()
	}
}

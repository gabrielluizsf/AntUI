package css

import (
	"fmt"
	"strconv"
	"strings"
)

// parser walks CSS source text. It is deliberately forgiving: a selector,
// at-rule or property the engine does not understand is dropped with a
// warning, and only structural damage — an unterminated block, a rule with no
// opening brace — is an error. The tokenizer also never spins: every scan step
// either advances the index or stops, so a stylesheet of any kind can spend
// forever being read rather than hang the app.
type parser struct {
	src string
	i   int
}

func (p *parser) eof() bool { return p.i >= len(p.src) }

func (p *parser) peek() byte {
	if p.eof() {
		return 0
	}
	return p.src[p.i]
}

func (p *parser) next() byte {
	c := p.src[p.i]
	p.i++
	return c
}

// skipSpace advances past whitespace and /* comments */.
func (p *parser) skipSpace() {
	for !p.eof() {
		c := p.src[p.i]
		if c == '/' && p.i+1 < len(p.src) && p.src[p.i+1] == '*' {
			end := strings.Index(p.src[p.i+2:], "*/")
			if end < 0 {
				p.i = len(p.src)
				return
			}
			p.i += end + 4
			continue
		}
		if isSpace(c) {
			p.i++
			continue
		}
		return
	}
}

// readHeader reads a prelude (a selector list or an at-rule condition) up to
// the '{' that opens the block, returning the text it passed over. When semi
// is true, a top-level ';' also ends the header, which a stray statement uses
// to drop itself. Strings, /* comments */ and parentheses are kept intact.
func (p *parser) readHeader(semi bool) (string, bool, error) {
	start := p.i
	depth := 0
	for !p.eof() {
		c := p.src[p.i]
		switch {
		case c == '/' && p.i+1 < len(p.src) && p.src[p.i+1] == '*':
			if end := strings.Index(p.src[p.i+2:], "*/"); end >= 0 {
				p.i += end + 4
			} else {
				p.i = len(p.src)
			}
		case c == '"' || c == '\'':
			if end := quoteEnd(p.src, p.i); end >= 0 {
				p.i = end + 1
			} else {
				return "", false, fmtErrf("unterminated string")
			}
		case c == '(':
			depth++
			p.i++
		case c == ')':
			if depth > 0 {
				depth--
			}
			p.i++
		case depth == 0 && c == '{':
			return p.src[start:p.i], false, nil
		case depth == 0 && semi && c == ';':
			return p.src[start:p.i], true, nil
		default:
			p.i++
		}
	}
	return "", false, fmtErrf("unterminated rule, missing '{'")
}

// readDecl reads one declaration inside a block, stopping at the ';' or the
// '}' that ends it. Strings, comments and parentheses are kept intact, so a
// value like url("a;b.png") reads as one unit.
func (p *parser) readDecl() (string, error) {
	start := p.i
	depth := 0
	for !p.eof() {
		c := p.src[p.i]
		switch {
		case c == '/' && p.i+1 < len(p.src) && p.src[p.i+1] == '*':
			if end := strings.Index(p.src[p.i+2:], "*/"); end >= 0 {
				p.i += end + 4
			} else {
				p.i = len(p.src)
			}
		case c == '"' || c == '\'':
			if end := quoteEnd(p.src, p.i); end >= 0 {
				p.i = end + 1
			} else {
				return "", fmtErrf("unterminated string")
			}
		case c == '(':
			depth++
			p.i++
		case c == ')':
			if depth > 0 {
				depth--
			}
			p.i++
		case depth == 0 && (c == ';' || c == '}'):
			return p.src[start:p.i], nil
		default:
			p.i++
		}
	}
	return "", fmtErrf("unterminated declaration")
}

// parseRule reads one selector block: `.a:hover, button { … }`, optionally
// constrained by an enclosing media query.
func (p *parser) parseRule(sh *Sheet, m Media) error {
	selText, semi, err := p.readHeader(true)
	if err != nil {
		return err
	}
	if p.peek() == '{' {
		p.next()
	}
	if semi {
		// A stray selector with no block: dropped, like a browser drops it.
		return nil
	}
	selectors, warns := parseSelectors(selText)
	sh.Warn = append(sh.Warn, warns...)
	rule := &Rule{Media: m, Selectors: selectors, Order: sh.order}
	sh.order++
	for {
		p.skipSpace()
		if p.eof() {
			return fmtErrf("unterminated rule, missing '}'")
		}
		if p.peek() == '}' {
			p.next()
			sh.rules = append(sh.rules, rule)
			return nil
		}
		text, err := p.readDecl()
		if err != nil {
			return err
		}
		if p.peek() == ';' {
			p.next()
		}
		if d, ok := parseDeclaration(text); ok {
			rule.Decls = append(rule.Decls, d)
		}
	}
}

// parseAtRule handles every at-rule: @media is evaluated against the window
// width, the rule-group at-rules (@supports, @layer, @container) ignore their
// condition and apply the rules inside, and every other at-rule is consumed
// wholesale and dropped with a warning. Whatever a stylesheet throws, the
// parser walks on.
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
		prelude, _, err := p.readHeader(false)
		if err != nil {
			return err
		}
		if p.peek() == '{' {
			p.next()
		}
		inner, warns := parseMedia(prelude)
		sh.Warn = append(sh.Warn, warns...)
		// A @media inside a @media must also live inside the outer one.
		m := outer.and(inner)
		return p.parseRuleGroup(sh, m)

	case "supports", "layer", "container":
		// The condition decides which browsers and windows a rule fits; a
		// canvas does not support feature-testing, so the rules inside apply
		// as if the condition held. A statement form (@layer a, b;) drops it.
		_, semi, err := p.readHeader(true)
		if err != nil {
			return err
		}
		if semi {
			if p.peek() == ';' {
				p.next()
			}
			return nil
		}
		if p.peek() == '{' {
			p.next()
		}
		return p.parseRuleGroup(sh, outer)

	case "keyframes", "-webkit-keyframes", "-moz-keyframes", "-o-keyframes":
		prelude, _, err := p.readHeader(false)
		if err != nil {
			return err
		}
		if p.peek() == '{' {
			p.next()
		}
		return p.parseKeyframes(sh, strings.TrimSpace(prelude))

	default:
		// Statement and declaration at-rules — @charset, @import, @namespace,
		// @font-face, @page, @property, vendor and future ones — name things
		// no box can paint: fonts, faces, fallbacks, namespaces. Their
		// prelude and body are consumed, and the rule is reported.
		sh.Warn = append(sh.Warn, fmtErrf("ignoring @%s", word).Error())
		return p.skipAtRuleBody()
	}
}

// parseRuleGroup reads the body of a rule-group at-rule — @media, @supports,
// @layer, @container — which contains whole rules and nested at-rules (a
// @media inside a @media, a rule directly inside @supports), until the
// closing '}'.
func (p *parser) parseRuleGroup(sh *Sheet, m Media) error {
	for {
		p.skipSpace()
		if p.eof() {
			return fmtErrf("unterminated block, missing '}'")
		}
		switch {
		case p.peek() == '}':
			p.next()
			return nil
		case p.peek() == ';':
			p.next()
		case p.peek() == '@':
			if err := p.parseAtRule(sh, m); err != nil {
				return err
			}
		default:
			if err := p.parseRule(sh, m); err != nil {
				return err
			}
		}
	}
}

// skipAtRuleBody swallows the rest of an unknown at-rule: a ';' statement or
// a balanced { … } block, honouring strings and comments.
func (p *parser) skipAtRuleBody() error {
	for !p.eof() {
		c := p.peek()
		switch {
		case c == '"' || c == '\'':
			if end := quoteEnd(p.src, p.i); end >= 0 {
				p.i = end + 1
			} else {
				p.i = len(p.src)
			}
		case c == '/':
			if p.i+1 < len(p.src) && p.src[p.i+1] == '*' {
				if end := strings.Index(p.src[p.i+2:], "*/"); end >= 0 {
					p.i += end + 4
					continue
				}
				p.i = len(p.src)
				return nil
			}
			p.next()
		case c == ';':
			p.next()
			return nil
		case c == '{':
			return p.skipBlock()
		case c == '}':
			return nil
		default:
			p.next()
		}
	}
	return nil
}

// skipBlock swallows a balanced { … } block.
func (p *parser) skipBlock() error {
	depth := 0
	for !p.eof() {
		c := p.src[p.i]
		switch {
		case c == '"' || c == '\'':
			if end := quoteEnd(p.src, p.i); end >= 0 {
				p.i = end + 1
			} else {
				p.i = len(p.src)
			}
		case c == '/':
			if p.i+1 < len(p.src) && p.src[p.i+1] == '*' {
				if end := strings.Index(p.src[p.i+2:], "*/"); end >= 0 {
					p.i += end + 4
					continue
				}
				p.i = len(p.src)
				return nil
			}
			p.next()
		case c == '{':
			depth++
			p.next()
		case c == '}':
			depth--
			p.next()
			if depth == 0 {
				return nil
			}
		default:
			p.next()
		}
	}
	return fmtErrf("unterminated block")
}

// parseSelectors splits a comma-separated selector list into one Selector per
// part. Every part is read tolerantly — an unsupported part marks its Selector
// as never-matching and warns instead of failing the sheet.
func parseSelectors(text string) ([]Selector, []string) {
	var out []Selector
	var warns []string
	for _, part := range splitTopLevel(text) {
		s, w := parseSelector(part)
		warns = append(warns, w...)
		out = append(out, s)
	}
	return out, warns
}

func splitTopLevel(text string) []string {
	var out []string
	depth := 0
	start := 0
	for i := 0; i < len(text); i++ {
		switch text[i] {
		case '(', '[':
			depth++
		case ')', ']':
			if depth > 0 {
				depth--
			}
		case '"', '\'':
			if end := quoteEnd(text, i); end >= 0 {
				i = end
			}
		case ',':
			if depth == 0 {
				out = append(out, text[start:i])
				start = i + 1
			}
		}
	}
	out = append(out, text[start:])
	trimmed := out[:0]
	for _, o := range out {
		if t := strings.TrimSpace(o); t != "" {
			trimmed = append(trimmed, t)
		}
	}
	return trimmed
}

// parseSelector reads one selector: a tag and/or classes and pseudo-classes,
// in any order. Whatever the flat widget model cannot evaluate — combinators,
// attribute selectors, ids, pseudo-elements, structural pseudo-classes — is
// recorded as [Selector.Unsupported] with a warning, so the rule parses and
// then never matches rather than failing the sheet.
func parseSelector(text string) (Selector, []string) {
	var sel Selector
	var warns []string
	warn := func(f string, a ...any) {
		warns = append(warns, fmtErrf(f, a...).Error())
	}
	t := stripComments(text)
	t = strings.TrimSpace(t)
	if t == "" {
		return Selector{Unsupported: "empty"}, nil
	}
	i := 0
	sawToken := false
	combinator := false
	for i < len(t) {
		switch c := t[i]; {
		case isSpace(c):
			j := i
			for j < len(t) && isSpace(t[j]) {
				j++
			}
			if j < len(t) && sawToken {
				combinator = true
			}
			i = j
		case c == '.':
			i++
			for i < len(t) && isSpace(t[i]) {
				i++
			}
			name, n := readName(t[i:])
			if name == "" {
				warn("bad class in selector %q", text)
				sel.Unsupported = markUnsupported(sel.Unsupported, "class")
				i += max(n, 1)
				sawToken = true
				continue
			}
			sel.Classes = append(sel.Classes, strings.ToLower(name))
			i += n
			sawToken = true
		case c == ':':
			if i+1 < len(t) && t[i+1] == ':' {
				// A pseudo-element (::before, ::placeholder, ::-webkit-*).
				i += 2
				name, n := readName(t[i:])
				i += n
				warn("ignoring pseudo-element :%s in %q", name, text)
				sel.Unsupported = markUnsupported(sel.Unsupported, "pseudo-element")
				if i < len(t) && t[i] == '(' {
					i = scanParen(t, i)
				}
				sawToken = true
				continue
			}
			i++
			name, n := readName(t[i:])
			i += n
			switch strings.ToLower(name) {
			case "hover":
				sel.State |= StateHover
			case "focus":
				sel.State |= StateFocus
			case "active", "pressed":
				sel.State |= StateActive
			case "checked":
				sel.State |= StateChecked
			case "root":
				// :root is the document root; on a window the template paints
				// that as the body element, so the rule targets body.
				sel.Tag = "body"
			default:
				// A structural or functional pseudo-class (:nth-child(2n+1),
				// :not(.x), :is(a, b), :first-child …). The parenthesised
				// argument list is balanced and consumed; the selector never
				// matches the flat widget model.
				warn("ignoring pseudo-class :%s in %q", name, text)
				sel.Unsupported = markUnsupported(sel.Unsupported, ":"+name)
				if i < len(t) && t[i] == '(' {
					i = scanParen(t, i)
				}
			}
			sawToken = true
		case c == '[':
			warn("ignoring attribute selector in %q", text)
			sel.Unsupported = markUnsupported(sel.Unsupported, "attribute")
			i = max(attrEnd(t, i), i+1)
			sawToken = true
		case c == '#':
			_, n := readName(t[i+1:])
			warn("ignoring id selector in %q", text)
			sel.Unsupported = markUnsupported(sel.Unsupported, "id")
			i += 1 + n
			sawToken = true
		case c == '*':
			sel.All = true
			i++
			sawToken = true
		case c == '>', c == '+', c == '~', c == '|':
			combinator = true
			i++
		default:
			name, n := readName(t[i:])
			if n == 0 {
				// An unrecognised byte (an @ inside a selector, a stray
				// glyph): skip it. The index always advances: no hang.
				i++
				continue
			}
			if sel.Tag != "" || sel.All {
				combinator = true
			} else {
				sel.Tag = strings.ToLower(name)
			}
			i += n
			sawToken = true
		}
	}
	if combinator {
		warn("ignoring combinators in %q", text)
		sel.Unsupported = markUnsupported(sel.Unsupported, "combinators")
	}
	return sel, warns
}

// markUnsupported records the first reason a selector cannot be evaluated,
// keeping the most specific one.
func markUnsupported(prev, next string) string {
	if prev != "" {
		return prev
	}
	return next
}

// readName reads a CSS identifier from s: ASCII letters, digits, underscore,
// hyphen, high bytes, or an escape sequence that resolves to one character.
// It returns the parsed name with escapes resolved, and the bytes consumed.
func readName(s string) (string, int) {
	var b strings.Builder
	i := 0
	for i < len(s) {
		c := s[i]
		switch {
		case c == '_' || c == '-' || isIdentByte(c) || c >= 0x80:
			b.WriteByte(c)
			i++
		case c == '\\':
			r, n := readEscape(s, i)
			if n == 0 {
				return b.String(), i
			}
			b.WriteString(r)
			i += n
		default:
			return b.String(), i
		}
	}
	return b.String(), i
}

// readEscape resolves the escape sequence starting at the '\' in s. A
// backslash followed by up to six hex digits reads a code point (with one
// optional trailing whitespace as its terminator); any other backslash reads
// the literal next character.
func readEscape(s string, i int) (string, int) {
	// s[i] is '\\'.
	if i+1 >= len(s) {
		return "", 0
	}
	n := i + 1
	if !isHex(s[n]) {
		return string(s[n]), 2
	}
	start := n
	for n < len(s) && n < start+6 && isHex(s[n]) {
		n++
	}
	v, err := strconv.ParseUint(s[start:n], 16, 32)
	if err != nil {
		return "", 0
	}
	if v > 0x10FFFF {
		v = 0xFFFD
	}
	if n < len(s) && isSpace(s[n]) {
		n++ // the single optional whitespace terminator of a hex escape
	}
	return string(rune(v)), n - i
}

// attrEnd returns the index just past the ']' closing the attribute selector
// that starts at the '[' in s, or the end of the string when it never closes,
// in which case the selector parses to "never matches" and warns.
func attrEnd(s string, i int) int {
	for i < len(s) {
		switch s[i] {
		case '"', '\'':
			if end := quoteEnd(s, i); end >= 0 {
				i = end + 1
				continue
			}
			return len(s)
		case ']':
			return i + 1
		}
		i++
	}
	return len(s)
}

// scanParen returns the index just past the ')' matching the '(' at i, or the
// end of the string when the group never closes.
func scanParen(s string, i int) int {
	end := parenEnd(s, i)
	if end < 0 {
		return len(s)
	}
	return end + 1
}

// parseMedia reads an @media prelude into a Media. The full grammar is
// walked — not/only, and/or, comma-separated query lists — while only the
// width features constrain the rule; everything else warns and is treated as
// satisfied, so a canvas never drops a rule it merely cannot measure.
func parseMedia(prelude string) (Media, []string) {
	var m Media
	var warns []string
	for _, branch := range splitTopLevel(prelude) {
		qs, w := parseMediaQuery(branch)
		warns = append(warns, w...)
		m.Queries = append(m.Queries, qs...)
	}
	return m, warns
}

// parseMediaQuery parses one comma-separated media query into its alternative
// width windows, or-separated groups each becoming a query of their own.
func parseMediaQuery(branch string) ([]MediaQuery, []string) {
	var out []MediaQuery
	var warns []string
	negated := false
	cur := MediaQuery{}
	flush := func() {
		if !cur.HasMin && !cur.HasMax {
			// A negation of nothing we can evaluate is left unconstrained so
			// the rule still applies.
			cur.Negated = false
		}
		out = append(out, cur)
	}
	for _, tok := range splitQueryTokens(branch) {
		t := strings.ToLower(strings.TrimSpace(tok))
		if t == "" {
			continue
		}
		switch t {
		case "not":
			negated = true
			cur.Negated = true
		case "only":
			// A media-type qualifier; the type itself is beyond the canvas.
		case "and":
			// Conjunction: keep filling the current query.
		case "or":
			flush()
			cur = MediaQuery{Negated: negated}
		default:
			if t[0] == '(' {
				body := strings.TrimSuffix(strings.TrimPrefix(t, "("), ")")
				kv := strings.SplitN(body, ":", 2)
				if len(kv) != 2 {
					warns = append(warns, fmtErrf("ignoring media condition %q", t).Error())
					continue
				}
				key := strings.TrimSpace(kv[0])
				val := strings.TrimSpace(kv[1])
				n, ok := mediaPx(val)
				switch key {
				case "min-width":
					if ok {
						cur.MinWidth, cur.HasMin = n, true
					} else {
						warns = append(warns, fmtErrf("ignoring media width %q", val).Error())
					}
				case "max-width":
					if ok {
						cur.MaxWidth, cur.HasMax = n, true
					} else {
						warns = append(warns, fmtErrf("ignoring media width %q", val).Error())
					}
				default:
					warns = append(warns, fmtErrf("ignoring media condition %q", t).Error())
				}
			} else {
				switch t {
				case "screen", "all":
					// The window kind the canvas always is.
				default:
					warns = append(warns, fmtErrf("ignoring media type %q", t).Error())
				}
			}
		}
	}
	flush()
	return out, warns
}

// splitQueryTokens breaks a media query into its words and (…)-groups,
// keeping parenthesised groups — including nested parentheses — whole.
func splitQueryTokens(s string) []string {
	var out []string
	i := 0
	for i < len(s) {
		for i < len(s) && isSpace(s[i]) {
			i++
		}
		if i >= len(s) {
			break
		}
		start := i
		if s[i] == '(' {
			i = max(scanParen(s, i), i+1)
		} else {
			for i < len(s) && !isSpace(s[i]) && s[i] != '(' {
				i++
			}
		}
		if i == start {
			i++
			continue
		}
		out = append(out, s[start:i])
	}
	return out
}

// mediaPx reads a media width like "720px" or "720" into pixels. Fractions
// round; junk returns ok=false and the callers warn.
func mediaPx(s string) (int, bool) {
	t := strings.TrimSpace(strings.ToLower(strings.TrimSuffix(strings.TrimSpace(s), "px")))
	var f float64
	if _, err := fmt.Sscanf(t, "%g", &f); err != nil {
		return 0, false
	}
	return int(f + 0.5), true
}

// stripComments removes every /* … */ comment from a string. A comment in a
// file any old position — selector, value, media prelude — carries nothing, so
// it is just deleted.
func stripComments(s string) string {
	for {
		a := strings.Index(s, "/*")
		if a < 0 {
			return s
		}
		b := strings.Index(s[a+2:], "*/")
		if b < 0 {
			return s[:a]
		}
		end := a + 2 + b + 2
		s = s[:a] + s[end:]
	}
}

// parseDeclaration reads "property: value". A missing ':' is not a
// declaration. A trailing !important is recorded, not parsed as a value.
func parseDeclaration(text string) (Declaration, bool) {
	text = stripComments(text)
	idx := strings.IndexByte(text, ':')
	if idx < 0 {
		return Declaration{}, false
	}
	prop := strings.TrimSpace(text[:idx])
	val := strings.TrimSpace(text[idx+1:])
	if prop == "" {
		return Declaration{}, false
	}
	val, important := splitImportant(val)
	return Declaration{Prop: prop, Raw: val, Important: important}, true
}

// splitImportant peels a trailing "!important" (in any spacing) off a value.
func splitImportant(raw string) (val string, important bool) {
	v := strings.TrimSpace(raw)
	if i := strings.LastIndexByte(v, '!'); i >= 0 {
		rest := strings.TrimSpace(strings.ToLower(v[i+1:]))
		if rest == "important" {
			return strings.TrimSpace(v[:i]), true
		}
	}
	return v, false
}

// scanFloat reads a floating point number with fmt-equivalent strictness.
func scanFloat(s string, dst *float64) (int, error) {
	t := strings.TrimSuffix(s, "%")
	v, err := strconv.ParseFloat(strings.TrimSpace(t), 64)
	if err != nil {
		return 0, err
	}
	*dst = v
	if strings.HasSuffix(strings.TrimSpace(s), "%") {
		*dst = v / 100
	}
	return len(t), nil
}

// isSpace reports whether c is CSS whitespace.
func isSpace(c byte) bool {
	switch c {
	case ' ', '\t', '\n', '\r', '\f':
		return true
	}
	return false
}

func isIdentByte(c byte) bool {
	return c == '_' || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
}

func isHex(c byte) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}

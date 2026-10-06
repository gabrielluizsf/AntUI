package css

import (
	"fmt"
	"os"
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

	// path is the file this text came from, cleaned, and empty when it came
	// from no file at all — which is also what leaves a relative @import
	// with nowhere to resolve to.
	path string

	// seen is every file this parse has read or is reading. An @import of
	// one of them would read it a second time, and two files importing each
	// other would never end, so the whole tree of files shares one list.
	seen map[string]bool

	// depth is how many @imports deep this text was read — the sheet's own
	// file is 0. It bounds a chain in which every file is a different one
	// and the seen list alone would let it run on.
	depth int
}

// maxImportDepth is how many @imports may nest: one sheet reading another,
// that reading another, and so on. A chain longer than this is reported and
// left unread rather than followed to wherever it goes.
const maxImportDepth = 16

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
	if err := p.declBlock(func(text string) error {
		if d, ok := parseDeclaration(text); ok {
			rule.Decls = append(rule.Decls, d)
		}
		return nil
	}); err != nil {
		return err
	}
	sh.rules = append(sh.rules, rule)
	return nil
}

// declBlock reads the declarations inside a '{' … '}' block, one at a time,
// handing each to fn until the closing '}' — which the caller has not been
// given yet. It is the body both a rule and a @font-face share.
func (p *parser) declBlock(fn func(text string) error) error {
	for {
		p.skipSpace()
		if p.eof() {
			return fmtErrf("unterminated rule, missing '}'")
		}
		if p.peek() == '}' {
			p.next()
			return nil
		}
		text, err := p.readDecl()
		if err != nil {
			return err
		}
		if p.peek() == ';' {
			p.next()
		}
		if err := fn(text); err != nil {
			return err
		}
	}
}

// parseAll reads the whole of one stylesheet's text: rules and at-rules, one
// statement at a time, until the text runs out. m is the media the text sits
// inside, which for a file read on its own is none at all.
func (p *parser) parseAll(sh *Sheet, m Media) error {
	for {
		p.skipSpace()
		if p.eof() {
			return nil
		}
		switch {
		case p.peek() == '@':
			if e := p.parseAtRule(sh, m); e != nil {
				return e
			}
		case p.peek() == ';':
			p.next()
		default:
			if e := p.parseRule(sh, m); e != nil {
				return e
			}
		}
	}
}

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

	case "supports":
		// The condition tests the engine: the properties it has, the values
		// it reads, the selectors it can evaluate. A block whose test fails is
		// not part of the sheet — that is the test doing its job — and one
		// whose condition cannot be read is reported and dropped the same way.
		prelude, semi, err := p.readHeader(true)
		if err != nil {
			return err
		}
		if semi {
			// A @supports with no block gates nothing, and there is nothing
			// inside it to keep.
			if p.peek() == ';' {
				p.next()
			}
			return nil
		}
		ok, warns := supportsCondition(prelude)
		sh.Warn = append(sh.Warn, warns...)
		if !ok {
			return p.skipBlock()
		}
		if p.peek() == '{' {
			p.next()
		}
		return p.parseRuleGroup(sh, outer)

	case "layer", "container":
		// A cascade layer and a containment context are names for where a rule
		// sits in a sheet a browser folds; a canvas has neither to fold, so
		// the rules inside apply as if the name were the only one. A statement
		// form (@layer a, b;) drops it.
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

	case "import":
		return p.parseImport(sh, outer)

	case "font-face":
		return p.parseFontFace(sh)

	default:
		// Statement and declaration at-rules — @charset, @namespace,
		// @font-face, @page, @property, vendor and future ones — name things
		// no box can paint: fonts, faces, fallbacks, namespaces. Their
		// prelude and body are consumed, and the rule is reported.
		sh.Warn = append(sh.Warn, fmtErrf("ignoring @%s", word).Error())
		return p.skipAtRuleBody()
	}
}

// parseImport reads another stylesheet in where the @import stands, so its
// rules land in this sheet at that point, carrying the media the statement is
// inside and the media query the statement itself names. Every reason the file
// cannot be read is a warning rather than an error: a stylesheet that imports
// something this canvas cannot open still styles everything else in it.
func (p *parser) parseImport(sh *Sheet, outer Media) error {
	prelude, semi, err := p.readHeader(true)
	if err != nil {
		// The statement never ended, and it is the last thing left to read:
		// either the file was cut off before its ';', or what follows is
		// inside a string that never closes. Both end the sheet there.
		sh.Warn = append(sh.Warn, fmtErrf("ignoring an unfinished @import").Error())
		p.i = len(p.src)
		return nil
	}
	if !semi {
		// An @import has no block, but the one written all the same is
		// swallowed so the rest of the sheet still is.
		if p.peek() == '{' {
			p.next()
			if err := p.skipBlock(); err != nil {
				return err
			}
		}
		sh.Warn = append(sh.Warn, fmtErrf("ignoring an @import with a block").Error())
		return nil
	}
	if p.peek() == ';' {
		p.next()
	}

	ref, mediaText, warns := splitImport(prelude)
	sh.Warn = append(sh.Warn, warns...)
	if ref == "" {
		return nil
	}
	what := fmt.Sprintf("@import of %q", ref)
	path, warn := p.localFile(ref, what)
	if warn != "" {
		sh.Warn = append(sh.Warn, warn)
		return nil
	}
	if p.depth >= maxImportDepth {
		sh.Warn = append(sh.Warn, fmtErrf("ignoring %s: more than %d files deep", what, maxImportDepth).Error())
		return nil
	}
	if p.seen == nil {
		p.seen = make(map[string]bool)
	}
	if p.seen[path] {
		sh.Warn = append(sh.Warn, fmtErrf("ignoring %s: already read", what).Error())
		return nil
	}
	p.seen[path] = true
	data, err := os.ReadFile(path)
	if err != nil {
		sh.Warn = append(sh.Warn, fmtErrf("ignoring %s: %v", what, err).Error())
		return nil
	}
	m, mediaWarns := parseMedia(mediaText)
	sh.Warn = append(sh.Warn, mediaWarns...)
	sub := &parser{src: string(data), path: path, seen: p.seen, depth: p.depth + 1}
	return sub.parseAll(sh, outer.and(m))
}

// splitImport takes an @import prelude apart: the URL it names, the media
// query that gates the whole of it, and the layer() and supports() clauses
// between the two — which this engine does not fold, and reports rather than
// drop the import along with them.
func splitImport(prelude string) (ref, media string, warns []string) {
	s := strings.TrimSpace(prelude)
	switch {
	case len(s) >= 4 && strings.EqualFold(s[:4], "url("):
		end := parenEnd(s, 3)
		if end < 0 {
			return "", "", []string{fmtErrf("ignoring an @import whose url( never ends").Error()}
		}
		ref = unquote(strings.TrimSpace(s[4:end]))
		s = strings.TrimSpace(s[end+1:])
	case len(s) > 0 && (s[0] == '"' || s[0] == '\''):
		end := quoteEnd(s, 0)
		if end < 0 {
			return "", "", []string{fmtErrf("ignoring an @import with an unterminated url").Error()}
		}
		ref = s[1:end]
		s = strings.TrimSpace(s[end+1:])
	default:
		return "", "", []string{fmtErrf("ignoring an @import with no file name").Error()}
	}
	// A query and a fragment are for a server to answer; localFile strips
	// them when it comes to open the file.
	if strings.TrimSpace(ref) == "" {
		return "", "", []string{fmtErrf("ignoring an @import with no file name").Error()}
	}
	for s != "" {
		s = strings.TrimLeft(s, " \t\r\n")
		clause := clauseWord(s)
		if clause == "" {
			break
		}
		s = s[len(clause):]
		if strings.HasPrefix(s, "(") {
			end := parenEnd(s, 0)
			if end < 0 {
				warns = append(warns, fmtErrf("ignoring an @import whose %s( never ends", clause).Error())
				return ref, "", warns
			}
			s = s[end+1:]
		}
		warns = append(warns, fmtErrf("ignoring the %s() clause of @import", clause).Error())
	}
	return ref, strings.TrimSpace(s), warns
}

// clauseWord reads the layer or supports keyword an @import may carry after
// its URL, or "" when what comes next is something else — a media type, or
// the end of the statement. s is already trimmed.
func clauseWord(s string) string {
	for _, clause := range []string{"layer", "supports"} {
		if len(s) < len(clause) || !strings.EqualFold(s[:len(clause)], clause) {
			continue
		}
		if len(s) == len(clause) {
			return clause
		}
		switch s[len(clause)] {
		case '(', ' ', '\t', '\r', '\n':
			return clause
		}
	}
	return ""
}

// isRemoteRef reports whether an @import names something outside this
// machine: another host, a data: document, anything with a scheme behind it.
// Reading a stylesheet means reading a file beside this one; a canvas that
// went out to fetch one would be a network client as well as a renderer.
func isRemoteRef(ref string) bool {
	if strings.HasPrefix(ref, "//") {
		return true
	}
	// A colon before any '/' is a scheme, unless it is a single letter —
	// that is a drive root, and reads as a path like any other.
	if i := strings.Index(ref, ":"); i > 0 && i != 1 && !strings.ContainsAny(ref[:i], `/\`) {
		return true
	}
	return false
}

// unquote drops the quotes around a url() body, if it has any.
func unquote(s string) string {
	if len(s) >= 2 && (s[0] == '"' || s[0] == '\'') && s[len(s)-1] == s[0] {
		return s[1 : len(s)-1]
	}
	return s
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
// walked — not/only, and/or, comma-separated query lists — while the
// features the viewport answers constrain the rule (its two edges, the
// window's shape, the display's density and the system's color scheme);
// everything else warns and is treated as satisfied, so a canvas never drops
// a rule it merely cannot measure.
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
		if !cur.constrained() {
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
				// measure writes one edge of the window, or says the number
				// was not a number and leaves the edge unset.
				measure := func(dst *int, has *bool, edge string) {
					if ok {
						*dst, *has = n, true
						return
					}
					warns = append(warns, fmtErrf("ignoring media %s %q", edge, val).Error())
				}
				switch key {
				case "min-width":
					measure(&cur.MinWidth, &cur.HasMin, "width")
				case "max-width":
					measure(&cur.MaxWidth, &cur.HasMax, "width")
				case "min-height":
					measure(&cur.MinHeight, &cur.HasMinHeight, "height")
				case "max-height":
					measure(&cur.MaxHeight, &cur.HasMaxHeight, "height")
				case "orientation":
					switch val {
					case "portrait":
						cur.HasOrientation, cur.Portrait = true, true
					case "landscape":
						cur.HasOrientation, cur.Portrait = true, false
					default:
						warns = append(warns, fmtErrf("ignoring media orientation %q", val).Error())
					}
				case "resolution", "min-resolution", "max-resolution":
					dpi, ok := mediaDpi(val)
					if !ok {
						warns = append(warns, fmtErrf("ignoring media resolution %q", val).Error())
						continue
					}
					switch key {
					case "resolution":
						cur.MinDpi, cur.MaxDpi, cur.HasMinDpi, cur.HasMaxDpi = dpi, dpi, true, true
					case "min-resolution":
						cur.MinDpi, cur.HasMinDpi = dpi, true
					case "max-resolution":
						cur.MaxDpi, cur.HasMaxDpi = dpi, true
					}
				case "prefers-color-scheme":
					switch val {
					case "light":
						cur.HasScheme, cur.Dark = true, false
					case "dark":
						cur.HasScheme, cur.Dark = true, true
					default:
						warns = append(warns, fmtErrf("ignoring media color-scheme %q", val).Error())
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

// mediaDpi reads a resolution — 96dpi, 40dpcm, 2dppx — and answers it in dots
// per inch, the unit every query is compared in. A bare number, a value with
// no digits in it, or a unit the spec does not give a resolution returns
// ok=false and the caller warns.
func mediaDpi(s string) (float64, bool) {
	t := strings.TrimSpace(strings.ToLower(s))
	i := 0
	for i < len(t) && (t[i] == '.' || (t[i] >= '0' && t[i] <= '9')) {
		i++
	}
	v, err := strconv.ParseFloat(t[:i], 64)
	if err != nil {
		return 0, false
	}
	switch strings.TrimSpace(t[i:]) {
	case "dpi":
	case "dpcm":
		v *= 2.54
	case "dppx":
		v *= 96
	default:
		return 0, false
	}
	return roundDpi(v), true
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

package css

import "strings"

// supportsCondition reads the prelude of an @supports rule and answers it
// against the engine: whether it has the properties, the values and the
// selectors the condition tests. A condition the reader cannot make sense of
// is reported and answered false — the block it guards leaves the sheet, which
// is what a browser does with a test it cannot parse — while a test that
// simply fails is answered false in silence: that is the test working.
func supportsCondition(prelude string) (bool, []string) {
	r := &supportsReader{src: stripComments(prelude)}
	ok := r.condition()
	if r.bad || !r.atEnd() {
		return false, []string{fmtErrf("ignoring @supports condition %q", strings.TrimSpace(prelude)).Error()}
	}
	return ok, nil
}

// supportsReader walks an @supports prelude. A condition is an or-list of
// and-lists of negations, each of which is either a parenthesised group — a
// declaration when a colon splits it, a nested condition when it does not — or
// one of the condition's own functions.
type supportsReader struct {
	src string
	i   int
	bad bool
}

func (r *supportsReader) eof() bool { return r.i >= len(r.src) }

func (r *supportsReader) skipSpace() {
	for !r.eof() && isSpace(r.src[r.i]) {
		r.i++
	}
}

// atEnd reports the reader having consumed the whole prelude.
func (r *supportsReader) atEnd() bool {
	r.skipSpace()
	return r.eof()
}

// word takes w when it is the next thing, and leaves the reader past the
// spaces before it when it is not, so a word that does not match costs
// nothing. Keywords are read without case, the way the media grammar reads
// them.
func (r *supportsReader) word(w string) bool {
	r.skipSpace()
	rest := r.src[r.i:]
	if len(rest) < len(w) || !strings.EqualFold(rest[:len(w)], w) {
		return false
	}
	if r.i+len(w) < len(r.src) && isIdentByte(r.src[r.i+len(w)]) {
		return false // a longer word that merely starts with w
	}
	r.i += len(w)
	return true
}

// condition reads the or-list at the top. Both sides are always read, even
// when the left one already settled the answer, so a malformed right-hand side
// is still reported instead of left behind for atEnd to call trailing garbage.
func (r *supportsReader) condition() bool {
	ok := r.conjunction()
	for r.word("or") {
		if r.conjunction() {
			ok = true
		}
	}
	return ok
}

func (r *supportsReader) conjunction() bool {
	ok := r.negation()
	for r.word("and") {
		if !r.negation() {
			ok = false
		}
	}
	return ok
}

// negation reads a "not" — which binds to what follows it, tighter than and
// and or — and lets a second one in, so `not not (x)` is the same test as
// `(x)`.
func (r *supportsReader) negation() bool {
	if r.word("not") {
		return !r.negation()
	}
	return r.primary()
}

func (r *supportsReader) primary() bool {
	r.skipSpace()
	if r.eof() {
		r.bad = true
		return false
	}
	if r.src[r.i] == '(' {
		inner, ok := r.group()
		if !ok {
			r.bad = true
			return false
		}
		return r.inside(inner)
	}
	name, n := readName(r.src[r.i:])
	if name == "" {
		r.bad = true
		return false
	}
	r.i += n
	inner, ok := r.group()
	if !ok {
		r.bad = true
		return false
	}
	return r.function(name, inner)
}

// inside answers what a pair of parentheses holds: a declaration when a colon
// splits it at the top level, a nested condition when it does not.
func (r *supportsReader) inside(inner string) bool {
	if d, ok := splitSupportDecl(inner); ok {
		return supportsDecl(d.Prop, d.Raw)
	}
	sub := &supportsReader{src: inner}
	ok := sub.condition()
	if sub.bad || !sub.atEnd() {
		r.bad = true
		return false
	}
	return ok
}

// function answers one of the condition's own functions. Anything else is a
// function the reader has no rule for, which leaves the whole condition
// unreadable rather than guessed at.
func (r *supportsReader) function(name, inner string) bool {
	switch strings.ToLower(name) {
	case "selector":
		return supportsSelector(inner)
	case "font-format":
		// Only TrueType outlines are read from a file, under either name
		// the spec gives them; the woff containers and the CFF outlines of
		// an OpenType file are formats no @font-face here can open.
		f := unquote(strings.ToLower(strings.TrimSpace(inner)))
		return f == "truetype" || f == "ttf"
	case "font-tech":
		// One technology is promised: the colour bitmaps of a CBDT
		// table, which an @font-face reads and every frame draws. The
		// rest — the COLR, SVG and sbix colour formats, variations,
		// layout features — are not read, so no stylesheet is told they
		// are.
		f := unquote(strings.ToLower(strings.TrimSpace(inner)))
		return f == "color-cbdt"
	}
	r.bad = true
	return false
}

// group reads one parenthesised group and returns what is inside it, with the
// strings and the nesting inside it left alone.
func (r *supportsReader) group() (string, bool) {
	r.skipSpace()
	if r.eof() || r.src[r.i] != '(' {
		return "", false
	}
	depth := 0
	start := 0
	for !r.eof() {
		c := r.src[r.i]
		if c == '"' || c == '\'' {
			end := quoteEnd(r.src, r.i)
			if end < 0 {
				return "", false
			}
			r.i = end + 1
			continue
		}
		switch c {
		case '(':
			if depth == 0 {
				start = r.i + 1
			}
			depth++
		case ')':
			depth--
			if depth == 0 {
				r.i++
				return r.src[start : r.i-1], true
			}
		}
		r.i++
	}
	return "", false
}

// splitSupportDecl reads "property: value" when a colon at the top level of
// the group splits it there. Inside parens or quotes a colon belongs to what
// it is written in — url(http://x), "a:b", selector(a:hover) — and a group
// whose first colon is one of those is not a declaration at all.
func splitSupportDecl(text string) (Declaration, bool) {
	top := -1
	depth := 0
	for i := 0; i < len(text); i++ {
		c := text[i]
		switch {
		case c == '"' || c == '\'':
			end := quoteEnd(text, i)
			if end < 0 {
				return Declaration{}, false
			}
			i = end
		case c == '(':
			depth++
		case c == ')':
			if depth > 0 {
				depth--
			}
		case c == ':' && depth == 0:
			top = i
		}
		if top >= 0 {
			break
		}
	}
	if top < 0 || strings.IndexByte(text, ':') != top {
		return Declaration{}, false
	}
	return parseDeclaration(text)
}

// supportsDecl reports whether the engine would take a declaration: a
// property it has, holding a value it reads. Both halves come from applyDecl
// itself — the fold the cascade runs — so @supports can only claim what the
// engine will really do with the declaration.
func supportsDecl(prop, raw string) bool {
	prop = strings.ToLower(strings.TrimSpace(prop))
	raw = strings.TrimSpace(raw)
	if prop == "" || raw == "" {
		return false
	}
	if strings.HasPrefix(prop, "--") {
		// A custom property is the author's to declare: the engine holds it
		// for var() without being asked what it holds.
		return true
	}
	if strings.HasPrefix(prop, "-") {
		// A vendor prefix the engine has no case for, which applyDecl says
		// out loud whenever it meets one in a stylesheet.
		return false
	}
	if !supportedProperty(prop) {
		return false
	}
	_, set := foldSupport(prop, raw)
	return set[prop]
}

// supportedProperty reports whether the engine has the property at all. It
// asks with a value no grammar reads, where the only answer that can come
// back is about the property's own name: a value the engine cannot read is
// dropped in silence, and only a name it has no case for is reported.
func supportedProperty(prop string) bool {
	warns, _ := foldSupport(prop, unreadableValue)
	return len(warns) == 0
}

// unreadableValue is what the knowledge probe throws at a property: no
// grammar accepts it, so nothing but the property's name can answer.
const unreadableValue = "?"

// foldSupport runs one declaration through the cascade's own fold on a
// scratch style. There is no sheet, no inheritance and a fixed measuring
// context, so the answer depends on the engine and on nothing else — least of
// all on the window the next frame will draw at.
func foldSupport(prop, raw string) ([]string, map[string]bool) {
	var st Style
	st.Custom = make(map[string]string, 4)
	st.inherit = make(map[string]bool, 4)
	set := map[string]bool{}
	ctx := Units{Width: 1000, Height: 1000, Font: DefaultFontSize, Root: DefaultFontSize}
	warns := applyDecl(&st, set, st.Custom, nil, Declaration{Prop: prop, Raw: raw}, &ctx)
	return warns, set
}

// supportsSelector answers selector() with the engine's own selector reader:
// a selector it would keep is one it can evaluate, and one it marks
// unsupported — a combinator, an attribute, an id, a pseudo-element, a
// structural pseudo-class — is one it would drop from the sheet. A selector
// with nothing in it at all selects no widget the engine knows how to name,
// which is how parseSelector reads the empty one.
func supportsSelector(text string) bool {
	sels, _ := parseSelectors(text)
	if len(sels) == 0 {
		return false
	}
	for _, s := range sels {
		if s.Unsupported != "" {
			return false
		}
		if s.Tag == "" && !s.All && s.State == 0 && len(s.Classes) == 0 {
			return false
		}
	}
	return true
}

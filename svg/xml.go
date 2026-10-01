package svg

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// element is one tag of a document: what it is called, what it was given, what
// was written inside it, and what was written inside those. A drawing is a tree
// of these, and the root of it is the element that has no one above.
type element struct {
	// Name is the tag with any namespace taken off, so a document that spells
	// its tags `svg:circle` and one that spells them `circle` read the same.
	Name string
	// Attr holds what the tag was given, keys lowercased because attribute
	// names are not case sensitive.
	Attr map[string]string
	// Kids are the elements written inside this one, in the order they came.
	Kids []*element
	// Text is the writing between the tags, with the whitespace at the two ends
	// of it taken out, because whitespace around a drawing is there to make
	// it readable and says nothing about where it is. It is the whole of the
	// writing, put together; see Parts for the order it was written in.
	Text string
	// Parts are what was written inside the tag, in the order it was written:
	// the writing and the elements one after another. A `<text>` holds both,
	// and reading it as one string beside a list of children would put the
	// writing on either side of a `<tspan>` in the wrong order.
	Parts []part
}

// part is one piece of what was written inside a tag: either the writing
// between two elements, or an element. Exactly one of the two is there.
type part struct {
	Text string
	Elem *element
}

// attr reads one attribute, answering the empty string for one that was not
// there, so a caller can ask without checking twice.
//
// The name is matched without regard to case, because the names a drawing uses
// are written the way XML has them — `viewBox` and `preserveAspectRatio` are
// camel case, and everything else is not — and a drawing that spelled one of them
// in a case this package did not expect would otherwise be read as though it
// had said nothing at all.
func (e *element) attr(name string) string { return e.Attr[strings.ToLower(name)] }

// hasAttr reports whether the attribute was written at all, which is a
// different question from whether it was written empty: `fill=""` is not the
// same as no fill, and only the first is an instruction.
func (e *element) hasAttr(name string) bool {
	_, ok := e.Attr[strings.ToLower(name)]
	return ok
}

// find returns the first descendant with the given tag name, or nil. It is how
// a `<defs>` or a `<style>` block is reached without walking the tree by hand.
func (e *element) find(name string) *element {
	if e.Name == name {
		return e
	}
	for _, k := range e.Kids {
		if got := k.find(name); got != nil {
			return got
		}
	}
	return nil
}

// parseDocument reads a whole XML document and answers its root element. An
// error is only ever structural — a tag that never closes, a string left open —
// because a drawing that does not close what it opens cannot be read at all,
// while anything else inside it is either drawn or reported as a warning.
func parseDocument(src string) (*element, error) {
	p := &xmlParser{src: src}
	return p.document()
}

// xmlParser walks the text of a document once, handing back what it finds. It
// is not a general XML parser: the parts of the language a drawing never uses,
// entities beyond the five XML defines and character references, DTDs and
// processing instructions, are read past rather than read into.
type xmlParser struct {
	src string
	pos int
}

// document reads from the opening angle bracket to the end of the root element.
func (p *xmlParser) document() (*element, error) {
	for {
		p.skipSpace()
		if p.at("<?") {
			p.through("?>")
			continue
		}
		if p.at("<!--") {
			p.through("-->")
			continue
		}
		if p.at("<!") {
			// A doctype or another declaration: it ends at the first > that is
			// not inside a pair of brackets, which is where the internal subset
			// of a DTD would otherwise run past it.
			p.skipDecl()
			continue
		}
		break
	}
	if p.eof() {
		return nil, fmt.Errorf("antui/svg: the drawing is empty, there is no <svg> in it")
	}
	if !p.at("<") {
		return nil, fmt.Errorf("antui/svg: the drawing starts with %q, where a tag was expected", p.peek(12))
	}
	root, err := p.element()
	if err != nil {
		return nil, err
	}
	return root, nil
}

// element reads one tag and everything inside it, up to its own closing tag.
func (p *xmlParser) element() (*element, error) {
	p.through("<")
	start := p.pos
	for !p.eof() && !isSpace(p.src[p.pos]) && p.src[p.pos] != '>' && p.src[p.pos] != '/' {
		p.pos++
	}
	name := p.src[start:p.pos]
	if name == "" {
		return nil, fmt.Errorf("antui/svg: a tag has no name at offset %d", start)
	}
	e := &element{Name: localName(name), Attr: map[string]string{}}
	if err := p.attributes(e); err != nil {
		return nil, err
	}
	if e.Attr[""] != "" {
		// A lone slash before the closing bracket is the older way of writing a
		// tag that closes itself, and it means the same thing.
		return e, nil
	}
	if err := p.children(e); err != nil {
		return nil, err
	}
	return e, nil
}

// attributes reads everything up to the closing bracket of the tag. A value may
// be in single quotes, double quotes, or no quotes at all, as XML allows the
// last when it has no space in it.
func (p *xmlParser) attributes(e *element) error {
	for {
		p.skipSpace()
		if p.eof() {
			return fmt.Errorf("antui/svg: <%s> is never closed", e.Name)
		}
		if p.at("/>") {
			p.pos += 2
			e.Attr[""] = "/"
			return nil
		}
		if p.at(">") {
			p.pos++
			return nil
		}
		name, ok := p.attrName()
		if !ok {
			// A slash with nothing before it is a tag that closes itself,
			// written the older way as `<rect/>` with no space in between.
			if p.at("/>") {
				p.pos += 2
				e.Attr[""] = "/"
				return nil
			}
			return fmt.Errorf("antui/svg: <%s> has something in it that is not an attribute", e.Name)
		}
		p.skipSpace()
		if !p.at("=") {
			// An attribute written on its own has the value its own name, which
			// is how `hidden` and `required` are written in a drawing.
			e.Attr[strings.ToLower(name)] = name
			continue
		}
		p.pos++
		p.skipSpace()
		v, err := p.attrValue()
		if err != nil {
			return fmt.Errorf("antui/svg: <%s %s>: %w", e.Name, name, err)
		}
		e.Attr[strings.ToLower(name)] = v
	}
}

// attrName reads the name of one attribute, stopping at the equals sign or the
// space after it.
func (p *xmlParser) attrName() (string, bool) {
	start := p.pos
	for !p.eof() {
		c := p.src[p.pos]
		if c == '=' || isSpace(c) || c == '>' || c == '/' {
			break
		}
		p.pos++
	}
	if p.pos == start {
		return "", false
	}
	return p.src[start:p.pos], true
}

// attrValue reads the value of one attribute, unwrapping its quotes and turning
// what is inside them into the characters it stands for.
func (p *xmlParser) attrValue() (string, error) {
	if p.eof() {
		return "", fmt.Errorf("the value is missing")
	}
	quote := p.src[p.pos]
	if quote != '"' && quote != '\'' {
		// An unquoted value runs to the first space, as XML has it — and not
		// into the slash of a `/>`, which closes the tag rather than being part
		// of what is written on it.
		start := p.pos
		for !p.eof() && !isSpace(p.src[p.pos]) && p.src[p.pos] != '>' {
			if p.src[p.pos] == '/' && p.pos+1 < len(p.src) && p.src[p.pos+1] == '>' {
				break
			}
			p.pos++
		}
		return unescape(p.src[start:p.pos]), nil
	}
	p.pos++
	start := p.pos
	for !p.eof() && p.src[p.pos] != quote {
		p.pos++
	}
	if p.eof() {
		return "", fmt.Errorf("the %c that opened the value is never closed", quote)
	}
	v := p.src[start:p.pos]
	p.pos++
	return unescape(v), nil
}

// children reads what is written inside a tag until the tag that closes it.
// Text and elements are kept apart, and the whitespace around a child element
// is left out of the text, so an element that only holds formatting whitespace
// holds nothing. Each piece of writing is also kept where it was written,
// because the writing of a `<text>` runs on across the elements between it.
func (p *xmlParser) children(e *element) error {
	var text strings.Builder
	keep := func() {
		if s := collapseSpace(text.String()); s != "" {
			e.Parts = append(e.Parts, part{Text: s})
		}
		text.Reset()
	}
	for {
		if p.eof() {
			return fmt.Errorf("antui/svg: <%s> is never closed", e.Name)
		}
		if p.at("</") {
			p.pos += 2
			start := p.pos
			for !p.eof() && p.src[p.pos] != '>' {
				p.pos++
			}
			if p.eof() {
				return fmt.Errorf("antui/svg: <%s> is never closed", e.Name)
			}
			// The tag that closes this one has to be the one this one is: a
			// drawing that opens a group and closes it with something else has
			// no shape to paint, and reading it as though it were right would
			// draw something other than what the file says.
			if closing := localName(p.src[start:p.pos]); closing != e.Name {
				return fmt.Errorf("antui/svg: <%s> is closed by </%s>", e.Name, closing)
			}
			p.pos++
			keep()
			e.Text = textOf(e)
			return nil
		}
		if p.at("<!--") {
			p.through("-->")
			continue
		}
		if p.at("<![CDATA[") {
			p.pos += len("<![CDATA[")
			start := p.pos
			// The section ends where `]]>` begins, so the text is what is
			// between the two rather than up to the end of them.
			if i := strings.Index(p.src[p.pos:], "]]>"); i >= 0 {
				text.WriteString(p.src[start : start+i])
				p.pos += i + len("]]>")
			} else {
				text.WriteString(p.src[start:])
				p.pos = len(p.src)
			}
			continue
		}
		if p.at("<?") {
			p.through("?>")
			continue
		}
		if p.at("<!") {
			p.skipDecl()
			continue
		}
		if p.at("<") {
			keep()
			kid, err := p.element()
			if err != nil {
				return err
			}
			e.Kids = append(e.Kids, kid)
			e.Parts = append(e.Parts, part{Elem: kid})
			continue
		}
		text.WriteByte(p.src[p.pos])
		p.pos++
	}
}

// textOf is all the writing inside an element, put together, with the
// whitespace at the two ends of it taken off, because whitespace around a
// drawing is there to make it readable and says nothing about where it is.
func textOf(e *element) string {
	var out strings.Builder
	for _, p := range e.Parts {
		out.WriteString(p.Text)
	}
	return strings.TrimSpace(out.String())
}

// skipSpace passes over whitespace and nothing else.
func (p *xmlParser) skipSpace() {
	for !p.eof() && isSpace(p.src[p.pos]) {
		p.pos++
	}
}

// through moves past the next occurrence of what, and past the end of the
// document if it is not there, which is what a comment or a CDATA section that
// was never finished wants: the rest of the file is not worth reading.
func (p *xmlParser) through(what string) {
	if i := strings.Index(p.src[p.pos:], what); i >= 0 {
		p.pos += i + len(what)
		return
	}
	p.pos = len(p.src)
}

// skipDecl passes over a `<!…>` declaration, counting brackets so that the
// internal subset of a DTD, which can hold a `>` of its own, does not end it
// early.
func (p *xmlParser) skipDecl() {
	depth := 0
	for !p.eof() {
		switch p.src[p.pos] {
		case '[':
			depth++
		case ']':
			depth--
		case '>':
			if depth <= 0 {
				p.pos++
				return
			}
		}
		p.pos++
	}
}

// at reports whether the text at the cursor starts with what.
func (p *xmlParser) at(what string) bool { return strings.HasPrefix(p.src[p.pos:], what) }

// eof reports whether the cursor has run off the end.
func (p *xmlParser) eof() bool { return p.pos >= len(p.src) }

// peek is a short look at the text ahead, for a message about what was found
// where a tag was expected.
func (p *xmlParser) peek(n int) string {
	end := min(p.pos+n, len(p.src))
	return p.src[p.pos:end]
}

// localName takes the namespace off a tag: the `svg:` of `svg:circle` and the
// whole of `{http://www.w3.org/2000/svg}circle` both come back as `circle`.
func localName(name string) string {
	if i := strings.LastIndexByte(name, ':'); i >= 0 {
		return name[i+1:]
	}
	if i := strings.IndexByte(name, '}'); i >= 0 {
		return name[i+1:]
	}
	return name
}

// isSpace reports whether c is one of the characters XML counts as whitespace.
func isSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}

// collapseSpace turns every run of whitespace into a single space, so that text
// wrapped across lines in the file is the same text as text written on one.
// The spaces at the two ends are kept: a space at the edge of a piece of
// writing is what sits between it and the piece written next to it, and only
// whoever reads all of it together knows whether the writing wants it, which
// is what taking them off the whole of it in [textOf] does.
func collapseSpace(s string) string {
	var out strings.Builder
	out.Grow(len(s))
	space := false
	for i := 0; i < len(s); {
		c := s[i]
		if isSpace(c) {
			space = true
			i++
			continue
		}
		if space {
			out.WriteByte(' ')
		}
		space = false
		_, size := utf8.DecodeRuneInString(s[i:])
		out.WriteString(s[i : i+size])
		i += size
	}
	if space {
		out.WriteByte(' ')
	}
	return out.String()
}

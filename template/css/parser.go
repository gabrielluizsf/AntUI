package css

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

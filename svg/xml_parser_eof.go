package svg

// eof reports whether the cursor has run off the end.
func (p *xmlParser) eof() bool { return p.pos >= len(p.src) }

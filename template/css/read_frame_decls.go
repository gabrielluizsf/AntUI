package css

// readFrameDecls reads one keyframe's declaration block, returning the
// declarations it held. A missing '}' ends the block and the parse, exactly
// as a rule that never closes would.
func (p *parser) readFrameDecls() []Declaration {
	var decls []Declaration
	for {
		p.skipSpace()
		if p.eof() {
			return decls
		}
		if p.peek() == '}' {
			p.next()
			return decls
		}
		text, err := p.readDecl()
		if err != nil {
			return decls
		}
		if p.peek() == ';' {
			p.next()
		}
		if d, ok := parseDeclaration(text); ok {
			decls = append(decls, d)
		}
	}
}

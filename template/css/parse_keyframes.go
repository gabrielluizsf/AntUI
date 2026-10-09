package css

import (
	"strings"
)

// parseKeyframes reads the body of an @keyframes block: a series of frames
// introduced by from, to or a percentage followed by a declaration block.
// A frame with an unrecognised main is consumed and skipped, and a malformed
// body ends the rule the way any structural damage does.
func (p *parser) parseKeyframes(sh *Sheet, name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		// A nameless @keyframes is nothing a rule can name; the block is
		// slurped and dropped.
		return p.skipBlock()
	}
	kf := &Keyframes{Name: name}
	for {
		p.skipSpace()
		if p.eof() {
			return fmtErrf("unterminated @keyframes %q, missing '}'", name)
		}
		if p.peek() == '}' {
			p.next()
			sh.addKeyframes(kf)
			return nil
		}
		header, semi, err := p.readHeader(true)
		if err != nil {
			return err
		}
		off, ok := parseKeyframeOffset(header)
		if !ok {
			if semi {
				if p.peek() == ';' {
					p.next()
				}
			} else if p.peek() == '{' {
				if err := p.skipBlock(); err != nil {
					return err
				}
			}
			continue
		}
		if p.peek() == '{' {
			p.next()
		}
		kf.Frames = append(kf.Frames, Keyframe{Offset: off, Decls: p.readFrameDecls()})
	}
}

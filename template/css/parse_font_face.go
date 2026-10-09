package css

import (
	"strings"
)

// parseFontFace reads a @font-face block: the family the file answers to,
// which file it is, and the weight and slant it was cut for. Everything the
// block says that this engine cannot use is reported and passed over, and a
// file that will not open leaves the family with no face at all rather than
// failing the stylesheet that asked for it.
func (p *parser) parseFontFace(sh *Sheet) error {
	prelude, _, err := p.readHeader(false)
	if err != nil {
		return err
	}
	if p.peek() == '{' {
		p.next()
	}
	if s := strings.TrimSpace(prelude); s != "" {
		sh.Warn = append(sh.Warn, fmtErrf("ignoring %q before a @font-face block", s).Error())
	}

	family, src, weight, style, err := p.faceDecls(sh)
	if err != nil {
		return err
	}

	family = normalizeFamily(family)
	if family == "" {
		sh.Warn = append(sh.Warn, fmtErrf("ignoring a @font-face with no font-family").Error())
		return nil
	}
	if strings.TrimSpace(src) == "" {
		sh.Warn = append(sh.Warn, fmtErrf("ignoring the @font-face for %q: no src", family).Error())
		return nil
	}

	ff := &FontFace{Family: family, Weight: FontWeightNormal}
	applyFaceWeightStyle(sh, ff, weight, style, family)
	if ff.face = p.fontSource(sh, family, src); ff.face == nil {
		return nil
	}
	sh.fonts = append(sh.fonts, ff)
	sh.fontMu.Lock()
	sh.fontCache = nil
	sh.fontMu.Unlock()
	return nil
}
